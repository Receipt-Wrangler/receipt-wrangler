package handlers

import (
	"archive/zip"
	"bytes"
	"crypto/sha256"
	"encoding/csv"
	"encoding/hex"
	"encoding/json"
	"io"
	"mime"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"testing"
	"time"

	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/permissions"
	"receipt-wrangler/api/internal/repositories"
	"receipt-wrangler/api/internal/services"
	"receipt-wrangler/api/internal/utils"

	"github.com/shopspring/decimal"
)

// The claim packet tests replace Chromium with a renderer that returns the HTML
// itself, so "packet.pdf" holds exactly what the PDF would print and can be
// asserted on as text. TestClaimPacket_RealPdf runs the real pipeline once.

type packetManifest struct {
	FormatVersion int    `json:"formatVersion"`
	GeneratedAt   string `json:"generatedAt"`
	GeneratedBy   string `json:"generatedBy"`
	Claim         struct {
		Id         uint   `json:"id"`
		Name       string `json:"name"`
		ReceiptIds []uint `json:"receiptIds"`
	} `json:"claim"`
	ReceiptIds []uint `json:"receiptIds"`
	Files      []struct {
		Path         string `json:"path"`
		Size         int64  `json:"size"`
		Sha256       string `json:"sha256"`
		ReceiptId    uint   `json:"receiptId"`
		OriginalName string `json:"originalName"`
	} `json:"files"`
	MissingImages []struct {
		ReceiptId    uint   `json:"receiptId"`
		FileId       uint   `json:"fileId"`
		OriginalName string `json:"originalName"`
	} `json:"missingImages"`
}

type unpackedPacket struct {
	files    map[string][]byte
	order    []string
	manifest packetManifest
}

func (packet unpackedPacket) html() string {
	return string(packet.files["packet.pdf"])
}

// everything is every text file in the packet, for "is this anywhere" checks.
func (packet unpackedPacket) everything() string {
	return packet.html() + string(packet.files["receipts.csv"]) + string(packet.files["items.csv"]) +
		string(packet.files["manifest.json"]) + strings.Join(packet.order, "\n")
}

func stubPacketPdf(t *testing.T) {
	t.Helper()
	restore := services.SetClaimPacketPdfRendererForTests(func(html string) ([]byte, error) {
		return []byte(html), nil
	})
	t.Cleanup(restore)
}

func requestPacket(t *testing.T, userId uint, claimId uint) *httptest.ResponseRecorder {
	t.Helper()
	w, r := claimRequest(http.MethodGet, userId, "", claimParams(claimId))
	GetClaimPacket(w, r)
	return w
}

func downloadPacket(t *testing.T, userId uint, claimId uint) unpackedPacket {
	t.Helper()
	w := requestPacket(t, userId, claimId)
	if w.Code != http.StatusOK {
		t.Fatalf("status %d, body %s", w.Code, w.Body.String())
	}
	return unpackPacket(t, w.Body.Bytes())
}

func unpackPacket(t *testing.T, body []byte) unpackedPacket {
	t.Helper()
	reader, err := zip.NewReader(bytes.NewReader(body), int64(len(body)))
	if err != nil {
		t.Fatalf("not a zip: %v", err)
	}
	packet := unpackedPacket{files: map[string][]byte{}}
	for _, file := range reader.File {
		opened, err := file.Open()
		if err != nil {
			t.Fatalf("open %s: %v", file.Name, err)
		}
		content, err := io.ReadAll(opened)
		opened.Close()
		if err != nil {
			t.Fatalf("read %s: %v", file.Name, err)
		}
		packet.files[file.Name] = content
		packet.order = append(packet.order, file.Name)
	}
	if err := json.Unmarshal(packet.files["manifest.json"], &packet.manifest); err != nil {
		t.Fatalf("manifest: %v", err)
	}
	return packet
}

type packetReceiptSeed struct {
	name       string
	amount     string
	date       time.Time
	paidBy     uint
	categories []models.Category
	tags       []models.Tag
	createdBy  uint
}

