import { Component, NO_ERRORS_SCHEMA, provideZonelessChangeDetection } from "@angular/core";
import { ComponentFixture, TestBed } from "@angular/core/testing";
import { FormArray, FormBuilder, FormControl, FormGroup } from "@angular/forms";
import { MatDialog } from "@angular/material/dialog";
import { UntilDestroy } from "@ngneat/until-destroy";
import { ReportReceiptsDialogComponent } from "../dialogs/report-receipts-dialog/report-receipts-dialog.component";
import { ReportBuilderValue, toReportRequestCommand } from "../models/report-command.mapper";
import { buildReportForm } from "../models/report-form.factory";
import { ReportPreviewPanelComponent } from "./report-preview-panel.component";

// buildReportForm wires untilDestroyed subscriptions, so it needs an
// @UntilDestroy()-decorated context.
@UntilDestroy()
@Component({ selector: "app-noop", template: "", standalone: false })
class NoopComponent {}

describe("ReportPreviewPanelComponent", () => {
  let fixture: ComponentFixture<ReportPreviewPanelComponent>;
  let component: ReportPreviewPanelComponent;
  let dialog: { open: jest.Mock };
  let form: FormGroup;

  beforeEach(async () => {
    dialog = { open: jest.fn() };
    await TestBed.configureTestingModule({
      declarations: [ReportPreviewPanelComponent, NoopComponent],
      providers: [provideZonelessChangeDetection(), FormBuilder, { provide: MatDialog, useValue: dialog }],
      schemas: [NO_ERRORS_SCHEMA],
    }).compileComponents();

    const host = TestBed.createComponent(NoopComponent).componentInstance;
    form = buildReportForm(TestBed.inject(FormBuilder), host);
    (form.get("scope") as FormArray).push(new FormControl("3", { nonNullable: true }));
    form.get("period.dateField")!.setValue("createdAt");

    fixture = TestBed.createComponent(ReportPreviewPanelComponent);
    component = fixture.componentInstance;
    fixture.componentRef.setInput("form", form);
  });

  // The drill-in asks the server for the receipts this exact command covers, so it
  // lists what the preview counted rather than re-deriving the period itself.
  it("opens the drill-in with the command the preview runs", () => {
    fixture.componentRef.setInput("receiptCount", 4);

    component.openReceipts();

    expect(dialog.open).toHaveBeenCalledTimes(1);
    const [dialogComponent, config] = dialog.open.mock.calls[0];
    const value = form.getRawValue() as ReportBuilderValue;
    expect(dialogComponent).toBe(ReportReceiptsDialogComponent);
    expect(config.data.command).toEqual(toReportRequestCommand(value));
    expect(config.data.command.period.dateField).toBe("createdAt");
    expect(config.data.period).toEqual(value.period);
    expect(config.data.receiptCount).toBe(4);
  });

  it("does not open the drill-in when the report covers no receipts", () => {
    fixture.componentRef.setInput("receiptCount", 0);

    component.openReceipts();

    expect(dialog.open).not.toHaveBeenCalled();
  });
});
