package repositories

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"os"
	"strings"
	"testing"

	"gopkg.in/gographics/imagick.v3/imagick"
)

// ---------- PDF rasterization (ConvertPdfToJpg) ----------
//
// The happy-path, invalid-bytes and DPI tests live in files_test.go. These pin what the
// page-at-a-time rewrite has to get right: transparency, colorspace, page geometry, the
// height limit and the temp files it spills.

// testPdfPage is one page of a hand-built PDF, sized in points (1/72 in), so a page renders
// to exactly its own size at 72 DPI.
type testPdfPage struct {
	width, height float64
	rotate        int
	// content is the raw content stream. Font /F1 (Helvetica) is always available, and
	// ExtGState /GS1 (50% fill opacity) when translucent is set.
	content     string
	translucent bool
}

// buildTestPdf assembles a minimal valid PDF, so a test can pin the page geometry and
// transparency features of a real PDF without a binary fixture.
func buildTestPdf(t *testing.T, pages ...testPdfPage) []byte {
	t.Helper()

	objects := []string{"<< /Type /Catalog /Pages 2 0 R >>", ""}
	kids := make([]string, 0, len(pages))
	for _, page := range pages {
		pageId := len(objects) + 1
		kids = append(kids, fmt.Sprintf("%d 0 R", pageId))

		resources := "/Font << /F1 << /Type /Font /Subtype /Type1 /BaseFont /Helvetica >> >>"
		if page.translucent {
			resources += " /ExtGState << /GS1 << /Type /ExtGState /ca 0.5 >> >>"
		}
		rotate := ""
		if page.rotate != 0 {
			rotate = fmt.Sprintf(" /Rotate %d", page.rotate)
		}

		objects = append(objects,
			fmt.Sprintf("<< /Type /Page /Parent 2 0 R /MediaBox [0 0 %g %g]%s /Contents %d 0 R /Resources << %s >> >>",
				page.width, page.height, rotate, pageId+1, resources),
			fmt.Sprintf("<< /Length %d >>\nstream\n%s\nendstream", len(page.content), page.content),
		)
	}
	objects[1] = fmt.Sprintf("<< /Type /Pages /Kids [%s] /Count %d >>", strings.Join(kids, " "), len(pages))

	var pdf bytes.Buffer
	pdf.WriteString("%PDF-1.4\n")
	offsets := make([]int, len(objects))
	for i, object := range objects {
		offsets[i] = pdf.Len()
		fmt.Fprintf(&pdf, "%d 0 obj\n%s\nendobj\n", i+1, object)
	}
	xref := pdf.Len()
	fmt.Fprintf(&pdf, "xref\n0 %d\n0000000000 65535 f \n", len(objects)+1)
	for _, offset := range offsets {
		fmt.Fprintf(&pdf, "%010d 00000 n \n", offset)
	}
	fmt.Fprintf(&pdf, "trailer << /Size %d /Root 1 0 R >>\nstartxref\n%d\n%%%%EOF\n", len(objects)+1, xref)
	return pdf.Bytes()
}

// cornerMarker fills a 40pt square at the page's top-left corner, so a test can find where
// each page landed in the stacked image.
func cornerMarker(pageHeight float64, rgb string) string {
	return fmt.Sprintf("%s rg 0 %g 40 40 re f", rgb, pageHeight-40)
}

// newSolidImage returns a w x h image filled with color ("none" is fully transparent).
func newSolidImage(t *testing.T, width, height uint, fill string) *imagick.MagickWand {
	t.Helper()
	background := imagick.NewPixelWand()
	defer background.Destroy()
	if !background.SetColor(fill) {
		t.Fatalf("SetColor(%q) failed", fill)
	}

	mw := imagick.NewMagickWand()
	t.Cleanup(mw.Destroy)
	if err := mw.NewImage(width, height, background); err != nil {
		t.Fatalf("NewImage: %v", err)
	}
	return mw
}

// pdfFromImage wraps an image in a one-page PDF at 72 DPI, so one pixel is one point and the
// page renders back to the image's own size at 72 DPI.
func pdfFromImage(t *testing.T, mw *imagick.MagickWand) []byte {
	t.Helper()
	if err := mw.SetImageResolution(72, 72); err != nil {
		t.Fatalf("SetImageResolution: %v", err)
	}
	if err := mw.SetImageUnits(imagick.RESOLUTION_PIXELS_PER_INCH); err != nil {
		t.Fatalf("SetImageUnits: %v", err)
	}
	if err := mw.SetImageFormat("pdf"); err != nil {
		t.Fatalf("SetImageFormat: %v", err)
	}
	pdf, err := mw.GetImageBlob()
	if err != nil {
		t.Fatalf("GetImageBlob: %v", err)
	}
	return pdf
}

