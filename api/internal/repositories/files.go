package repositories

import (
	"archive/zip"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"receipt-wrangler/api/internal/constants"
	config "receipt-wrangler/api/internal/env"
	"receipt-wrangler/api/internal/logging"
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/utils"
	"regexp"
	"strings"

	"github.com/gabriel-vasile/mimetype"
	"gopkg.in/gographics/imagick.v3/imagick"
	"gorm.io/gorm"
)

type FileRepository struct {
	BaseRepository
}

func NewFileRepository(tx *gorm.DB) FileRepository {
	repository := FileRepository{BaseRepository: BaseRepository{
		DB: GetDB(),
		TX: tx,
	}}
	return repository
}

func (repository FileRepository) BuildFilePath(receiptId string, receiptImageId string, receiptImageFileName string) (string, error) {
	db := repository.GetDB()
	var receipt models.Receipt

	err := db.Model(models.Receipt{}).Where("id = ?", receiptId).Select("group_id").Find(&receipt).Error
	if err != nil {
		return "", err
	}

	groupPath, err := repository.BuildGroupPath(receipt.GroupId, "")
	if err != nil {
		return "", err
	}

	fileName := utils.BuildFileName(receiptId, receiptImageId, receiptImageFileName)
	path := filepath.Join(groupPath, fileName)

	// The file name component (which can originate from an upload or email
	// attachment) must not escape the data directory either (CWE-22).
	if err := utils.AssertWithinDataDir(path); err != nil {
		return "", err
	}

	return path, nil
}

func (repository FileRepository) BuildGroupPath(groupId uint, alternateGroupName string) (string, error) {
	db := repository.GetDB()
	var groupNameToUse string

	if len(alternateGroupName) > 0 {
		groupNameToUse = alternateGroupName
	} else {
		var group models.Group
		err := db.Model(models.Group{}).Where("id = ?", groupId).Select("name").Find(&group).Error
		if err != nil {
			return "", err
		}

		groupNameToUse = group.Name
	}

	strGroupId := utils.UintToString(groupId)
	groupPath, err := utils.BuildGroupPathString(strGroupId, groupNameToUse)
	if err != nil {
		return "", err
	}

	return groupPath, nil
}

func (repository FileRepository) GetBytesForFileData(fileData models.FileData) ([]byte, error) {
	path, err := repository.BuildFilePath(utils.UintToString(fileData.ReceiptId), utils.UintToString(fileData.ID), fileData.Name)

	if err != nil {
		return nil, err
	}

	fileBytes, err := utils.ReadDataFile(path)
	if err != nil {
		return nil, err
	}

	resultBytes, err := repository.GetBytesFromImageBytes(fileBytes)
	if err != nil {
		return nil, err
	}

	return resultBytes, nil
}

func (repository FileRepository) GetBytesFromImageBytes(imageData []byte) ([]byte, error) {
	var bytes []byte
	validatedType, err := repository.ValidateFileType(imageData)
	if err != nil {
		return nil, err
	}

	if validatedType == constants.ApplicationPdf {
		bytes, err = repository.ConvertPdfToJpg(imageData)
		if err != nil {
			return nil, err
		}
	} else if validatedType == constants.ImageHeic || validatedType == constants.ImageHeif {
		// HEIF is the container HEIC is a profile of; a mif1-branded file sniffs as
		// HEIF and is read the same way.
		bytes, err = repository.ConvertHeicToJpg(imageData)
		if err != nil {
			return nil, err
		}
	} else {
		bytes = imageData
	}

	return bytes, nil
}

// BuildDisplayImageString turns raw file bytes into a data URI a browser can
// render in an <img>: a PDF rasterizes, HEIC and HEIF transcode, anything else
// passes through.
//
// Every browser-facing image response goes through this. Calling
// BuildEncodedImageString on unconverted bytes is the trap it closes — GetFileType
// relabels a PDF as image/jpeg WITHOUT converting it, so the result is a
// well-formed data URI carrying bytes the browser cannot decode: a broken image
// with no error anywhere.
//
// Reading the bytes stays the caller's job, because that half is a trust boundary
// and legitimately differs — data/ goes through the containment-checked
// utils.ReadDataFile, temp/ through os.ReadFile plus AssertWithinTempDirectory,
// and a multipart upload has no path at all. This half is identical everywhere.
func (repository FileRepository) BuildDisplayImageString(raw []byte) (string, error) {
	converted, err := repository.GetBytesFromImageBytes(raw)
	if err != nil {
		return "", err
	}

	return repository.BuildEncodedImageString(converted)
}

