import { CurrencyPipe } from "@angular/common";
import { provideZonelessChangeDetection } from "@angular/core";
import { ComponentFixture, TestBed } from "@angular/core/testing";
import { NoopAnimationsModule } from "@angular/platform-browser/animations";
import { NgxsModule } from "@ngxs/store";
import { Group, ReceiptStatus, ReceiptSummary } from "../../open-api";
import { PipesModule } from "../../pipes";
import { SystemSettingsState } from "../../store/system-settings.state";
import { ReceiptTotalsComponent } from "./receipt-totals.component";

function buildSummary(overrides: Partial<ReceiptSummary> = {}): ReceiptSummary {
  return {
    enabled: true,
    configurationGroupId: 1,
    overall: {
      status: "" as ReceiptStatus,
      receiptCount: 3,
      total: "122.24",
      customFieldTotals: [],
    },
    statuses: [],
    ...overrides,
  } as ReceiptSummary;
}

function buildGroup(id: number, name: string): Group {
  return { id, name } as Group;
}

describe("ReceiptTotalsComponent", () => {
  let component: ReceiptTotalsComponent;
  let fixture: ComponentFixture<ReceiptTotalsComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [
        ReceiptTotalsComponent,
        NoopAnimationsModule,
        PipesModule,
        NgxsModule.forRoot([SystemSettingsState]),
      ],
      providers: [provideZonelessChangeDetection(), CurrencyPipe],
    }).compileComponents();

    fixture = TestBed.createComponent(ReceiptTotalsComponent);
    component = fixture.componentInstance;
  });

  function text(): string {
    return fixture.nativeElement.textContent ?? "";
  }

  it("renders nothing without a summary", async () => {
    fixture.detectChanges();
    await fixture.whenStable();

    expect(fixture.nativeElement.querySelector("[data-testid='receipt-totals']")).toBeNull();
  });

  // The off state is a well-formed 200 from the server, not an error, so the component's job is to
  // render nothing rather than to show an empty block.
  it("renders nothing when the summary is disabled", async () => {
    fixture.componentRef.setInput("summary", buildSummary({ enabled: false }));
    fixture.detectChanges();
    await fixture.whenStable();

    expect(fixture.nativeElement.querySelector("[data-testid='receipt-totals']")).toBeNull();
  });

  it("renders the overall row with its count and total", async () => {
    fixture.componentRef.setInput("summary", buildSummary());
    fixture.detectChanges();
    await fixture.whenStable();

    const row = fixture.nativeElement.querySelector("[data-testid='receipt-totals-row-overall']");
    expect(row).not.toBeNull();
    expect(row.textContent).toContain("All Receipts");
    expect(row.textContent).toContain("3 receipts");
    // Through customCurrency, so money reads the same as everywhere else in the app.
    expect(row.textContent).toContain("122.24");
  });

  it("singularizes a one-receipt row", async () => {
    fixture.componentRef.setInput(
      "summary",
      buildSummary({
        overall: {
          status: "" as ReceiptStatus,
          receiptCount: 1,
          total: "10.00",
          customFieldTotals: [],
        },
      })
    );
    fixture.detectChanges();
    await fixture.whenStable();

    expect(text()).toContain("1 receipt)");
    expect(text()).not.toContain("1 receipts");
  });

  // The rule that keeps the block's shape steady as a filter narrows.
  it("renders a configured status that matched nothing as a zero row", async () => {
    fixture.componentRef.setInput(
      "summary",
      buildSummary({
        statuses: [
          {
            status: ReceiptStatus.Open,
            receiptCount: 2,
            total: "22.00",
            customFieldTotals: [],
          },
          {
            status: ReceiptStatus.Resolved,
            receiptCount: 0,
            total: "0",
            customFieldTotals: [],
          },
        ],
      })
    );
    fixture.detectChanges();
    await fixture.whenStable();

    const emptyRow = fixture.nativeElement.querySelector(
      "[data-testid='receipt-totals-row-RESOLVED']"
    );
    expect(emptyRow).not.toBeNull();
    expect(emptyRow.textContent).toContain("0 receipts");
    // Muted, so a legitimate zero does not read as a bug.
    expect(emptyRow.classList).toContain("receipt-totals__row--empty");
  });

  it("renders a column per currency custom field total", async () => {
    fixture.componentRef.setInput(
      "summary",
      buildSummary({
        overall: {
          status: "" as ReceiptStatus,
          receiptCount: 1,
          total: "11.30",
          customFieldTotals: [
            { customFieldId: 1, name: "HST", total: "1.30" },
            { customFieldId: 2, name: "Subtotal", total: "10.00" },
          ],
        },
      })
    );
    fixture.detectChanges();
    await fixture.whenStable();

    expect(text()).toContain("HST");
    expect(text()).toContain("Subtotal");
  });

  describe("the configuration group chips", () => {
    // A picker with one option is not a choice; the caption says the same thing without one.
    it("renders a caption rather than chips for a single group", async () => {
      fixture.componentRef.setInput("summary", buildSummary());
      fixture.componentRef.setInput("configGroups", [buildGroup(1, "Alpha")]);
      fixture.detectChanges();
      await fixture.whenStable();

      expect(fixture.nativeElement.querySelector("mat-chip-listbox")).toBeNull();
      expect(text()).toContain("Using the summary configuration from Alpha");
    });

    it("renders chips and emits the picked group", async () => {
      const picked: number[] = [];
      component.configGroupSelected.subscribe((id) => picked.push(id));

      fixture.componentRef.setInput("summary", buildSummary());
      fixture.componentRef.setInput("configGroups", [buildGroup(1, "Alpha"), buildGroup(2, "Beta")]);
      fixture.componentRef.setInput("selectedConfigGroupId", 1);
      fixture.detectChanges();
      await fixture.whenStable();

      // The clickable element is the inner button[matchipaction]; mat-chip-option's own host is
      // role="presentation", so clicking it does nothing.
      const beta: HTMLElement = fixture.nativeElement.querySelector(
        "[data-testid='receipt-totals-config-group-2'] button[matchipaction]"
      );
      expect(beta).not.toBeNull();

      // The listbox's [value] alone drives selection; there is no per-option [selected].
      const alpha = fixture.nativeElement.querySelector(
        "[data-testid='receipt-totals-config-group-1']"
      );
      expect(alpha.classList).toContain("mat-mdc-chip-selected");

      beta.click();
      fixture.detectChanges();
      await fixture.whenStable();

      expect(picked).toEqual([2]);
    });
  });
});
