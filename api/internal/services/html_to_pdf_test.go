package services

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"

	"gopkg.in/gographics/imagick.v3/imagick"
)

func TestHtmlToPdfService_Render_BasicHtml(t *testing.T) {
	service := NewHtmlToPdfService(nil)

	html := `<!DOCTYPE html>
<html>
<head><title>Test Receipt</title></head>
<body>
  <h1>Receipt #1234</h1>
  <p>Total: $12.34</p>
</body>
</html>`

	pdfBytes, taskCmd, err := service.Render(html)
	if err != nil {
		t.Fatalf("Render returned error: %v", err)
	}

	if len(pdfBytes) == 0 {
		t.Fatal("expected non-empty PDF bytes, got empty slice")
	}

	if !bytes.HasPrefix(pdfBytes, []byte("%PDF-")) {
		t.Fatalf("expected PDF bytes to start with %%PDF-, got: %q", string(pdfBytes[:min(20, len(pdfBytes))]))
	}

	if taskCmd.Status != "SUCCEEDED" {
		t.Errorf("expected system task status SUCCEEDED, got %s", taskCmd.Status)
	}

	if taskCmd.Type != "HTML_TO_PDF" {
		t.Errorf("expected system task type HTML_TO_PDF, got %s", taskCmd.Type)
	}

	if taskCmd.EndedAt == nil {
		t.Error("expected EndedAt to be set on success")
	}
}

func TestHtmlToPdfService_Render_EmptyHtmlFails(t *testing.T) {
	service := NewHtmlToPdfService(nil)

	pdfBytes, taskCmd, err := service.Render("")
	if err == nil {
		t.Fatal("expected error for empty HTML, got nil")
	}

	if pdfBytes != nil {
		t.Errorf("expected nil PDF bytes on error, got %d bytes", len(pdfBytes))
	}

	if taskCmd.Status != "FAILED" {
		t.Errorf("expected system task status FAILED, got %s", taskCmd.Status)
	}

	if !strings.Contains(taskCmd.ResultDescription, "empty") {
		t.Errorf("expected ResultDescription to mention empty, got %q", taskCmd.ResultDescription)
	}
}

// Security regression for the email HTML->PDF local-file-read (file://) bug.
// Attacker-controlled email HTML must NOT be able to pull local files into the
// rendered PDF via a file:// sub-resource. Renders HTML containing a file://
// iframe pointing at a canary file and asserts the canary text is absent from
// the PDF while the legitimate body text is present. Uses ghostscript to extract
// PDF text; skips (does not fail) where gs is unavailable.
func TestHtmlToPdfService_Render_BlocksLocalFileRead(t *testing.T) {
	gs, err := exec.LookPath("gs")
	if err != nil {
		t.Skip("ghostscript (gs) not available; skipping PDF text extraction")
	}

	const canary = "TOPSECRETCANARY9f3a2b"
	const marker = "LEGITIMATEBODYMARKER"
	canaryPath := filepath.Join(t.TempDir(), "secret.txt")
	if err := os.WriteFile(canaryPath, []byte(canary), 0644); err != nil {
		t.Fatalf("failed to write canary file: %v", err)
	}

	html := `<!DOCTYPE html><html><body><h1>` + marker + `</h1>` +
		`<iframe src="file://` + canaryPath + `" width="600" height="200"></iframe>` +
		`<img src="file://` + canaryPath + `">` +
		`</body></html>`

	service := NewHtmlToPdfService(nil)
	pdfBytes, taskCmd, err := service.Render(html)
	if err != nil {
		t.Fatalf("render failed: %v (%s)", err, taskCmd.ResultDescription)
	}

	pdfPath := filepath.Join(t.TempDir(), "out.pdf")
	if err := os.WriteFile(pdfPath, pdfBytes, 0644); err != nil {
		t.Fatalf("failed to write pdf: %v", err)
	}
	txtPath := filepath.Join(t.TempDir(), "out.txt")
	cmd := exec.Command(gs, "-q", "-dNOPAUSE", "-dBATCH", "-sDEVICE=txtwrite", "-o", txtPath, pdfPath)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Skipf("gs extraction failed (%v): %s", err, string(out))
	}
	extracted, err := os.ReadFile(txtPath)
	if err != nil {
		t.Fatalf("failed to read extracted text: %v", err)
	}
	text := string(extracted)

	if strings.Contains(text, canary) {
		t.Fatalf("SECURITY: local file contents leaked into the PDF: %q", text)
	}
	// Sanity: the render actually produced our document (so the negative result
	// above isn't just an empty/failed render).
	if !strings.Contains(text, marker) {
		t.Fatalf("expected the legitimate body marker in the PDF, got: %q", text)
	}
}

