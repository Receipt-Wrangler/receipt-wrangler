package services

import (
	"archive/zip"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/gabriel-vasile/mimetype"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"

	"receipt-wrangler/api/internal/constants"
	"receipt-wrangler/api/internal/logging"
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/reporting"
	"receipt-wrangler/api/internal/reporting/render"
	"receipt-wrangler/api/internal/repositories"
	"receipt-wrangler/api/internal/structs"
	"receipt-wrangler/api/internal/utils"
)

// The claim packet (GET /claim/{claimId}/packet): one zip a finance clerk can
// read without this app, holding packet.pdf, images/ (the untouched originals),
// receipts.csv, items.csv and manifest.json (every other file's SHA-256).
//
// It covers exactly what the caller may see of the claim: the receipts
// ClaimService already filters (paid-by visibility, member isolation), with the
// same category/tag stripping and member masking GET /receipt applies.

// The bounds that keep one packet inside a 1 GB host. A packet over either is a
// 400 before anything heavy runs.
const (
	// maxClaimPacketReceipts bounds the receipts sections in the PDF.
	maxClaimPacketReceipts = 500
	// MaxClaimPacketImages bounds the image files in a packet. Each one is read,
	// rasterized (a PDF at the PdfDpi setting) and downscaled in turn, and the
	// downscaled copies all sit in the HTML Chromium renders, so this is what
	// bounds both the request's duration and the render's memory.
	MaxClaimPacketImages = 60
	// maxClaimPacketEmbeddedBytes bounds the downscaled JPEG bytes embedded in
	// the PDF's HTML (base64 adds a third). Typical receipt photos come to
	// 150-400 KB each, so this only trips on unusually large or many-paged images.
	maxClaimPacketEmbeddedBytes = 24 << 20
)

// claimPacketFormatVersion is manifest.json's formatVersion.
const claimPacketFormatVersion = 1

// ClaimPacketTooLargeError is a packet over one of the bounds above; handlers
// answer it 400 with its message.
type ClaimPacketTooLargeError struct {
	Message string
}

func (err ClaimPacketTooLargeError) Error() string {
	return err.Message
}

// claimPacketSlot lets one packet build at a time in this process. Two at once
// would double the peak (two rasterizations, two Chromium renders), which a 1 GB
// host cannot afford; a second request waits rather than fails.
var claimPacketSlot = make(chan struct{}, 1)

// claimPacketPdfRenderer turns the packet's HTML into PDF bytes. A variable so
// tests can replace Chromium (SetClaimPacketPdfRendererForTests).
var claimPacketPdfRenderer = func(html string) ([]byte, error) {
	pdf, _, err := NewHtmlToPdfService(nil).Render(html)
	return pdf, err
}

// SetClaimPacketPdfRendererForTests replaces the HTML-to-PDF step and returns a
// function restoring the real one. A test renderer can return the HTML itself, so
// a test can assert on what the PDF would print without Chromium.
func SetClaimPacketPdfRendererForTests(renderer func(html string) ([]byte, error)) func() {
	previous := claimPacketPdfRenderer
	claimPacketPdfRenderer = renderer
	return func() { claimPacketPdfRenderer = previous }
}

// ClaimPacket is a built packet: a zip in temp/ the caller streams and removes.
type ClaimPacket struct {
	Path     string
	Filename string
	Size     int64
}

// Remove deletes the packet's temp file. A flat file in temp/, so the hourly
// sweeper reclaims it if the process dies first.
func (packet ClaimPacket) Remove() {
	if packet.Path != "" {
		os.Remove(packet.Path)
	}
}

type ClaimPacketService struct {
	BaseService
}

func NewClaimPacketService(tx *gorm.DB) ClaimPacketService {
	return ClaimPacketService{BaseService: BaseService{
		DB: repositories.GetDB(),
		TX: tx,
	}}
}

// claimPacketManifest is manifest.json.
type claimPacketManifest struct {
	FormatVersion int                       `json:"formatVersion"`
	Note          string                    `json:"note"`
	GeneratedAt   string                    `json:"generatedAt"`
	GeneratedBy   string                    `json:"generatedBy"`
	TimeZone      string                    `json:"timeZone"`
	Claim         structs.ClaimView         `json:"claim"`
	ReceiptIds    []uint                    `json:"receiptIds"`
	Files         []claimPacketManifestFile `json:"files"`
	MissingImages []claimPacketMissingImage `json:"missingImages"`
}

