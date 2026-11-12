package services

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"gopkg.in/gographics/imagick.v3/imagick"

	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/utils"
)

func packetTestReceipt(id uint, amount string, categories ...models.Category) models.Receipt {
	return models.Receipt{
		BaseModel:  models.BaseModel{ID: id},
		Amount:     decimal.RequireFromString(amount),
		Categories: categories,
	}
}

func TestFoldPacketCategories(t *testing.T) {
	food := models.Category{BaseModel: models.BaseModel{ID: 1}, Name: "food"}
	travel := models.Category{BaseModel: models.BaseModel{ID: 2}, Name: "Travel"}
	apples := models.Category{BaseModel: models.BaseModel{ID: 3}, Name: "Apples"}

	receipts := []models.Receipt{
		packetTestReceipt(10, "0.10", travel, food), // multi-category
		packetTestReceipt(11, "0.20", food),
		packetTestReceipt(12, "100.05"),                     // uncategorised
		packetTestReceipt(13, "1.01", travel, travel),       // a duplicate is one category
		packetTestReceipt(14, "0.07", apples, food, travel), // multi-category
	}

	folds, total, multi := foldPacketCategories(receipts)

	type want struct {
		name    string
		indexes []int
		total   string
	}
	wants := []want{
		{"Apples", []int{4}, "0.07"},
		{"food", []int{0, 1, 4}, "0.37"},
		{"Travel", []int{0, 3, 4}, "1.18"},
		{uncategorisedLabel, []int{2}, "100.05"},
	}
	if len(folds) != len(wants) {
		t.Fatalf("got %d groups, want %d: %+v", len(folds), len(wants), folds)
	}
	sum := decimal.Zero
	for i, w := range wants {
		fold := folds[i]
		if fold.Name != w.name || fold.Subtotal.StringFixed(2) != w.total || !equalInts(fold.ReceiptIndexes, w.indexes) {
			utils.PrintTestError(t, fold, w)
		}
		sum = sum.Add(fold.Subtotal)
	}
	if total.StringFixed(2) != "101.43" {
		utils.PrintTestError(t, total.StringFixed(2), "101.43")
	}
	if multi != 2 {
		utils.PrintTestError(t, multi, 2)
	}
	// Each multi-category receipt adds its amount once per extra category.
	if sum.StringFixed(2) != "101.67" {
		utils.PrintTestError(t, sum.StringFixed(2), "101.67")
	}
}

func TestFoldPacketCategories_Empty(t *testing.T) {
	folds, total, multi := foldPacketCategories(nil)
	if len(folds) != 0 || !total.IsZero() || multi != 0 {
		utils.PrintTestError(t, folds, "nothing")
	}
}

func equalInts(left, right []int) bool {
	if len(left) != len(right) {
		return false
	}
	for i := range left {
		if left[i] != right[i] {
			return false
		}
	}
	return true
}

func TestPacketImageName(t *testing.T) {
	receipt := models.Receipt{
		BaseModel: models.BaseModel{ID: 412},
		Name:      "Costco #42 / Gas & Food",
		// Stored as midnight UTC; a zone must never shift the day.
		Date: time.Date(2026, 10, 3, 0, 0, 0, 0, time.UTC),
	}

	cases := []struct {
		file models.FileData
		want string
	}{
		{models.FileData{Name: "IMG_1.HEIC", FileType: "image/jpeg"}, "2026-10-03_Costco_42_Gas_Food_412_2.jpg"},
		{models.FileData{Name: "scan", FileType: "application/pdf"}, "2026-10-03_Costco_42_Gas_Food_412_2.pdf"},
		{models.FileData{Name: "photo.PNG", FileType: ""}, "2026-10-03_Costco_42_Gas_Food_412_2.png"},
		{models.FileData{Name: "weird.ex/e..", FileType: "nonsense"}, "2026-10-03_Costco_42_Gas_Food_412_2.bin"},
	}
	for _, c := range cases {
		if got := packetImageName(receipt, 2, c.file); got != c.want {
			utils.PrintTestError(t, got, c.want)
		}
	}

	receipt.Name = "Café 東京"
	if got := packetImageName(receipt, 1, models.FileData{FileType: "image/png"}); got != "2026-10-03_Caf_412_1.png" {
		utils.PrintTestError(t, got, "2026-10-03_Caf_412_1.png")
	}
	receipt.Name = "東京"
	if got := packetImageName(receipt, 1, models.FileData{FileType: "image/png"}); got != "2026-10-03_receipt_412_1.png" {
		utils.PrintTestError(t, got, "2026-10-03_receipt_412_1.png")
	}
	receipt.Name = strings.Repeat("a", 100)
	if got := packetImageName(receipt, 1, models.FileData{FileType: "image/png"}); got != "2026-10-03_"+strings.Repeat("a", 40)+"_412_1.png" {
		utils.PrintTestError(t, got, "a name cut to 40 characters")
	}
}

