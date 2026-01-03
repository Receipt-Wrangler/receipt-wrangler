import { AfterViewInit, Component, OnInit, TemplateRef, ViewChild, ViewEncapsulation, } from "@angular/core";
import { MatDialog } from "@angular/material/dialog";
import { PageEvent } from "@angular/material/paginator";
import { Sort } from "@angular/material/sort";
import { MatTableDataSource } from "@angular/material/table";
import { ActivatedRoute, Router } from "@angular/router";
import { UntilDestroy, untilDestroyed } from "@ngneat/until-destroy";
import { Select, Store } from "@ngxs/store";
import { map, Observable, take, tap } from "rxjs";
import { fadeInOut } from "src/animations";
import { ReceiptFilterService } from "src/services/receipt-filter.service";
import { ConfirmationDialogComponent } from "src/shared-ui/confirmation-dialog/confirmation-dialog.component";
import { ResetReceiptFilter, SetColumnConfig, SetPage, SetPageSize, SetReceiptFilterData, } from "src/store/receipt-table.actions";
import { ReceiptTableState } from "src/store/receipt-table.state";
import { TableColumn } from "src/table/table-column.interface";
import { TableComponent } from "src/table/table/table.component";
import { DEFAULT_DIALOG_CONFIG, DEFAULT_HOST_CLASS } from "../../constants";
import { ReceiptTableColumnConfig } from "../../interfaces";
import {
  BulkStatusUpdateCommand,
  Category,
  Group,
  GroupRole,
  GroupsService,
  PagedDataDataInner,
  Receipt,
  ReceiptPagedRequestCommand,
  ReceiptService,
  ReceiptStatus,
  Tag,
} from "../../open-api";
import { GroupRolePipe } from "../../pipes/group-role.pipe";
import { SnackbarService } from "../../services";
import { ReceiptExportService } from "../../services/receipt-export.service";
import { ReceiptFilterComponent } from "../../shared-ui/receipt-filter/receipt-filter.component";
import { GroupState } from "../../store";
import { applyFormCommand } from "../../utils/index";
import { buildReceiptFilterForm } from "../../utils/receipt-filter";
import { BulkStatusUpdateComponent } from "../bulk-resolve-dialog/bulk-status-update-dialog.component";
import { ColumnConfigurationDialogComponent } from "../column-configuration-dialog/column-configuration-dialog.component";

@UntilDestroy()
@Component({
  selector: "app-receipts-table",
  templateUrl: "./receipts-table.component.html",
  styleUrls: ["./receipts-table.component.scss"],
  providers: [GroupRolePipe],
  animations: [fadeInOut],
  encapsulation: ViewEncapsulation.None,
  host: DEFAULT_HOST_CLASS,
  standalone: false
})
export class ReceiptsTableComponent implements OnInit, AfterViewInit {
  constructor(
    private activatedRoute: ActivatedRoute,
    private groupPipe: GroupRolePipe,
    private groupsService: GroupsService,
    private matDialog: MatDialog,
    private receiptExportService: ReceiptExportService,
    private receiptFilterService: ReceiptFilterService,
    private receiptService: ReceiptService,
    private router: Router,
    private snackbarService: SnackbarService,
    private store: Store,
  ) {}

  @ViewChild("createdAtCell") createdAtCell!: TemplateRef<any>;

  @ViewChild("dateCell") dateCell!: TemplateRef<any>;

  @ViewChild("nameCell") nameCell!: TemplateRef<any>;

  @ViewChild("paidByCell") paidByCell!: TemplateRef<any>;

  @ViewChild("amountCell") amountCell!: TemplateRef<any>;

  @ViewChild("categoryCell") categoryCell!: TemplateRef<any>;

  @ViewChild("tagCell") tagCell!: TemplateRef<any>;

  @ViewChild("statusCell") statusCell!: TemplateRef<any>;

  @ViewChild("resolvedDateCell") resolvedDateCell!: TemplateRef<any>;