func seedPacketReceipt(t *testing.T, groupId uint, seed packetReceiptSeed) models.Receipt {
	t.Helper()
	receipt := models.Receipt{
		Name:         seed.name,
		Amount:       decimal.RequireFromString(seed.amount),
		Date:         seed.date,
		GroupId:      groupId,
		PaidByUserID: seed.paidBy,
		Status:       models.OPEN,
		Categories:   seed.categories,
		Tags:         seed.tags,
	}
	if seed.createdBy != 0 {
		createdBy := seed.createdBy
		receipt.CreatedBy = &createdBy
		receipt.CreatedByString = "Creator Number Three"
	}
	if err := repositories.GetDB().Create(&receipt).Error; err != nil {
		t.Fatalf("seed receipt: %v", err)
	}
	return receipt
}

func seedCategory(t *testing.T, name string) models.Category {
	t.Helper()
	category := models.Category{Name: name}
	if err := repositories.GetDB().Create(&category).Error; err != nil {
		t.Fatalf("seed category: %v", err)
	}
	return category
}

func sha256Hex(content []byte) string {
	sum := sha256.Sum256(content)
	return hex.EncodeToString(sum[:])
}

func storedImagePath(t *testing.T, image models.FileData) string {
	t.Helper()
	path, err := repositories.NewFileRepository(nil).BuildFilePath(utils.UintToString(image.ReceiptId), utils.UintToString(image.ID), image.Name)
	if err != nil {
		t.Fatalf("BuildFilePath: %v", err)
	}
	return path
}

func removeStoredImages(t *testing.T, images ...models.FileData) {
	t.Helper()
	// The paths are resolved now: by cleanup time the deferred teardown has
	// emptied the database the path is built from.
	paths := make([]string, 0, len(images))
	for _, image := range images {
		paths = append(paths, storedImagePath(t, image))
	}
	t.Cleanup(func() {
		for _, path := range paths {
			// utils.WriteFile writes with mode 777 DECIMAL, which is read-only on
			// Windows, where a read-only file cannot be removed.
			os.Chmod(path, 0o644)
			utils.RemoveDataPath(path)
		}
	})
}

func readCsv(t *testing.T, content []byte) [][]string {
	t.Helper()
	rows, err := csv.NewReader(bytes.NewReader(content)).ReadAll()
	if err != nil {
		t.Fatalf("csv: %v", err)
	}
	return rows
}

