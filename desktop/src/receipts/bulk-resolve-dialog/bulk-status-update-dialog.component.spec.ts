import { CUSTOM_ELEMENTS_SCHEMA } from "@angular/core";
import { ComponentFixture, TestBed } from "@angular/core/testing";
import { ReactiveFormsModule } from "@angular/forms";
import { MatDialogRef } from "@angular/material/dialog";
import { ReceiptStatus } from "../../open-api";
import { PipesModule } from "../../pipes";
import { BulkStatusUpdateComponent } from "./bulk-status-update-dialog.component";

describe("BulkStatusUpdateComponent", () => {
  let component: BulkStatusUpdateComponent;
  let fixture: ComponentFixture<BulkStatusUpdateComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [BulkStatusUpdateComponent],
      imports: [ReactiveFormsModule, PipesModule],
      providers: [
        {
          provide: MatDialogRef,
          useValue: {
            close: (arg: any) => {},
          },
        },
      ],
      schemas: [CUSTOM_ELEMENTS_SCHEMA],
    }).compileComponents();

    fixture = TestBed.createComponent(BulkStatusUpdateComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it("should create", () => {
    expect(component).toBeTruthy();
  });

  it("should init form correctly", () => {
    component.ngOnInit();

    expect(component.form.value).toEqual({
      comment: "",
      status: ReceiptStatus.Resolved,
    });
  });

  it("should close dialog with undefined", () => {
    const spy = jest.spyOn(component.matDialogRef, "close");
    component.cancelButtonClicked();

    expect(spy).toHaveBeenCalledWith(undefined);
  });

  it("should close dialog with form value", () => {
    const spy = jest.spyOn(component.matDialogRef, "close");
    component.form.patchValue({
      comment: "resolved",
      status: ReceiptStatus.NeedsAttention,
    });
    component.submitButtonClicked();

    expect(spy).toHaveBeenCalledWith({
      comment: "resolved",
      status: ReceiptStatus.NeedsAttention,
    });
  });
});
