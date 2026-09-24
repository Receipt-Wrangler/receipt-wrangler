package services

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/fetch"
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
// process. Network and file:// resource loads are blocked by default for
// security; inline data: URIs remain allowed. In the default mode the HTML is
// injected into an about:blank page via Page.setDocumentContent, so the document
// has no local-file origin and cannot read files off disk. Set
// CHROMIUM_ALLOW_EXTERNAL_RESOURCES=true to permit remote loads (logos / product
// imagery): the HTML is then served from an ephemeral loopback HTTP server and
// navigated to, giving it a real http origin so remote sub-resources load. In
// that mode sub-resource requests are intercepted and only PUBLIC hosts are
// allowed — loopback/private/link-local targets (e.g. 127.0.0.1:*,
// 169.254.169.254, RFC1918) and file:// are denied, so attacker-controlled email
// HTML cannot use the render as an SSRF into internal services.
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
	printAction := chromedp.ActionFunc(func(ctx context.Context) error {
		buf, _, err := page.PrintToPDF().WithPrintBackground(true).Do(ctx)
		if err != nil {
			return err
		}
		pdfBuf = buf
		return nil
	})

	var actions []chromedp.Action
	if env.GetChromiumAllowExternalResources() {
		// External resources are opt-in (remote logos / product imagery). Serve the
		// HTML from an ephemeral loopback HTTP server and navigate to it so the
		// document has a real http origin — an about:blank document filled via
		// Page.setDocumentContent never fetches http(s) sub-resources at all, so
		// remote images only load via a navigated origin.
		//
		// The email body is attacker-controlled, so allowing arbitrary network loads
		// is an SSRF vector: <iframe src="http://169.254.169.254/…"> can bake cloud
		// metadata into the PDF, and http://127.0.0.1:<port> reaches loopback
		// services. A URL blocklist can't express "allow our own loopback page but
		// deny every OTHER internal host", nor catch a hostname that resolves to an
		// internal IP — so requests are intercepted via the Fetch domain and each is
		// allowed/denied by its resolved address (isRequestAllowed). file:// is
		// denied there too. chromedp.Navigate waits for the load event, so allowed
		// remote images finish loading before PrintToPDF.
		server, addr, serveErr := startLoopbackHtmlServer(html)
		if serveErr != nil {
			endTime := time.Now()
			systemTaskCommand.Status = models.SYSTEM_TASK_FAILED
			systemTaskCommand.EndedAt = &endTime
			systemTaskCommand.ResultDescription = serveErr.Error()
			return nil, systemTaskCommand, serveErr
		}
		defer server.Close()

		// Register the interception handler BEFORE Navigate so the main navigation
		// request is intercepted too. The callback runs on chromedp's event loop and
		// must not block, so it hands each paused request to a goroutine.
		chromedp.ListenTarget(browserCtx, func(ev interface{}) {
			if paused, ok := ev.(*fetch.EventRequestPaused); ok {
				go resolveInterceptedRequest(browserCtx, paused, addr)
			}
		})

		actions = []chromedp.Action{
			fetch.Enable().WithPatterns([]*fetch.RequestPattern{{URLPattern: "*"}}),
			chromedp.Navigate("http://" + addr + "/"),
			printAction,
		}
	} else {
		// Default: block external network resources AND file:// loads. Receipt
		// emails contain attacker-controllable URLs and we run chromium with
		// --no-sandbox, so disallowing these removes an SSRF / tracking-pixel
		// surface. The HTML is injected into an about:blank page instead of a
		// file:// URL, so it has no local-file origin — <iframe src="file:///...">/
		// <img> cannot exfiltrate local files into the PDF, and the multi-MB
		// data-URL truncation cap is avoided. Inline data: URIs remain allowed and
		// load synchronously, so no explicit load wait is needed before printing.
		actions = []chromedp.Action{
			network.Enable(),
			network.SetBlockedURLs(blockedExternalUrlPatterns),
			chromedp.Navigate("about:blank"),
			chromedp.ActionFunc(func(ctx context.Context) error {
				frameTree, err := page.GetFrameTree().Do(ctx)
				if err != nil {
					return err
				}
				return page.SetDocumentContent(frameTree.Frame.ID, html).Do(ctx)
			}),
			printAction,
		}
	}
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

