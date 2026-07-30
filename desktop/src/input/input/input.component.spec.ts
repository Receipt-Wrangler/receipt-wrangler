import { CommonModule } from "@angular/common";
import { provideZonelessChangeDetection } from "@angular/core";
import { ComponentFixture, TestBed } from "@angular/core/testing";
import { FormControl, ReactiveFormsModule } from "@angular/forms";
import { MatButtonModule } from "@angular/material/button";
import { MatFormFieldModule } from "@angular/material/form-field";
import { MatIconModule } from "@angular/material/icon";
import { MatInputModule } from "@angular/material/input";
import { MatTooltipModule } from "@angular/material/tooltip";
import { NoopAnimationsModule } from "@angular/platform-browser/animations";
import { NgxsModule } from "@ngxs/store";
import { NgxMaskDirective, provideNgxMask } from "ngx-mask";
import { PasswordGeneratorService } from "../../services/password-generator.service";
import { SystemSettingsState } from "../../store/system-settings.state";
import { InputComponent } from "./input.component";

const GENERATED_PASSWORD = "Fixed-Pass-1";

describe("InputComponent", () => {
  let component: InputComponent;
  let fixture: ComponentFixture<InputComponent>;
  let passwordGeneratorService: { generateAndCopy: jest.Mock };

  beforeEach(async () => {
    passwordGeneratorService = {
      generateAndCopy: jest.fn().mockReturnValue(GENERATED_PASSWORD),
    };

    await TestBed.configureTestingModule({
      declarations: [InputComponent],
      imports: [
        CommonModule,
        MatButtonModule,
        MatFormFieldModule,
        MatIconModule,
        MatInputModule,
        MatTooltipModule,
        NoopAnimationsModule,
        ReactiveFormsModule,
        NgxMaskDirective,
        NgxsModule.forRoot([SystemSettingsState]),
      ],
      providers: [
        provideNgxMask(),
        provideZonelessChangeDetection(),
        {
          provide: PasswordGeneratorService,
          useValue: passwordGeneratorService,
        },
      ],
    }).compileComponents();

    fixture = TestBed.createComponent(InputComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it("should create", () => {
    expect(component).toBeTruthy();
  });

  describe("generate password button", () => {
    function generateButton(): HTMLButtonElement | null {
      return fixture.nativeElement.querySelector(
        '[data-testid="password-generate"]'
      );
    }

    function iconButtonCount(): number {
      return fixture.nativeElement.querySelectorAll("button[mat-icon-button]")
        .length;
    }

    it("is not rendered by default", () => {
      expect(generateButton()).toBeNull();
      expect(iconButtonCount()).toEqual(0);
    });

    it("is not rendered for a field that only shows the visibility eye", async () => {
      fixture.componentRef.setInput("showVisibilityEye", true);
      await fixture.whenStable();

      expect(generateButton()).toBeNull();
      expect(iconButtonCount()).toEqual(1);
    });

    it("renders alongside the visibility eye when enabled", async () => {
      fixture.componentRef.setInput("showVisibilityEye", true);
      fixture.componentRef.setInput("showGeneratePassword", true);
      await fixture.whenStable();

      expect(generateButton()).not.toBeNull();
      expect(iconButtonCount()).toEqual(2);
    });

    it("masks the field on init when only the generate button is enabled", async () => {
      fixture.componentRef.setInput("showGeneratePassword", true);
      await fixture.whenStable();

      expect(component.type).toEqual("password");
    });

    it("masks the field when the generate button is enabled later", async () => {
      fixture.componentRef.setInput("showGeneratePassword", false);
      await fixture.whenStable();
      expect(component.type).toEqual("text");

      // A parent flipping the flag on after init still turns this into a
      // password field — otherwise typed passwords would sit in plain text.
      fixture.componentRef.setInput("showGeneratePassword", true);
      await fixture.whenStable();

      expect(component.type).toEqual("password");
    });

    it("does not re-mask a revealed password when an unrelated input changes", async () => {
      fixture.componentRef.setInput("showGeneratePassword", true);
      await fixture.whenStable();

      component.generatePassword();
      expect(component.type).toEqual("text");

      fixture.componentRef.setInput("label", "New Password");
      await fixture.whenStable();

      expect(component.type).toEqual("text");
    });

    it("fills the control with a generated password and reveals it", async () => {
      fixture.componentRef.setInput("inputFormControl", new FormControl(""));
      fixture.componentRef.setInput("showGeneratePassword", true);
      await fixture.whenStable();

      generateButton()?.click();
      await fixture.whenStable();

      expect(passwordGeneratorService.generateAndCopy).toHaveBeenCalled();
      expect(component.inputFormControl.value).toEqual(GENERATED_PASSWORD);
      expect(component.inputFormControl.dirty).toEqual(true);
      expect(component.type).toEqual("text");
      expect(
        fixture.nativeElement.querySelector("input").getAttribute("type")
      ).toEqual("text");
    });

    it("can be re-masked with the visibility eye after generating", () => {
      component.type = "password";

      component.generatePassword();
      expect(component.type).toEqual("text");

      component.toggleVisibility();
      expect(component.type).toEqual("password");
    });

    it("does nothing for a disabled control", async () => {
      fixture.componentRef.setInput(
        "inputFormControl",
        new FormControl({ value: "", disabled: true })
      );
      fixture.componentRef.setInput("showGeneratePassword", true);
      await fixture.whenStable();

      expect(generateButton()?.disabled).toEqual(true);

      component.generatePassword();

      expect(passwordGeneratorService.generateAndCopy).not.toHaveBeenCalled();
      expect(component.inputFormControl.value).toEqual("");
    });

    it("does nothing for a readonly field", () => {
      fixture.componentRef.setInput("readonly", true);

      component.generatePassword();

      expect(passwordGeneratorService.generateAndCopy).not.toHaveBeenCalled();
    });
  });
});
