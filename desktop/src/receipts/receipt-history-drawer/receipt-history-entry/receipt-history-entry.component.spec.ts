import { provideHttpClient } from "@angular/common/http";
import { provideHttpClientTesting } from "@angular/common/http/testing";
import { provideZonelessChangeDetection } from "@angular/core";
import { ComponentFixture, TestBed } from "@angular/core/testing";
import { NgxsModule, Store } from "@ngxs/store";
import { of, throwError } from "rxjs";
import { ReceiptHistoryChange, ReceiptHistoryEntry, ReceiptService, User } from "../../../open-api";
import { SystemSettingsState } from "../../../store/system-settings.state";
import { UserState } from "../../../store/user.state";
import { SetUsers } from "../../../store/user.state.actions";
import { ReceiptHistoryEntryComponent } from "./receipt-history-entry.component";

const maya = { id: 2, displayName: "Maya Chen", username: "maya", defaultAvatarColor: "#a78bfa" } as User;

function buildEntry(overrides: Partial<ReceiptHistoryEntry> = {}): ReceiptHistoryEntry {
  return {
    systemTaskId: 9188,
    action: "UPDATED",
    recordedAt: "2026-10-13T10:02:51Z",
    actor: { type: "USER", userId: 2 },
    hasSnapshots: true,
    beforeState: "COMPLETE",
    changes: [],
    ...overrides,
  };
}

function textChange(field: string, from: string, to: string): ReceiptHistoryChange {
  return { field, kind: "CHANGED", from: { type: "TEXT", value: from }, to: { type: "TEXT", value: to } };
}