// startLoopbackHtmlServer serves html at "/" on a random loopback port so the
// renderer can navigate to it with a real http origin (required for remote
// sub-resources to load; an about:blank document set via Page.setDocumentContent
// never fetches them). Used only when external resources are enabled. The
// listener is bound to 127.0.0.1 and the caller closes the returned server as
// soon as the render finishes; only "/" is served (any other path 404s).
func startLoopbackHtmlServer(html string) (*http.Server, string, error) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return nil, "", err
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = io.WriteString(w, html)
	})

	server := &http.Server{Handler: mux}
	go func() { _ = server.Serve(listener) }()

	return server, listener.Addr().String(), nil
}

// resolveInterceptedRequest continues or blocks one Fetch-intercepted request
// according to the SSRF policy (isRequestAllowed). It runs in its own goroutine
// so the ListenTarget callback never blocks the CDP event loop. Errors from
// Continue/Fail are ignored: the browser context may already be cancelled (the
// render finished or timed out), which is not actionable here.
func resolveInterceptedRequest(ctx context.Context, paused *fetch.EventRequestPaused, loopbackAddr string) {
	executor := cdp.WithExecutor(ctx, chromedp.FromContext(ctx).Target)

	var action chromedp.Action
	if isRequestAllowed(paused.Request.URL, loopbackAddr) {
		action = fetch.ContinueRequest(paused.RequestID)
	} else {
		logging.LogStd(logging.LOG_LEVEL_INFO, "HTML to PDF blocked non-public sub-resource: ", paused.Request.URL)
		action = fetch.FailRequest(paused.RequestID, network.ErrorReasonBlockedByClient)
	}
	_ = action.Do(executor)
}

// isRequestAllowed decides whether a Fetch-intercepted sub-resource request may
// proceed in external-resources mode. Inline data: URIs and our own loopback page
// origin are always allowed; other http(s) requests are allowed only to
// non-internal hosts (public logos / imagery); every other scheme (file, ws, …)
// is denied. A parse failure denies (fail closed).
func isRequestAllowed(rawURL string, loopbackAddr string) bool {
	u, err := url.Parse(rawURL)
	if err != nil {
		return false
	}
	switch u.Scheme {
	case "data":
		return true
	case "http", "https":
		if u.Host == loopbackAddr {
			// Our own served page and its same-origin sub-resources.
			return true
		}
		return !hostIsInternal(u.Hostname())
	default:
		return false
	}
}

// hostIsInternal reports whether a request host points at a loopback, private,
// link-local (incl. 169.254.169.254 cloud metadata) or unspecified address.
// Hostnames are resolved and treated as internal if ANY resolved address is, or
// if resolution fails (fail closed). NOTE: the address is checked at decision
// time; a host that rebinds to an internal IP between this check and chromium's
// own resolution could still slip through (documented DNS-rebinding residual).
func hostIsInternal(host string) bool {
	if host == "" || strings.EqualFold(host, "localhost") {
		return true
	}
	if ip := net.ParseIP(host); ip != nil {
		return ipIsInternal(ip)
	}
	ips, err := net.LookupIP(host)
	if err != nil || len(ips) == 0 {
		return true
	}
	for _, ip := range ips {
		if ipIsInternal(ip) {
			return true
		}
	}
	return false
}

// ipIsInternal reports whether an IP is one the render must never reach: loopback
// (127/8, ::1), RFC1918 / IPv6 ULA (via IsPrivate), link-local unicast/multicast
// (169.254/16 incl. cloud metadata, fe80::/10) or unspecified (0.0.0.0, ::).
func ipIsInternal(ip net.IP) bool {
	return ip.IsLoopback() ||
		ip.IsPrivate() ||
		ip.IsLinkLocalUnicast() ||
		ip.IsLinkLocalMulticast() ||
		ip.IsUnspecified()
}
