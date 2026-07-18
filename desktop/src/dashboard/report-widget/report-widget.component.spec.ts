import { provideHttpClient, withInterceptorsFromDi } from "@angular/common/http";
import { provideHttpClientTesting } from "@angular/common/http/testing";
import { CUSTOM_ELEMENTS_SCHEMA, provideZonelessChangeDetection } from "@angular/core";
import { ComponentFixture, TestBed } from "@angular/core/testing";
import { By } from "@angular/platform-browser";
import { provideRouter } from "@angular/router";
import { Subject, of, throwError } from "rxjs";
import { ReportPreviewResponse, Widget, WidgetType } from "../../open-api";
import { ReportRunnerService } from "../../reports/services/report-runner.service";
import { ReportWidgetComponent } from "./report-widget.component";

describe("ReportWidgetComponent", () => {
  let fixture: ComponentFixture<ReportWidgetComponent>;
  let component: ReportWidgetComponent;
  let runner: { renderTemplate: jest.Mock; downloadTemplateById: jest.Mock };

  const widget = (reportTemplateId?: number): Widget => ({
    id: 1,
    name: "My Report",
    widgetType: WidgetType.Report,
    configuration: reportTemplateId === undefined ? {} : { reportTemplateId },
  });

  async function setup(w: Widget): Promise<void> {
    fixture = TestBed.createComponent(ReportWidgetComponent);
    component = fixture.componentInstance;
    fixture.componentRef.setInput("widget", w);
    TestBed.flushEffects();
    await fixture.whenStable();
  }

  beforeEach(async () => {
    runner = { renderTemplate: jest.fn(), downloadTemplateById: jest.fn() };

    await TestBed.configureTestingModule({
      imports: [ReportWidgetComponent],
      schemas: [CUSTOM_ELEMENTS_SCHEMA],
      providers: [
        provideZonelessChangeDetection(),
        { provide: ReportRunnerService, useValue: runner },
        provideRouter([]),
        provideHttpClient(withInterceptorsFromDi()),
        provideHttpClientTesting(),
      ],
    }).compileComponents();
  });

  it("renders the returned report HTML in a sandboxed iframe", async () => {
    const response: ReportPreviewResponse = { html: "<p>report</p>", receiptCount: 3, allowedActions: ["read"] };
    runner.renderTemplate.mockReturnValue(of(response));

    await setup(widget(5));

    expect(runner.renderTemplate).toHaveBeenCalledWith(5);
    expect(component.html()).toBe("<p>report</p>");
    expect(component.isLoading()).toBe(false);
    const iframe = fixture.debugElement.query(By.css("iframe"));
    expect(iframe).toBeTruthy();
    expect(iframe.attributes["sandbox"]).toBe("allow-same-origin");
  });

  it("shows the loading spinner while the render is in flight", async () => {
    runner.renderTemplate.mockReturnValue(new Subject<ReportPreviewResponse>()); // never emits

    await setup(widget(5));

    expect(component.isLoading()).toBe(true);
    expect(fixture.debugElement.query(By.css('[data-testid="report-widget-loading"]'))).toBeTruthy();
    expect(fixture.debugElement.query(By.css("iframe"))).toBeNull();
  });

  it("renders whatever HTML comes back for the restricted case and hides the download button", async () => {
    const restricted: ReportPreviewResponse = { html: "<html>restricted</html>", receiptCount: 0, allowedActions: [] };
    runner.renderTemplate.mockReturnValue(of(restricted));

    await setup(widget(5));

    // No special "restricted" branch — the widget just renders the HTML it was given.
    expect(component.html()).toBe("<html>restricted</html>");
    expect(fixture.debugElement.query(By.css("iframe"))).toBeTruthy();
    expect(component.canDownload()).toBe(false);
    expect(fixture.debugElement.query(By.css('[data-testid="report-widget-download"]'))).toBeNull();
  });

  it("shows an error state and does not call the API when no template is pinned", async () => {
    await setup(widget(undefined));

    expect(runner.renderTemplate).not.toHaveBeenCalled();
    expect(component.hasError()).toBe(true);
    expect(fixture.debugElement.query(By.css('[data-testid="report-widget-error"]'))).toBeTruthy();
  });

  it("shows an error state when the render request fails", async () => {
    runner.renderTemplate.mockReturnValue(throwError(() => new Error("boom")));

    await setup(widget(5));

    expect(component.hasError()).toBe(true);
    expect(component.isLoading()).toBe(false);
    expect(fixture.debugElement.query(By.css('[data-testid="report-widget-error"]'))).toBeTruthy();
    expect(fixture.debugElement.query(By.css("iframe"))).toBeNull();
  });

  it("ignores a slow render for a swapped-out template (stale result cannot win)", async () => {
    const first = new Subject<ReportPreviewResponse>(); // template 5 — never resolves until later
    const second: ReportPreviewResponse = { html: "<p>seven</p>", receiptCount: 1, allowedActions: ["read", "generate"] };
    runner.renderTemplate.mockImplementation((id: number) => (id === 5 ? first : of(second)));

    await setup(widget(5));
    expect(component.isLoading()).toBe(true);

    // Swap the pinned template while template 5's request is still pending.
    fixture.componentRef.setInput("widget", widget(7));
    TestBed.flushEffects();
    await fixture.whenStable();

    // The new template resolved synchronously.
    expect(component.html()).toBe("<p>seven</p>");
    expect(component.allowedActions()).toEqual(["read", "generate"]);

    // The stale template-5 request now emits — its subscription was cancelled on the
    // swap, so it must not overwrite the current report or its permissions.
    first.next({ html: "<p>five</p>", receiptCount: 9, allowedActions: [] });
    TestBed.flushEffects();
    await fixture.whenStable();

    expect(component.html()).toBe("<p>seven</p>");
    expect(component.allowedActions()).toEqual(["read", "generate"]);
  });

  it("shows the download button only when allowedActions include generate, and downloads on click", async () => {
    const response: ReportPreviewResponse = { html: "<p>report</p>", receiptCount: 2, allowedActions: ["read", "generate"] };
    runner.renderTemplate.mockReturnValue(of(response));
    runner.downloadTemplateById.mockReturnValue(of(new Blob()));

    await setup(widget(7));

    expect(component.canDownload()).toBe(true);
    expect(fixture.debugElement.query(By.css('[data-testid="report-widget-download"]'))).toBeTruthy();

    component.download();
    expect(runner.downloadTemplateById).toHaveBeenCalledWith(7);
  });

  it("resets the downloading flag when the download fails", async () => {
    const response: ReportPreviewResponse = { html: "<p>report</p>", receiptCount: 2, allowedActions: ["read", "generate"] };
    runner.renderTemplate.mockReturnValue(of(response));
    runner.downloadTemplateById.mockReturnValue(throwError(() => new Error("download failed")));

    await setup(widget(7));

    component.download();
    await fixture.whenStable();

    expect(component.isDownloading()).toBe(false);
  });
});
