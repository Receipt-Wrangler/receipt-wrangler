import { Component, OnInit, ViewEncapsulation, computed, signal, viewChildren, viewChild } from "@angular/core";
import { FormArray, FormBuilder, FormGroup, Validators } from "@angular/forms";
import { MatDialogRef } from "@angular/material/dialog";
import { UntilDestroy, untilDestroyed } from "@ngneat/until-destroy";
import { Store } from "@ngxs/store";
import { catchError, of, take, tap } from "rxjs";
import { ReceiptFilterComponent } from "src/shared-ui/receipt-filter/receipt-filter.component";
import { FormOption } from "../../interfaces/form-option.interface";
import { BaseFormComponent } from "../../form/index";
import { Category, ChartGrouping, Dashboard, DashboardService, PagedRequestCommand, Permission, ReportService, ReportTemplate, SortDirection, Tag, Widget, WidgetType } from "../../open-api";
import { SnackbarService } from "../../services";
import { EditableListComponent } from "../../shared-ui/editable-list/editable-list.component";
import { AuthState, GroupState } from "../../store";
import { buildReceiptFilterForm } from "../../utils/receipt-filter";
import { chartGroupingOptions } from "../constants/chart-grouping-options";
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
  public readonly receiptFilterComponents = viewChildren(ReceiptFilterComponent);

  public readonly widgetList = viewChild.required(EditableListComponent);

  public headerText: string = "";

  public dashboard?: Dashboard;

  public isAddingWidget: boolean = false;

  public originalWidgets: Widget[] = [];

  private readonly selectedGroupId = this.store.selectSignal(
    GroupState.selectedGroupId
  );

  // Filter options come from the selected group's grant-filtered AppData catalog
  // (same source as the receipts table), not the admin-only global endpoints.
  // Derived so they refresh if the selected group changes while the dialog is open.
  public readonly categories = computed<Category[]>(() => {
    const groupId = Number(this.selectedGroupId());
    return Number.isNaN(groupId)
      ? []
      : this.store.selectSnapshot(AuthState.groupCategories(groupId));
  });

  public readonly tags = computed<Tag[]>(() => {
    const groupId = Number(this.selectedGroupId());
    return Number.isNaN(groupId)
      ? []
      : this.store.selectSnapshot(AuthState.groupTags(groupId));
  });

  // Only offer the Report widget type to users who can view report templates —
  // otherwise the template picker would be empty. Resolved through the AuthState
  // selector, exactly like the report builder / templates-list gate.
  private readonly canSelectReport = this.store.selectSignal(
    AuthState.hasAnyAppPermission([
      Permission.AppReportsRead,
      Permission.AppReportsReadAll,
    ])
  );

  public readonly availableWidgetTypeOptions = computed<FormOption[]>(() =>
    this.canSelectReport()
      ? widgetTypeOptions
      : widgetTypeOptions.filter((option) => option.value !== WidgetType.Report)
  );

  // The templates the user may view, as determined by the API (the list endpoint is
  // server-filtered to the caller's visible templates) — the Report widget's picker.
  public readonly reportTemplateOptions = signal<{ id: number; name: string }[]>([]);

  public get widgets(): FormArray {
    return this.form.get("widgets") as FormArray;
  }

  constructor(
    private dashboardService: DashboardService,
    private formBuilder: FormBuilder,
    private store: Store,
    private snackbarService: SnackbarService,
    private reportService: ReportService,
    private matDialogRef: MatDialogRef<DashboardFormComponent>
  ) {
    super();
  }


  public ngOnInit(): void {
    this.originalWidgets = Array.from(this.dashboard?.widgets ?? []);
    this.initForm();
    if (this.canSelectReport()) {
      this.loadReportTemplates();
    }
  }

  // Populate the Report widget's template picker from the server-filtered list, so
  // the choices are exactly the templates the user may view. A read failure yields
  // an empty picker rather than an error (the user can still add other widgets).
  private loadReportTemplates(): void {
    const command: PagedRequestCommand = {
      page: 1,
      pageSize: 200,
      orderBy: "name",
      sortDirection: SortDirection.Asc,
    };
    this.reportService
      .getReportTemplates(command)
      .pipe(
        take(1),
        catchError(() => of({ data: [], totalCount: 0 })),
        tap((paged) => {
          const templates = (paged.data ?? []) as unknown as ReportTemplate[];
          this.reportTemplateOptions.set(
            templates.map((template) => ({ id: template.id, name: template.name }))
          );
        }),
        untilDestroyed(this)
      )
      .subscribe();
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
      case WidgetType.PieChart:
        formGroup = this.formBuilder.group({
          name: [widget.name, Validators.required],
          widgetType: [widget.widgetType, Validators.required],
          configuration: this.buildPieChartConfigForm(widget.configuration),
        });
        break;
      case WidgetType.Report:
        formGroup = this.formBuilder.group({
          name: [widget.name, Validators.required],
          widgetType: [widget.widgetType, Validators.required],
          configuration: this.buildReportConfigForm(widget.configuration),
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
            formGroup.removeControl("configuration");
            formGroup.addControl("configuration", buildReceiptFilterForm({}, this));
          } else if (widgetType === WidgetType.PieChart) {
            formGroup.removeControl("configuration");
            formGroup.addControl("configuration", this.buildPieChartConfigForm({}));
          } else if (widgetType === WidgetType.Report) {
            formGroup.removeControl("configuration");
            formGroup.addControl("configuration", this.buildReportConfigForm({}));
          } else {
            formGroup.removeControl("configuration");
          }
        }),
      ).subscribe();

    return formGroup;
  }

  private buildPieChartConfigForm(config: any): FormGroup {
    return this.formBuilder.group({
      chartGrouping: [config?.chartGrouping ?? ChartGrouping.Categories, Validators.required],
      filter: buildReceiptFilterForm(config?.filter ?? {}, this),
    });
  }

  private buildReportConfigForm(config: any): FormGroup {
    return this.formBuilder.group({
      reportTemplateId: [config?.reportTemplateId ?? null, Validators.required],
    });
  }

  public submit(): void {
    const canSubmit = this.form.valid && this.widgetList().getCurrentRowOpen() === undefined;

    if (this.widgetList().getCurrentRowOpen() !== undefined) {
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
      this.widgetList().closeRow();
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
    this.widgetList().openLastRow(this.widgets.length - 1);
    this.isAddingWidget = true;
  }

  public cancelWidgetEdit(): void {
    if (this.isAddingWidget) {
      this.widgets.removeAt(this.widgets.length - 1);
      this.widgetList().closeRow();
      this.isAddingWidget = false;
    } else {
      const widget = this.originalWidgets[this.widgetList().getCurrentRowOpen() as number];

      const widgetList = this.widgetList();
      if (widget.widgetType === WidgetType.FilteredReceipts) {
        this.patchFilterConfig(widgetList.getCurrentRowOpen() as number);
      } else {
        this.widgets.at(widgetList.getCurrentRowOpen() as number).patchValue(widget);
      }
      widgetList.closeRow();
    }
  }

  public filterSubmitted(): void {
    if (this.isAddingWidget) {
      const widget = this.widgets.at(this.widgets.length - 1) as FormGroup;
      if (widget.valid) {
        const form = this.receiptFilterComponents().at(-1)!.parentForm;
        widget.get("configuration")?.patchValue(form.value);
        this.originalWidgets.push(widget.value);

        this.widgetList().closeRow();
        this.isAddingWidget = false;
      }
    } else {
      const widget = this.widgets.at(
        this.widgetList().getCurrentRowOpen() as number
      ) as FormGroup;

      if (widget.valid) {
        const form = this.receiptFilterComponents().at(0)!.parentForm;
        widget.get("configuration")?.patchValue(form.value);
        const widgetList = this.widgetList();
        this.originalWidgets.splice(
          widgetList.getCurrentRowOpen() as number,
          1,
          widget.value
        );

        widgetList.closeRow();
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
  protected readonly chartGroupingOptions = chartGroupingOptions;
}
