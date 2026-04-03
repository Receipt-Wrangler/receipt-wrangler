import { Component, Input, OnInit, TemplateRef, input, output } from "@angular/core";
import { FormArray, FormBuilder, FormControl, FormGroup, } from "@angular/forms";
import { MatDialogRef } from "@angular/material/dialog";
import { UntilDestroy } from "@ngneat/until-destroy";
import { Store } from "@ngxs/store";
import { endOfDay, startOfMonth } from "date-fns";
import { forkJoin, take, tap } from "rxjs";
import { RECEIPT_STATUS_OPTIONS } from "src/constants";
import { SetReceiptFilter } from "src/store/receipt-table.actions";
import { buildCustomFieldFilterFormGroup, listenForBetweenOperation } from "src/utils/receipt-filter";
import { FormCommand } from "../../form/index";
import { Category, CategoryService, CustomField, CustomFieldService, CustomFieldType, FilterOperation, Tag, TagService } from "../../open-api";
import { OperationsPipe } from "./operations.pipe";

@UntilDestroy()
@Component({
  selector: "app-receipt-filter",
  templateUrl: "./receipt-filter.component.html",
  styleUrls: ["./receipt-filter.component.scss"],
  standalone: false
})
export class ReceiptFilterComponent implements OnInit {
  @Input() public headerText: string = "";

  public readonly footerTemplate = input<TemplateRef<any>>();

  public readonly isOpen = input<boolean>(true);

  @Input() public previewTemplate?: TemplateRef<any>;

  public readonly previewTemplateContext = input<any>();

  public readonly inDialog = input<boolean>(true);

  @Input() public parentForm: FormGroup = new FormGroup({});

  @Input() public basePath: string = "";

  public readonly formCommand = output<FormCommand>();

  public readonly formInitialized = output<FormGroup>();

  public receiptStatusOptions = RECEIPT_STATUS_OPTIONS;

  public categories: Category[] = [];

  public tags: Tag[] = [];

  public customFields: CustomField[] = [];

  public booleanOptions = [
    { value: true, displayValue: "True" },
    { value: false, displayValue: "False" },
  ];

  public startOfMonthFormControl = new FormControl(startOfMonth(new Date()));

  public endOfTodayFormControl = new FormControl(endOfDay(new Date()));

  private operationsPipe = new OperationsPipe();

  constructor(
    private store: Store,
    private dialogRef: MatDialogRef<ReceiptFilterComponent>,
    private categoryService: CategoryService,
    private customFieldService: CustomFieldService,
    private tagService: TagService,
  ) {}

  public ngOnInit(): void {
    this.startOfMonthFormControl.disable();
    this.endOfTodayFormControl.disable();

    forkJoin([
      this.categoryService.getAllCategories(),
      this.tagService.getAllTags(),
      this.customFieldService.getAllCustomFields(),
    ])
      .pipe(
        take(1),
        tap(([categories, tags, customFields]) => {
          this.categories = categories;
          this.tags = tags;
          this.customFields = customFields;
          this.setupAutoOperationSelection();
          this.setupCustomFieldAutoOperationSelection();
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

    const customFieldsArray = this.getCustomFieldsArray();
    while (customFieldsArray.length > 0) {
      customFieldsArray.removeAt(0);
    }
  }

  public submitButtonClicked(): void {
    const filter = this.parentForm.value;

    if (this.parentForm.valid) {
      filter.customFields = (filter.customFields || []).filter((cf: any) =>
        cf.customFieldId != null && cf.operation != null && cf.value != null
      );

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

  public getCustomFieldsArray(): FormArray {
    return this.parentForm.get(`${this.basePath}customFields`) as FormArray;
  }

  public addCustomFieldFilter(): void {
    const group = buildCustomFieldFilterFormGroup();
    const index = this.getCustomFieldsArray().length;
    this.getCustomFieldsArray().push(group);

    group.get("customFieldId")?.valueChanges.subscribe(() => {
      this.onCustomFieldSelected(index);
    });
  }

  public removeCustomFieldFilter(index: number): void {
    this.getCustomFieldsArray().removeAt(index);
  }

  public getSelectedCustomField(index: number): CustomField | undefined {
    const id = this.getCustomFieldsArray().at(index).get("customFieldId")?.value;
    return this.customFields.find(cf => cf.id === id);
  }

  public getAvailableCustomFields(index: number): CustomField[] {
    const usedIds = new Set<number>();
    const array = this.getCustomFieldsArray();
    for (let i = 0; i < array.length; i++) {
      if (i !== index) {
        const id = array.at(i).get("customFieldId")?.value;
        if (id != null) {
          usedIds.add(id);
        }
      }
    }
    return this.customFields.filter(cf => !usedIds.has(cf.id!));
  }

  public getFilterType(customField: CustomField): string {
    switch (customField.type) {
      case CustomFieldType.Text: return "text";
      case CustomFieldType.Date: return "date";
      case CustomFieldType.Currency: return "currency";
      case CustomFieldType.Boolean: return "boolean";
      case CustomFieldType.Select: return "list";
      default: return "text";
    }
  }

  public onCustomFieldSelected(index: number): void {
    const group = this.getCustomFieldsArray().at(index) as FormGroup;
    const formBuilder = new FormBuilder();
    const selectedField = this.getSelectedCustomField(index);

    group.get("operation")?.setValue(null);
    group.removeControl("value");

    if (selectedField && this.getFilterType(selectedField) === "list") {
      group.addControl("value", formBuilder.array([]));
    } else {
      group.addControl("value", formBuilder.control(null));
    }

    if (selectedField) {
      const filterType = this.getFilterType(selectedField);
      if (filterType === "date" || filterType === "currency") {
        listenForBetweenOperation(this.parentForm, `${this.basePath}customFields.${index}`, this);
      }
    }

    this.setupCustomFieldAutoOperationSelectionForIndex(index);
  }

  private setupCustomFieldAutoOperationSelection(): void {
    const array = this.getCustomFieldsArray();
    for (let i = 0; i < array.length; i++) {
      this.setupCustomFieldAutoOperationSelectionForIndex(i);

      const selectedField = this.getSelectedCustomField(i);
      if (selectedField) {
        const filterType = this.getFilterType(selectedField);
        if (filterType === "date" || filterType === "currency") {
          listenForBetweenOperation(this.parentForm, `${this.basePath}customFields.${i}`, this);
        }
      }

      const group = array.at(i);
      group.get("customFieldId")?.valueChanges.subscribe(() => {
        this.onCustomFieldSelected(i);
      });
    }
  }

  private setupCustomFieldAutoOperationSelectionForIndex(index: number): void {
    const group = this.getCustomFieldsArray().at(index);
    const valueControl = group?.get("value");
    const operationControl = group?.get("operation");
    const selectedField = this.getSelectedCustomField(index);

    if (!valueControl || !operationControl || !selectedField) {
      return;
    }

    const type = this.getFilterType(selectedField);

    valueControl.valueChanges.subscribe(value => {
      const hasValue = this.hasFieldValue(value, type);

      if (hasValue) {
        if (!operationControl.value) {
          const operations = this.operationsPipe.transform(type, false);
          if (operations.length > 0) {
            operationControl.setValue(operations[0]);
          }
        }
      } else {
        operationControl.setValue(null);
      }
    });
  }

  protected readonly FilterOperation = FilterOperation;
}
