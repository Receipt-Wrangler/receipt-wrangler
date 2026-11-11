import {
  Component,
  ElementRef,
  HostListener,
  OnDestroy,
  OnInit,
  ViewEncapsulation,
  input,
  output,
  viewChild
} from "@angular/core";
import { FormGroup } from "@angular/forms";
import { Subject, takeUntil } from "rxjs";
import { FormMode } from "src/enums/form-mode.enum";
import { InputComponent } from "../../input";
import { Category, Group, Item, Tag } from "../../open-api";
import { KeyboardShortcutService } from "../../services/keyboard-shortcut.service";
import { buildItemForm } from "../utils/form.utils";
import { KEYBOARD_SHORTCUT_ACTIONS, DISPLAY_SHORTCUTS } from "../../constants/keyboard-shortcuts.constant";
import { CardComponent } from "../../shared-ui/card/card.component";
import { CategoryAutocompleteComponent } from "../../category-autocomplete/category-autocomplete.component";
import { TagAutocompleteComponent } from "../../tag-autocomplete/tag-autocomplete.component";
import { ButtonComponent } from "../../button/button/button.component";
import { FormGetPipe } from "../../pipes/form-get.pipe";
import { InputReadonlyPipe } from "../../pipes/input-readonly.pipe";

@Component({
  selector: "app-item-add-form",
  templateUrl: "./item-add-form.component.html",
  styleUrls: ["./item-add-form.component.scss"],
  encapsulation: ViewEncapsulation.None,
  imports: [CardComponent, InputComponent, CategoryAutocompleteComponent, TagAutocompleteComponent, ButtonComponent, FormGetPipe, InputReadonlyPipe]
})
export class ItemAddFormComponent implements OnInit, OnDestroy {
  public readonly addForm = viewChild.required<ElementRef>("addForm");

  public readonly nameInput = viewChild.required<InputComponent>("nameInput");

  public readonly amountInput = viewChild.required<InputComponent>("amountInput");

  public readonly categoryInput = viewChild.required<ElementRef>("categoryInput");

  public readonly tagInput = viewChild.required<ElementRef>("tagInput");

  public readonly categories = input<Category[]>([]);
  public readonly tags = input<Tag[]>([]);
  public readonly selectedGroup = input<Group>();
  public readonly mode = input<FormMode>(FormMode.add);
  public readonly receiptId = input<string>();

  public readonly itemAdded = output<Item>();
  public readonly cancelled = output<void>();
  public readonly submitAndContinue = output<Item>();
  public readonly submitAndFinish = output<Item>();

  public formMode = FormMode;
  public newItemFormGroup!: FormGroup;
  public rapidAddMode: boolean = false;
  public displayShortcuts = DISPLAY_SHORTCUTS;
  public showKeyboardHint: boolean = false;

  private destroy$ = new Subject<void>();

  constructor(private keyboardShortcutService: KeyboardShortcutService) {}

  public ngOnInit(): void {
    this.initializeForm();
    this.setupKeyboardShortcuts();
    
    // Auto-focus name field after view init
    setTimeout(() => {
      this.focusNameField();
    }, 50);
  }

  public ngOnDestroy(): void {
    this.destroy$.next();
    this.destroy$.complete();
  }

  @HostListener("document:keydown", ["$event"])
  public handleKeyboardShortcut(event: KeyboardEvent): void {
    // Only handle shortcuts when this form is active
    if (this.mode() === FormMode.view) {
      return;
    }

    // Let the service handle the keyboard event
    this.keyboardShortcutService.handleKeyboardEvent(event);
  }

  private initializeForm(): void {
    this.newItemFormGroup = buildItemForm(
      undefined,
      this.receiptId(),
      false,
      false
    );
  }

  private setupKeyboardShortcuts(): void {
    // Subscribe to keyboard shortcut events
    this.keyboardShortcutService.shortcutTriggered
      .pipe(takeUntil(this.destroy$))
      .subscribe(shortcutEvent => {
        this.handleShortcutAction(shortcutEvent.action);
      });

    // Subscribe to show hint observable
    this.keyboardShortcutService.showHint
      .pipe(takeUntil(this.destroy$))
      .subscribe(showHint => {
        this.showKeyboardHint = showHint;
      });
  }

  private handleShortcutAction(action: string): void {
    switch (action) {
      case KEYBOARD_SHORTCUT_ACTIONS.SUBMIT_AND_CONTINUE:
        if (this.newItemFormGroup.valid) {
          this.onSubmitAndContinue();
        }
        break;
      case KEYBOARD_SHORTCUT_ACTIONS.SUBMIT_AND_FINISH:
        if (this.newItemFormGroup.valid) {
          this.onSubmitAndFinish();
        }
        break;
      case KEYBOARD_SHORTCUT_ACTIONS.CANCEL:
        this.onCancel();
        break;
    }
  }

  public onSubmitAndContinue(): void {
    if (this.newItemFormGroup.valid) {
      const newItem = this.newItemFormGroup.value as Item;
      newItem.chargedToUserId = undefined;
      
      this.submitAndContinue.emit(newItem);
      
      // Reset form for rapid add mode
      this.rapidAddMode = true;
      this.initializeForm();
      
      // Re-focus name field
      setTimeout(() => {
        this.focusNameField();
      }, 50);
    }
  }

  public onSubmitAndFinish(): void {
    if (this.newItemFormGroup.valid) {
      const newItem = this.newItemFormGroup.value as Item;
      newItem.chargedToUserId = undefined;
      
      this.submitAndFinish.emit(newItem);
    }
  }

  public onCancel(): void {

    this.cancelled.emit();
  }

  public onNameEnter(event: Event): void {
    event.preventDefault();
    this.focusAmountField();
  }

  public onAmountEnter(event: Event): void {
    event.preventDefault();
    
    // If categories are hidden, submit directly
    if (this.selectedGroup()?.groupReceiptSettings?.hideItemCategories) {
      if (this.newItemFormGroup.valid) {
        this.onSubmitAndContinue();
      }
      return;
    }
    
    this.focusCategoryField();
  }

  public onCategoryEnter(event: Event): void {
    event.preventDefault();
    
    if (this.selectedGroup()?.groupReceiptSettings?.hideItemTags) {
      this.onSubmitAndContinue();
      return;
    }
    
    this.focusTagField();
  }

  public onTagEnter(event: Event): void {
    event.preventDefault();
    this.onSubmitAndContinue();
  }

  private focusNameField(): void {
    const nameInput = this.nameInput();
    const nativeInput = nameInput?.nativeInput();
    if (nativeInput?.nativeElement) {
      (nativeInput.nativeElement as HTMLInputElement).focus();
    }
  }

  private focusAmountField(): void {
    const amountInput = this.amountInput();
    const nativeInput = amountInput?.nativeInput();
    if (nativeInput?.nativeElement) {
      (nativeInput.nativeElement as HTMLInputElement).focus();
    }
  }

  private focusCategoryField(): void {
    const categoryInput = this.categoryInput();
    if (categoryInput?.nativeElement) {
      const input = categoryInput.nativeElement.querySelector('input');
      if (input) {
        input.focus();
      }
    }
  }

  private focusTagField(): void {
    const tagInput = this.tagInput();
    if (tagInput?.nativeElement) {
      const input = tagInput.nativeElement.querySelector('input');
      if (input) {
        input.focus();
      }
    }
  }
}