func TestClaimPacket_ZipLayoutManifestHashesAndContent(t *testing.T) {
	defer tearDownGenericHandlerTest()
	stubPacketPdf(t)
	repositories.CreateTestGroupWithUsers()
	grantGroupPerms(t, 1, 1, permissions.GroupReceiptsRead)

	food := seedCategory(t, "Food")
	travel := seedCategory(t, "Travel")
	claim := seedClaimRow(t, 1, 1, "PCS Rome / TDY")
	repositories.GetDB().Model(&models.Claim{}).Where("id = ?", claim.ID).Updates(map[string]interface{}{
		"reference_number": "TA-77", "status": models.ClaimStatusSubmitted,
	})
	seedLedgerEntry(t, claim.ID, 1, models.ClaimLedgerEntryPayment, "10.00")

	// 0.10 + 0.20 is the classic float trap; the subtotals must still be exact.
	both := seedPacketReceipt(t, 1, packetReceiptSeed{name: "Taxi & dinner", amount: "0.10",
		date: time.Date(2026, 9, 2, 0, 0, 0, 0, time.UTC), paidBy: 1, categories: []models.Category{food, travel}})
	foodOnly := seedPacketReceipt(t, 1, packetReceiptSeed{name: "Groceries", amount: "0.20",
		date: time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC), paidBy: 1, categories: []models.Category{food}})
	none := seedPacketReceipt(t, 1, packetReceiptSeed{name: "Visa fee", amount: "100.05",
		date: time.Date(2026, 9, 3, 0, 0, 0, 0, time.UTC), paidBy: 1})
	linkReceipts(t, claim.ID, both.ID, foodOnly.ID, none.ID)

	first := addStoredImage(t, both.ID, "IMG_0001.JPG")
	second := addStoredImage(t, both.ID, "scan.jpg")
	removeStoredImages(t, first, second)
	addStoredComment(t, both.ID, 1, "Approved by post management")

	field := models.CustomField{Name: "Justification", Type: models.TEXT}
	repositories.GetDB().Create(&field)
	justification := "Official travel to Rome"
	repositories.GetDB().Create(&models.CustomFieldValue{ReceiptId: both.ID, CustomFieldId: field.ID, StringValue: &justification})

	w := requestPacket(t, 1, claim.ID)
	assertStatus(t, w, http.StatusOK)
	if got := w.Header().Get("Content-Type"); got != "application/zip" {
		utils.PrintTestError(t, got, "application/zip")
	}
	_, params, err := mime.ParseMediaType(w.Header().Get("Content-Disposition"))
	if err != nil {
		t.Fatalf("content disposition: %v", err)
	}
	today := time.Now().In(repositories.GetAppLocation()).Format(time.DateOnly)
	if want := "PCS_Rome_TDY_packet_" + today + ".zip"; params["filename"] != want {
		utils.PrintTestError(t, params["filename"], want)
	}
	if w.Header().Get("Content-Length") != utils.UintToString(uint(w.Body.Len())) {
		utils.PrintTestError(t, w.Header().Get("Content-Length"), w.Body.Len())
	}

	packet := unpackPacket(t, w.Body.Bytes())

	imageOne := "2026-09-02_Taxi_dinner_" + utils.UintToString(both.ID) + "_1.jpg"
	imageTwo := "2026-09-02_Taxi_dinner_" + utils.UintToString(both.ID) + "_2.jpg"
	wantFiles := []string{"images/" + imageOne, "images/" + imageTwo, "items.csv", "manifest.json", "packet.pdf", "receipts.csv"}
	gotFiles := append([]string{}, packet.order...)
	sort.Strings(gotFiles)
	if strings.Join(gotFiles, ",") != strings.Join(wantFiles, ",") {
		utils.PrintTestError(t, gotFiles, wantFiles)
	}

	// The originals are byte-identical to what is on disk.
	for name, image := range map[string]models.FileData{imageOne: first, imageTwo: second} {
		onDisk, err := os.ReadFile(storedImagePath(t, image))
		if err != nil {
			t.Fatalf("read original: %v", err)
		}
		if !bytes.Equal(onDisk, packet.files["images/"+name]) {
			t.Errorf("images/%s differs from the stored file", name)
		}
	}

	// Every file but the manifest is listed, with a hash of exactly its bytes.
	if len(packet.manifest.Files) != len(packet.files)-1 {
		utils.PrintTestError(t, len(packet.manifest.Files), len(packet.files)-1)
	}
	for _, file := range packet.manifest.Files {
		content, ok := packet.files[file.Path]
		if !ok {
			t.Errorf("manifest lists %s, which is not in the zip", file.Path)
			continue
		}
		if file.Sha256 != sha256Hex(content) || file.Size != int64(len(content)) {
			t.Errorf("%s: manifest says %s/%d, zip holds %s/%d", file.Path, file.Sha256, file.Size, sha256Hex(content), len(content))
		}
		if strings.HasPrefix(file.Path, "images/") && (file.ReceiptId != both.ID || file.OriginalName == "") {
			t.Errorf("%s: not attributed to its receipt: %+v", file.Path, file)
		}
	}
	if packet.manifest.Claim.Id != claim.ID || packet.manifest.GeneratedBy != "test" || packet.manifest.FormatVersion != 1 {
		utils.PrintTestError(t, packet.manifest, "the claim, generated by the caller")
	}
	if _, err := time.Parse(time.RFC3339, packet.manifest.GeneratedAt); err != nil {
		t.Errorf("generatedAt %q is not RFC 3339", packet.manifest.GeneratedAt)
	}
	wantIds := []uint{foodOnly.ID, both.ID, none.ID} // by date
	if utils.UintToString(uint(len(packet.manifest.ReceiptIds))) != "3" || packet.manifest.ReceiptIds[0] != wantIds[0] ||
		packet.manifest.ReceiptIds[1] != wantIds[1] || packet.manifest.ReceiptIds[2] != wantIds[2] {
		utils.PrintTestError(t, packet.manifest.ReceiptIds, wantIds)
	}

	// receipts.csv has the export's columns plus Images.
	rows := readCsv(t, packet.files["receipts.csv"])
	header := rows[0]
	if header[len(header)-1] != "Images" || header[0] != "Id" {
		utils.PrintTestError(t, header, "the export's columns, then Images")
	}
	for _, row := range rows[1:] {
		want := ""
		if row[0] == utils.UintToString(both.ID) {
			want = imageOne + "," + imageTwo
		}
		if row[len(row)-1] != want {
			utils.PrintTestError(t, row[len(row)-1], want)
		}
	}

	html := packet.html()
	for _, want := range []string{
		"PCS Rome / TDY", "TA-77", "Submitted",
		"Subtotal: Food</td><td class=\"num\">$0.30", // 0.10 + 0.20, exactly
		"Subtotal: Travel</td><td class=\"num\">$0.10",
		"Subtotal: Uncategorised</td><td class=\"num\">$100.05",
		"$100.35",                          // the grand total, each receipt once
		"add up to $100.45, more than the", // the multi-category note
		"Taxi &amp; dinner (#" + utils.UintToString(both.ID) + ") *",
		"Official travel to Rome", "Approved by post management",
		"images/" + imageOne, "SHA-256 " + sha256Hex(packet.files["images/"+imageOne]),
		"data:image/jpeg;base64,",
		"Payment", "$10.00",
	} {
		if !strings.Contains(html, want) {
			t.Errorf("the PDF does not print %q", want)
		}
	}
	if strings.Index(html, "Subtotal: Food") > strings.Index(html, "Subtotal: Travel") ||
		strings.Index(html, "Subtotal: Travel") > strings.Index(html, "Subtotal: Uncategorised") {
		t.Errorf("categories are not in name order with Uncategorised last")
	}
}

