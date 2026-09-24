package services

import (
	"bytes"
	"context"
	"errors"
	"time"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"gorm.io/gorm"
	"receipt-wrangler/api/internal/commands"
	"receipt-wrangler/api/internal/env"
	"receipt-wrangler/api/internal/logging"
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/repositories"
)

// blockedExternalUrlPatterns matches all common network schemes so chromium
// refuses to load remote resources referenced from the rendered HTML, plus
// file:// so attacker-controlled email HTML cannot read local files (e.g.
// <iframe src="file:///etc/passwd"> or another group's receipt image under
// data/). The HTML is loaded via Page.setDocumentContent into an about:blank
// page (see Render), so there is no legitimate file:// load to allow. Inline
// data: URIs (base64 content commonly embedded in receipt emails) are NOT in
// this list and remain allowed.
var blockedExternalUrlPatterns = []string{
	"http://*",
	"https://*",
	"ws://*",
	"wss://*",
	"ftp://*",
	"file://*",
}

const htmlToPdfTimeout = 30 * time.Second

type HtmlToPdfService struct {
	BaseService
}

func NewHtmlToPdfService(tx *gorm.DB) HtmlToPdfService {
	return HtmlToPdfService{
		BaseService: BaseService{
			DB: repositories.GetDB(),
			TX: tx,
		},
	}
}

// Render converts the given HTML to a PDF using a fresh headless Chromium
// process. The HTML is injected into an about:blank page via
// Page.setDocumentContent (NOT navigated to as a file://), so the document has
// no local-file origin and cannot read files off disk. Network and file://
// resource loads are blocked by default for security; inline data: URIs remain
// allowed. Set CHROMIUM_ALLOW_EXTERNAL_RESOURCES=true to permit remote loads if
// you need logos or product imagery from URLs.
func (service HtmlToPdfService) Render(html string) ([]byte, commands.UpsertSystemTaskCommand, error) {
	startTime := time.Now()
	systemTaskCommand := commands.UpsertSystemTaskCommand{
		Type:                 models.HTML_TO_PDF,
		Status:               models.SYSTEM_TASK_SUCCEEDED,
		AssociatedEntityType: models.NOOP_ENTITY_TYPE,
		StartedAt:            startTime,
	}

	if len(html) == 0 {
		endTime := time.Now()
		systemTaskCommand.Status = models.SYSTEM_TASK_FAILED
		systemTaskCommand.EndedAt = &endTime
		systemTaskCommand.ResultDescription = "html content is empty"
		return nil, systemTaskCommand, errors.New("html content is empty")
	}

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.ExecPath(env.GetChromiumPath()),
		chromedp.Headless,
		chromedp.DisableGPU,
		chromedp.Flag("disable-javascript", true),
	)
	// Default behavior is --no-sandbox because the supported docker images
	// run as root, where chromium's sandbox refuses to start. Operators
	// running the API as a non-root user can opt back into the sandbox via
	// the CHROMIUM_SANDBOX env var.
	if !env.GetChromiumSandboxEnabled() {
		opts = append(opts, chromedp.NoSandbox)
	}

	allocCtx, cancelAlloc := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancelAlloc()

	browserCtx, cancelBrowser := chromedp.NewContext(allocCtx)
	defer cancelBrowser()

	timeoutCtx, cancelTimeout := context.WithTimeout(browserCtx, htmlToPdfTimeout)
	defer cancelTimeout()

	var pdfBuf []byte
	actions := []chromedp.Action{}
	// Default behavior is to block external network resources AND file:// loads:
	// receipt emails contain attacker-controllable URLs and we run chromium with
	// --no-sandbox, so disallowing these removes an SSRF / tracking-pixel surface
	// and (with the about:blank injection below) a local-file-read surface. Opt
	// back into network loads via CHROMIUM_ALLOW_EXTERNAL_RESOURCES; file:// stays
	// blocked regardless since nothing legitimate needs it.
	if env.GetChromiumAllowExternalResources() {
		actions = append(actions,
			network.Enable(),
			network.SetBlockedURLs([]string{"file://*"}),
		)
	} else {
		actions = append(actions,
			network.Enable(),
			network.SetBlockedURLs(blockedExternalUrlPatterns),
		)
	}
	actions = append(actions,
		// Inject the HTML into a blank page instead of navigating to a file://
		// URL. This avoids both the data-URL length cap (chromium silently
		// truncates multi-MB data: URLs, which large receipt emails can exceed)
		// AND a local-file origin — an about:blank document cannot read file://
		// sub-resources, so <iframe src="file:///...">/<img> can't exfiltrate
		// local files into the PDF.
		chromedp.Navigate("about:blank"),
		chromedp.ActionFunc(func(ctx context.Context) error {
			frameTree, err := page.GetFrameTree().Do(ctx)
			if err != nil {
				return err
			}
			if err := page.SetDocumentContent(frameTree.Frame.ID, html).Do(ctx); err != nil {
				return err
			}
			buf, _, err := page.PrintToPDF().WithPrintBackground(true).Do(ctx)
			if err != nil {
				return err
			}
			pdfBuf = buf
			return nil
		}),
	)
	err := chromedp.Run(timeoutCtx, actions...)

	endTime := time.Now()
	systemTaskCommand.EndedAt = &endTime
	elapsed := endTime.Sub(startTime)

	if err != nil {
		systemTaskCommand.Status = models.SYSTEM_TASK_FAILED
		systemTaskCommand.ResultDescription = err.Error()
		logging.LogStd(logging.LOG_LEVEL_ERROR, "HTML to PDF render failed: ", err.Error())
		return nil, systemTaskCommand, err
	}

	if !bytes.HasPrefix(pdfBuf, []byte("%PDF-")) {
		err = errors.New("chromedp returned non-PDF bytes")
		systemTaskCommand.Status = models.SYSTEM_TASK_FAILED
		systemTaskCommand.ResultDescription = err.Error()
		return nil, systemTaskCommand, err
	}

	systemTaskCommand.ResultDescription = "rendered " + elapsed.String()
	logging.LogStd(logging.LOG_LEVEL_INFO, "HTML to PDF render took: ", elapsed)
	return pdfBuf, systemTaskCommand, nil
}
