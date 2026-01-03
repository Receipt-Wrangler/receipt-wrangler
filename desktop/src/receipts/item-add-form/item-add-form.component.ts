import {
  Component,
  ElementRef,
  EventEmitter,
  HostListener,
  Input,
  OnDestroy,
  OnInit,
  Output,
  ViewChild,
  ViewEncapsulation,
} from "@angular/core";
import { FormGroup } from "@angular/forms";
import { Subject, takeUntil } from "rxjs";
import { FormMode } from "src/enums/form-mode.enum";
import { InputComponent } from "../../input";
import { Category, Group, Item, Tag } from "../../open-api";
import { KeyboardShortcutService } from "../../services/keyboard-shortcut.service";
import { buildItemForm } from "../utils/form.utils";
import { KEYBOARD_SHORTCUT_ACTIONS, DISPLAY_SHORTCUTS } from "../../constants/keyboard-shortcuts.constant";

@Component({
  selector: "app-item-add-form",
  templateUrl: "./item-add-form.component.html",
  styleUrls: ["./item-add-form.component.scss"],
  encapsulation: ViewEncapsulation.None,
  standalone: false
})
export class ItemAddFormComponent implements OnInit, OnDestroy {
  @ViewChild("addForm")
  public addForm!: ElementRef;

  @ViewChild("nameInput")
  public nameInput!: InputComponent;

  @ViewChild("amountInput")
  public amountInput!: InputComponent;

  @ViewChild("categoryInput")
  public categoryInput!: ElementRef;

  @ViewChild("tagInput")
  public tagInput!: ElementRef;

  @Input() public categories: Category[] = [];
  @Input() public tags: Tag[] = [];
  @Input() public selectedGroup?: Group;
  @Input() public mode: FormMode = FormMode.add;
  @Input() public receiptId?: string;

  @Output() public itemAdded = new EventEmitter<Item>();
  @Output() public cancelled = new EventEmitter<void>();
  @Output() public submitAndContinue = new EventEmitter<Item>();
  @Output() public submitAndFinish = new EventEmitter<Item>();

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
    if (this.mode === FormMode.view) {
      return;
    }

    // Let the service handle the keyboard event
    this.keyboardShortcutService.handleKeyboardEvent(event);
  }

  private initializeForm(): void {
    this.newItemFormGroup = buildItemForm(
      undefined,
      this.receiptId,
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
    if (this.selectedGroup?.groupReceiptSettings?.hideItemCategories) {
      if (this.newItemFormGroup.valid) {
        this.onSubmitAndContinue();
      }
      return;
    }
    
    this.focusCategoryField();
  }

  public onCategoryEnter(event: Event): void {
    event.preventDefault();
    
    if (this.selectedGroup?.groupReceiptSettings?.hideItemTags) {
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
    if (this.nameInput?.nativeInput?.nativeElement) {
      (this.nameInput.nativeInput.nativeElement as HTMLInputElement).focus();
    }
  }

  private focusAmountField(): void {
    if (this.amountInput?.nativeInput?.nativeElement) {
      (this.amountInput.nativeInput.nativeElement as HTMLInputElement).focus();
    }
  }

  private focusCategoryField(): void {
    if (this.categoryInput?.nativeElement) {
      const input = this.categoryInput.nativeElement.querySelector('input');
      if (input) {
        input.focus();
      }
    }
  }

  private focusTagField(): void {
    if (this.tagInput?.nativeElement) {
      const input = this.tagInput.nativeElement.querySelector('input');
      if (input) {
        input.focus();
      }
    }
  }
}