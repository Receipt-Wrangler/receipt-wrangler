import { CurrencyPipe, DatePipe } from "@angular/common";
import { AfterViewInit, Component, computed, OnInit, signal, TemplateRef, ViewEncapsulation, viewChild } from "@angular/core";
import { MatDialog } from "@angular/material/dialog";
import { PageEvent } from "@angular/material/paginator";
import { Sort } from "@angular/material/sort";
import { MatTableDataSource } from "@angular/material/table";
import { ActivatedRoute, Router } from "@angular/router";
import { UntilDestroy, untilDestroyed } from "@ngneat/until-destroy";
import { Store } from "@ngxs/store";
import { catchError, EMPTY, map, Subject, switchMap, take, tap } from "rxjs";
import { fadeInOut } from "src/animations";
import { ReceiptFilterService } from "src/services/receipt-filter.service";
import { ConfirmationDialogComponent } from "src/shared-ui/confirmation-dialog/confirmation-dialog.component";
import { ResetReceiptFilter, SetColumnConfig, SetPage, SetPageSize, SetQuickDateField, SetReceiptFilterData, SetReceiptFilterField, SetSummaryConfigGroupId, } from "src/store/receipt-table.actions";
import { DEFAULT_RECEIPT_ORDER_BY, DEFAULT_RECEIPT_SORT_DIRECTION, ReceiptTableState, } from "src/store/receipt-table.state";
import { TableColumn } from "src/table/table-column.interface";
import { TableComponent } from "src/table/table/table.component";
import { DEFAULT_DIALOG_CONFIG, DEFAULT_HOST_CLASS, RECEIPT_DATE_FILTER_FIELDS, ReceiptDateFilterFieldKey } from "../../constants";
import { ReceiptTableColumnConfig } from "../../interfaces";
import {
  BulkStatusUpdateCommand,
  Category,
  CustomField,
  FilterOperation,
  Group,
  GroupsService,
  PagedDataDataInner,
  Permission,
  Receipt,
  ReceiptPagedRequestFilter,
  ReceiptService,
  ReceiptStatus,
  ReceiptSummary,
  Tag,
} from "../../open-api";
import { SnackbarService } from "../../services";
import { ReceiptExportService } from "../../services/receipt-export.service";
import { ReceiptFilterComponent } from "../../shared-ui/receipt-filter/receipt-filter.component";
import { canReadCustomFieldCatalog } from "../../resolvers/custom-field.resolver";
import { AuthState, GroupState, UserState } from "../../store";
import { CustomCurrencyPipe } from "../../pipes/custom-currency.pipe";
import {
  applyFormCommand,
  customFieldColumnDef,
  mergeCustomFieldColumns,
  RECEIPT_COLUMN_DISPLAY_NAMES,
} from "../../utils/index";
import { FilterMonth, monthFilterEntry, monthFromFilterEntry } from "../../utils/receipt-date-filter";
import { buildReceiptFilterForm } from "../../utils/receipt-filter";
import { buildReceiptFilterChips, ReceiptFilterChip } from "../../utils/receipt-filter-chips";
import { isFilterEntryActive } from "../../utils/receipt-filter-entry";
import { resolveSummaryConfigGroup, summaryConfigGroups } from "../../utils/receipt-summary";
import { openQuickScanDialog } from "../quick-scan-dialog/open-quick-scan-dialog";
import { BulkStatusUpdateComponent } from "../bulk-resolve-dialog/bulk-status-update-dialog.component";
import { ColumnConfigurationDialogComponent } from "../column-configuration-dialog/column-configuration-dialog.component";

/**
 * A receipts-table column, carrying the custom field definition when the column
 * is one. The shared table hands the whole column to a cell template, which is
 * how a single template can render every custom field.
 */
interface ReceiptTableColumn extends TableColumn {
  customField?: CustomField;
}