func (repository FileRepository) IsImage(imageData []byte) (bool, error) {
	validatedFileType, err := repository.ValidateFileType(imageData)
	if err != nil {
		return false, err
	}

	isImage, err := regexp.Match(constants.AnyImage, []byte(validatedFileType))
	if err != nil {
		return false, err
	}

	return isImage, nil
}

func (repository FileRepository) IsPdf(imageData []byte) (bool, error) {
	validatedFileType, err := repository.ValidateFileType(imageData)
	if err != nil {
		return false, err
	}

	isPdf, err := regexp.Match(constants.ApplicationPdf, []byte(validatedFileType))
	if err != nil {
		return false, err
	}

	return isPdf, nil
}

func (repository FileRepository) ValidateFileType(bytes []byte) (string, error) {
	fileType := utils.GetMimeType(bytes).String()
	acceptedFileTypes := []string{constants.AnyImage, constants.ApplicationPdf}

	for _, acceptedFileType := range acceptedFileTypes {
		matched, _ := regexp.Match(acceptedFileType, []byte(fileType))

		if matched {
			return fileType, nil
		}
	}

	return "", errors.New("invalid file type")
}

func (repository FileRepository) ValidateJsonFileType(bytes []byte) (string, error) {
	fileType := mimetype.Detect(bytes).String()
	acceptedFileTypes := []string{constants.ApplicationJson}

	for _, acceptedFileType := range acceptedFileTypes {
		matched, _ := regexp.Match(acceptedFileType, []byte(fileType))

		if matched {
			return fileType, nil
		}
	}

	return "", errors.New("invalid file type")
}

func (repository FileRepository) ConvertHeicToJpg(bytes []byte) ([]byte, error) {
	mw := imagick.NewMagickWand()
	defer mw.Destroy()

	if err := mw.ReadImageBlob(bytes); err != nil {
		return nil, err
	}

	if err := mw.SetImageFormat("jpeg"); err != nil {
		return nil, err
	}

	if err := mw.SetCompressionQuality(95); err != nil {
		return nil, err
	}

	return mw.GetImageBlob()
}

// pdfRasterizationDpi returns the DPI used to rasterize PDFs into images. It is
// configured via the PdfDpi system setting; when unset (0) or out of the
// supported range it falls back to defaultPdfDpi. 300 DPI is a common
// document-scanning resolution that preserves small receipt text for OCR and
// vision models while keeping output sizes and memory use reasonable.
const (
	defaultPdfDpi = 300.0
	minPdfDpi     = 72.0
	maxPdfDpi     = 1200.0
)

func (repository FileRepository) pdfRasterizationDpi() float64 {
	systemSettingsRepository := NewSystemSettingsRepository(repository.DB)
	systemSettings, err := systemSettingsRepository.GetSystemSettings()
	if err != nil {
		return defaultPdfDpi
	}

	dpi := float64(systemSettings.PdfDpi)
	if dpi < minPdfDpi || dpi > maxPdfDpi {
		return defaultPdfDpi
	}
	return dpi
}

const (
	// maxJpegDimension is the largest width or height a JPEG can hold (libjpeg's
	// JPEG_MAX_DIMENSION).
	maxJpegDimension = 65500

	// imageMagickHeightResource is MagickCore's HeightResource, the policy.xml "height"
	// limit. imagick.v3 does not export it, so it is named by its position in MagickCore's
	// ResourceType enum (resource_.h); TestImageMagickHeightResource_IsTheHeightLimit pins it.
	imageMagickHeightResource = imagick.ResourceType(4)
)

// rasterizedPage is one PDF page, rendered, flattened onto white and spilled to temp/.
type rasterizedPage struct {
	path          string
	width, height uint
}