// A receipt hidden from the caller by paid-by visibility is absent from every
// file of the packet.
func TestClaimPacket_HiddenReceiptIsAbsentEverywhere(t *testing.T) {
	defer tearDownGenericHandlerTest()
	stubPacketPdf(t)
	services.ClearRolePermissionCacheForTests()
	services.ClearGroupRoleGrantCacheForTests()
	repositories.CreateTestGroupWithUsers()

	role, err := repositories.NewRoleRepository(nil).CreateGroupRole("Packet PaidBy Restricted", "",
		[]string{permissions.GroupReceiptsRead}, nil, nil, []uint{1}, false, false)
	if err != nil {
		t.Fatalf("role: %v", err)
	}
	repositories.GetDB().Model(&models.GroupMember{}).Where("user_id = ? AND group_id = ?", 2, 1).Update("group_role_id", role.ID)

	claim := seedClaimRow(t, 1, 1, "shared")
	visible := seedPacketReceipt(t, 1, packetReceiptSeed{name: "VisibleLunch", amount: "10.00", date: time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC), paidBy: 1})
	hidden := seedPacketReceipt(t, 1, packetReceiptSeed{name: "HiddenHotel", amount: "90.00", date: time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC), paidBy: 3})
	linkReceipts(t, claim.ID, visible.ID, hidden.ID)
	hiddenImage := addStoredImage(t, hidden.ID, "hidden-hotel.jpg")
	visibleImage := addStoredImage(t, visible.ID, "lunch.jpg")
	removeStoredImages(t, hiddenImage, visibleImage)
	addStoredComment(t, hidden.ID, 3, "HiddenComment")

	packet := downloadPacket(t, 2, claim.ID)

	everything := packet.everything()
	for _, leak := range []string{"HiddenHotel", "hidden-hotel.jpg", "HiddenComment", "90.00"} {
		if strings.Contains(everything, leak) {
			t.Errorf("the packet leaks %q", leak)
		}
	}
	if len(packet.manifest.ReceiptIds) != 1 || packet.manifest.ReceiptIds[0] != visible.ID {
		utils.PrintTestError(t, packet.manifest.ReceiptIds, []uint{visible.ID})
	}
	if len(packet.manifest.Claim.ReceiptIds) != 1 {
		utils.PrintTestError(t, packet.manifest.Claim.ReceiptIds, "only the visible receipt")
	}
	images := 0
	for _, name := range packet.order {
		if strings.HasPrefix(name, "images/") {
			images++
			if !strings.Contains(name, "VisibleLunch") {
				t.Errorf("unexpected image %s", name)
			}
		}
	}
	if images != 1 {
		utils.PrintTestError(t, images, 1)
	}
	if !strings.Contains(packet.html(), "VisibleLunch") {
		t.Errorf("the visible receipt is missing")
	}
}