type claimPacketManifestFile struct {
	Path         string `json:"path"`
	Size         int64  `json:"size"`
	Sha256       string `json:"sha256"`
	ReceiptId    uint   `json:"receiptId,omitempty"`
	OriginalName string `json:"originalName,omitempty"`
}

// claimPacketMissingImage is an image the receipt records but whose file is not
// on the server, so it could not be included.
type claimPacketMissingImage struct {
	ReceiptId    uint   `json:"receiptId"`
	FileId       uint   `json:"fileId"`
	OriginalName string `json:"originalName"`
}

const claimPacketManifestNote = "sha256 is computed over each file's bytes exactly as stored in this zip. " +
	"manifest.json does not list itself. images/ holds the original files byte for byte as they are " +
	"stored on the server; the copies printed in packet.pdf are downscaled."

// BuildClaimPacket builds the packet for one claim the caller may see, as a zip in
// temp/. ErrClaimNotFound for a claim they may not see, ClaimPacketTooLargeError
// over a bound. The caller streams the file and calls Remove.
func (service ClaimPacketService) BuildClaimPacket(ctx context.Context, userId uint, claimId uint) (ClaimPacket, error) {
	claim, err := NewClaimService(service.TX).GetClaimForUser(userId, claimId)
	if err != nil {
		return ClaimPacket{}, err
	}

	receipts, err := service.loadPacketReceipts(userId, claim)
	if err != nil {
		return ClaimPacket{}, err
	}
	if err := checkClaimPacketBounds(receipts); err != nil {
		return ClaimPacket{}, err
	}

	user, err := repositories.NewUserRepository(service.TX).GetUserById(userId)
	if err != nil {
		return ClaimPacket{}, err
	}
	settings, err := repositories.NewSystemSettingsRepository(service.TX).GetSystemSettings()
	if err != nil {
		return ClaimPacket{}, err
	}
	location := repositories.ResolveAppLocation(settings.TimeZone)
	now := time.Now().In(location)

	select {
	case claimPacketSlot <- struct{}{}:
		defer func() { <-claimPacketSlot }()
	case <-ctx.Done():
		return ClaimPacket{}, ctx.Err()
	}

	path, err := repositories.NewFileRepository(service.TX).BuildTempFilePath("zip")
	if err != nil {
		return ClaimPacket{}, err
	}
	file, err := os.Create(path)
	if err != nil {
		return ClaimPacket{}, err
	}
	built := false
	defer func() {
		if !built {
			file.Close()
			os.Remove(path)
		}
	}()

	builder := claimPacketBuilder{
		service:   service,
		claim:     claim,
		receipts:  receipts,
		location:  location,
		now:       now,
		username:  user.Username,
		currency:  packetCurrencyFormat(settings),
		zip:       newPacketZip(file, now),
		manifest:  claimPacketManifest{Files: []claimPacketManifestFile{}, MissingImages: []claimPacketMissingImage{}},
		imageName: make(map[uint][]string, len(receipts)),
	}
	if err := builder.build(); err != nil {
		return ClaimPacket{}, err
	}
	if err := file.Close(); err != nil {
		return ClaimPacket{}, err
	}

	info, err := os.Stat(path)
	if err != nil {
		return ClaimPacket{}, err
	}
	built = true
	return ClaimPacket{
		Path:     path,
		Filename: claimPacketFilename(claim.Name, now),
		Size:     info.Size(),
	}, nil
}