// usePdfConversion points temp/ at a fresh directory and sets the rasterization DPI.
func usePdfConversion(t *testing.T, dpi int) FileRepository {
	t.Helper()
	t.Setenv("BASE_PATH", t.TempDir())
	t.Cleanup(TruncateTestDb)
	setPdfDpiSetting(t, dpi)
	return NewFileRepository(nil)
}

func convertPdf(t *testing.T, repository FileRepository, pdf []byte) image.Image {
	t.Helper()
	out, err := repository.ConvertPdfToJpg(pdf)
	if err != nil {
		t.Fatalf("ConvertPdfToJpg: %v", err)
	}
	img, format, err := image.Decode(bytes.NewReader(out))
	if err != nil {
		t.Fatalf("decode output: %v", err)
	}
	if format != "jpeg" {
		t.Fatalf("output format = %q, want jpeg", format)
	}
	return img
}

func rgbAt(img image.Image, x, y int) color.RGBA {
	return color.RGBAModel.Convert(img.At(x, y)).(color.RGBA)
}

// assertColorNear checks each channel within tolerance; JPEG is lossy, so exact matches
// are not available.
func assertColorNear(t *testing.T, img image.Image, x, y int, want color.RGBA, tolerance uint8) {
	t.Helper()
	got := rgbAt(img, x, y)
	near := func(a, b uint8) bool {
		if a > b {
			return a-b <= tolerance
		}
		return b-a <= tolerance
	}
	if !near(got.R, want.R) || !near(got.G, want.G) || !near(got.B, want.B) {
		t.Errorf("pixel (%d,%d) = rgb(%d,%d,%d), want rgb(%d,%d,%d) ±%d",
			x, y, got.R, got.G, got.B, want.R, want.G, want.B, tolerance)
	}
}

func assertSize(t *testing.T, img image.Image, width, height int) {
	t.Helper()
	if got := img.Bounds().Size(); got.X != width || got.Y != height {
		t.Fatalf("output is %dx%d, want %dx%d", got.X, got.Y, width, height)
	}
}

func assertTempDirectoryEmpty(t *testing.T, repository FileRepository) {
	t.Helper()
	entries, err := os.ReadDir(repository.GetTempDirectoryPath())
	if err != nil && !os.IsNotExist(err) {
		t.Fatalf("read temp dir: %v", err)
	}
	for _, entry := range entries {
		t.Errorf("left %s behind in temp/", entry.Name())
	}
}

// limitImageMagickHeight lowers ImageMagick's height limit for one test.
func limitImageMagickHeight(t *testing.T, limit uint64) {
	t.Helper()
	previous := imagick.GetResourceLimit(imageMagickHeightResource)
	if !imagick.SetResourceLimit(imageMagickHeightResource, limit) {
		t.Fatalf("SetResourceLimit(height, %d) refused", limit)
	}
	// GetResourceLimit reports "unlimited" as -1; converting back restores the same value.
	t.Cleanup(func() { imagick.SetResourceLimit(imageMagickHeightResource, uint64(previous)) })
}

var (
	rgbWhite = color.RGBA{R: 255, G: 255, B: 255}
	rgbRed   = color.RGBA{R: 255}
	rgbGreen = color.RGBA{G: 255}
	rgbBlue  = color.RGBA{B: 255}
)

// The black-block bug: an image with a transparent background came out as a black page,
// because the alpha channel was dropped rather than flattened onto white.
func TestConvertPdfToJpg_TransparentAreasRenderWhite(t *testing.T) {
	repository := usePdfConversion(t, 72)

	page := newSolidImage(t, 200, 200, "none")
	if err := page.CompositeImage(newSolidImage(t, 80, 80, "blue"), imagick.COMPOSITE_OP_OVER, true, 60, 60); err != nil {
		t.Fatalf("CompositeImage: %v", err)
	}

	img := convertPdf(t, repository, pdfFromImage(t, page))

	assertSize(t, img, 200, 200)
	assertColorNear(t, img, 10, 10, rgbWhite, 10)
	assertColorNear(t, img, 190, 190, rgbWhite, 10)
	assertColorNear(t, img, 100, 100, rgbBlue, 40)
}

