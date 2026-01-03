import { Component, OnInit, QueryList, ViewChild, ViewChildren, ViewEncapsulation } from "@angular/core";
import { FormArray, FormBuilder, FormGroup, Validators } from "@angular/forms";
import { MatDialogRef } from "@angular/material/dialog";
import { UntilDestroy, untilDestroyed } from "@ngneat/until-destroy";
import { Store } from "@ngxs/store";
import { take, tap } from "rxjs";
import { ReceiptFilterComponent } from "src/shared-ui/receipt-filter/receipt-filter.component";
import { BaseFormComponent } from "../../form/index";
import { Dashboard, DashboardService, Widget, WidgetType } from "../../open-api";
import { SnackbarService } from "../../services";
import { EditableListComponent } from "../../shared-ui/editable-list/editable-list.component";
import { GroupState } from "../../store";
import { buildReceiptFilterForm } from "../../utils/receipt-filter";
import { widgetTypeOptions } from "../constants/widget-options";

@UntilDestroy()
@Component({
    selector: "app-dashboard-form",
    templateUrl: "./dashboard-form.component.html",
    styleUrls: ["./dashboard-form.component.scss"],
    encapsulation: ViewEncapsulation.None,
    standalone: false
})
export class DashboardFormComponent extends BaseFormComponent implements OnInit {
  @ViewChildren(ReceiptFilterComponent)
  public receiptFilterComponents!: QueryList<ReceiptFilterComponent>;

  @ViewChild(EditableListComponent)
  public widgetList!: EditableListComponent;

  public headerText: string = "";

  public dashboard?: Dashboard;

  public isAddingWidget: boolean = false;

  public originalWidgets: Widget[] = [];

  public get widgets(): FormArray {
    return this.form.get("widgets") as FormArray;
  }

  constructor(
    private dashboardService: DashboardService,
    private formBuilder: FormBuilder,
    private store: Store,
    private snackbarService: SnackbarService,
    private matDialogRef: MatDialogRef<DashboardFormComponent>
  ) {
    super();
  }


  public ngOnInit(): void {
    this.originalWidgets = Array.from(this.dashboard?.widgets ?? []);
    this.initForm();
  }

  public initForm(): void {
    this.form = this.formBuilder.group({
      name: [this.dashboard?.name ?? "", Validators.required],
      groupId: [
        this.store.selectSnapshot(GroupState.selectedGroupId),
        Validators.required,
      ],
      widgets: this.formBuilder.array(
        this.dashboard?.widgets?.map((w) => this.buildWidgetFormGroup(w)) ?? []
      ),
    });
  }

  private buildWidgetFormGroup(widget: Widget): FormGroup {
    let formGroup: FormGroup;
    switch (widget.widgetType) {
      case WidgetType.FilteredReceipts:
        formGroup = this.formBuilder.group({
          name: [widget.name, Validators.required],
          widgetType: [widget.widgetType, Validators.required],
          configuration: buildReceiptFilterForm(widget.configuration, this),
        });
        break;
      default:
        formGroup = this.formBuilder.group({
          name: [widget.name, Validators.required],
          widgetType: [widget.widgetType, Validators.required],
        });
        break;
    }

    formGroup.get("widgetType")
      ?.valueChanges
      .pipe(
        untilDestroyed(this),
        tap((widgetType: WidgetType) => {
          if (widgetType === WidgetType.FilteredReceipts) {
            formGroup.addControl("configuration", buildReceiptFilterForm({}, this));
          } else {
            formGroup.removeControl("configuration");
          }
        }),
      ).subscribe();

    return formGroup;
  }

  public submit(): void {
    const canSubmit = this.form.valid && this.widgetList.getCurrentRowOpen() === undefined;

    if (this.widgetList.getCurrentRowOpen() !== undefined) {
      this.snackbarService.error(
        "Please finish editing the open widget before submitting"
      );
      return;
    }

    if (canSubmit && !this.dashboard) {
      this.dashboardService
        .createDashboard(this.form.value)
        .pipe(
          take(1),
          tap((dashboard) => {
            this.snackbarService.success("Dashboard successfully created");
            this.matDialogRef.close(dashboard);
          })
        )
        .subscribe();
    } else if (canSubmit && this.dashboard) {
      this.dashboardService
        .updateDashboard(this.dashboard.id, this.form.value)
        .pipe(
          take(1),
          tap((dashboard) => {
            this.snackbarService.success("Dashboard successfully updated");
            this.matDialogRef.close(dashboard);
          })
        )
        .subscribe();
    }
  }

  public submitWidget(index: number): void {
    const widgetFormGroup = (this.widgets.at(index) as FormGroup);
    const widget = widgetFormGroup.value;

    if (!widgetFormGroup.valid) {
      widgetFormGroup.markAllAsTouched();
      return;
    }

    if (widget["widgetType"] === WidgetType.FilteredReceipts) {
      this.filterSubmitted();
    } else {
      this.widgetList.closeRow();
    }
  }

  public cancelButtonClicked(): void {
    this.matDialogRef.close(undefined);
  }

  public addWidget(): void {
    const formGroup = this.buildWidgetFormGroup({
      name: "",
      widgetType: undefined,
    } as Widget);
    this.widgets.push(formGroup);
    this.widgetList.openLastRow();
    this.isAddingWidget = true;
  }

  public cancelWidgetEdit(): void {
    if (this.isAddingWidget) {
      this.widgets.removeAt(this.widgets.length - 1);
      this.widgetList.closeRow();
      this.isAddingWidget = false;
    } else {
      const widget = this.originalWidgets[this.widgetList.getCurrentRowOpen() as number];

      if (widget.widgetType === WidgetType.FilteredReceipts) {
        this.patchFilterConfig(this.widgetList.getCurrentRowOpen() as number);
      } else {
        this.widgets.at(this.widgetList.getCurrentRowOpen() as number).patchValue(widget);
      }
      this.widgetList.closeRow();
    }
  }

  public filterSubmitted(): void {
    if (this.isAddingWidget) {
      const widget = this.widgets.at(this.widgets.length - 1) as FormGroup;
      if (widget.valid) {
        const form = this.receiptFilterComponents.last.parentForm;
        widget.get("configuration")?.patchValue(form.value);
        this.originalWidgets.push(widget.value);

        this.widgetList.closeRow();
        this.isAddingWidget = false;
      }
    } else {
      const widget = this.widgets.at(
        this.widgetList.getCurrentRowOpen() as number
      ) as FormGroup;

      if (widget.valid) {
        const form = this.receiptFilterComponents.first.parentForm;
        widget.get("configuration")?.patchValue(form.value);
        this.originalWidgets.splice(
          this.widgetList.getCurrentRowOpen() as number,
          1,
          widget.value
        );

        this.widgetList.closeRow();
      }
    }
  }

  private patchFilterConfig(index: number): void {
    if (this.widgets.at(index)) {
      const originalWidget =
        this.originalWidgets[index];

      (this.widgets.at(index) as FormGroup).removeControl("configuration");
      (this.widgets.at(index) as FormGroup).addControl(
        "configuration",
        buildReceiptFilterForm(originalWidget.configuration, this)
      );
    }
  }

  public removeWidget(index: number): void {
    this.widgets.removeAt(index);
    this.originalWidgets.splice(index, 1);
  }

  protected readonly WidgetType = WidgetType;
  protected readonly widgetTypeOptions = widgetTypeOptions;
}