// ConvertPdfToJpg rasterizes every page of a PDF at the PdfDpi system setting and stacks
// them top to bottom into a single JPEG.
//
// It holds at most one rendered page in memory, plus the final image: each page is rendered
// alone, flattened and spilled to temp/, and the final image is built at its exact size once
// every page has been measured. Reading the whole PDF and appending held every page and a
// full-height copy at once, enough to get the API OOM-killed on a 1 GB host. Ping cannot
// measure pages up front: ImageMagick reports the document's largest page for every page.
// See api/CLAUDE.md → "PDF rasterization".
func (repository FileRepository) ConvertPdfToJpg(pdfBytes []byte) ([]byte, error) {
	// PDFs are vector documents with no inherent pixel resolution. ImageMagick
	// rasterizes them at a default of 72 DPI unless a resolution is set BEFORE
	// the page is read, which produces low-resolution images that degrade OCR
	// and vision-model accuracy. Set the rasterization density first; it is
	// configurable via the PdfDpi system setting (default 300).
	dpi := repository.pdfRasterizationDpi()

	pdfPath, err := repository.BuildTempFilePath("pdf")
	if err != nil {
		return nil, err
	}
	defer os.Remove(pdfPath)
	if err := os.WriteFile(pdfPath, pdfBytes, 0o600); err != nil {
		return nil, err
	}

	pageCount, err := countPdfPages(pdfPath)
	if err != nil {
		return nil, err
	}

	pages, err := repository.rasterizePdfPages(pdfPath, pageCount, dpi)
	if err != nil {
		return nil, err
	}
	defer removeRasterizedPages(pages)

	return stackPagesAsJpeg(pages, dpi)
}

// maxPdfImageHeight is the tallest stacked image ConvertPdfToJpg can produce.
func maxPdfImageHeight() uint {
	return pdfImageHeightLimit(imagick.GetResourceLimit(imageMagickHeightResource))
}

// pdfImageHeightLimit is ImageMagick's height limit or the JPEG format's own, whichever is
// lower. The binding casts MagickSizeType to int64, so "unlimited" arrives as -1.
func pdfImageHeightLimit(imageMagickLimit int64) uint {
	if imageMagickLimit > 0 && imageMagickLimit < maxJpegDimension {
		return uint(imageMagickLimit)
	}
	return maxJpegDimension
}

// countPdfPages pings the PDF for its page count only; see ConvertPdfToJpg for why ping's
// page sizes are not used.
func countPdfPages(pdfPath string) (int, error) {
	mw := imagick.NewMagickWand()
	defer mw.Destroy()

	// The pdf: prefix pins the decoder, so the bytes are never sniffed as another format.
	if err := mw.PingImage("pdf:" + pdfPath); err != nil {
		return 0, err
	}

	pageCount := int(mw.GetNumberImages())
	if pageCount == 0 {
		return 0, errors.New("pdf has no pages")
	}
	return pageCount, nil
}

// rasterizePdfPages renders each page to its own temp file, failing as soon as the pages
// would stack taller than the final image may be. Without that check ImageMagick fails to
// build the stacked image, and the old code went on to write the nil wand it got back,
// which aborted the whole process. On error it removes every page it spilled.
func (repository FileRepository) rasterizePdfPages(pdfPath string, pageCount int, dpi float64) ([]rasterizedPage, error) {
	maxHeight := maxPdfImageHeight()
	pages := make([]rasterizedPage, 0, pageCount)
	var totalHeight uint

	for i := 0; i < pageCount; i++ {
		page, err := repository.rasterizePdfPage(pdfPath, i, dpi)
		if err != nil {
			removeRasterizedPages(pages)
			return nil, err
		}
		pages = append(pages, page)

		totalHeight += page.height
		if totalHeight > maxHeight {
			removeRasterizedPages(pages)
			return nil, fmt.Errorf(
				"pdf is too tall to convert: %d of its %d pages already stack to %dpx at %g DPI, over the %dpx limit; lower the PDF DPI system setting",
				i+1, pageCount, totalHeight, dpi, maxHeight,
			)
		}
	}

	return pages, nil
}

