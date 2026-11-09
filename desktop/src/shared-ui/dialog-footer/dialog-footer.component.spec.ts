import { CUSTOM_ELEMENTS_SCHEMA } from "@angular/core";
import { ComponentFixture, TestBed } from "@angular/core/testing";
import { By } from "@angular/platform-browser";
import { NgxsModule } from "@ngxs/store";
import { ButtonModule } from "../../button";
import { CancelButtonComponent } from "../cancel-button/cancel-button.component";
import { DialogFooterComponent } from "./dialog-footer.component";
import { provideRouter } from "@angular/router";

describe("DialogFooterComponent", () => {
  let component: DialogFooterComponent;
  let fixture: ComponentFixture<DialogFooterComponent>;

  const query = (testId: string) =>
    fixture.nativeElement.querySelector(`[data-testid='${testId}']`);

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [ButtonModule, NgxsModule.forRoot([]), DialogFooterComponent],
      providers: [provideRouter([])],
      schemas: [CUSTOM_ELEMENTS_SCHEMA],
    }).compileComponents();

    fixture = TestBed.createComponent(DialogFooterComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it("should create", () => {
    expect(component).toBeTruthy();
  });

  it("renders submit and cancel by default", () => {
    expect(query("dialog-submit-button")).not.toBeNull();
    expect(query("dialog-cancel-button")).not.toBeNull();
  });

  // A dialog with nothing to submit keeps the cancel button as its only action.
  it("renders only the cancel button when the submit button is hidden", () => {
    fixture.componentRef.setInput("showSubmitButton", false);
    fixture.componentRef.setInput("cancelButtonTooltip", "Close");
    fixture.detectChanges();

    expect(query("dialog-submit-button")).toBeNull();
    const cancel: CancelButtonComponent = fixture.debugElement
      .query(By.css("[data-testid='dialog-cancel-button']"))
      .componentInstance;
    expect(cancel.tooltip()).toBe("Close");

    const cancelled = jest.fn();
    component.cancelClicked.subscribe(cancelled);
    cancel.clicked.emit();

    expect(cancelled).toHaveBeenCalled();
  });
});