// The other shape of the black-block bug: a vector fill at 50% opacity (a highlight) takes
// Ghostscript through its transparency compositor.
func TestConvertPdfToJpg_TranslucentFillBlendsWithWhite(t *testing.T) {
	repository := usePdfConversion(t, 72)

	pdf := buildTestPdf(t, testPdfPage{
		width: 300, height: 200, translucent: true,
		content: "0 0 0 rg BT /F1 20 Tf 20 160 Td (RECEIPT) Tj ET " +
			"/GS1 gs 1 1 0 rg 20 20 260 60 re f",
	})

	img := convertPdf(t, repository, pdf)

	assertSize(t, img, 300, 200)
	// Half-opaque yellow over white.
	assertColorNear(t, img, 150, 150, color.RGBA{R: 255, G: 255, B: 128}, 40)
	// Between the heading and the highlight: page background.
	assertColorNear(t, img, 150, 90, rgbWhite, 10)
	assertColorNear(t, img, 280, 15, rgbWhite, 10)
}

// A CMYK PDF used to come out as a CMYK JPEG. Pages are composited onto an sRGB canvas, so
// they must be converted first or their channels land as garbage.
func TestConvertPdfToJpg_CmykPageBecomesRgb(t *testing.T) {
	repository := usePdfConversion(t, 72)

	page := newSolidImage(t, 100, 100, "red")
	if err := page.TransformImageColorspace(imagick.COLORSPACE_CMYK); err != nil {
		t.Fatalf("TransformImageColorspace: %v", err)
	}

	img := convertPdf(t, repository, pdfFromImage(t, page))

	if _, isCmyk := img.(*image.CMYK); isCmyk {
		t.Errorf("output is a CMYK JPEG, want RGB or grayscale")
	}
	assertColorNear(t, img, 50, 50, rgbRed, 40)
}

// Pages differ in size (a phone scanner's auto-crop does this). ImageMagick's ping reports
// the largest page for every page, so each page must be measured from its real render.
func TestConvertPdfToJpg_StacksMixedPageSizesInOrder(t *testing.T) {
	repository := usePdfConversion(t, 72)

	pdf := buildTestPdf(t,
		testPdfPage{width: 612, height: 792, content: cornerMarker(792, "1 0 0")},
		testPdfPage{width: 300, height: 400, content: cornerMarker(400, "0 1 0")},
		testPdfPage{width: 612, height: 200, content: cornerMarker(200, "0 0 1")},
	)

	img := convertPdf(t, repository, pdf)

	assertSize(t, img, 612, 792+400+200)
	assertColorNear(t, img, 20, 20, rgbRed, 40)
	assertColorNear(t, img, 20, 792+20, rgbGreen, 40)
	assertColorNear(t, img, 20, 792+400+20, rgbBlue, 40)
	// The narrower middle page is padded with white, not left black or transparent.
	assertColorNear(t, img, 450, 792+200, rgbWhite, 10)
}

func TestConvertPdfToJpg_RotatedPageIsNotClipped(t *testing.T) {
	repository := usePdfConversion(t, 72)

	pdf := buildTestPdf(t, testPdfPage{width: 200, height: 300, rotate: 90, content: cornerMarker(300, "1 0 0")})

	img := convertPdf(t, repository, pdf)

	assertSize(t, img, 300, 200)
}

// Tesseract reads the image's density, and the appended output used to carry the render DPI.
func TestConvertPdfToJpg_KeepsTheRenderResolution(t *testing.T) {
	repository := usePdfConversion(t, 150)
	pdf := buildTestPdf(t, testPdfPage{width: 72, height: 72, content: cornerMarker(72, "1 0 0")})

	out, err := repository.ConvertPdfToJpg(pdf)
	if err != nil {
		t.Fatalf("ConvertPdfToJpg: %v", err)
	}

	mw := imagick.NewMagickWand()
	defer mw.Destroy()
	if err := mw.ReadImageBlob(out); err != nil {
		t.Fatalf("ReadImageBlob: %v", err)
	}
	if w, h := mw.GetImageWidth(), mw.GetImageHeight(); w != 150 || h != 150 {
		t.Errorf("output is %dx%d, want 150x150", w, h)
	}
	x, y, err := mw.GetImageResolution()
	if err != nil {
		t.Fatalf("GetImageResolution: %v", err)
	}
	if x != 150 || y != 150 {
		t.Errorf("resolution = %gx%g, want 150x150", x, y)
	}
	if depth := mw.GetImageDepth(); depth != 8 {
		t.Errorf("depth = %d, want 8", depth)
	}
}

