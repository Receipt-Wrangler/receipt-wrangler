import { CUSTOM_ELEMENTS_SCHEMA } from "@angular/core";
import { ComponentFixture, TestBed } from "@angular/core/testing";
import { NgxsModule } from "@ngxs/store";
import { ButtonModule } from "../../button";
import { DialogFooterComponent } from "./dialog-footer.component";

describe("DialogFooterComponent", () => {
  let component: DialogFooterComponent;
  let fixture: ComponentFixture<DialogFooterComponent>;

  const query = (testId: string) =>
    fixture.nativeElement.querySelector(`[data-testid='${testId}']`);

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [DialogFooterComponent],
      imports: [ButtonModule, NgxsModule.forRoot([])],
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
    const cancel = query("dialog-cancel-button");
    expect(cancel.tooltip).toBe("Close");

    const cancelled = jest.fn();
    component.cancelClicked.subscribe(cancelled);
    cancel.dispatchEvent(new CustomEvent("clicked"));

    expect(cancelled).toHaveBeenCalled();
  });
});
