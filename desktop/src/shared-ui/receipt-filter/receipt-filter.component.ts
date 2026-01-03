import { Component, EventEmitter, Input, OnInit, Output, TemplateRef, } from "@angular/core";
import { FormControl, FormGroup, } from "@angular/forms";
import { MatDialogRef } from "@angular/material/dialog";
import { Store } from "@ngxs/store";
import { endOfDay, startOfMonth } from "date-fns";
import { forkJoin, take, tap } from "rxjs";
import { RECEIPT_STATUS_OPTIONS } from "src/constants";
import { SetReceiptFilter } from "src/store/receipt-table.actions";
import { FormCommand } from "../../form/index";
import { Category, CategoryService, FilterOperation, Tag, TagService } from "../../open-api";
import { OperationsPipe } from "./operations.pipe";

@Component({
  selector: "app-receipt-filter",
  templateUrl: "./receipt-filter.component.html",
  styleUrls: ["./receipt-filter.component.scss"],
  standalone: false
})
export class ReceiptFilterComponent implements OnInit {
  @Input() public headerText: string = "";

  @Input() public footerTemplate?: TemplateRef<any>;

  @Input() public isOpen: boolean = true;

  @Input() public previewTemplate?: TemplateRef<any>;

  @Input() public previewTemplateContext?: any;

  @Input() public inDialog: boolean = true;

  @Input() public parentForm: FormGroup = new FormGroup({});

  @Input() public basePath: string = "";

  @Output() public formCommand: EventEmitter<FormCommand> = new EventEmitter<FormCommand>();

  @Output() public formInitialized: EventEmitter<FormGroup> =
    new EventEmitter<FormGroup>();

  public receiptStatusOptions = RECEIPT_STATUS_OPTIONS;

  public categories: Category[] = [];

  public tags: Tag[] = [];

  public startOfMonthFormControl = new FormControl(startOfMonth(new Date()));

  public endOfTodayFormControl = new FormControl(endOfDay(new Date()));

  private operationsPipe = new OperationsPipe();

  constructor(
    private store: Store,
    private dialogRef: MatDialogRef<ReceiptFilterComponent>,
    private categoryService: CategoryService,
    private tagService: TagService
  ) {}

  public ngOnInit(): void {
    this.startOfMonthFormControl.disable();
    this.endOfTodayFormControl.disable();

    forkJoin([
      this.categoryService.getAllCategories(),
      this.tagService.getAllTags(),
    ])
      .pipe(
        take(1),
        tap(([categories, tags]) => {
          this.categories = categories;
          this.tags = tags;
          this.setupAutoOperationSelection();
        })
      )
      .subscribe();
  }

  public resetFilter(): void {
    this.formCommand.emit({
      path: `${this.basePath}`,
      command: "reset",
    });
    this.formCommand.emit({
      path: `${this.basePath}paidBy.value`,
      command: "clear",
    });
    this.formCommand.emit({
      path: `${this.basePath}categories.value`,
      command: "clear",
    });
    this.formCommand.emit({
      path: `${this.basePath}tags.value`,
      command: "clear",
    });
    this.formCommand.emit({
      path: `${this.basePath}status.value`,
      command: "clear",
    });
  }

  public submitButtonClicked(): void {
    const filter = this.parentForm.value;

    if (this.parentForm.valid) {
      this.store
        .dispatch(new SetReceiptFilter(filter))
        .pipe(
          take(1),
          tap(() => {
            this.dialogRef.close(true);
          })
        )
        .subscribe();
    } else {
      this.parentForm.markAllAsTouched();
    }
  }

  public cancelButtonClicked(): void {
    this.dialogRef.close(false);
  }

  private setupAutoOperationSelection(): void {
    // List of all filter fields
    const fieldsToWatch = [
      { fieldName: "date", type: "date" },
      { fieldName: "name", type: "text" },
      { fieldName: "paidBy", type: "users" },
      { fieldName: "amount", type: "number" },
      { fieldName: "categories", type: "list" },
      { fieldName: "tags", type: "list" },
      { fieldName: "status", type: "list" },
      { fieldName: "resolvedDate", type: "date" },
      { fieldName: "createdAt", type: "date" }
    ];

    fieldsToWatch.forEach(({ fieldName, type }) => {
      const valueControl = this.parentForm.get(`${this.basePath}${fieldName}.value`);
      const operationControl = this.parentForm.get(`${this.basePath}${fieldName}.operation`);

      if (valueControl && operationControl) {
        valueControl.valueChanges.subscribe(value => {
          const hasValue = this.hasFieldValue(value, type);

          if (hasValue) {
            // Set first operation if none is selected
            if (!operationControl.value) {
              const operations = this.operationsPipe.transform(type, false);
              if (operations.length > 0) {
                operationControl.setValue(operations[0]);
              }
            }
          } else {
            // Clear operation if field is empty
            operationControl.setValue(null);
          }
        });
      }
    });
  }

  private hasFieldValue(value: any, type: string): boolean {
    if (value === null || value === undefined) {
      return false;
    }

    if (type === "list" || type === "users") {
      return Array.isArray(value) && value.length > 0;
    }

    if (typeof value === "string") {
      return value.trim().length > 0;
    }

    return value !== "";
  }

  protected readonly FilterOperation = FilterOperation;
}