func TestConvertPdfToJpg_LeavesNoTempFiles(t *testing.T) {
	t.Run("after a conversion", func(t *testing.T) {
		repository := usePdfConversion(t, 72)
		pdf := buildTestPdf(t,
			testPdfPage{width: 100, height: 100, content: cornerMarker(100, "1 0 0")},
			testPdfPage{width: 100, height: 100, content: cornerMarker(100, "0 1 0")},
		)

		if _, err := repository.ConvertPdfToJpg(pdf); err != nil {
			t.Fatalf("ConvertPdfToJpg: %v", err)
		}
		assertTempDirectoryEmpty(t, repository)
	})

	t.Run("after a pdf that will not render", func(t *testing.T) {
		repository := usePdfConversion(t, 72)

		if _, err := repository.ConvertPdfToJpg([]byte("%PDF-1.4\nthis is not really a pdf\n")); err == nil {
			t.Fatal("expected an error for a corrupt pdf")
		}
		assertTempDirectoryEmpty(t, repository)
	})
}

func TestConvertPdfToJpg_RejectsDocumentsTallerThanTheLimit(t *testing.T) {
	repository := usePdfConversion(t, 72)
	limitImageMagickHeight(t, 1000)

	page := testPdfPage{width: 100, height: 600, content: cornerMarker(600, "1 0 0")}
	_, err := repository.ConvertPdfToJpg(buildTestPdf(t, page, page, page))

	if err == nil {
		t.Fatal("expected an error for a pdf taller than the limit")
	}
	if !strings.Contains(err.Error(), "too tall") || !strings.Contains(err.Error(), "1000") {
		t.Errorf("error %q should say the pdf is too tall and name the limit", err)
	}
	assertTempDirectoryEmpty(t, repository)
}

// imagick.v3 has no constant for MagickCore's HeightResource, so the code names it by its
// position in the enum. Prove that position really is the height limit.
func TestImageMagickHeightResource_IsTheHeightLimit(t *testing.T) {
	limitImageMagickHeight(t, 100)

	white := imagick.NewPixelWand()
	defer white.Destroy()
	white.SetColor("white")

	within := imagick.NewMagickWand()
	defer within.Destroy()
	if err := within.NewImage(1, 50, white); err != nil {
		t.Errorf("a 50px-tall image under a 100px limit failed: %v", err)
	}

	over := imagick.NewMagickWand()
	defer over.Destroy()
	if err := over.NewImage(1, 200, white); err == nil {
		t.Error("a 200px-tall image under a 100px limit was allowed; imageMagickHeightResource is not the height limit")
	}
}

// Driven with plain values: policy.xml's height limit is a ceiling that SetResourceLimit
// cannot raise, so the JPEG-cap branch is unreachable through ImageMagick itself.
func TestPdfImageHeightLimit(t *testing.T) {
	tests := map[string]struct {
		imageMagickLimit int64
		want             uint
	}{
		"policy.xml's 32,000px limit applies":       {imageMagickLimit: 32000, want: 32000},
		"a limit above the jpeg maximum is capped":  {imageMagickLimit: 1_000_000, want: maxJpegDimension},
		"exactly the jpeg maximum":                  {imageMagickLimit: maxJpegDimension, want: maxJpegDimension},
		"unlimited (-1) falls back to the jpeg max": {imageMagickLimit: -1, want: maxJpegDimension},
		"zero falls back to the jpeg max":           {imageMagickLimit: 0, want: maxJpegDimension},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			if got := pdfImageHeightLimit(test.imageMagickLimit); got != test.want {
				t.Errorf("pdfImageHeightLimit(%d) = %d, want %d", test.imageMagickLimit, got, test.want)
			}
		})
	}
}

func TestMaxPdfImageHeight_FollowsImageMagicksLimit(t *testing.T) {
	limitImageMagickHeight(t, 1000)
	if got := maxPdfImageHeight(); got != 1000 {
		t.Errorf("maxPdfImageHeight() = %d under a 1000px ImageMagick limit, want 1000", got)
	}
}