// Render must embed inline images. The email pipeline feeds user/attacker email
// BodyHtml (which commonly embeds images as data: URIs — the only image kind that
// survives the default external-resource block) straight into Render, and the
// resulting PDF is rasterized for OCR/vision, so a dropped image loses content.
// This proves a data: image is painted into the PDF. Rasterizing a PDF needs the
// ghostscript delegate, so skip (do not fail) where gs is unavailable.
func TestHtmlToPdfService_Render_EmbedsDataImage(t *testing.T) {
	if _, err := exec.LookPath("gs"); err != nil {
		t.Skip("ghostscript (gs) not available; skipping PDF rasterization")
	}

	dataURI := "data:image/png;base64," + base64.StdEncoding.EncodeToString(solidRedPng(t))
	html := `<!DOCTYPE html><html><body style="margin:0">` +
		`<img src="` + dataURI + `" style="width:400px;height:400px">` +
		`</body></html>`

	if !hasRedPixel(rasterizePdfPage(t, renderOrFatal(t, html))) {
		t.Fatal("expected the embedded data: image to be painted into the PDF, but found no red pixels")
	}
}

// In external-resources mode the render navigates a loopback http origin and
// intercepts sub-resource requests. This proves the interception ALLOWS the page
// through (the render produces a PDF at all → the main navigation was continued)
// and still embeds inline data: images.
func TestHtmlToPdfService_Render_ExternalMode_EmbedsDataImage(t *testing.T) {
	if _, err := exec.LookPath("gs"); err != nil {
		t.Skip("ghostscript (gs) not available; skipping PDF rasterization")
	}
	t.Setenv("CHROMIUM_ALLOW_EXTERNAL_RESOURCES", "true")

	dataURI := "data:image/png;base64," + base64.StdEncoding.EncodeToString(solidRedPng(t))
	html := `<!DOCTYPE html><html><body style="margin:0">` +
		`<img src="` + dataURI + `" style="width:400px;height:400px">` +
		`</body></html>`

	if !hasRedPixel(rasterizePdfPage(t, renderOrFatal(t, html))) {
		t.Fatal("expected the data: image in the PDF (external mode): interception must continue the page load and allow data: URIs")
	}
}

// C1 regression: in external-resources mode, attacker-controlled email HTML must
// NOT be able to make the render fetch an internal/loopback host (SSRF) or bake
// its response into the PDF. Points an iframe + img at a mock "internal" server on
// 127.0.0.1 and asserts it is NEVER reached and its secret never appears in the
// PDF. (The render's own loopback PAGE server is a different host:port and is the
// only loopback origin the policy allows.)
func TestHtmlToPdfService_Render_ExternalMode_BlocksInternalSSRF(t *testing.T) {
	gs, err := exec.LookPath("gs")
	if err != nil {
		t.Skip("ghostscript (gs) not available; skipping PDF text extraction")
	}
	t.Setenv("CHROMIUM_ALLOW_EXTERNAL_RESOURCES", "true")

	const secret = "INTERNAL_SSRF_CANARY_5e1f8a"
	pngBytes := solidRedPng(t)
	var hits int32
	internal := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&hits, 1)
		if r.URL.Path == "/pixel" {
			w.Header().Set("Content-Type", "image/png")
			_, _ = w.Write(pngBytes)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprintf(w, "<html><body><h1>%s</h1></body></html>", secret)
	}))
	defer internal.Close()

	html := `<!DOCTYPE html><html><body>` +
		`<iframe src="` + internal.URL + `/iframe" width="600" height="200"></iframe>` +
		`<img src="` + internal.URL + `/pixel">` +
		`</body></html>`

	pdfBytes := renderOrFatal(t, html)

	pdfPath := filepath.Join(t.TempDir(), "out.pdf")
	if err := os.WriteFile(pdfPath, pdfBytes, 0644); err != nil {
		t.Fatalf("failed to write pdf: %v", err)
	}
	txtPath := filepath.Join(t.TempDir(), "out.txt")
	if out, err := exec.Command(gs, "-q", "-dNOPAUSE", "-dBATCH", "-sDEVICE=txtwrite", "-o", txtPath, pdfPath).CombinedOutput(); err != nil {
		t.Skipf("gs extraction failed (%v): %s", err, string(out))
	}
	extracted, err := os.ReadFile(txtPath)
	if err != nil {
		t.Fatalf("failed to read extracted text: %v", err)
	}

	if strings.Contains(string(extracted), secret) {
		t.Fatalf("SECURITY: internal response leaked into the PDF: %q", string(extracted))
	}
	if h := atomic.LoadInt32(&hits); h != 0 {
		t.Fatalf("SECURITY: internal server was reached %d time(s); SSRF not blocked", h)
	}
}