@UntilDestroy()
@Component({
  selector: "app-receipts-table",
  templateUrl: "./receipts-table.component.html",
  styleUrls: ["./receipts-table.component.scss"],
  animations: [fadeInOut],
  encapsulation: ViewEncapsulation.None,
  host: DEFAULT_HOST_CLASS,
  // The chip labels are built in TS rather than the template, so the formatting
  // pipes are injected. CustomCurrencyPipe is declared in PipesModule and needs
  // CurrencyPipe, neither of which is providedIn: "root".
  providers: [CurrencyPipe, CustomCurrencyPipe, DatePipe],
  standalone: false
})
export class ReceiptsTableComponent implements OnInit, AfterViewInit {
  constructor(
    private activatedRoute: ActivatedRoute,
    private groupsService: GroupsService,
    private matDialog: MatDialog,
    private receiptExportService: ReceiptExportService,
    private receiptFilterService: ReceiptFilterService,
    private receiptService: ReceiptService,
    private router: Router,
    private snackbarService: SnackbarService,
    private store: Store,
    private customCurrencyPipe: CustomCurrencyPipe,
    private datePipe: DatePipe,
  ) {
    this.listenForRefreshRequests();
    this.listenForSummaryRequests();
  }

  readonly createdAtCell = viewChild.required<TemplateRef<any>>("createdAtCell");

  readonly dateCell = viewChild.required<TemplateRef<any>>("dateCell");

  readonly nameCell = viewChild.required<TemplateRef<any>>("nameCell");

  readonly paidByCell = viewChild.required<TemplateRef<any>>("paidByCell");

  readonly amountCell = viewChild.required<TemplateRef<any>>("amountCell");

  readonly categoryCell = viewChild.required<TemplateRef<any>>("categoryCell");

  readonly tagCell = viewChild.required<TemplateRef<any>>("tagCell");

  readonly statusCell = viewChild.required<TemplateRef<any>>("statusCell");

  readonly resolvedDateCell = viewChild.required<TemplateRef<any>>("resolvedDateCell");

  readonly firstCommentCell = viewChild.required<TemplateRef<any>>("firstCommentCell");

  readonly customFieldCell = viewChild.required<TemplateRef<any>>("customFieldCell");

  readonly actionsCell = viewChild.required<TemplateRef<any>>("actionsCell");

  readonly table = viewChild.required(TableComponent);

  public page = this.store.selectSignal(ReceiptTableState.page);

  public pageSize = this.store.selectSignal(ReceiptTableState.pageSize);

  public filter = this.store.selectSignal(ReceiptTableState.filterData);

  public columnConfig = this.store.selectSignal(ReceiptTableState.columnConfig);

  public selectedGroupId = this.store.selectSignal(GroupState.selectedGroupId);

  private readonly refreshRequested = new Subject<void>();

  /**
   * The summary rides its own stream rather than the table's, for two reasons: a
   * forkJoin would make the table wait on the slower unpaged aggregate on the
   * app's hottest screen, and the totals do not change when only the page or the
   * sort does — see getFilteredReceiptsPage().
   *
   * It mirrors refreshRequested exactly otherwise, switchMap and inner catchError
   * included, so a newer request supersedes an in-flight one and an error cannot
   * kill later refreshes.
   */
  private readonly summaryRequested = new Subject<void>();

  public readonly summary = signal<ReceiptSummary | undefined>(undefined);

  private groups = this.store.selectSignal(GroupState.groups);

  private persistedSummaryConfigGroupId = this.store.selectSignal(
    ReceiptTableState.summaryConfigGroupId
  );

  /**
   * The groups whose summary configuration the viewer may choose between — only
   * ever more than one on the synthetic "All" group, which spans several groups
   * and has no configuration of its own.
   */
  public readonly summaryConfigGroupOptions = computed(() => {
    const group = this.groups().find((candidate) => candidate.id === Number(this.groupId));
    if (!group?.isAllGroup) {
      return [];
    }

    return summaryConfigGroups(this.groups());
  });