// loadPacketReceipts loads the claim's visible receipts in full and applies what
// GET /receipt applies before anything is printed: category/tag grants and member
// masking. The claim view's ids are already the visible ones; the group and
// paid-by checks are repeated against the loaded rows so a receipt that moved or
// changed payer since cannot slip in. Sorted by date, then id.
func (service ClaimPacketService) loadPacketReceipts(userId uint, claim structs.ClaimView) ([]models.Receipt, error) {
	if len(claim.ReceiptIds) == 0 {
		return []models.Receipt{}, nil
	}

	ids := make([]string, 0, len(claim.ReceiptIds))
	for _, id := range claim.ReceiptIds {
		ids = append(ids, utils.UintToString(id))
	}

	associations := append(append([]string{}, constants.FULL_RECEIPT_ASSOCIATIONS...),
		"Categories", "Tags", "Comments.User", "ReceiptItems.ChargedToUser")
	loaded, err := repositories.NewReceiptRepository(service.TX).GetReceiptsByIds(ids, associations)
	if err != nil {
		return nil, err
	}

	inGroup := make([]models.Receipt, 0, len(loaded))
	for _, receipt := range loaded {
		if receipt.GroupId == claim.GroupId {
			inGroup = append(inGroup, receipt)
		}
	}

	permissionService := NewPermissionService(service.TX)
	receipts, err := permissionService.FilterReceiptsByPaidBy(userId, inGroup)
	if err != nil {
		return nil, err
	}
	if err := permissionService.FilterReceiptCategoriesTags(userId, receipts); err != nil {
		return nil, err
	}
	if err := permissionService.MaskReceiptsForMemberVisibility(userId, receipts); err != nil {
		return nil, err
	}

	sort.SliceStable(receipts, func(i, j int) bool {
		if !receipts[i].Date.Equal(receipts[j].Date) {
			return receipts[i].Date.Before(receipts[j].Date)
		}
		return receipts[i].ID < receipts[j].ID
	})
	for i := range receipts {
		receipt := &receipts[i]
		sort.SliceStable(receipt.ImageFiles, func(a, b int) bool { return receipt.ImageFiles[a].ID < receipt.ImageFiles[b].ID })
		sort.SliceStable(receipt.Comments, func(a, b int) bool {
			if !receipt.Comments[a].CreatedAt.Equal(receipt.Comments[b].CreatedAt) {
				return receipt.Comments[a].CreatedAt.Before(receipt.Comments[b].CreatedAt)
			}
			return receipt.Comments[a].ID < receipt.Comments[b].ID
		})
		// items.csv names each item's receipt.
		for j := range receipt.ReceiptItems {
			receipt.ReceiptItems[j].Receipt = models.Receipt{Name: receipt.Name, Date: receipt.Date}
		}
	}

	return receipts, nil
}

func checkClaimPacketBounds(receipts []models.Receipt) error {
	if len(receipts) > maxClaimPacketReceipts {
		return ClaimPacketTooLargeError{Message: fmt.Sprintf(
			"This claim has %d receipts; a packet can include at most %d", len(receipts), maxClaimPacketReceipts)}
	}

	images := 0
	for _, receipt := range receipts {
		images += len(receipt.ImageFiles)
	}
	if images > MaxClaimPacketImages {
		return ClaimPacketTooLargeError{Message: fmt.Sprintf(
			"This claim has %d images; a packet can include at most %d", images, MaxClaimPacketImages)}
	}
	return nil
}

// claimPacketBuilder writes one packet's zip. Its zip entries go out in this
// order: images/, receipts.csv, items.csv, packet.pdf, manifest.json. The images
// come first because each original is read once, hashed into the zip and
// downscaled for the PDF in the same pass, so no two originals are ever in
// memory together and the PDF can print each image's hash.
type claimPacketBuilder struct {
	service   ClaimPacketService
	claim     structs.ClaimView
	receipts  []models.Receipt
	location  *time.Location
	now       time.Time
	username  string
	currency  reporting.CurrencyFormat
	zip       *packetZip
	manifest  claimPacketManifest
	imageName map[uint][]string

	embeddedBytes int
}

func (builder *claimPacketBuilder) build() error {
	sections := make([]packetReceiptSection, 0, len(builder.receipts))
	for _, receipt := range builder.receipts {
		section, err := builder.receiptSection(receipt)
		if err != nil {
			return err
		}
		sections = append(sections, section)
	}

	if err := builder.writeCsvs(); err != nil {
		return err
	}

	document := builder.document(sections)
	html, err := renderClaimPacketHtml(document, builder.embeddedBytes)
	if err != nil {
		return err
	}
	pdf, err := claimPacketPdfRenderer(html)
	if err != nil {
		return err
	}
	if _, err := builder.zip.add("packet.pdf", pdf, zip.Deflate); err != nil {
		return err
	}

	return builder.writeManifest()
}