// TestIsRequestAllowed is the hermetic (no-browser) guard for the SSRF policy.
// Public hosts load; loopback/private/link-local/metadata hosts, non-http(s)
// schemes, and everything but the render's own loopback page are refused.
func TestIsRequestAllowed(t *testing.T) {
	const loopbackAddr = "127.0.0.1:5000" // the render's own page origin

	tests := map[string]struct {
		rawURL string
		want   bool
	}{
		"our loopback page":         {"http://127.0.0.1:5000/", true},
		"our loopback sub-resource": {"http://127.0.0.1:5000/logo.png", true},
		"data URI":                  {"data:image/png;base64,AAAA", true},
		"public IPv4":               {"https://93.184.216.34/logo.png", true},
		"public DNS root IPv4":      {"http://8.8.8.8/x", true},
		"cloud metadata":            {"http://169.254.169.254/latest/meta-data/", false},
		"other loopback port":       {"http://127.0.0.1:6379/", false},
		"private 10/8":              {"http://10.0.0.5/", false},
		"private 192.168/16":        {"http://192.168.1.1/", false},
		"private 172.16/12":         {"http://172.16.0.1/", false},
		"loopback IPv6":             {"http://[::1]/", false},
		"link-local IPv6":           {"http://[fe80::1]/", false},
		"localhost name":            {"http://localhost:8081/api", false},
		"unspecified":               {"http://0.0.0.0/", false},
		"file scheme":               {"file:///etc/passwd", false},
		"ftp scheme":                {"ftp://93.184.216.34/x", false},
		"missing scheme":            {"://bad", false},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			if got := isRequestAllowed(tt.rawURL, loopbackAddr); got != tt.want {
				t.Fatalf("isRequestAllowed(%q) = %v, want %v", tt.rawURL, got, tt.want)
			}
		})
	}
}