  /**
   * Whose configuration is actually in use. On a real group it is that group; on
   * "All" it is the viewer's persisted pick, falling back to the first enabled
   * group when that pick is stale (the group turned its summary off, the user
   * left it, or the state predates the key).
   */
  public readonly summaryConfigGroupId = computed(() => {
    const options = this.summaryConfigGroupOptions();
    if (options.length === 0) {
      const group = this.groups().find((candidate) => candidate.id === Number(this.groupId));
      return group?.isAllGroup ? undefined : group?.id;
    }

    return resolveSummaryConfigGroup(options, this.persistedSummaryConfigGroupId())?.id;
  });

  private users = this.store.selectSignal(UserState.users);

  private receiptFilter = computed(
    () => this.filter()?.filter as ReceiptPagedRequestFilter | undefined
  );

  /** The date field the quick date control writes to. */
  public quickDateField = this.store.selectSignal(ReceiptTableState.quickDateField);

  public readonly dateFilterFields = RECEIPT_DATE_FILTER_FIELDS;

  public quickDateFieldLabel = computed(
    () =>
      RECEIPT_DATE_FILTER_FIELDS.find((field) => field.key === this.quickDateField())?.label ??
      "Receipt Date"
  );

  /** The entry the quick date control currently owns. */
  private quickDateEntry = computed(() => this.receiptFilter()?.[this.quickDateField()]);

  /**
   * The month the quick date control is showing, or null when its field is
   * unset or holds something a month cannot express.
   */
  public stepperMonth = computed(() => monthFromFilterEntry(this.quickDateEntry()));

  public stepperLabel = computed(() => {
    const month = this.stepperMonth();
    if (month) {
      return this.datePipe.transform(new Date(month.year, month.month, 1), "LLLL y") ?? "";
    }

    // A filter the stepper cannot describe still has to be visible as a
    // filter — the chip below spells out what it actually is.
    return isFilterEntryActive(this.quickDateEntry()) ? "Custom" : "All time";
  });

  public filterChips = computed<ReceiptFilterChip[]>(() => {
    const categories = this.categories();
    const tags = this.tags();
    const groups = this.groups();
    const users = this.users();

    return buildReceiptFilterChips(
      this.receiptFilter(),
      {
        categories,
        tags,
        // Not groupsWithoutAll: a group filter set on the All Groups view is
        // persisted, so it can outlive the view that offers the control and
        // still has to name itself here.
        groups,
        users,
        formatDate: (value) => this.datePipe.transform(value as string) ?? "",
        formatCurrency: (value) => this.customCurrencyPipe.transform(value as number),
      }
    );
  });

  private numFiltersAppliedRaw = this.store.selectSignal(ReceiptTableState.numFiltersApplied);

  public numFiltersApplied = computed(() => {
    const num = this.numFiltersAppliedRaw();
    return num > 0 ? num : undefined;
  });

  public customFields = signal<CustomField[]>([]);

  /** Whether `customFields` is the real catalog rather than a permission stub. */
  public customFieldsAvailable = signal<boolean>(true);

  public categories = signal<Category[]>([]);

  public tags = signal<Tag[]>([]);

  public groupId: string = "0";

  public dataSource = signal(new MatTableDataSource<PagedDataDataInner>([]));

  public displayedColumns = signal<string[]>([]);

  public columns = signal<TableColumn[]>([]);

  public totalCount = signal(0);

  public selectedReceiptIds = signal<number[]>([]);

  public firstSort: boolean = true;

  public canEdit: boolean = false;

  public canCreate: boolean = false;

  public canQuickScan: boolean = false;

  public canPollEmail: boolean = false;

  public headerText: string = "";

  public group?: Group;

  protected readonly Permission = Permission;