// What is printed about a visible receipt is what GET /receipt shows: categories
// outside the caller's grants are stripped, and in a group with isolated members
// another member's comments and authorship are dropped.
func TestClaimPacket_GrantFilteringAndMemberMaskingApply(t *testing.T) {
	defer tearDownGenericHandlerTest()
	stubPacketPdf(t)
	services.ClearRolePermissionCacheForTests()
	services.ClearGroupRoleGrantCacheForTests()
	repositories.CreateTestGroupWithUsers()
	repositories.GetDB().Model(&models.Group{}).Where("id = ?", 1).Update("isolate_members", true)

	allowed := seedCategory(t, "AllowedCategory")
	restricted := seedCategory(t, "SecretCategory")
	role, err := repositories.NewRoleRepository(nil).CreateGroupRole("Packet Category Restricted", "",
		[]string{permissions.GroupReceiptsRead}, []uint{allowed.ID}, nil, nil, false, false)
	if err != nil {
		t.Fatalf("role: %v", err)
	}
	repositories.GetDB().Model(&models.GroupMember{}).Where("user_id = ? AND group_id = ?", 2, 1).Update("group_role_id", role.ID)

	claim := seedClaimRow(t, 1, 2, "mine")
	receipt := seedPacketReceipt(t, 1, packetReceiptSeed{name: "Conference", amount: "50.00",
		date: time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC), paidBy: 2,
		categories: []models.Category{allowed, restricted}, createdBy: 3})
	linkReceipts(t, claim.ID, receipt.ID)
	addStoredComment(t, receipt.ID, 3, "PeerCommentText")
	addStoredComment(t, receipt.ID, 2, "OwnCommentText")

	packet := downloadPacket(t, 2, claim.ID)

	everything := packet.everything()
	for _, leak := range []string{"SecretCategory", "PeerCommentText", "Creator Number Three"} {
		if strings.Contains(everything, leak) {
			t.Errorf("the packet leaks %q", leak)
		}
	}
	for _, want := range []string{"AllowedCategory", "OwnCommentText", "Conference"} {
		if !strings.Contains(packet.html(), want) {
			t.Errorf("the PDF does not print %q", want)
		}
	}
	// With its only other category stripped, the receipt is single-category.
	if strings.Contains(packet.html(), "listed under each of them") {
		t.Errorf("a stripped category still made the receipt multi-category")
	}
}

func TestClaimPacket_Permissions(t *testing.T) {
	defer tearDownGenericHandlerTest()
	stubPacketPdf(t)
	repositories.CreateTestGroupWithUsers()
	grantGroupPerms(t, 1, 1, permissions.GroupReceiptsRead)
	grantGroupPerms(t, 2, 1, permissions.GroupReceiptsUpdate) // no read
	grantGroupPerms(t, 4, 2, permissions.GroupReceiptsRead)   // another group
	claim := seedClaimRow(t, 1, 1, "permissions")

	assertStatus(t, requestPacket(t, 1, claim.ID), http.StatusOK)

	noRead := requestPacket(t, 2, claim.ID)
	otherGroup := requestPacket(t, 4, claim.ID)
	missing := requestPacket(t, 1, 9999)
	for name, w := range map[string]*httptest.ResponseRecorder{"no read": noRead, "other group": otherGroup, "missing": missing} {
		if w.Code != http.StatusForbidden {
			t.Errorf("%s: status %d, want 403", name, w.Code)
		}
		if w.Body.String() != missing.Body.String() {
			t.Errorf("%s: body %q differs from a missing claim's %q", name, w.Body.String(), missing.Body.String())
		}
		if w.Header().Get("Content-Disposition") != "" {
			t.Errorf("%s: a refusal must not look like a download", name)
		}
	}
}