// receiptSection reads, stores and embeds each of the receipt's images, and
// lays out its section of the PDF.
func (builder *claimPacketBuilder) receiptSection(receipt models.Receipt) (packetReceiptSection, error) {
	section := buildPacketReceiptSection(receipt, builder.location, builder.money)

	names := make([]string, 0, len(receipt.ImageFiles))
	for index, file := range receipt.ImageFiles {
		name := packetImageName(receipt, index+1, file)
		image := packetImage{Name: name, OriginalName: file.Name}

		original, err := builder.readOriginal(file)
		if errors.Is(err, fs.ErrNotExist) {
			logging.LogStd(logging.LOG_LEVEL_ERROR, fmt.Sprintf("claim packet: image %d of receipt %d is missing on disk", file.ID, receipt.ID))
			image.Note = "The file for this image is missing on the server, so it is not in this packet."
			builder.manifest.MissingImages = append(builder.manifest.MissingImages, claimPacketMissingImage{
				ReceiptId: receipt.ID, FileId: file.ID, OriginalName: file.Name,
			})
			section.Images = append(section.Images, image)
			continue
		}
		if err != nil {
			return packetReceiptSection{}, err
		}

		entry, err := builder.zip.add("images/"+name, original, zip.Store)
		if err != nil {
			return packetReceiptSection{}, err
		}
		builder.zip.annotateLast(receipt.ID, file.Name)
		image.Sha256 = entry.Sha256
		names = append(names, name)

		slices, err := builder.displaySlices(original)
		if err != nil {
			logging.LogStd(logging.LOG_LEVEL_ERROR, fmt.Sprintf("claim packet: image %d of receipt %d could not be rendered: %s", file.ID, receipt.ID, err.Error()))
			image.Note = "This image could not be shown here; the original is included as images/" + name + "."
		}
		for _, slice := range slices {
			builder.embeddedBytes += len(slice)
			if builder.embeddedBytes > maxClaimPacketEmbeddedBytes {
				return packetReceiptSection{}, ClaimPacketTooLargeError{Message: fmt.Sprintf(
					"This claim's images are too large to print together in one packet (over %d MB once downscaled); split the claim",
					maxClaimPacketEmbeddedBytes>>20)}
			}
			image.Slices = append(image.Slices, template.URL(utils.BuildDataURI("image/jpeg", slice)))
		}
		section.Images = append(section.Images, image)
	}

	builder.imageName[receipt.ID] = names
	section.Row.Images = names
	return section, nil
}

// readOriginal reads a stored image through the data-directory helpers. A missing
// file is fs.ErrNotExist; anything else fails the packet.
func (builder *claimPacketBuilder) readOriginal(file models.FileData) ([]byte, error) {
	path, err := repositories.NewFileRepository(builder.service.TX).BuildFilePath(
		utils.UintToString(file.ReceiptId), utils.UintToString(file.ID), file.Name)
	if err != nil {
		return nil, err
	}
	return utils.ReadDataFile(path)
}

// displaySlices is the receipt image view's conversion (a PDF rasterizes, HEIC
// transcodes; see FileRepository.BuildDisplayImageString) followed by the
// packet's downscale.
func (builder *claimPacketBuilder) displaySlices(original []byte) ([][]byte, error) {
	displayable, err := repositories.NewFileRepository(builder.service.TX).GetBytesFromImageBytes(original)
	if err != nil {
		return nil, err
	}
	return packetDisplaySlices(displayable)
}

func (builder *claimPacketBuilder) writeCsvs() error {
	csvService := NewReceiptCsvService()
	csvService.Location = builder.location
	csvService.ImageNames = builder.imageName
	result, err := csvService.BuildReceiptCsv(builder.receipts)
	if err != nil {
		return err
	}
	if _, err := builder.zip.add("receipts.csv", result.ReceiptCsvBytes, zip.Deflate); err != nil {
		return err
	}
	_, err = builder.zip.add("items.csv", result.ReceiptItemCsvBytes, zip.Deflate)
	return err
}