  public ngOnInit(): void {
    this.groupId = this.store
      .selectSnapshot(GroupState.selectedGroupId)
      ?.toString();
    this.setGroup();
    this.setCanEdit();

    this.setHeaderText();

    // Filter options come from the selected group's AppData catalog (filtered to
    // the user's grants), so a restricted user can't filter by a hidden one.
    const numericGroupId = Number(this.groupId);
    this.categories.set(
      Number.isNaN(numericGroupId)
        ? []
        : this.store.selectSnapshot(AuthState.groupCategories(numericGroupId))
    );
    this.tags.set(
      Number.isNaN(numericGroupId)
        ? []
        : this.store.selectSnapshot(AuthState.groupTags(numericGroupId))
    );

    // Resolved per route, and empty for a user without app.custom-fields.read -
    // which is the permission gate: no catalog, no custom field columns.
    this.customFields.set(
      this.activatedRoute.snapshot.data["customFields"] ?? []
    );

    // ...but that empty list means "not permitted to look", not "none exist", and
    // the two must not be confused when reconciling - see reconcileColumnConfig.
    this.customFieldsAvailable.set(canReadCustomFieldCatalog(this.store));
    this.reconcileColumnConfig();

    this.getInitialData();
    // getInitialData stays on its own one-shot subscription because it owns the single
    // setColumns() call; the summary has no such constraint and just rides its stream.
    this.summaryRequested.next();
  }

  /**
   * Brings the persisted column configuration back in line with the custom
   * fields that exist now, before anything reads it.
   *
   * Both halves matter on the *first* load. The configuration and the sort are
   * persisted to localStorage and are shared by every account using the browser,
   * so the table can start up holding a column for a deleted custom field, or one
   * only an administrator can see. A stale column would make mat-table render an
   * id it has no definition for; a stale sort would make the very first request
   * ask the API to order by a column that no longer exists.
   */
  private reconcileColumnConfig(): void {
    const persisted = this.store.selectSnapshot(
      ReceiptTableState.columnConfig
    );
    const reconciled = mergeCustomFieldColumns(
      persisted,
      this.customFields(),
      this.customFieldsAvailable()
    );

    const changed =
      reconciled.length !== persisted.length ||
      reconciled.some(
        (column, index) =>
          column.matColumnDef !== persisted[index].matColumnDef ||
          column.visible !== persisted[index].visible ||
          column.order !== persisted[index].order
      );

    if (changed) {
      this.store.dispatch(new SetColumnConfig(reconciled));
    }

    const filterData = this.store.selectSnapshot(ReceiptTableState.filterData);
    const isKnownColumn = reconciled.some(
      (column) => column.matColumnDef === filterData.orderBy
    );

    if (!isKnownColumn) {
      this.store.dispatch(
        new SetReceiptFilterData({
          ...filterData,
          orderBy: DEFAULT_RECEIPT_ORDER_BY,
          sortDirection: DEFAULT_RECEIPT_SORT_DIRECTION,
        })
      );
    }
  }

  private setGroup(): void {
    this.group = this.store.selectSnapshot(GroupState.getGroupById(this.groupId));
  }

  private getInitialData(): void {
    this.receiptFilterService
      .getPagedReceiptsForGroups(this.groupId)
      .pipe(
        take(1),
        tap((pagedData) => {
          this.dataSource.set(new MatTableDataSource<PagedDataDataInner>(pagedData.data));
          this.totalCount.set(pagedData.totalCount);
          this.setColumns();
        })
      )
      .subscribe();
  }

  private setCanEdit(): void {
    const groupId = Number.parseInt(this.groupId);
    this.canEdit = this.store.selectSnapshot(
      AuthState.hasGroupPermission(groupId, Permission.GroupReceiptsUpdate)
    );
    this.canCreate = this.store.selectSnapshot(
      AuthState.hasGroupPermission(groupId, Permission.GroupReceiptsCreate)
    );
    this.canQuickScan = this.store.selectSnapshot(
      AuthState.hasGroupPermission(groupId, Permission.GroupReceiptsQuickScan)
    );
    this.canPollEmail = this.store.selectSnapshot(
      AuthState.hasGroupPermission(groupId, Permission.GroupEmailPoll)
    );
  }