func TestClaimPacketFilename(t *testing.T) {
	now := time.Date(2026, 11, 12, 23, 30, 0, 0, time.FixedZone("EST", -5*3600))
	if got := claimPacketFilename(`PCS "Rome" / TDY`, now); got != "PCS_Rome_TDY_packet_2026-11-12.zip" {
		utils.PrintTestError(t, got, "PCS_Rome_TDY_packet_2026-11-12.zip")
	}
	if got := claimPacketFilename("???", now); got != "claim_packet_2026-11-12.zip" {
		utils.PrintTestError(t, got, "claim_packet_2026-11-12.zip")
	}
}

func TestHumanizeCode(t *testing.T) {
	for input, want := range map[string]string{"NEEDS_ATTENTION": "Needs attention", "PAYMENT": "Payment", "": "", "x": "X"} {
		if got := humanizeCode(input); got != want {
			utils.PrintTestError(t, got, want)
		}
	}
}

func TestPacketLayoutFor(t *testing.T) {
	cases := []struct {
		name          string
		width, height uint
		want          packetImageLayout
		slices        int
	}{
		{"small image is never enlarged", 400, 300, packetImageLayout{400, 300, 300}, 1},
		{"portrait photo fits the longest side", 3000, 4000, packetImageLayout{1200, 1600, 1600}, 1},
		{"landscape photo fits the longest side", 4000, 3000, packetImageLayout{1600, 1200, 1200}, 1},
		// A three-page Letter PDF at 300 DPI, stacked.
		{"stacked pdf keeps a readable width and is sliced", 2550, 9900, packetImageLayout{1230, 4775, 1599}, 3},
		{"a long till receipt is sliced at its own width", 600, 4000, packetImageLayout{600, 4000, 780}, 6},
		{"an extreme strip is capped at the slice limit", 100, 20000, packetImageLayout{100, 20000, 1667}, 12},
	}
	for _, c := range cases {
		got := packetLayoutFor(c.width, c.height)
		if got != c.want || got.sliceCount() != c.slices {
			t.Errorf("%s: got %+v (%d slices), want %+v (%d slices)", c.name, got, got.sliceCount(), c.want, c.slices)
		}
		if got.Width > packetImageMaxSide || got.SliceHeight > packetImageMaxSide && c.slices < packetMaxSlicesPerImage {
			t.Errorf("%s: a slice exceeds %d px: %+v", c.name, packetImageMaxSide, got)
		}
	}
}

func newPacketTestImage(t *testing.T, width, height uint, format string, transparent bool) []byte {
	t.Helper()
	background := imagick.NewPixelWand()
	defer background.Destroy()
	if transparent {
		background.SetColor("transparent")
	} else {
		background.SetColor("#3366cc")
	}
	mw := imagick.NewMagickWand()
	defer mw.Destroy()
	if err := mw.NewImage(width, height, background); err != nil {
		t.Fatalf("new image: %v", err)
	}
	if err := mw.SetImageFormat(format); err != nil {
		t.Fatalf("format: %v", err)
	}
	blob, err := mw.GetImageBlob()
	if err != nil && len(blob) == 0 {
		t.Fatalf("blob: %v", err)
	}
	return blob
}

