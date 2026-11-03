import { provideHttpClient } from "@angular/common/http";
import { provideHttpClientTesting } from "@angular/common/http/testing";
import { provideZonelessChangeDetection } from "@angular/core";
import { ComponentFixture, TestBed } from "@angular/core/testing";
import { MatDialog } from "@angular/material/dialog";
import { NoopAnimationsModule } from "@angular/platform-browser/animations";
import { provideRouter } from "@angular/router";
import { NgxsModule } from "@ngxs/store";
import { NEVER, of, throwError } from "rxjs";
import { ReceiptHistory, ReceiptHistoryEntry, ReceiptService } from "../../open-api";
import { SystemSettingsState } from "../../store/system-settings.state";
import { UserState } from "../../store/user.state";
import { ReceiptHistoryDrawerComponent } from "./receipt-history-drawer.component";

function buildEntry(systemTaskId: number, recordedAt: string, action = "UPDATED"): ReceiptHistoryEntry {
  return { systemTaskId, action, recordedAt, actor: { type: "SYSTEM" }, hasSnapshots: false, changes: [] };
}

const HISTORY: ReceiptHistory = {
  entries: [
    buildEntry(3, "2026-10-13T18:00:00Z", "RESOLVED"),
    buildEntry(2, "2026-10-13T09:00:00Z"),
    buildEntry(1, "2026-10-12T09:14:08Z", "CREATED"),
  ],
  changeCount: 5,
  saveCount: 3,
  contributorCount: 1,
};