func TestClaimPacket_ImageCountBoundIs400(t *testing.T) {
	defer tearDownGenericHandlerTest()
	stubPacketPdf(t)
	repositories.CreateTestGroupWithUsers()
	grantGroupPerms(t, 1, 1, permissions.GroupReceiptsRead)
	claim := seedClaimRow(t, 1, 1, "too many")
	receipt := seedPacketReceipt(t, 1, packetReceiptSeed{name: "Bulk", amount: "1.00", date: time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC), paidBy: 1})
	linkReceipts(t, claim.ID, receipt.ID)

	// Rows alone: the bound is checked before any file is read.
	rows := make([]models.FileData, 0, services.MaxClaimPacketImages+1)
	for i := 0; i <= services.MaxClaimPacketImages; i++ {
		rows = append(rows, models.FileData{Name: "x.jpg", FileType: "image/jpeg", ReceiptId: receipt.ID})
	}
	if err := repositories.GetDB().Create(&rows).Error; err != nil {
		t.Fatalf("seed files: %v", err)
	}

	w := requestPacket(t, 1, claim.ID)
	assertStatus(t, w, http.StatusBadRequest)
	if !strings.Contains(w.Body.String(), "at most") {
		utils.PrintTestError(t, w.Body.String(), "a message naming the bound")
	}

	// At the bound it passes (the absent files are noted, not fatal).
	repositories.GetDB().Where("id = ?", rows[0].ID).Delete(&models.FileData{})
	assertStatus(t, requestPacket(t, 1, claim.ID), http.StatusOK)
}

// An image whose file is gone from disk is noted in the PDF and the manifest,
// and the rest of the packet is still produced.
func TestClaimPacket_MissingImageFileIsNotedNotFatal(t *testing.T) {
	defer tearDownGenericHandlerTest()
	stubPacketPdf(t)
	repositories.CreateTestGroupWithUsers()
	grantGroupPerms(t, 1, 1, permissions.GroupReceiptsRead)
	claim := seedClaimRow(t, 1, 1, "missing file")
	receipt := seedPacketReceipt(t, 1, packetReceiptSeed{name: "Lost", amount: "3.00", date: time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC), paidBy: 1})
	linkReceipts(t, claim.ID, receipt.ID)
	gone := models.FileData{Name: "gone.jpg", FileType: "image/jpeg", ReceiptId: receipt.ID}
	repositories.GetDB().Create(&gone)

	packet := downloadPacket(t, 1, claim.ID)
	if len(packet.manifest.MissingImages) != 1 || packet.manifest.MissingImages[0].FileId != gone.ID {
		utils.PrintTestError(t, packet.manifest.MissingImages, "the missing image")
	}
	for _, name := range packet.order {
		if strings.HasPrefix(name, "images/") {
			t.Errorf("unexpected image %s", name)
		}
	}
	if !strings.Contains(packet.html(), "missing on the server") {
		t.Errorf("the PDF does not say the image is missing")
	}
}

// The real pipeline once: Chromium renders packet.pdf, and its text is read back
// with Ghostscript where available.
func TestClaimPacket_RealPdf(t *testing.T) {
	defer tearDownGenericHandlerTest()
	repositories.CreateTestGroupWithUsers()
	grantGroupPerms(t, 1, 1, permissions.GroupReceiptsRead)
	claim := seedClaimRow(t, 1, 1, "RealPdfClaim")
	receipt := seedPacketReceipt(t, 1, packetReceiptSeed{name: "RealPdfReceipt", amount: "12.34", date: time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC), paidBy: 1})
	linkReceipts(t, claim.ID, receipt.ID)
	image := addStoredImage(t, receipt.ID, "photo.jpg")
	removeStoredImages(t, image)

	packet := downloadPacket(t, 1, claim.ID)
	pdf := packet.files["packet.pdf"]
	if !bytes.HasPrefix(pdf, []byte("%PDF-")) {
		t.Fatalf("packet.pdf is not a PDF")
	}

	gs, err := exec.LookPath("gs")
	if err != nil {
		t.Skip("ghostscript (gs) not available; skipping PDF text extraction")
	}
	directory := t.TempDir()
	pdfPath := filepath.Join(directory, "packet.pdf")
	textPath := filepath.Join(directory, "packet.txt")
	if err := os.WriteFile(pdfPath, pdf, 0644); err != nil {
		t.Fatalf("write: %v", err)
	}
	if out, err := exec.Command(gs, "-q", "-dNOPAUSE", "-dBATCH", "-sDEVICE=txtwrite", "-o", textPath, pdfPath).CombinedOutput(); err != nil {
		t.Skipf("gs extraction failed (%v): %s", err, string(out))
	}
	text, err := os.ReadFile(textPath)
	if err != nil {
		t.Fatalf("read text: %v", err)
	}
	for _, want := range []string{"RealPdfClaim", "RealPdfReceipt", "Subtotal: Uncategorised", "$12.34"} {
		if !strings.Contains(string(text), want) {
			t.Errorf("the PDF text lacks %q", want)
		}
	}
}
