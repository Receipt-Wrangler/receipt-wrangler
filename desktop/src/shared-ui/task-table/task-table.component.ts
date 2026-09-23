import { AfterViewInit, Component, computed, Inject, OnInit, signal, TemplateRef, input, viewChild } from "@angular/core";
import { takeUntilDestroyed } from "@angular/core/rxjs-interop";
import { MatDialog } from "@angular/material/dialog";
import { PageEvent } from "@angular/material/paginator";
import { Sort } from "@angular/material/sort";
import { MatTableDataSource } from "@angular/material/table";
import { catchError, EMPTY, Subject, switchMap, take, tap } from "rxjs";
import { DEFAULT_DIALOG_CONFIG } from "../../constants/dialog.constant";
import { AssociatedEntityType, GetSystemTaskCommand, SystemTask, SystemTaskPagedRequestFilter, SystemTaskService, SystemTaskType } from "../../open-api";
import { BaseTableService } from "../../services/base-table.service";
import { changedTopLevelKeys, parseReceiptUpdateDescription, ReceiptUpdateSnapshots } from "../../utils/receipt-update-description";
import { toSystemTaskWireFilter } from "../../utils/system-task-filter";
import { TABLE_SERVICE_INJECTION_TOKEN } from "../../services/injection-tokens/table-service";
import { TableColumn } from "../../table/table-column.interface";
import { downloadFile, filenameFromContentDisposition } from "../../utils/file";
import { DescriptionViewerDialogComponent, DescriptionViewerDialogData } from "../description-viewer-dialog/description-viewer-dialog.component";
import { ReceiptUpdateDiffDialogComponent, ReceiptUpdateDiffDialogData } from "../receipt-update-diff-dialog/receipt-update-diff-dialog.component";
import { SourceFileViewerDialogComponent, SourceFileViewerDialogData } from "../source-file-viewer-dialog/source-file-viewer-dialog.component";

interface ReceiptUpdateRow {
  snapshots: ReceiptUpdateSnapshots;
  changedKeys: string[];
}

@Component({
  selector: "app-task-table",
  templateUrl: "./task-table.component.html",
  styleUrl: "./task-table.component.scss",
  standalone: false
})
export class TaskTableComponent implements OnInit, AfterViewInit {
  public readonly typeCell = viewChild.required<TemplateRef<any>>("typeCell");

  public readonly startedAtCell = viewChild.required<TemplateRef<any>>("startedAtCell");

  public readonly endedAtCell = viewChild.required<TemplateRef<any>>("endedAtCell");

  public readonly statusCell = viewChild.required<TemplateRef<any>>("statusCell");

  public readonly resultDescriptionCell = viewChild.required<TemplateRef<any>>("resultDescriptionCell");

  public readonly ranByUserIdCell = viewChild.required<TemplateRef<any>>("ranByUserIdCell");

  public readonly sourceFileCell = viewChild.required<TemplateRef<any>>("sourceFileCell");

  /**
   * Opt-in, so the component's other hosts are unchanged. The server decides per
   * row whether a source file is reachable (models.SystemTask.hasSourceFile), so
   * this only says whether the host wants the column at all.
   */
  public readonly showSourceFileActions = input(false);

  public readonly associatedEntityType = input<AssociatedEntityType>();

  public readonly associatedEntityId = input<number>();

  public readonly expandedRowTemplate = input<TemplateRef<any>>();

  /**
   * Reads the system task filter to send with each request.
   *
   * A **function, not the filter itself**: the System Tasks page dispatches a
   * filter change and calls `getTableData()` in the same synchronous turn,
   * before change detection has pushed a new input value in — so a value input
   * would send the *previous* filter, and clearing a chip would leave the
   * cleared condition applied. Called at request time, it always reads current
   * state.
   *
   * Only that page binds it. The two embedded task tables leave it undefined,
   * so the key is omitted from the request and the API's zero-value filter adds
   * no predicates.
   */
  public readonly filterProvider = input<() => SystemTaskPagedRequestFilter | undefined>();

  public displayedColumns: string[] = [];

  public columns: TableColumn[] = [];

  public dataSource = signal(new MatTableDataSource<SystemTask>([]));

  public totalCount = signal(0);

  /**
   * RECEIPT_UPDATED rows whose description parses, keyed by task id. Parsed
   * once per page rather than per change-detection pass, since each one holds
   * two whole receipts. A row missing from the map (a failed update stores its
   * plain error text) renders like any other task.
   */
  public readonly receiptUpdates = computed(() => {
    const updates = new Map<number, ReceiptUpdateRow>();
    for (const task of this.dataSource().data) {
      if (task.type !== SystemTaskType.ReceiptUpdated) {
        continue;
      }
      const snapshots = parseReceiptUpdateDescription(task.resultDescription);
      if (snapshots) {
        updates.set(task.id, { snapshots, changedKeys: changedTopLevelKeys(snapshots.before, snapshots.after) });
      }
    }
    return updates;
  });

  public rowExpandable: (row: SystemTask) => boolean = (systemTask) => (systemTask?.childSystemTasks?.length || 0) > 0;

  // Descriptions shorter than this render inline without a "View More"
  // button; longer values get truncated and opened in a dialog on demand.
  public readonly descriptionInlineMaxLength = 120;

  private readonly refreshRequests = new Subject<void>();

  constructor(
    @Inject(TABLE_SERVICE_INJECTION_TOKEN) public tableService: BaseTableService,
    private systemTaskService: SystemTaskService,
    private dialog: MatDialog,
  ) {
    this.listenForRefreshRequests();
  }

