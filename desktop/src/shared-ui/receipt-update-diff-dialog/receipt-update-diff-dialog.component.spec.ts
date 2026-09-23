import { CommonModule } from "@angular/common";
import { CUSTOM_ELEMENTS_SCHEMA, provideZonelessChangeDetection } from "@angular/core";
import { ComponentFixture, TestBed } from "@angular/core/testing";
import { MAT_DIALOG_DATA, MatDialogRef } from "@angular/material/dialog";
import { CollapsedDiffRow } from "../../utils/line-diff";
import { ReceiptUpdateSnapshots } from "../../utils/receipt-update-description";
import { ReceiptUpdateDiffDialogComponent, ReceiptUpdateDiffDialogData } from "./receipt-update-diff-dialog.component";

describe("ReceiptUpdateDiffDialogComponent", () => {
  let fixture: ComponentFixture<ReceiptUpdateDiffDialogComponent>;
  let component: ReceiptUpdateDiffDialogComponent;
  let dialogRef: { close: jest.Mock };

  // Enough unchanged keys around the edit that "Changes only" has lines to hide.
  const before = {
    id: 7,
    createdAt: "2026-09-01T10:00:00Z",
    updatedAt: "2026-09-22T10:00:00Z",
    name: "Costco",
    amount: "42.1",
    date: "2026-09-20T00:00:00Z",
    resolvedDate: null,
    paidByUserId: 1,
    status: "OPEN",
    groupId: 2,
    categories: [],
    tags: [{ id: 1, name: "Groceries" }],
    imageFiles: [],
    receiptItems: [],
    comments: [],
    customFields: [],
  };
  const after = {
    ...before,
    updatedAt: "2026-09-23T10:00:00Z",
    name: "Costco Wholesale",
    tags: [{ id: 1, name: "Groceries" }, { id: 2, name: "Bulk" }],
  };

  const create = async (snapshots: ReceiptUpdateSnapshots) => {
    dialogRef = { close: jest.fn() };
    await TestBed.configureTestingModule({
      declarations: [ReceiptUpdateDiffDialogComponent],
      imports: [CommonModule],
      schemas: [CUSTOM_ELEMENTS_SCHEMA],
      providers: [
        provideZonelessChangeDetection(),
        { provide: MatDialogRef, useValue: dialogRef },
        { provide: MAT_DIALOG_DATA, useValue: { snapshots } satisfies ReceiptUpdateDiffDialogData },
      ],
    }).compileComponents();

    fixture = TestBed.createComponent(ReceiptUpdateDiffDialogComponent);
    component = fixture.componentInstance;
    await fixture.whenStable();
  };

  const renderedRows = (): HTMLElement[] =>
    Array.from(fixture.nativeElement.querySelectorAll("[data-testid='receipt-diff-row']"));

  const rowsOfKind = (kind: string) => renderedRows().filter((row) => row.dataset["kind"] === kind);

  beforeEach(async () => {
    await create({ before, after });
  });

  it("shows the old value on the left and the new value on the right of a changed line", () => {
    const nameRow = rowsOfKind("changed").find((row) => row.textContent?.includes("Costco"))!;
    const [left, right] = Array.from(nameRow.querySelectorAll<HTMLElement>(".receipt-diff__code"));

    expect(left.textContent).toBe('  "name": "Costco",');
    expect(right.textContent).toBe('  "name": "Costco Wholesale",');
    expect(right.querySelector("mark")?.textContent).toBe(" Wholesale");
    expect(left.querySelector("mark")).toBeNull();
  });

  it("shows an added line on the right only", () => {
    const added = rowsOfKind("added");

    expect(added.length).toBeGreaterThan(0);
    for (const row of added) {
      const [left] = Array.from(row.querySelectorAll<HTMLElement>(".receipt-diff__code"));
      expect(left.classList).toContain("receipt-diff__empty");
      expect(left.textContent).toBe("");
    }
    expect(added.map((row) => row.textContent).join("")).toContain('"name": "Bulk"');
  });

  it("shows every line by default", () => {
    expect(component.view()).toBe("all");
    expect(fixture.nativeElement.querySelector("[data-testid='receipt-diff-collapsed']")).toBeNull();
    expect(rowsOfKind("equal").length).toBeGreaterThan(10);
  });

  it("collapses unchanged lines away from the changes in the changes-only view", async () => {
    component.setView("changes");
    await fixture.whenStable();

    const collapsed = component.displayRows().filter((row): row is CollapsedDiffRow => row.kind === "collapsed");
    expect(collapsed.length).toBeGreaterThan(0);
    expect(fixture.nativeElement.querySelector("[data-testid='receipt-diff-collapsed']")?.textContent)
      .toContain("unchanged lines");
    // Every change is still there.
    const changes = component.displayRows().filter((row) => row.kind !== "equal" && row.kind !== "collapsed");
    expect(changes.length).toBe(component.tabs[1].count);
  });

  it("counts added and removed lines", () => {
    // updatedAt and name each swap one line; the new tag adds four more lines
    // ("},", "{", its id, its name), its closing brace matching the old one.
    expect(component.removedLineCount).toBe(2);
    expect(component.addedLineCount).toBe(6);
    expect(component.tabs[1].count).toBe(renderedRows().filter((row) => row.dataset["kind"] !== "equal").length);
  });

  it("titles the dialog with the receipt's new name", () => {
    expect(component.headerText).toBe("Receipt update: Costco Wholesale");
  });

  it("copies the pair as readable JSON rather than the double-encoded description", () => {
    expect(JSON.parse(component.copyText)).toEqual({ before, after });
    expect(component.copyText).not.toContain('\\"');
  });

  it("closes", () => {
    component.close();

    expect(dialogRef.close).toHaveBeenCalled();
  });
});
