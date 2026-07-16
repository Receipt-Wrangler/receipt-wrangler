import {
  AfterViewInit,
  Component,
  DestroyRef,
  TemplateRef,
  inject,
  signal,
  viewChild,
} from "@angular/core";
import { takeUntilDestroyed } from "@angular/core/rxjs-interop";
import { MatDialog } from "@angular/material/dialog";
import { MatTableDataSource } from "@angular/material/table";
import { Router } from "@angular/router";
import { Store } from "@ngxs/store";
import { EMPTY, catchError, finalize, take, tap } from "rxjs";
import { DEFAULT_DIALOG_CONFIG, DEFAULT_HOST_CLASS } from "../../constants";
import { PagedTableInterface } from "../../interfaces/paged-table.interface";
import { Permission, ReportTemplate } from "../../open-api";
import { SnackbarService } from "../../services";
import { BaseTableService } from "../../services/base-table.service";
import { BaseTableComponent } from "../../shared-ui/base-table/base-table.component";
import { BreadcrumbItem } from "../../shared-ui/breadcrumb/breadcrumb-item.interface";
import { ConfirmationDialogComponent } from "../../shared-ui/confirmation-dialog/confirmation-dialog.component";
import { GroupState } from "../../store";
import {
  columnCount,
  detailSummary,
  formatChips,
  groupingSummary,
  scopeNames,
} from "../models/report-template-summary";
import { ReportRunnerService } from "../services/report-runner.service";
import { ReportTemplateTableService } from "./report-template-table.service";

/**
 * The Reports landing page: a paged table of saved report templates. Each row
 * opens in the builder (edit/open), and — gated by the matching app.reports.*
 * permission — generates, duplicates, or deletes. The paging/sort mechanics come
 * from BaseTableComponent + ReportTemplateTableService; the derived Scope/Grouping/
 * Detail/Formats columns are read out of each template's stored configuration.
 */
@Component({
  selector: "app-report-template-list",
  templateUrl: "./report-template-list.component.html",
  styleUrls: ["./report-template-list.component.scss"],
  host: DEFAULT_HOST_CLASS,
  providers: [{ provide: BaseTableService, useClass: ReportTemplateTableService }],
  standalone: false,
})
export class ReportTemplateListComponent extends BaseTableComponent<ReportTemplate> implements AfterViewInit {
  private readonly router = inject(Router);
  private readonly store = inject(Store);
  private readonly matDialog = inject(MatDialog);
  private readonly snackbar = inject(SnackbarService);
  private readonly runner = inject(ReportRunnerService);
  private readonly destroyRef = inject(DestroyRef);

  private readonly nameCell = viewChild.required<TemplateRef<any>>("nameCell");
  private readonly scopeCell = viewChild.required<TemplateRef<any>>("scopeCell");
  private readonly groupingCell = viewChild.required<TemplateRef<any>>("groupingCell");
  private readonly detailCell = viewChild.required<TemplateRef<any>>("detailCell");
  private readonly formatsCell = viewChild.required<TemplateRef<any>>("formatsCell");
  private readonly updatedCell = viewChild.required<TemplateRef<any>>("updatedCell");
  private readonly actionsCell = viewChild.required<TemplateRef<any>>("actionsCell");

  protected readonly Permission = Permission;

  // Guards the empty-state against a first-paint flash: it renders only once a load
  // has completed and returned nothing.
  public readonly loaded = signal(false);
  // The template currently generating (per-row spinner + single-flight guard).
  public readonly generatingId = signal<number | null>(null);

  private readonly groups = this.store.selectSignal(GroupState.groupsWithoutAll);

  public readonly crumbs: BreadcrumbItem[] = [{ label: "Reports" }];

  constructor(public override baseTableService: BaseTableService) {
    super(baseTableService);
  }

  public ngAfterViewInit(): void {
    this.setColumns();
    this.getTableData();
  }