  private setHeaderText(): void {
    const group = this.store.selectSnapshot(
      GroupState.getGroupById(this.groupId)
    );
    if (group) {
      if (group.name.toLowerCase().includes("receipt")) {
        this.headerText = group.name;
      } else {
        this.headerText = `${group.name} Receipts`;
      }
    }
  }

  public ngAfterViewInit(): void {
    this.setSelectedReceiptIdsObservable();
  }

  private setSelectedReceiptIdsObservable(): void {
    this.table()?.selection?.changed
      .pipe(
        untilDestroyed(this),
        map((event) => (event.source.selected as Receipt[]).map((r) => r.id)),
        tap((ids) => this.selectedReceiptIds.set(ids))
      )
      .subscribe();
  }

  private setColumns(): void {
    const currentColumnConfig = this.store.selectSnapshot(ReceiptTableState.columnConfig);

    const allColumns = [
      {
        columnHeader: "Added At",
        matColumnDef: "created_at",
        template: this.createdAtCell(),
        sortable: true,
      },
      {
        columnHeader: "Receipt Date",
        matColumnDef: "date",
        template: this.dateCell(),
        sortable: true,
      },
      {
        columnHeader: "Name",
        matColumnDef: "name",
        template: this.nameCell(),
        sortable: true,
      },
      {
        columnHeader: "Paid By",
        matColumnDef: "paid_by_user_id",
        template: this.paidByCell(),
        sortable: true,
      },
      {
        columnHeader: "Amount",
        matColumnDef: "amount",
        template: this.amountCell(),
        sortable: true,
      },
      {
        columnHeader: "Categories",
        matColumnDef: "categories",
        template: this.categoryCell(),
        sortable: false,
      },
      {
        columnHeader: "Tags",
        matColumnDef: "tags",
        template: this.tagCell(),
        sortable: false,
      },
      {
        columnHeader: "Status",
        matColumnDef: "status",
        template: this.statusCell(),
        sortable: true,
      },
      {
        columnHeader: "Resolved Date",
        matColumnDef: "resolved_date",
        template: this.resolvedDateCell(),
        sortable: true,
      },
      {
        // The receipt's first comment, which the API both returns
        // (firstComment) and sorts on under this same key.
        columnHeader: RECEIPT_COLUMN_DISPLAY_NAMES["first_comment"],
        matColumnDef: "first_comment",
        template: this.firstCommentCell(),
        sortable: true,
      },
    ] as ReceiptTableColumn[];

    // One column per custom field, after the built-in ones. They share a single
    // cell template, which reads the definition off the column it is rendering.
    allColumns.push(
      ...this.customFields().map((customField) => ({
        columnHeader: customField.name,
        matColumnDef: customFieldColumnDef(customField.id),
        template: this.customFieldCell(),
        sortable: true,
        customField,
      }))
    );

    // Filter and order columns based on configuration
    const visibleColumnConfigs = currentColumnConfig
      .filter(config => config.visible)
      .sort((a, b) => a.order - b.order);

    const columns = visibleColumnConfigs
      .map(config => allColumns.find(col => col.matColumnDef === config.matColumnDef))
      .filter(col => col !== undefined) as ReceiptTableColumn[];

    if (this.canEdit) {
      columns.push({
        columnHeader: "Actions",
        matColumnDef: "actions",
        template: this.actionsCell(),
        sortable: false,
      });
    }

    // Derived from the columns that actually resolved, never from the stored
    // configuration: mat-table throws on a displayed id it has no definition for,
    // and a configuration naming a since-deleted custom field is ordinary.
    const displayColumns = ["select", ...columns.map(col => col.matColumnDef)];

    const filter = this.store.selectSnapshot(ReceiptTableState.filterData);
    const orderByIndex = columns.findIndex(
      (c) => c.matColumnDef === filter.orderBy
    );

    if (orderByIndex >= 0) {
      columns[orderByIndex].defaultSortDirection = filter.sortDirection;
    } else if (columns.length > 0) {
      columns[0].defaultSortDirection = "desc";
    }

    this.columns.set(columns);
    this.displayedColumns.set(displayColumns);
  }