  @ViewChild("actionsCell") actionsCell!: TemplateRef<any>;

  @ViewChild(TableComponent) table!: TableComponent;

  @Select(ReceiptTableState.page) public page!: Observable<number>;

  @Select(ReceiptTableState.pageSize) public pageSize!: Observable<number>;

  @Select(ReceiptTableState.filterData) public filter!: Observable<ReceiptPagedRequestCommand>;

  @Select(ReceiptTableState.columnConfig) public columnConfig!: Observable<ReceiptTableColumnConfig[]>;

  @Select(GroupState.selectedGroupId)
  public selectedGroupId!: Observable<string>;

  public numFiltersApplied!: Observable<number | undefined>;

  public categories: Category[] = [];

  public tags: Tag[] = [];

  public groupId: string = "0";

  public groupRole = GroupRole;

  public dataSource: MatTableDataSource<PagedDataDataInner> =
    new MatTableDataSource<PagedDataDataInner>([]);

  public displayedColumns: string[] = [];

  public columns: TableColumn[] = [];

  public totalCount: number = 0;

  public selectedReceiptIds: number[] = [];

  public firstSort: boolean = true;

  public canEdit: boolean = false;

  public headerText: string = "";

  public group?: Group;

  public ngOnInit(): void {
    this.groupId = this.store
      .selectSnapshot(GroupState.selectedGroupId)
      ?.toString();
    this.setGroup();
    this.setCanEdit();

    this.setHeaderText();
    this.numFiltersApplied = this.store
      .select(ReceiptTableState.numFiltersApplied)
      .pipe(
        map((num) => {
          if (num > 0) {
            return num;
          } else {
            return undefined;
          }
        })
      );

    const data = this.activatedRoute.snapshot.data;
    this.categories = data["categories"];
    this.tags = data["tags"];
    this.getInitialData();
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
          this.dataSource = new MatTableDataSource<PagedDataDataInner>(pagedData.data);
          this.totalCount = pagedData.totalCount;
          this.setColumns();
        })
      )
      .subscribe();
  }

  private setCanEdit(): void {
    this.canEdit = this.groupPipe.transform(this.groupId, GroupRole.Editor);
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
    this.table?.selection?.changed
      .pipe(
        untilDestroyed(this),
        map((event) => (event.source.selected as Receipt[]).map((r) => r.id)),
        tap((ids) => (this.selectedReceiptIds = ids))
      )
      .subscribe();
  }

  private setColumns(): void {
    const currentColumnConfig = this.store.selectSnapshot(ReceiptTableState.columnConfig);

    const allColumns = [
      {
        columnHeader: "Added At",
        matColumnDef: "created_at",
        template: this.createdAtCell,
        sortable: true,
      },
      {
        columnHeader: "Receipt Date",
        matColumnDef: "date",
        template: this.dateCell,
        sortable: true,
      },
      {
        columnHeader: "Name",
        matColumnDef: "name",
        template: this.nameCell,
        sortable: true,
      },
      {
        columnHeader: "Paid By",
        matColumnDef: "paid_by_user_id",
        template: this.paidByCell,
        sortable: true,
      },
      {
        columnHeader: "Amount",
        matColumnDef: "amount",
        template: this.amountCell,
        sortable: true,
      },
      {
        columnHeader: "Categories",
        matColumnDef: "categories",
        template: this.categoryCell,
        sortable: false,
      },
      {
        columnHeader: "Tags",
        matColumnDef: "tags",
        template: this.tagCell,
        sortable: false,
      },
      {
        columnHeader: "Status",
        matColumnDef: "status",
        template: this.statusCell,
        sortable: true,
      },
      {
        columnHeader: "Resolved Date",
        matColumnDef: "resolved_date",
        template: this.resolvedDateCell,
        sortable: true,
      },
    ] as TableColumn[];

    // Filter and order columns based on configuration
    const visibleColumnConfigs = currentColumnConfig
      .filter(config => config.visible)
      .sort((a, b) => a.order - b.order);

    const columns = visibleColumnConfigs
      .map(config => allColumns.find(col => col.matColumnDef === config.matColumnDef))
      .filter(col => col !== undefined) as TableColumn[];

    const displayColumns = ["select", ...visibleColumnConfigs.map(config => config.matColumnDef)];

    if (this.canEdit) {
      columns.push({
        columnHeader: "Actions",
        matColumnDef: "actions",
        template: this.actionsCell,
        sortable: false,
      });
      displayColumns.push("actions");
    }

    const filter = this.store.selectSnapshot(ReceiptTableState.filterData);
    const orderByIndex = columns.findIndex(
      (c) => c.matColumnDef === filter.orderBy
    );

    if (orderByIndex >= 0) {
      columns[orderByIndex].defaultSortDirection = filter.sortDirection;
    } else if (columns.length > 0) {
      columns[0].defaultSortDirection = "desc";
    }

    this.columns = columns;
    this.displayedColumns = displayColumns;
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

      this.getFilteredReceipts();
    }
    this.firstSort = false;
  }

  public filterButtonClicked(): void {
    const filter = this.store.selectSnapshot(ReceiptTableState.filterData).filter as any;

    const dialogRef = this.matDialog.open(ReceiptFilterComponent, {
      minWidth: "75%",
      maxWidth: "100%",
      data: {
        categories: this.categories,
        tags: this.tags,
      },
    });

    dialogRef.componentInstance.parentForm = buildReceiptFilterForm(filter, this);
    dialogRef.componentInstance.headerText = "Filter Receipts";
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

  public resetFilterButtonClicked(): void {
    this.store.dispatch(new ResetReceiptFilter());
    this.getFilteredReceipts();
  }

  public configureColumnsButtonClicked(): void {
    const currentColumnConfig = this.store.selectSnapshot(ReceiptTableState.columnConfig);

    const dialogRef = this.matDialog.open(ColumnConfigurationDialogComponent, {
      ...DEFAULT_DIALOG_CONFIG,
      data: { currentColumns: currentColumnConfig }
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

  public getFilteredReceipts(): void {
    this.receiptFilterService
      .getPagedReceiptsForGroups(this.groupId.toString())
      .pipe(
        take(1),
        tap((pagedData) => {
          this.dataSource.data = pagedData.data;
          this.totalCount = pagedData.totalCount;
        })
      )
      .subscribe();
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
                  this.dataSource.data = this.dataSource.data.filter(
                    (r) => r.id !== row.id
                  );
                  this.snackbarService.success("Receipt successfully deleted");
                })
              )
              .subscribe();
          }
        })
      )
      .subscribe();
  }

  public duplicateReceipt(id: string): void {
    this.receiptService
      .duplicateReceipt(Number.parseInt(id))
      .pipe(
        tap((r: Receipt) => {
          this.snackbarService.success("Receipt successfully duplicated");
          this.router.navigateByUrl(`/receipts/${r.id}/view`);
        })
      )
      .subscribe();
  }

  public updatePageData(pageEvent: PageEvent): void {
    const newPage = pageEvent.pageIndex + 1;
    this.store.dispatch(new SetPage(newPage));
    this.store.dispatch(new SetPageSize(pageEvent.pageSize));

    this.getFilteredReceipts();
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
            if (this.table.selection.hasValue() && commentForm) {
              const receiptIds = (
                this.table.selection.selected as Receipt[]
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
                    let newReceipts = Array.from(this.dataSource.data);
                    receipts.forEach((r) => {
                      const receiptInTable = newReceipts.find(
                        (nr) => r.id === nr.id
                      ) as any as Receipt;
                      if (receiptInTable) {
                        receiptInTable.status = r.status;
                        receiptInTable.resolvedDate = r.resolvedDate;
                      }
                    });
                    this.dataSource.data = newReceipts;
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
    const receiptIds = this.dataSource.data.map(data => data.id);
    this.receiptExportService.exportReceiptsById(receiptIds);
  }
}