// Security regression for the local-file-read (file://) bug, exercised on the
// NEW external-resources path (loopback http origin). Even with remote resources
// allowed, attacker-controlled email HTML must not pull local files into the PDF
// via a file:// sub-resource. Mirrors TestHtmlToPdfService_Render_BlocksLocalFileRead
// but with CHROMIUM_ALLOW_EXTERNAL_RESOURCES enabled, so it guards the loopback
// path specifically.
func TestHtmlToPdfService_Render_ExternalMode_BlocksLocalFileRead(t *testing.T) {
	gs, err := exec.LookPath("gs")
	if err != nil {
		t.Skip("ghostscript (gs) not available; skipping PDF text extraction")
	}
	t.Setenv("CHROMIUM_ALLOW_EXTERNAL_RESOURCES", "true")

	const canary = "EXTERNALCANARY7c1d4e"
	const marker = "EXTERNALBODYMARKER"
	canaryPath := filepath.Join(t.TempDir(), "secret.txt")
	if err := os.WriteFile(canaryPath, []byte(canary), 0644); err != nil {
		t.Fatalf("failed to write canary file: %v", err)
	}

	html := `<!DOCTYPE html><html><body><h1>` + marker + `</h1>` +
		`<iframe src="file://` + canaryPath + `" width="600" height="200"></iframe>` +
		`<img src="file://` + canaryPath + `">` +
		`</body></html>`

	pdfBytes := renderOrFatal(t, html)

	pdfPath := filepath.Join(t.TempDir(), "out.pdf")
	if err := os.WriteFile(pdfPath, pdfBytes, 0644); err != nil {
		t.Fatalf("failed to write pdf: %v", err)
	}
	txtPath := filepath.Join(t.TempDir(), "out.txt")
	cmd := exec.Command(gs, "-q", "-dNOPAUSE", "-dBATCH", "-sDEVICE=txtwrite", "-o", txtPath, pdfPath)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Skipf("gs extraction failed (%v): %s", err, string(out))
	}
	extracted, err := os.ReadFile(txtPath)
	if err != nil {
		t.Fatalf("failed to read extracted text: %v", err)
	}
	text := string(extracted)

	if strings.Contains(text, canary) {
		t.Fatalf("SECURITY: local file contents leaked into the PDF (external mode): %q", text)
	}
	if !strings.Contains(text, marker) {
		t.Fatalf("expected the legitimate body marker in the PDF, got: %q", text)
	}
}

// solidRedPng returns the PNG bytes of a small solid pure-red image, generated
// in-test so there is no binary fixture to maintain.
func solidRedPng(t *testing.T) []byte {
	t.Helper()
	src := image.NewRGBA(image.Rect(0, 0, 32, 32))
	red := color.RGBA{R: 255, G: 0, B: 0, A: 255}
	for y := 0; y < 32; y++ {
		for x := 0; x < 32; x++ {
			src.Set(x, y, red)
		}
	}
	var buf bytes.Buffer
	if err := png.Encode(&buf, src); err != nil {
		t.Fatalf("failed to encode source png: %v", err)
	}
	return buf.Bytes()
}

// renderOrFatal renders html to a PDF and fails the test unless it gets PDF bytes
// back. The image is scaled to 400x400 by the callers so it covers a large,
// easy-to-find region of the rasterized page.
func renderOrFatal(t *testing.T, html string) []byte {
	t.Helper()
	pdfBytes, taskCmd, err := NewHtmlToPdfService(nil).Render(html)
	if err != nil {
		t.Fatalf("render failed: %v (%s)", err, taskCmd.ResultDescription)
	}
	if !bytes.HasPrefix(pdfBytes, []byte("%PDF-")) {
		t.Fatalf("expected PDF bytes, got: %q", string(pdfBytes[:min(20, len(pdfBytes))]))
	}
	return pdfBytes
}

// rasterizePdfPage renders the first page of a PDF to an image via ImageMagick
// (which delegates PDF reading to ghostscript). PNG is used as the intermediate
// format so colors are preserved exactly (unlike lossy JPEG).
func rasterizePdfPage(t *testing.T, pdf []byte) image.Image {
	t.Helper()
	mw := imagick.NewMagickWand()
	defer mw.Destroy()

	if err := mw.SetResolution(150, 150); err != nil {
		t.Fatalf("SetResolution: %v", err)
	}
	if err := mw.ReadImageBlob(pdf); err != nil {
		t.Fatalf("ReadImageBlob: %v", err)
	}
	mw.SetIteratorIndex(0)
	if err := mw.SetImageFormat("png"); err != nil {
		t.Fatalf("SetImageFormat: %v", err)
	}
	blob, err := mw.GetImageBlob()
	if err != nil {
		t.Fatalf("GetImageBlob: %v", err)
	}
	img, err := png.Decode(bytes.NewReader(blob))
	if err != nil {
		t.Fatalf("failed to decode rasterized page: %v", err)
	}
	return img
}

// hasRedPixel reports whether the image contains a strongly-red pixel. RGBA()
// returns 16-bit channel values, so shift down to 8-bit before comparing.
func hasRedPixel(img image.Image) bool {
	b := img.Bounds()
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			r, g, bl, _ := img.At(x, y).RGBA()
			if r>>8 > 200 && g>>8 < 80 && bl>>8 < 80 {
				return true
			}
		}
	}
	return false
}
