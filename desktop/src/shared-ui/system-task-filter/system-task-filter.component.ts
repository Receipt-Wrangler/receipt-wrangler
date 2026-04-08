import { Component, Input, OnInit, output } from "@angular/core";
import { FormControl, FormGroup } from "@angular/forms";
import { MatDialogRef } from "@angular/material/dialog";
import { UntilDestroy, untilDestroyed } from "@ngneat/until-destroy";
import { Store } from "@ngxs/store";
import { endOfDay, startOfMonth } from "date-fns";
import { take, tap } from "rxjs";
import { FormCommand } from "../../form/index";
import { FilterOperation, SystemTaskStatus, SystemTaskType } from "../../open-api";
import { SetSystemTaskFilter } from "../../store/system-task-table.state.actions";
import { OperationsPipe } from "../receipt-filter/operations.pipe";
import { SystemTaskTypePipe } from "../task-table/system-task-type.pipe";

@UntilDestroy()
@Component({
  selector: "app-system-task-filter",
  templateUrl: "./system-task-filter.component.html",
  styleUrls: ["./system-task-filter.component.scss"],
  standalone: false
})
export class SystemTaskFilterComponent implements OnInit {
  @Input() public headerText: string = "";

  @Input() public parentForm: FormGroup = new FormGroup({});

  public readonly formCommand = output<FormCommand>();

  public startOfMonthFormControl = new FormControl(startOfMonth(new Date()));
  public endOfTodayFormControl = new FormControl(endOfDay(new Date()));

  public systemTaskTypeOptions: { value: string; displayValue: string }[] = [];

  public systemTaskStatusOptions: { value: string; displayValue: string }[] = [
    { value: SystemTaskStatus.Succeeded, displayValue: "Succeeded" },
    { value: SystemTaskStatus.Failed, displayValue: "Failed" },
  ];

  private operationsPipe = new OperationsPipe();
  private systemTaskTypePipe = new SystemTaskTypePipe();

  constructor(
    private store: Store,
    private dialogRef: MatDialogRef<SystemTaskFilterComponent>,
  ) {}

  public ngOnInit(): void {
    this.startOfMonthFormControl.disable();
    this.endOfTodayFormControl.disable();
    this.buildTypeOptions();
    this.setupAutoOperationSelection();
  }

  private buildTypeOptions(): void {
    // Only show types that appear as parent tasks (exclude child-only types)
    const visibleTypes: SystemTaskType[] = [
      SystemTaskType.MagicFill,
      SystemTaskType.QuickScan,
      SystemTaskType.EmailUpload,
      SystemTaskType.EmailRead,
      SystemTaskType.SystemEmailConnectivityCheck,
      SystemTaskType.ReceiptProcessingSettingsConnectivityCheck,
      SystemTaskType.PromptGenerated,
      SystemTaskType.ReceiptUpdated,
      SystemTaskType.ApiKeyDeleted,
    ];

    this.systemTaskTypeOptions = visibleTypes.map((type) => ({
      value: type,
      displayValue: this.systemTaskTypePipe.transform(type) || type,
    }));
  }

  public resetFilter(): void {
    this.formCommand.emit({
      path: "",
      command: "reset",
    });
    this.formCommand.emit({
      path: "type.value",
      command: "clear",
    });
    this.formCommand.emit({
      path: "status.value",
      command: "clear",
    });
    this.formCommand.emit({
      path: "ranByUserId.value",
      command: "clear",
    });
  }

  public submitButtonClicked(): void {
    const filter = this.parentForm.value;

    if (this.parentForm.valid) {
      this.store
        .dispatch(new SetSystemTaskFilter(filter))
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
    const fieldsToWatch = [
      { fieldName: "type", type: "list" },
      { fieldName: "status", type: "list" },
      { fieldName: "ranByUserId", type: "users" },
      { fieldName: "startedAt", type: "date" },
      { fieldName: "endedAt", type: "date" },
    ];

    fieldsToWatch.forEach(({ fieldName, type }) => {
      const valueControl = this.parentForm.get(`${fieldName}.value`);
      const operationControl = this.parentForm.get(`${fieldName}.operation`);

      if (valueControl && operationControl) {
        valueControl.valueChanges.pipe(
          untilDestroyed(this)
        ).subscribe(value => {
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