  public sort(sortState: Sort): void {
    if (!this.firstSort) {
      const filterData = this.store.selectSnapshot(
        ReceiptTableState.filterData
      );

      this.store.dispatch(
        new SetReceiptFilterData({
          page: filterData.page,
          pageSize: filterData.pageSize,
          orderBy: sortState.active,
          sortDirection: sortState.direction,
          filter: filterData.filter,
        })
      );

      this.getFilteredReceiptsPage();
    }
    this.firstSort = false;
  }

  public filterButtonClicked(): void {
    const filter = this.store.selectSnapshot(ReceiptTableState.filterData).filter as any;

    const dialogRef = this.matDialog.open(ReceiptFilterComponent, {
      minWidth: "75%",
      maxWidth: "100%",
    });

    dialogRef.componentInstance.categories = this.categories();
    dialogRef.componentInstance.tags = this.tags();
    dialogRef.componentInstance.parentForm = buildReceiptFilterForm(filter, this);
    dialogRef.componentInstance.headerText = "Filter Receipts";
    // The group filter is only meaningful on the "All groups" view; a
    // single-group view is already scoped to one group.
    dialogRef.componentInstance.showGroupFilter = this.group?.isAllGroup ?? false;
    const formCommandSubscription = dialogRef.componentInstance.formCommand.subscribe((formCommand) => {
      applyFormCommand(dialogRef.componentInstance.parentForm, formCommand);
    });

    dialogRef
      .afterClosed()
      .pipe(
        take(1),
        tap((applyFilter) => {
          if (applyFilter) {
            this.store.dispatch(new SetPage(1));
            this.getFilteredReceipts();
          }

          formCommandSubscription.unsubscribe();
        })
      )
      .subscribe();
  }

  public monthSelected(month: FilterMonth): void {
    // The quick control IS its field's filter, so it overwrites whatever was there.
    this.applyFilterField(this.quickDateField(), monthFilterEntry(month) as any);
  }

  public allTimeSelected(): void {
    this.applyFilterField(this.quickDateField(), null);
  }

  /**
   * Re-points the quick date control at another date field. Deliberately
   * non-destructive: it changes no condition — only which one the stepper
   * describes — so whatever the previous field held stays applied and keeps its
   * chip. Nothing to refetch.
   */
  public quickDateFieldSelected(field: ReceiptDateFilterFieldKey): void {
    this.store.dispatch(new SetQuickDateField(field));
  }

  public filterChipCleared(field: keyof ReceiptPagedRequestFilter): void {
    this.applyFilterField(field, null);
  }

  /**
   * The one write path for a single-field filter change, so the page reset and
   * the refetch can never be forgotten — narrowing a filter while on page 7
   * would otherwise land on an empty page.
   */
  private applyFilterField(
    field: keyof ReceiptPagedRequestFilter,
    entry: { operation: FilterOperation | null; value: unknown } | null
  ): void {
    this.store.dispatch(new SetReceiptFilterField(field, entry));
    this.store.dispatch(new SetPage(1));
    this.getFilteredReceipts();
  }

  public quickScanClicked(): void {
    openQuickScanDialog(this.matDialog)
      .pipe(
        take(1),
        tap(() => this.getFilteredReceipts())
      )
      .subscribe();
  }

  public exportAllReceipts(): void {
    this.receiptExportService.exportReceiptsFromFilter(this.groupId, this.filter());
  }

  public resetFilterButtonClicked(): void {
    this.store.dispatch(new ResetReceiptFilter());
    this.getFilteredReceipts();
  }

