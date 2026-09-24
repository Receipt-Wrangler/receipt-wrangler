package services

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
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