  public openSourceFileDialog(element: SystemTask): void {
    this.systemTaskService
      .getSystemTaskSourceFile(element.id)
      .pipe(
        take(1),
        tap((sourceFile) => {
          const data: SourceFileViewerDialogData = {
            encodedImage: sourceFile.encodedImage,
            name: sourceFile.name,
          };
          this.dialog.open(SourceFileViewerDialogComponent, { ...DEFAULT_DIALOG_CONFIG, data });
        })
      ).subscribe();
  }

  public downloadSourceFile(element: SystemTask): void {
    // Observed as a response so the original file name can be read off
    // Content-Disposition; a task row does not carry it.
    this.systemTaskService
      .downloadSystemTaskSourceFile(element.id, "response")
      .pipe(
        take(1),
        tap((response) => {
          if (!response.body) {
            return;
          }

          downloadFile(response.body, filenameFromContentDisposition(response.headers.get("Content-Disposition")));
        })
      ).subscribe();
  }

  public openDescriptionDialog(element: SystemTask): void {
    const data: DescriptionViewerDialogData = {
      description: element.resultDescription ?? "",
      headerText: "Description",
    };
    this.dialog.open(DescriptionViewerDialogComponent, {
      ...DEFAULT_DIALOG_CONFIG,
      data,
    });
  }

  public openReceiptUpdateDialog(snapshots: ReceiptUpdateSnapshots): void {
    const data: ReceiptUpdateDiffDialogData = { snapshots };
    this.dialog.open(ReceiptUpdateDiffDialogComponent, {
      ...DEFAULT_DIALOG_CONFIG,
      width: "90vw",
      maxWidth: "1400px",
      data,
    });
  }

  public ngOnInit(): void {
    this.getTableData();
  }

  public getTableData(): void {
    this.refreshRequests.next();
  }

  /**
   * One subscription for every refresh, so the last *request* wins rather than
   * the last *response* — clearing two filter chips in quick succession puts
   * two fetches in flight, and a slow earlier one would otherwise repaint the
   * table with a filter the user has already moved past.
   *
   * The catchError lives on the inner observable: an error surfacing *through*
   * switchMap would complete the outer subscription and silently kill every
   * later refresh.
   */
  private listenForRefreshRequests(): void {
    this.refreshRequests
      .pipe(
        switchMap(() => {
          const pagedCommand = this.tableService.getPagedRequestCommand();
          const getSystemTaskCommand: GetSystemTaskCommand = {
            page: pagedCommand.page,
            pageSize: pagedCommand.pageSize,
            orderBy: pagedCommand.orderBy,
            sortDirection: pagedCommand.sortDirection,
            associatedEntityId: this.associatedEntityId(),
            associatedEntityType: this.associatedEntityType(),
            // Normalized here rather than in the store: see toSystemTaskWireFilter.
            filter: toSystemTaskWireFilter(this.filterProvider()?.())
          };

          return this.systemTaskService.getPagedSystemTasks(getSystemTaskCommand)
            .pipe(catchError(() => EMPTY));
        }),
        tap((pagedData) => {
          this.dataSource.set(new MatTableDataSource<SystemTask>((pagedData.data as any[]) as SystemTask[]));
          this.totalCount.set(pagedData.totalCount);
        }),
        takeUntilDestroyed(),
      )
      .subscribe();
  }

  public ngAfterViewInit(): void {
    this.initTable();
  }

  private initTable(): void {
    this.setColumns();
  }

  private setColumns(): void {
    this.columns = [
      {
        columnHeader: "Type",
        matColumnDef: "type",
        template: this.typeCell(),
        sortable: true,
      },
      {
        columnHeader: "Started At",
        matColumnDef: "started_at",
        template: this.startedAtCell(),
        sortable: true,
      },
      {
        columnHeader: "Ended At",
        matColumnDef: "ended_at",
        template: this.endedAtCell(),
        sortable: true,
      },
      {
        columnHeader: "Status",
        matColumnDef: "status",
        template: this.statusCell(),
        sortable: true,
      },
      {
        columnHeader: "Description",
        matColumnDef: "result_description",
        template: this.resultDescriptionCell(),
        sortable: true,
      },
      {
        columnHeader: "Ran By",
        matColumnDef: "ran_by_user_id",
        template: this.ranByUserIdCell(),
        sortable: true,
      }
    ];

    this.displayedColumns = ["started_at", "ended_at", "type", "ran_by_user_id", "result_description", "status"];

    if (this.showSourceFileActions()) {
      this.columns.push({
        columnHeader: "Source File",
        matColumnDef: "source_file",
        template: this.sourceFileCell(),
        sortable: false,
      });
      // Before the expand push below: "expand" has to stay last, and mat-table
      // throws on a displayed id with no matching column definition.
      this.displayedColumns.push("source_file");
    }

    if (this.expandedRowTemplate()) {
      this.displayedColumns.push("expand");
    }
  }

  public sorted(sort: Sort): void {
    this.tableService.setOrderBy(sort);
    this.tableService.setSortDirection(sort.direction);

    this.getTableData();
  }

  public pageChanged(event: PageEvent): void {
    const newPage = event.pageIndex + 1;

    this.tableService.setPage(newPage);
    this.tableService.setPageSize(event.pageSize);

    this.getTableData();
  }

  protected readonly SystemTaskType = SystemTaskType;
}