  public configureColumnsButtonClicked(): void {
    const currentColumnConfig = this.store.selectSnapshot(ReceiptTableState.columnConfig);

    const dialogRef = this.matDialog.open(ColumnConfigurationDialogComponent, {
      ...DEFAULT_DIALOG_CONFIG,
      data: {
        currentColumns: currentColumnConfig,
        customFields: this.customFields(),
        customFieldsAvailable: this.customFieldsAvailable(),
      }
    });

    dialogRef
      .afterClosed()
      .pipe(
        take(1),
        tap((result: ReceiptTableColumnConfig[] | null) => {
          if (result) {
            this.store.dispatch(new SetColumnConfig(result));
            this.setColumns();
          }
        })
      )
      .subscribe();
  }

  /**
   * The filter changed, or receipts were mutated: refresh the page AND the totals.
   *
   * This keeps its name and every existing caller, so the safe behaviour is the
   * default — a new call site that forgets the distinction over-refreshes rather
   * than leaving stale figures on screen.
   */
  public getFilteredReceipts(): void {
    this.refreshRequested.next();
    this.summaryRequested.next();
  }

  /**
   * The page or the sort changed. Deliberately does NOT refresh the totals:
   * neither changes which receipts the filter matches, so the figures are already
   * right, and the summary aggregates the whole result set unpaged — refetching it
   * on every page click would be the most expensive no-op in the app.
   */
  private getFilteredReceiptsPage(): void {
    this.refreshRequested.next();
  }

  /**
   * Every refresh goes through one switchMap, so a newer request supersedes
   * whatever is in flight (and aborts its XHR) instead of racing it. Without
   * this the last *response* wins rather than the last *request* — and the
   * quick date arrows put those one click apart, so a slow earlier page could
   * repaint the table with a month the user has already stepped past.
   */
  private listenForRefreshRequests(): void {
    this.refreshRequested
      .pipe(
        untilDestroyed(this),
        switchMap(() =>
          this.receiptFilterService.getPagedReceiptsForGroups(this.groupId.toString()).pipe(
            // Keeps the outer subscription alive. An error surfacing through
            // switchMap would complete it and silently kill every later
            // refresh; the HTTP interceptor already reports the failure.
            catchError(() => EMPTY)
          )
        ),
        tap((pagedData) => {
          this.dataSource.set(new MatTableDataSource(pagedData.data));
          this.totalCount.set(pagedData.totalCount);
        })
      )
      .subscribe();
  }

  /**
   * Mirrors listenForRefreshRequests. The guard sits inside the switchMap rather
   * than at the call sites so a group with no summary configured costs no request
   * at all — which is what keeps this free for every install that has not opted in.
   */
  private listenForSummaryRequests(): void {
    this.summaryRequested
      .pipe(
        untilDestroyed(this),
        switchMap(() => {
          const configGroupId = this.summaryConfigGroupId();
          if (!configGroupId) {
            this.summary.set(undefined);
            return EMPTY;
          }

          return this.receiptFilterService
            .getReceiptSummaryForGroup(this.groupId.toString(), configGroupId)
            .pipe(
              // Same reasoning as the table's: an error through switchMap would
              // complete the outer subscription and kill every later refresh. The
              // last good figures stay on screen and the interceptor reports it.
              catchError(() => EMPTY)
            );
        }),
        tap((summary) => this.summary.set(summary))
      )
      .subscribe();
  }

  /**
   * Only the breakdown's shape changes, so the table is untouched — this must not
   * go through getFilteredReceipts().
   */
  public summaryConfigGroupSelected(groupId: number): void {
    this.store.dispatch(new SetSummaryConfigGroupId(groupId));
    this.summaryRequested.next();
  }