describe("ReceiptHistoryDrawerComponent", () => {
  let fixture: ComponentFixture<ReceiptHistoryDrawerComponent>;
  let component: ReceiptHistoryDrawerComponent;
  let getReceiptHistory: jest.SpyInstance;
  let closed: jest.Mock;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [
        ReceiptHistoryDrawerComponent,
        NoopAnimationsModule,
        NgxsModule.forRoot([UserState, SystemSettingsState]),
      ],
      providers: [
        provideZonelessChangeDetection(),
        provideHttpClient(),
        provideHttpClientTesting(),
        // app-button binds routerLink.
        provideRouter([]),
      ],
    }).compileComponents();

    getReceiptHistory = jest
      .spyOn(TestBed.inject(ReceiptService), "getReceiptHistory")
      .mockReturnValue(of(HISTORY) as any);
  });

  afterEach(() => {
    document.querySelectorAll("app-header, [data-testid='opener']").forEach((node) => node.remove());
  });

  // `whenStable` waits for the resource's request, so a load that never answers
  // is rendered with a plain change-detection pass instead.
  async function create({ waitForLoad = true } = {}): Promise<void> {
    fixture = TestBed.createComponent(ReceiptHistoryDrawerComponent);
    component = fixture.componentInstance;
    closed = jest.fn();
    component.closed.subscribe(closed);
    fixture.componentRef.setInput("receiptId", 1042);
    if (waitForLoad) {
      await fixture.whenStable();
    } else {
      fixture.detectChanges();
    }
  }

  function host(): HTMLElement {
    return fixture.nativeElement;
  }

  function query(testId: string): HTMLElement | null {
    return host().querySelector(`[data-testid="${testId}"]`);
  }

  function closeButton(): HTMLButtonElement {
    return query("receipt-history-close")!.querySelector("button")!;
  }

  it("loads the receipt's history and summarizes it", async () => {
    await create();

    expect(getReceiptHistory).toHaveBeenCalledWith(1042);
    expect(query("receipt-history-summary")!.textContent).toBe("5 changes · 3 saves · 1 contributor");
  });

  it("groups the entries by day, newest first", async () => {
    await create();

    const days = Array.from(host().querySelectorAll("[data-testid='receipt-history-day']"));
    expect(days.map((day) => day.textContent)).toEqual(["Tuesday, Oct 13", "Monday, Oct 12"]);
    const actions = Array.from(host().querySelectorAll("[data-testid='receipt-history-entry']")).map((entry) =>
      entry.getAttribute("data-action")
    );
    expect(actions).toEqual(["RESOLVED", "UPDATED", "CREATED"]);
  });

  it("is a labelled complementary region", async () => {
    await create();

    expect(host().getAttribute("role")).toBe("complementary");
    const title = host().querySelector(`#${host().getAttribute("aria-labelledby")}`);
    expect(title?.textContent).toBe("Receipt History");
    expect(closeButton().getAttribute("aria-label")).toBe("Close receipt history");
  });

  it("shows a loading state until the history arrives", async () => {
    getReceiptHistory.mockReturnValue(NEVER);
    await create({ waitForLoad: false });

    expect(query("receipt-history-loading")).not.toBeNull();
    expect(query("receipt-history-summary")).toBeNull();
  });

  it("shows an empty state for a receipt with no history", async () => {
    getReceiptHistory.mockReturnValue(
      of({ entries: [], changeCount: 0, saveCount: 0, contributorCount: 0 }) as any
    );
    await create();

    expect(query("receipt-history-empty")).not.toBeNull();
  });

  it("reports a failed load and retries", async () => {
    getReceiptHistory.mockReturnValueOnce(throwError(() => new Error("403")));
    await create();

    expect(query("receipt-history-error")!.textContent).toContain("Couldn't load this receipt's history.");
    expect(query("receipt-history-entry")).toBeNull();

    query("receipt-history-retry")!.querySelector("button")!.click();
    await fixture.whenStable();

    expect(getReceiptHistory).toHaveBeenCalledTimes(2);
    expect(query("receipt-history-error")).toBeNull();
    expect(host().querySelectorAll("[data-testid='receipt-history-entry']")).toHaveLength(3);
  });

  describe("closing", () => {
    it("closes from the close button", async () => {
      await create();

      closeButton().click();

      expect(closed).toHaveBeenCalledTimes(1);
    });

    it("closes on Escape", async () => {
      await create();

      document.dispatchEvent(new KeyboardEvent("keydown", { key: "Escape" }));

      expect(closed).toHaveBeenCalledTimes(1);
    });

    // A dialog, menu or autocomplete above the drawer consumes its own Escape.
    it("leaves an Escape another overlay handled alone", async () => {
      await create();

      const event = new KeyboardEvent("keydown", { key: "Escape", cancelable: true });
      event.preventDefault();
      document.dispatchEvent(event);

      expect(closed).not.toHaveBeenCalled();
    });

    it("stays open on Escape while a dialog is open", async () => {
      await create();
      jest.spyOn(TestBed.inject(MatDialog), "openDialogs", "get").mockReturnValue([{} as any]);

      document.dispatchEvent(new KeyboardEvent("keydown", { key: "Escape" }));

      expect(closed).not.toHaveBeenCalled();
    });
  });

  describe("focus", () => {
    function addOpener(): HTMLButtonElement {
      const opener = document.createElement("button");
      opener.setAttribute("data-testid", "opener");
      document.body.appendChild(opener);
      opener.focus();
      return opener;
    }

    it("moves focus to the close button on open", async () => {
      await create();

      expect(document.activeElement).toBe(closeButton());
    });

    it("hands focus back to the control that opened it", async () => {
      const opener = addOpener();
      await create();
      expect(document.activeElement).not.toBe(opener);

      fixture.destroy();

      expect(document.activeElement).toBe(opener);
    });

    it("leaves focus alone when the user has moved on", async () => {
      const opener = addOpener();
      await create();
      const elsewhere = document.createElement("input");
      elsewhere.setAttribute("data-testid", "opener");
      document.body.appendChild(elsewhere);
      elsewhere.focus();

      fixture.destroy();

      expect(document.activeElement).toBe(elsewhere);
      expect(opener).not.toBe(document.activeElement);
    });
  });

  describe("position", () => {
    it("sits below the app header and beside the sidebar once measured", async () => {
      const header = document.createElement("app-header");
      header.getBoundingClientRect = () => ({ top: 0, bottom: 89, left: 82, height: 89 }) as DOMRect;
      document.body.appendChild(header);

      await create();

      expect(host().style.top).toBe("89px");
      expect(host().style.maxWidth).toBe("calc(100vw - 82px)");
    });

    it("keeps the stylesheet's fallback without a header to measure", async () => {
      await create();

      expect(host().style.top).toBe("");
      expect(host().style.maxWidth).toBe("");
    });
  });
});