// rasterizePdfPage renders one page, flattens it onto white and spills it to temp/ as an
// 8-bit MIFF. It is its own function so the deferred Destroy frees the page before the next
// one renders: a defer in the caller's loop would hold every page until the end.
func (repository FileRepository) rasterizePdfPage(pdfPath string, index int, dpi float64) (rasterizedPage, error) {
	white, err := newWhitePixelWand()
	if err != nil {
		return rasterizedPage{}, err
	}
	defer white.Destroy()

	page := imagick.NewMagickWand()
	defer page.Destroy()

	if err := page.SetResolution(dpi, dpi); err != nil {
		return rasterizedPage{}, err
	}
	// The [index] suffix makes ImageMagick ask Ghostscript for this page alone.
	if err := page.ReadImage(fmt.Sprintf("pdf:%s[%d]", pdfPath, index)); err != nil {
		return rasterizedPage{}, err
	}

	// A CMYK PDF renders CMYK pages, and the stacked image is sRGB.
	if err := page.TransformImageColorspace(imagick.COLORSPACE_SRGB); err != nil {
		return rasterizedPage{}, err
	}

	// Ghostscript leaves anything the PDF does not paint transparent, with black underneath.
	// JPEG has no alpha, so without flattening onto white those areas come out as black
	// blocks. (The old SetImageAlphaChannel(0) was UndefinedAlphaChannel: a no-op.)
	if err := page.SetImageBackgroundColor(white); err != nil {
		return rasterizedPage{}, err
	}
	if err := page.SetImageAlphaChannel(imagick.ALPHA_CHANNEL_REMOVE); err != nil {
		return rasterizedPage{}, err
	}
	if err := page.SetImageDepth(8); err != nil {
		return rasterizedPage{}, err
	}

	path, err := repository.BuildTempFilePath("miff")
	if err != nil {
		return rasterizedPage{}, err
	}
	if err := page.WriteImage("miff:" + path); err != nil {
		os.Remove(path)
		return rasterizedPage{}, err
	}

	return rasterizedPage{path: path, width: page.GetImageWidth(), height: page.GetImageHeight()}, nil
}

// stackPagesAsJpeg builds the final image at its exact size and copies each spilled page
// into place, reading one at a time.
func stackPagesAsJpeg(pages []rasterizedPage, dpi float64) ([]byte, error) {
	var width, height uint
	for _, page := range pages {
		width = max(width, page.width)
		height += page.height
	}

	white, err := newWhitePixelWand()
	if err != nil {
		return nil, err
	}
	defer white.Destroy()

	// A narrower page leaves white to its right, as AppendImages did.
	canvas := imagick.NewMagickWand()
	defer canvas.Destroy()
	if err := canvas.NewImage(width, height, white); err != nil {
		return nil, fmt.Errorf("could not create a %dx%d image for the pdf: %w", width, height, err)
	}
	// Keep the render density, as the appended pages did: Tesseract reads it.
	if err := canvas.SetImageResolution(dpi, dpi); err != nil {
		return nil, err
	}

	y := 0
	for _, page := range pages {
		if err := copySpilledPage(canvas, page, y); err != nil {
			return nil, err
		}
		y += int(page.height)
	}

	if err := canvas.SetImageDepth(8); err != nil {
		return nil, err
	}
	if err := canvas.SetImageFormat("jpeg"); err != nil {
		return nil, err
	}

	jpg, err := canvas.GetImageBlob()
	if len(jpg) == 0 {
		if err == nil {
			err = errors.New("encoding the pdf as a jpeg produced no data")
		}
		return nil, err
	}
	if err != nil {
		// GetImageBlob reports warnings as errors too; with a JPEG in hand, the encode worked.
		logging.LogStd(logging.LOG_LEVEL_INFO, "ImageMagick warning while encoding a pdf as a jpeg: "+err.Error())
	}
	return jpg, nil
}

// copySpilledPage reads one spilled page and copies it onto the canvas at row y. Pages are
// opaque by now, so a straight copy is enough.
func copySpilledPage(canvas *imagick.MagickWand, page rasterizedPage, y int) error {
	mw := imagick.NewMagickWand()
	defer mw.Destroy()

	if err := mw.ReadImage("miff:" + page.path); err != nil {
		return err
	}
	return canvas.CompositeImage(mw, imagick.COMPOSITE_OP_COPY, true, 0, y)
}

func removeRasterizedPages(pages []rasterizedPage) {
	for _, page := range pages {
		os.Remove(page.path)
	}
}

func newWhitePixelWand() (*imagick.PixelWand, error) {
	white := imagick.NewPixelWand()
	if !white.SetColor("white") {
		white.Destroy()
		return nil, errors.New("could not create a white pixel wand")
	}
	return white, nil
}

