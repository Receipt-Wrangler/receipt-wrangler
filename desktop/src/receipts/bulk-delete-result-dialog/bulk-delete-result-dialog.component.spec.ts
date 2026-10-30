import { CUSTOM_ELEMENTS_SCHEMA } from "@angular/core";
import { ComponentFixture, TestBed } from "@angular/core/testing";
import { MAT_DIALOG_DATA, MatDialogRef } from "@angular/material/dialog";
import { BulkDeleteResultDialogComponent, BulkDeleteResultDialogData } from "./bulk-delete-result-dialog.component";

describe("BulkDeleteResultDialogComponent", () => {
  let component: BulkDeleteResultDialogComponent;
  let fixture: ComponentFixture<BulkDeleteResultDialogComponent>;

  const data: BulkDeleteResultDialogData = {
    deletedCount: 1,
    requestedCount: 3,
    failures: [
      { name: "Groceries", message: "You do not have permission to delete this receipt" },
      { name: "Receipt #9", message: "Error deleting receipt" },
    ],
  };

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [BulkDeleteResultDialogComponent],
      providers: [
        { provide: MatDialogRef, useValue: { close: jest.fn() } },
        { provide: MAT_DIALOG_DATA, useValue: data },
      ],
      schemas: [CUSTOM_ELEMENTS_SCHEMA],
    }).compileComponents();

    fixture = TestBed.createComponent(BulkDeleteResultDialogComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it("shows how many of the requested receipts were deleted", () => {
    const summary = fixture.nativeElement.querySelector("[data-testid='bulk-delete-result-summary']");

    expect(summary.textContent.trim()).toBe("Deleted 1 of 3 receipts.");
  });

  it("lists every failed receipt by name with its reason", () => {
    const items: HTMLElement[] = Array.from(
      fixture.nativeElement.querySelectorAll("[data-testid='bulk-delete-result-failures'] li")
    );

    expect(items.map((item) => item.textContent?.replace(/\s+/g, " ").trim())).toEqual([
      "Groceries: You do not have permission to delete this receipt",
      "Receipt #9: Error deleting receipt",
    ]);
  });

  it("closes", () => {
    component.close();

    expect(component.matDialogRef.close).toHaveBeenCalled();
  });
});
