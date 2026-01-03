import { CUSTOM_ELEMENTS_SCHEMA } from "@angular/core";
import { ComponentFixture, TestBed } from "@angular/core/testing";
import { ReactiveFormsModule } from "@angular/forms";
import { Subject } from "rxjs";
import { FormMode } from "src/enums/form-mode.enum";
import { PipesModule } from "src/pipes/pipes.module";
import { KEYBOARD_SHORTCUT_ACTIONS } from "../../constants/keyboard-shortcuts.constant";
import { Group } from "../../open-api";
import { KeyboardShortcutService } from "../../services/keyboard-shortcut.service";
import { ItemAddFormComponent } from "./item-add-form.component";

describe("ItemAddFormComponent", () => {
  let component: ItemAddFormComponent;
  let fixture: ComponentFixture<ItemAddFormComponent>;
  let keyboardShortcutService: jest.Mocked<KeyboardShortcutService>;
  let shortcutTriggered$: Subject<any>;
  let showHint$: Subject<boolean>;

  beforeEach(async () => {
    shortcutTriggered$ = new Subject();
    showHint$ = new Subject();

    const keyboardSpy = {
      handleKeyboardEvent: jest.fn(),
      shortcutTriggered: shortcutTriggered$.asObservable(),
      showHint: showHint$.asObservable()
    };

    await TestBed.configureTestingModule({
      declarations: [ItemAddFormComponent],
      schemas: [CUSTOM_ELEMENTS_SCHEMA],
      imports: [ReactiveFormsModule, PipesModule],
      providers: [
        { provide: KeyboardShortcutService, useValue: keyboardSpy }
      ]
    }).compileComponents();

    fixture = TestBed.createComponent(ItemAddFormComponent);
    component = fixture.componentInstance;
    keyboardShortcutService = TestBed.inject(KeyboardShortcutService) as jest.Mocked<KeyboardShortcutService>;
  });

  it("should create", () => {
    expect(component).toBeTruthy();
  });

  it("should initialize form on ngOnInit", () => {
    component.receiptId = "123";
    component.ngOnInit();

    expect(component.newItemFormGroup).toBeDefined();
    expect(component.newItemFormGroup.get("name")).toBeDefined();
    expect(component.newItemFormGroup.get("amount")).toBeDefined();
  });

  it("should setup keyboard shortcuts on init", () => {
    component.ngOnInit();

    // Emit a shortcut event
    shortcutTriggered$.next({ action: KEYBOARD_SHORTCUT_ACTIONS.SUBMIT_AND_CONTINUE, event: new KeyboardEvent("keydown") });

    // Should have subscribed to the service
    expect((component as any).showKeyboardHint).toBe(false);
  });

  it("should show keyboard hint when service emits", () => {
    component.ngOnInit();

    showHint$.next(true);
    expect((component as any).showKeyboardHint).toBe(true);

    showHint$.next(false);
    expect((component as any).showKeyboardHint).toBe(false);
  });

  describe("keyboard shortcut actions", () => {
    beforeEach(() => {
      component.ngOnInit();
      // Mock form as valid
      Object.defineProperty(component.newItemFormGroup, 'valid', {
        get: () => true
      });
      jest.spyOn(component, "onSubmitAndContinue");
      jest.spyOn(component, "onSubmitAndFinish");
      jest.spyOn(component, "onCancel");
    });

    it("should handle SUBMIT_AND_CONTINUE action", () => {
      shortcutTriggered$.next({ action: KEYBOARD_SHORTCUT_ACTIONS.SUBMIT_AND_CONTINUE, event: new KeyboardEvent("keydown") });

      expect(component.onSubmitAndContinue).toHaveBeenCalled();
    });

    it("should handle SUBMIT_AND_FINISH action", () => {
      shortcutTriggered$.next({ action: KEYBOARD_SHORTCUT_ACTIONS.SUBMIT_AND_FINISH, event: new KeyboardEvent("keydown") });

      expect(component.onSubmitAndFinish).toHaveBeenCalled();
    });

    it("should handle CANCEL action", () => {
      shortcutTriggered$.next({ action: KEYBOARD_SHORTCUT_ACTIONS.CANCEL, event: new KeyboardEvent("keydown") });

      expect(component.onCancel).toHaveBeenCalled();
    });
  });

  describe("form submission", () => {
    beforeEach(() => {
      component.ngOnInit();
      component.newItemFormGroup.patchValue({
        name: "Test Item",
        amount: 10.00
      });
    });

    it("should emit submitAndContinue event with item data", () => {
      jest.spyOn(component.submitAndContinue, "emit");

      component.onSubmitAndContinue();

      expect(component.submitAndContinue.emit).toHaveBeenCalledWith(
        expect.objectContaining({
          name: "Test Item",
          amount: 10.00,
          chargedToUserId: undefined
        })
      );
      expect(component.rapidAddMode).toBe(true);
    });

    it("should emit submitAndFinish event with item data", () => {
      jest.spyOn(component.submitAndFinish, "emit");

      component.onSubmitAndFinish();

      expect(component.submitAndFinish.emit).toHaveBeenCalledWith(
        expect.objectContaining({
          name: "Test Item",
          amount: 10.00,
          chargedToUserId: undefined
        })
      );
    });

    it("should not submit if form is invalid", () => {
      jest.spyOn(component.submitAndContinue, "emit");
      component.newItemFormGroup.patchValue({ name: "" }); // Make form invalid

      component.onSubmitAndContinue();

      expect(component.submitAndContinue.emit).not.toHaveBeenCalled();
    });
  });

  it("should emit cancelled event on cancel", () => {
    jest.spyOn(component.cancelled, "emit");

    component.onCancel();

    expect(component.cancelled.emit).toHaveBeenCalled();
  });

  describe("field navigation", () => {
    beforeEach(() => {
      component.ngOnInit();
      // Mock focus methods
      jest.spyOn(component as any, "focusAmountField");
      jest.spyOn(component as any, "focusCategoryField");
      jest.spyOn(component as any, "focusTagField");
      jest.spyOn(component, "onSubmitAndContinue");
    });

    it("should focus amount field on name enter", () => {
      const event = new Event("keydown");
      jest.spyOn(event, "preventDefault");

      component.onNameEnter(event);

      expect(event.preventDefault).toHaveBeenCalled();
      expect(component["focusAmountField"]).toHaveBeenCalled();
    });

    it("should submit directly if categories are hidden", () => {
      component.selectedGroup = {
        groupReceiptSettings: {
          hideItemCategories: true
        }
      } as Group;
      component.newItemFormGroup.patchValue({ name: "Test", amount: 10 });

      const event = new Event("keydown");
      component.onAmountEnter(event);

      expect(component.onSubmitAndContinue).toHaveBeenCalled();
    });

    it("should focus category field if categories are shown", () => {
      component.selectedGroup = {
        groupReceiptSettings: {
          hideItemCategories: false
        }
      } as Group;

      const event = new Event("keydown");
      component.onAmountEnter(event);

      expect(component["focusCategoryField"]).toHaveBeenCalled();
    });
  });

  it("should handle keyboard events in view mode", () => {
    component.mode = FormMode.view;
    const event = new KeyboardEvent("keydown");

    component.handleKeyboardShortcut(event);

    expect(keyboardShortcutService.handleKeyboardEvent).not.toHaveBeenCalled();
  });

  it("should handle keyboard events in edit mode", () => {
    component.mode = FormMode.edit;
    const event = new KeyboardEvent("keydown");

    component.handleKeyboardShortcut(event);

    expect(keyboardShortcutService.handleKeyboardEvent).toHaveBeenCalledWith(event);
  });

  it("should cleanup on destroy", () => {
    component.ngOnInit();
    jest.spyOn(component["destroy$"], "next");
    jest.spyOn(component["destroy$"], "complete");

    component.ngOnDestroy();

    expect(component["destroy$"].next).toHaveBeenCalled();
    expect(component["destroy$"].complete).toHaveBeenCalled();
  });
});