func (builder *claimPacketBuilder) writeManifest() error {
	receiptIds := make([]uint, 0, len(builder.receipts))
	for _, receipt := range builder.receipts {
		receiptIds = append(receiptIds, receipt.ID)
	}

	manifest := builder.manifest
	manifest.FormatVersion = claimPacketFormatVersion
	manifest.Note = claimPacketManifestNote
	manifest.GeneratedAt = builder.now.Format(time.RFC3339)
	manifest.GeneratedBy = builder.username
	manifest.TimeZone = builder.location.String()
	manifest.Claim = builder.claim
	manifest.ReceiptIds = receiptIds
	manifest.Files = builder.zip.files

	bytes, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		return err
	}
	if _, err := builder.zip.add("manifest.json", bytes, zip.Deflate); err != nil {
		return err
	}
	return builder.zip.writer.Close()
}

func (builder *claimPacketBuilder) money(value decimal.Decimal) string {
	return render.FormatCurrency(value, builder.currency)
}

// packetCurrencyFormat is the app's currency display, always with cents: a
// packet is evidence, and its subtotals have to visibly add up.
func packetCurrencyFormat(settings models.SystemSettings) reporting.CurrencyFormat {
	format := *currencyFormat(settings)
	format.HideDecimals = false
	return format
}

// packetZip writes zip entries and records each one's size and SHA-256 over the
// bytes written, for the manifest.
type packetZip struct {
	writer   *zip.Writer
	modified time.Time
	files    []claimPacketManifestFile
}

func newPacketZip(destination io.Writer, modified time.Time) *packetZip {
	return &packetZip{writer: zip.NewWriter(destination), modified: modified, files: []claimPacketManifestFile{}}
}

func (packet *packetZip) add(path string, content []byte, method uint16) (claimPacketManifestFile, error) {
	entry, err := packet.writer.CreateHeader(&zip.FileHeader{Name: path, Method: method, Modified: packet.modified})
	if err != nil {
		return claimPacketManifestFile{}, err
	}

	hash := sha256.New()
	written, err := io.MultiWriter(entry, hash).Write(content)
	if err != nil {
		return claimPacketManifestFile{}, err
	}

	file := claimPacketManifestFile{Path: path, Size: int64(written), Sha256: hex.EncodeToString(hash.Sum(nil))}
	if path != "manifest.json" {
		packet.files = append(packet.files, file)
	}
	return file, nil
}

// annotateLast records which receipt the last added file belongs to.
func (packet *packetZip) annotateLast(receiptId uint, originalName string) {
	last := &packet.files[len(packet.files)-1]
	last.ReceiptId = receiptId
	last.OriginalName = originalName
}

var packetNameUnsafe = regexp.MustCompile(`[^A-Za-z0-9]+`)

// packetSafeName reduces user text to ASCII letters and digits joined by "_",
// so it is safe in a zip entry and a Content-Disposition on every OS.
func packetSafeName(name string, fallback string, maxLength int) string {
	safe := strings.Trim(packetNameUnsafe.ReplaceAllString(name, "_"), "_")
	if len(safe) > maxLength {
		safe = strings.TrimRight(safe[:maxLength], "_")
	}
	if safe == "" {
		return fallback
	}
	return safe
}

// packetImageName is an image's name in images/:
// <receiptDate>_<receiptName>_<receiptId>_<n>.<ext>, so a folder listing sorts by
// date and every name is unique. n counts from 1 in the receipt's image order.
func packetImageName(receipt models.Receipt, n int, file models.FileData) string {
	return fmt.Sprintf("%s_%s_%d_%d.%s",
		receipt.Date.UTC().Format(time.DateOnly),
		packetSafeName(receipt.Name, "receipt", 40),
		receipt.ID,
		n,
		packetImageExtension(file),
	)
}

var packetExtensionSafe = regexp.MustCompile(`^[a-z0-9]{1,8}$`)

// packetImageExtension takes the extension from the type sniffed at upload, then
// from the uploaded name, then gives up with "bin".
func packetImageExtension(file models.FileData) string {
	if mime := mimetype.Lookup(file.FileType); mime != nil && mime.Extension() != "" {
		return strings.TrimPrefix(mime.Extension(), ".")
	}
	extension := strings.ToLower(strings.TrimPrefix(filepath.Ext(file.Name), "."))
	if packetExtensionSafe.MatchString(extension) {
		return extension
	}
	return "bin"
}

// claimPacketFilename is the download's name: the claim's name and today in the
// app zone.
func claimPacketFilename(claimName string, now time.Time) string {
	return packetSafeName(claimName, "claim", 60) + "_packet_" + now.Format(time.DateOnly) + ".zip"
}