func (repository FileRepository) WriteTempFile(data []byte) (string, error) {
	tempPath := repository.GetTempDirectoryPath()
	utils.MakeDirectory(tempPath)

	validatedFileType, err := repository.ValidateFileType(data)
	if err != nil {
		return "", err
	}

	parts := strings.Split(validatedFileType, "/")
	if len(parts) != 2 {
		return "", errors.New("malformed mime type")
	}

	fileType := parts[1]

	filePath, err := repository.BuildTempFilePath(fileType)
	if err != nil {
		return "", err
	}

	err = utils.WriteFile(filePath, data)
	if err != nil {
		os.Remove(filePath)
		return "", err
	}

	return filePath, nil
}

func (repository FileRepository) BuildTempFilePath(fileType string) (string, error) {
	tempPath := repository.GetTempDirectoryPath()
	// Ensure the temp directory exists before returning a path under it.
	// ConvertPdfToJpg spills its PDF and rendered pages here, and on a fresh
	// checkout/runtime the directory may not exist yet, which would fail the write.
	utils.MakeDirectory(tempPath)

	filename, err := utils.GetRandomString(10)
	if err != nil {
		return "", err
	}

	filePath := tempPath + "/" + filename
	filePath = filePath + "." + fileType
	return filePath, nil
}

func (repository FileRepository) GetFileType(bytes []byte) (string, error) {
	fileType, err := repository.ValidateFileType(bytes)
	if err != nil {
		return "", err
	}

	isPdf, err := repository.IsPdf(bytes)
	if err != nil {
		return "", err
	}

	if isPdf {
		fileType = "image/jpeg"
	}

	return fileType, nil
}

func (repository FileRepository) BuildEncodedImageString(bytes []byte) (string, error) {
	fileType, err := repository.GetFileType(bytes)
	if err != nil {
		return "", err
	}

	return utils.BuildDataURI(fileType, bytes), nil
}

func (repository FileRepository) CreateZipFromTempFiles(zipFilename string, filenames []string) (string, error) {
	tempPath := repository.GetTempDirectoryPath()
	zipPath := filepath.Join(tempPath, zipFilename)

	zipFile, err := os.Create(zipPath)
	if err != nil {
		return "", err
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	for _, filename := range filenames {
		filePath := filepath.Join(tempPath, filename)
		file, err := os.Open(filePath)
		if err != nil {
			return "", err
		}
		defer file.Close()

		info, err := file.Stat()
		if err != nil {
			return "", err
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return "", err
		}

		header.Name = filename
		header.Method = zip.Deflate

		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			return "", err
		}

		if _, err = io.Copy(writer, file); err != nil {
			return "", err
		}
	}

	return zipPath, nil
}

func (repository FileRepository) ZipFiles(filenames []string, fileContents [][]byte) ([]byte, error) {
	if len(filenames) != len(fileContents) {
		return nil, errors.New("number of filenames does not match number of file contents")
	}

	if len(filenames) == 0 {
		return nil, errors.New("no files to zip")
	}

	zipBuffer := new(bytes.Buffer)
	zipWriter := zip.NewWriter(zipBuffer)

	for i, filename := range filenames {
		file, err := zipWriter.Create(filename)
		if err != nil {
			return nil, fmt.Errorf("error creating file %s in zip: %w", filename, err)
		}

		_, err = file.Write(fileContents[i])
		if err != nil {
			return nil, fmt.Errorf("error writing content for file %s in zip: %w", filename, err)
		}
	}

	err := zipWriter.Close()
	if err != nil {
		return nil, fmt.Errorf("error closing zip writer: %w", err)
	}

	return zipBuffer.Bytes(), nil
}

func (repository FileRepository) GetTempDirectoryPath() string {
	return filepath.Join(config.GetBasePath(), "temp")
}

// AssertWithinTempDirectory returns an error when path resolves outside temp/.
//
// temp/ is deliberately exempt from the data-scoped helpers (see api/CLAUDE.md →
// "Filesystem Access & Path-Traversal Safety"): those resolve against data/ and
// would reject every temp path. That exemption holds for paths the server builds
// itself, but a temp path read back out of an asynq payload is attacker-adjacent
// — the payload is JSON in Redis — so anything sourced from one must come
// through here before it is opened or served.
func (repository FileRepository) AssertWithinTempDirectory(path string) error {
	return utils.AssertWithinDir(repository.GetTempDirectoryPath(), path)
}

func (repository FileRepository) GetTestJpgBytes() ([]byte, error) {
	path := filepath.Join(config.GetBasePath(), "testing", "test.jpg")

	return utils.ReadFile(path)
}