  public deleteReceipt(row: Receipt): void {
    const dialogRef = this.matDialog.open(ConfirmationDialogComponent);

    dialogRef.componentInstance.headerText = "Delete Receipt";
    dialogRef.componentInstance.dialogContent = `Are you sure you would like to delete the receipt ${row.name}? This action is irreversible.`;

    dialogRef
      .afterClosed()
      .pipe(
        take(1),
        tap((r) => {
          if (r) {
            this.receiptService
              .deleteReceiptById(row.id as number)
              .pipe(
                take(1),
                tap(() => {
                  this.dataSource.update(ds => new MatTableDataSource(ds.data.filter(
                    (r) => r.id !== row.id
                  )));
                  // The table row is patched out in place rather than refetched, so the
                  // totals have to be told separately or they keep counting the receipt.
                  this.summaryRequested.next();
                  this.snackbarService.success("Receipt successfully deleted");
                })
              )
              .subscribe();
          }
        })
      )
      .subscribe();
  }

  public duplicateReceipt(row: Receipt): void {
    const dialogRef = this.matDialog.open(ConfirmationDialogComponent);

    dialogRef.componentInstance.headerText = "Duplicate Receipt";
    dialogRef.componentInstance.dialogContent = `Are you sure you would like to duplicate the receipt ${row.name}?`;

    dialogRef
      .afterClosed()
      .pipe(
        take(1),
        tap((confirmed) => {
          if (confirmed) {
            this.receiptService
              .duplicateReceipt(row.id)
              .pipe(
                take(1),
                tap((r: Receipt) => {
                  this.snackbarService.success("Receipt successfully duplicated");
                  this.router.navigateByUrl(`/receipts/${r.id}/view`);
                })
              )
              .subscribe();
          }
        })
      )
      .subscribe();
  }

  public updatePageData(pageEvent: PageEvent): void {
    const newPage = pageEvent.pageIndex + 1;
    this.store.dispatch(new SetPage(newPage));
    this.store.dispatch(new SetPageSize(pageEvent.pageSize));

    this.getFilteredReceiptsPage();
  }

  public showStatusUpdateDialog(): void {
    const ref = this.matDialog.open(
      BulkStatusUpdateComponent,
      DEFAULT_DIALOG_CONFIG
    );

    ref
      .afterClosed()
      .pipe(
        take(1),
        tap(
          (
            commentForm:
              | {
              comment: string;
              status: ReceiptStatus;
            }
              | undefined
          ) => {
            const table = this.table();
            if (table.selection.hasValue() && commentForm) {
              const receiptIds = (
                table.selection.selected as Receipt[]
              ).map((r) => r.id as number);

              const bulkResolve: BulkStatusUpdateCommand = {
                comment: commentForm?.comment ?? "",
                status: commentForm?.status,
                receiptIds: receiptIds,
              };
              this.receiptService
                .bulkReceiptStatusUpdate(bulkResolve)
                .pipe(
                  take(1),
                  tap((receipts) => {
                    let newReceipts = Array.from(this.dataSource().data);
                    receipts.forEach((r) => {
                      const receiptInTable = newReceipts.find(
                        (nr) => r.id === nr.id
                      ) as any as Receipt;
                      if (receiptInTable) {
                        receiptInTable.status = r.status;
                        receiptInTable.resolvedDate = r.resolvedDate;
                      }
                    });
                    this.dataSource.set(new MatTableDataSource(newReceipts));
                    // Same in-place patch, and this one moves receipts between the
                    // summary's status rows — the case that looks like it needs no
                    // refresh and needs it most.
                    this.summaryRequested.next();
                  })
                )
                .subscribe();
            }
          }
        )
      )
      .subscribe();
  }

  public pollEmail(): void {
    const groupId = this.store.selectSnapshot(GroupState.selectedGroupId);

    this.groupsService
      .pollGroupEmail(groupId as any)
      .pipe(
        take(1),
        tap(() => {
          this.snackbarService.success("Email successfully poll successfully queued");
        }),
      )
      .subscribe();
  }

  public exportSelectedReceipts(): void {
    const receiptIds = this.dataSource().data.map(data => data.id);
    this.receiptExportService.exportReceiptsById(receiptIds);
  }
}
