import { CUSTOM_ELEMENTS_SCHEMA } from "@angular/core";
import { ComponentFixture, TestBed } from "@angular/core/testing";
import { FormControl, ReactiveFormsModule } from "@angular/forms";
import { MatAutocompleteModule } from "@angular/material/autocomplete";
import { MatFormField, MatFormFieldModule } from "@angular/material/form-field";
import { MatInputModule } from "@angular/material/input";
import { By } from "@angular/platform-browser";
import { NoopAnimationsModule } from "@angular/platform-browser/animations";
import { TextareaComponent } from "./textarea.component";

describe("TextareaComponent", () => {
  let component: TextareaComponent;
  let fixture: ComponentFixture<TextareaComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [TextareaComponent],
      imports: [
        MatFormFieldModule,
        MatInputModule,
        ReactiveFormsModule,
        NoopAnimationsModule,
        MatAutocompleteModule
      ],
      schemas: [CUSTOM_ELEMENTS_SCHEMA],
    }).compileComponents();

    fixture = TestBed.createComponent(TextareaComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it("should create", () => {
    expect(component).toBeTruthy();
  });

  describe("subscript sizing", () => {
    function subscriptSizing(): string {
      return fixture.debugElement.query(By.directive(MatFormField))
        .componentInstance.subscriptSizing;
    }

    it("stays fixed without a hint", () => {
      expect(subscriptSizing()).toEqual("fixed");
    });

    it("goes dynamic with a hint, so a wrapping hint cannot overlap what follows", () => {
      fixture.componentRef.setInput(
        "hint",
        "A hint long enough to wrap onto more than one line."
      );
      fixture.detectChanges();

      expect(subscriptSizing()).toEqual("dynamic");
    });
  });

  it("should set selection end to where word was inserted", () => {
    fixture.componentRef.setInput('trigger', "@");
    component.inputFormControl = new FormControl("hello @trigger world");
    component.lastKnownSelection = 6;
    fixture.detectChanges();

    // Mock the matAutocompleteTrigger closePanel
    jest.spyOn(component.matAutocompleteTrigger(), "closePanel").mockImplementation(() => {});

    // Mock the textarea viewChild with a fake nativeElement to track selectionEnd
    const fakeNativeElement = { selectionEnd: 6 };
    Object.defineProperty(component, 'textarea', {
      value: () => ({ nativeElement: fakeNativeElement }),
      configurable: true,
    });

    component.onOptionSelected();

    expect(fakeNativeElement.selectionEnd).toBe(15);
  });
});