  // Overrides the base fetch to also flip `loaded` (drives the empty state) on both
  // success and error, so a failed load doesn't hang on a blank page.
  public override getTableData(): void {
    this.baseTableService
      .getPagedData()
      .pipe(
        take(1),
        tap((pagedData) => {
          this.dataSource.set(new MatTableDataSource(pagedData.data as unknown as ReportTemplate[]));
          this.totalCount.set(pagedData.totalCount);
          this.loaded.set(true);
        }),
        catchError(() => {
          this.loaded.set(true);
          return EMPTY;
        }),
        takeUntilDestroyed(this.destroyRef),
      )
      .subscribe();
  }

  private setColumns(): void {
    this.columns = [
      { columnHeader: "Name", matColumnDef: "name", template: this.nameCell(), sortable: true },
      { columnHeader: "Scope", matColumnDef: "scope", template: this.scopeCell(), sortable: false },
      { columnHeader: "Grouping", matColumnDef: "grouping", template: this.groupingCell(), sortable: false },
      { columnHeader: "Detail", matColumnDef: "detail", template: this.detailCell(), sortable: false },
      { columnHeader: "Formats", matColumnDef: "formats", template: this.formatsCell(), sortable: false },
      { columnHeader: "Updated", matColumnDef: "updated_at", template: this.updatedCell(), sortable: true },
      { columnHeader: "", matColumnDef: "actions", template: this.actionsCell(), sortable: false },
    ];
    this.setInitialSortedColumn(
      this.baseTableService.getPagedRequestCommand() as PagedTableInterface,
      this.columns,
    );
    this.displayedColumns = ["name", "scope", "grouping", "detail", "formats", "updated_at", "actions"];
  }

  // ---- Derived display strings (read from each template's stored configuration) ----

  public columnCountFor(template: ReportTemplate): number {
    return columnCount(template.configuration);
  }

  public scopeSummary(template: ReportTemplate): string {
    const names = scopeNames(template.configuration, (id) => this.groupName(id));
    if (names.length === 0) {
      return "—";
    }
    if (names.length <= 2) {
      return names.join(", ");
    }
    return `${names.length} groups`;
  }

  public groupingSummaryFor(template: ReportTemplate): string {
    return groupingSummary(template.configuration);
  }

  public detailSummaryFor(template: ReportTemplate): string {
    return detailSummary(template.configuration);
  }

  public formatChipsFor(template: ReportTemplate): string[] {
    return formatChips(template.configuration);
  }

  private groupName(id: string): string | undefined {
    return this.groups().find((group) => group.id?.toString() === id)?.name;
  }

  // ---- Row + page actions ----

  public newReport(): void {
    this.router.navigate(["/reports/new"]);
  }

  public generate(template: ReportTemplate): void {
    if (this.generatingId() !== null) {
      return;
    }
    this.generatingId.set(template.id);
    this.runner
      .generateFromTemplate(template.configuration)
      .pipe(
        take(1),
        catchError(() => EMPTY), // the HTTP interceptor surfaces the failure toast
        finalize(() => this.generatingId.set(null)),
        takeUntilDestroyed(this.destroyRef),
      )
      .subscribe();
  }

  public duplicate(template: ReportTemplate): void {
    this.runner
      .duplicateTemplate(template.id)
      .pipe(
        take(1),
        tap(() => {
          this.snackbar.success("Template duplicated");
          this.getTableData();
        }),
        catchError(() => EMPTY),
        takeUntilDestroyed(this.destroyRef),
      )
      .subscribe();
  }

  public delete(template: ReportTemplate): void {
    const dialogRef = this.matDialog.open(ConfirmationDialogComponent, DEFAULT_DIALOG_CONFIG);
    dialogRef.componentInstance.headerText = "Delete Report Template";
    dialogRef.componentInstance.dialogContent = `Are you sure you want to delete "${template.name}"? This action is irreversible.`;

    dialogRef
      .afterClosed()
      .pipe(
        take(1),
        tap((confirmed) => {
          if (confirmed) {
            this.callDelete(template);
          }
        }),
        takeUntilDestroyed(this.destroyRef),
      )
      .subscribe();
  }

  private callDelete(template: ReportTemplate): void {
    this.runner
      .deleteTemplate(template.id)
      .pipe(
        take(1),
        tap(() => {
          this.snackbar.success("Template deleted");
          this.getTableData();
        }),
        catchError(() => EMPTY),
        takeUntilDestroyed(this.destroyRef),
      )
      .subscribe();
  }
}