describe("ReceiptHistoryEntryComponent", () => {
  let fixture: ComponentFixture<ReceiptHistoryEntryComponent>;
  let receiptService: ReceiptService;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [ReceiptHistoryEntryComponent, NgxsModule.forRoot([UserState, SystemSettingsState])],
      providers: [provideZonelessChangeDetection(), provideHttpClient(), provideHttpClientTesting()],
    }).compileComponents();

    TestBed.inject(Store).dispatch(new SetUsers([maya]));
    receiptService = TestBed.inject(ReceiptService);
    fixture = TestBed.createComponent(ReceiptHistoryEntryComponent);
    fixture.componentRef.setInput("receiptId", 1042);
  });

  async function render(entry: ReceiptHistoryEntry): Promise<void> {
    fixture.componentRef.setInput("entry", entry);
    await fixture.whenStable();
  }

  function query<T extends HTMLElement = HTMLElement>(testId: string): T | null {
    return fixture.nativeElement.querySelector(`[data-testid="${testId}"]`);
  }

  function queryAll(testId: string): HTMLElement[] {
    return Array.from(fixture.nativeElement.querySelectorAll(`[data-testid="${testId}"]`));
  }

  async function click(testId: string): Promise<void> {
    query(testId)!.click();
    await fixture.whenStable();
  }

  it("names the actor and what they did, with the action on the element", async () => {
    await render(buildEntry());

    expect(query("receipt-history-entry")!.getAttribute("data-action")).toBe("UPDATED");
    expect(query("receipt-history-actor")!.textContent).toBe("Maya Chen");
    expect(query("receipt-history-verb")!.textContent).toBe("edited this receipt");
    // Angular's CLDR data puts a narrow no-break space before AM/PM.
    expect(fixture.nativeElement.querySelector("time").textContent.replace(/\s/g, " ")).toBe("10:02 AM");
  });

  it("falls back on codes this client does not know", async () => {
    await render(
      buildEntry({
        action: "MERGED",
        actor: { type: "ROBOT" },
        changes: [{ field: "RECEIPT_NOTES", kind: "MOVED", from: { type: "MYSTERY", value: "a" }, to: { type: "MYSTERY", value: "b" } }],
      })
    );

    expect(query("receipt-history-verb")!.textContent).toBe("updated this receipt");
    expect(query("receipt-history-actor")!.textContent).toBe("System");
    const row = query("receipt-history-change")!;
    expect(row.textContent).toContain("Receipt Notes");
    expect(query("receipt-history-from")!.textContent).toBe("a");
    expect(query("receipt-history-to")!.textContent).toBe("b");
  });

  describe("collapsing", () => {
    const changes = ["A", "B", "C", "D", "E", "F"].map((name) => textChange("NAME", `${name}1`, `${name}2`));

    it("shows four changes and offers the rest", async () => {
      await render(buildEntry({ changes }));

      expect(queryAll("receipt-history-change")).toHaveLength(4);
      const more = query("receipt-history-show-more")!;
      expect(more.textContent).toBe("Show 2 more");
      expect(more.getAttribute("aria-expanded")).toBe("false");

      await click("receipt-history-show-more");
      expect(queryAll("receipt-history-change")).toHaveLength(6);
      expect(query("receipt-history-show-more")!.textContent).toBe("Show less");
      expect(query("receipt-history-show-more")!.getAttribute("aria-expanded")).toBe("true");

      await click("receipt-history-show-more");
      expect(queryAll("receipt-history-change")).toHaveLength(4);
    });

    it("offers nothing more at exactly four", async () => {
      await render(buildEntry({ changes: changes.slice(0, 4) }));

      expect(queryAll("receipt-history-change")).toHaveLength(4);
      expect(query("receipt-history-show-more")).toBeNull();
    });
  });

  it("boxes a failed save with its error", async () => {
    await render(
      buildEntry({
        action: "FAILED",
        hasSnapshots: false,
        errorMessage: "Your role requires at least one image on this group's receipts.",
      })
    );

    expect(query("receipt-history-verb")!.textContent).toBe("tried to save changes");
    expect(query("receipt-history-failed")!.textContent).toContain(
      "Save failed: Your role requires at least one image on this group's receipts."
    );
    expect(query("receipt-history-json-toggle")).toBeNull();
  });

  it("explains a rebuilt entry", async () => {
    await render(
      buildEntry({
        beforeState: "REBUILT",
        beforeSourceType: "RECEIPT_UPDATED",
        beforeSourceRecordedAt: "2026-10-16T08:03:17Z",
      })
    );

    const notice = query("receipt-history-notice")!;
    expect(notice.getAttribute("data-before-state")).toBe("REBUILT");
    expect(notice.textContent).toContain("compared against the previous update (Oct 16, 2026)");
  });

  it("warns about an incomplete entry", async () => {
    await render(buildEntry({ beforeState: "INCOMPLETE" }));

    const notice = query("receipt-history-notice")!;
    expect(notice.getAttribute("data-before-state")).toBe("INCOMPLETE");
    expect(notice.textContent).toContain("Some changes may be inaccurate.");
  });

  it("shows no notice or failure for a complete save", async () => {
    await render(buildEntry());

    expect(query("receipt-history-notice")).toBeNull();
    expect(query("receipt-history-failed")).toBeNull();
  });

  describe("View JSON", () => {
    it("fetches the snapshots once, and toggles them", async () => {
      const spy = jest
        .spyOn(receiptService, "getReceiptHistorySnapshots")
        .mockReturnValue(of({ before: '{"name":"Old"}', after: '{"name":"New"}' }) as any);
      await render(buildEntry());

      expect(query("receipt-history-json-toggle")!.textContent).toBe("View JSON");
      expect(spy).not.toHaveBeenCalled();

      await click("receipt-history-json-toggle");
      expect(spy).toHaveBeenCalledWith(1042, 9188);
      expect(JSON.parse(query("receipt-history-json")!.textContent!)).toEqual({
        before: { name: "Old" },
        after: { name: "New" },
      });
      expect(query("receipt-history-json")!.textContent).toContain('\n  "before": {\n    "name": "Old"');
      expect(query("receipt-history-json-toggle")!.textContent).toBe("Hide JSON");
      expect(query("receipt-history-json-toggle")!.getAttribute("aria-expanded")).toBe("true");

      await click("receipt-history-json-toggle");
      expect(query("receipt-history-json")).toBeNull();
      expect(query("receipt-history-json-toggle")!.textContent).toBe("View JSON");

      await click("receipt-history-json-toggle");
      expect(query("receipt-history-json")).not.toBeNull();
      expect(spy).toHaveBeenCalledTimes(1);
    });

    it("shows the created receipt alone for a creation entry", async () => {
      jest.spyOn(receiptService, "getReceiptHistorySnapshots").mockReturnValue(of({ after: '{"name":"New"}' }) as any);
      await render(buildEntry({ action: "CREATED" }));

      await click("receipt-history-json-toggle");

      expect(JSON.parse(query("receipt-history-json")!.textContent!)).toEqual({ name: "New" });
    });

    it("offers no JSON for an entry without snapshots", async () => {
      await render(buildEntry({ hasSnapshots: false }));

      expect(query("receipt-history-json-toggle")).toBeNull();
    });

    it("reports a failed load and retries on the next open", async () => {
      const spy = jest
        .spyOn(receiptService, "getReceiptHistorySnapshots")
        .mockReturnValueOnce(throwError(() => new Error("boom")))
        .mockReturnValueOnce(of({ after: "{}" }) as any);
      await render(buildEntry());

      await click("receipt-history-json-toggle");
      expect(fixture.nativeElement.textContent).toContain("Couldn't load the JSON for this entry.");
      expect(query("receipt-history-json")).toBeNull();

      await click("receipt-history-json-toggle");
      await click("receipt-history-json-toggle");
      expect(spy).toHaveBeenCalledTimes(2);
      expect(query("receipt-history-json")!.textContent).toBe("{}");
    });
  });

  describe("value formatting", () => {
    function change(type: string, value?: string, userId?: number): ReceiptHistoryChange {
      return { field: "CUSTOM_FIELD", label: type, kind: "ADDED", to: { type, value, userId } };
    }

    async function shown(...changes: ReceiptHistoryChange[]): Promise<string[]> {
      await render(buildEntry({ changes }));
      await click("receipt-history-show-more");
      return queryAll("receipt-history-to").map((chip) => chip.textContent ?? "");
    }

    it("formats each value type with the app's own formatting", async () => {
      const values = await shown(
        change("TEXT", "Costco"),
        change("CURRENCY", "1187.42"),
        change("DATE", "2026-10-11"),
        change("INSTANT", "2026-10-20T16:52:31Z"),
        change("USER", undefined, 2),
        change("USER", undefined, 99),
        change("STATUS", "NEEDS_ATTENTION"),
        change("BOOLEAN", "true"),
        change("SOMETHING_NEW", "raw")
      );

      expect(values.map((value) => value.replace(/\s/g, " "))).toEqual([
        "Costco",
        "$1,187.42",
        "Oct 11, 2026",
        "Oct 20, 2026, 4:52:31 PM",
        "Maya Chen",
        "Unknown user",
        "Needs Attention",
        "Yes",
        "raw",
      ]);
    });
  });

  it("strikes out the old value and shows the new one beside an arrow", async () => {
    await render(
      buildEntry({
        changes: [
          textChange("NAME", "COSTCO WHSE #482", "Costco Wholesale"),
          { field: "TAGS", kind: "ADDED", to: { type: "TEXT", value: "Bulk" } },
          { field: "IMAGE", kind: "REMOVED", from: { type: "TEXT", value: "costco_1011.jpg" } },
        ],
      })
    );

    const [name, tags, image] = queryAll("receipt-history-change");
    expect(name.getAttribute("data-field")).toBe("NAME");
    expect(name.getAttribute("data-kind")).toBe("CHANGED");
    expect(name.querySelector("del")!.textContent).toBe("COSTCO WHSE #482");
    expect(name.querySelector("ins")!.textContent).toBe("Costco Wholesale");
    expect(name.querySelector("mat-icon")!.textContent).toBe("arrow_forward");

    expect(tags.querySelector("del")).toBeNull();
    expect(tags.querySelector("mat-icon")).toBeNull();
    expect(tags.querySelector("ins")!.textContent).toBe("Bulk");

    expect(image.querySelector("ins")).toBeNull();
    expect(image.querySelector("del")!.textContent).toBe("costco_1011.jpg");
  });

  it("titles a clipped comment with its full text", async () => {
    const text = "Took the bag fee off, we brought our own. My share is the chicken + half the paper towels.";
    await render(
      buildEntry({
        action: "IMAGES_COMMENTS_UPDATED",
        changes: [{ field: "COMMENT", kind: "ADDED", userId: 2, to: { type: "TEXT", value: text } }],
      })
    );

    const chip = query("receipt-history-to")!;
    expect(chip.textContent!.startsWith("Maya Chen: “Took the bag fee off")).toBe(true);
    expect(chip.textContent!.endsWith("…”")).toBe(true);
    expect(chip.getAttribute("title")).toBe(`Maya Chen: “${text}”`);
  });
});