func TestPacketDisplaySlices(t *testing.T) {
	cases := []struct {
		name          string
		width, height uint
		format        string
		transparent   bool
		slices        int
		maxWidth      uint
	}{
		{"jpeg photo", 3000, 4000, "jpeg", false, 1, 1200},
		{"transparent png", 500, 400, "png", true, 1, 500},
		{"stacked pdf-like jpeg", 2550, 9900, "jpeg", false, 3, 1230},
	}
	for _, c := range cases {
		slices, err := packetDisplaySlices(newPacketTestImage(t, c.width, c.height, c.format, c.transparent))
		if err != nil {
			t.Fatalf("%s: %v", c.name, err)
		}
		if len(slices) != c.slices {
			t.Errorf("%s: %d slices, want %d", c.name, len(slices), c.slices)
		}
		for _, slice := range slices {
			if !bytes.HasPrefix(slice, []byte{0xFF, 0xD8, 0xFF}) {
				t.Errorf("%s: a slice is not a JPEG", c.name)
				continue
			}
			width, height, isJpeg, err := pingImage(slice)
			if err != nil || !isJpeg || width > c.maxWidth || height > packetImageMaxSide {
				t.Errorf("%s: slice %dx%d jpeg=%v err=%v", c.name, width, height, isJpeg, err)
			}
		}
	}

	// Transparency is flattened onto white, never black.
	slices, err := packetDisplaySlices(newPacketTestImage(t, 50, 50, "png", true))
	if err != nil || len(slices) != 1 {
		t.Fatalf("transparent: %v", err)
	}
	mw := imagick.NewMagickWand()
	defer mw.Destroy()
	if err := mw.ReadImageBlob(slices[0]); err != nil {
		t.Fatalf("read: %v", err)
	}
	pixel, err := mw.GetImagePixelColor(25, 25)
	if err != nil {
		t.Fatalf("pixel: %v", err)
	}
	defer pixel.Destroy()
	if pixel.GetRed() < 0.95 || pixel.GetGreen() < 0.95 || pixel.GetBlue() < 0.95 {
		t.Errorf("transparency came out as %v,%v,%v, not white", pixel.GetRed(), pixel.GetGreen(), pixel.GetBlue())
	}

	if _, err := packetDisplaySlices([]byte("not an image")); err == nil {
		t.Errorf("garbage bytes should fail, so the packet notes the image instead")
	}
}

func TestCheckClaimPacketBounds(t *testing.T) {
	withImages := func(count int) []models.Receipt {
		return []models.Receipt{{ImageFiles: make([]models.FileData, count)}}
	}
	if err := checkClaimPacketBounds(withImages(MaxClaimPacketImages)); err != nil {
		t.Errorf("at the bound: %v", err)
	}
	var tooLarge ClaimPacketTooLargeError
	err := checkClaimPacketBounds(withImages(MaxClaimPacketImages + 1))
	if err == nil || !strings.Contains(err.Error(), "at most") {
		utils.PrintTestError(t, err, "an image bound error")
	} else if _, ok := err.(ClaimPacketTooLargeError); !ok {
		utils.PrintTestError(t, err, tooLarge)
	}
	if err := checkClaimPacketBounds(make([]models.Receipt, maxClaimPacketReceipts+1)); err == nil {
		t.Errorf("over the receipt bound should fail")
	}
}

func TestReceiptCsv_ImagesColumnOnlyWhenAsked(t *testing.T) {
	receipts := []models.Receipt{{BaseModel: models.BaseModel{ID: 7}, Name: "r", Amount: decimal.NewFromInt(1), Status: models.OPEN}}

	plain := NewReceiptCsvService()
	plain.Location = time.UTC
	result, err := plain.BuildReceiptCsv(receipts)
	if err != nil {
		t.Fatalf("csv: %v", err)
	}
	if strings.Contains(string(result.ReceiptCsvBytes), "Images") {
		t.Errorf("the normal export gained an Images column")
	}

	packet := NewReceiptCsvService()
	packet.Location = time.UTC
	packet.ImageNames = map[uint][]string{7: {"a.jpg", "b.pdf"}}
	result, err = packet.BuildReceiptCsv(receipts)
	if err != nil {
		t.Fatalf("csv: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(string(result.ReceiptCsvBytes)), "\n")
	if !strings.HasSuffix(lines[0], ",Images") || !strings.HasSuffix(lines[1], `,"a.jpg,b.pdf"`) {
		utils.PrintTestError(t, lines, "an Images column listing a.jpg,b.pdf")
	}
}
