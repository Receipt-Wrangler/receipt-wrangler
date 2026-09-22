import { AfterViewInit, Component, Inject, OnInit, signal, TemplateRef, input, viewChild } from "@angular/core";
import { MatDialog } from "@angular/material/dialog";
import { PageEvent } from "@angular/material/paginator";
import { Sort } from "@angular/material/sort";
import { MatTableDataSource } from "@angular/material/table";
import { take, tap } from "rxjs";
import { DEFAULT_DIALOG_CONFIG } from "../../constants/dialog.constant";
import { AssociatedEntityType, GetSystemTaskCommand, SystemTask, SystemTaskService, SystemTaskType } from "../../open-api";
import { BaseTableService } from "../../services/base-table.service";
import { TABLE_SERVICE_INJECTION_TOKEN } from "../../services/injection-tokens/table-service";
import { TableColumn } from "../../table/table-column.interface";
import { downloadFile, filenameFromContentDisposition } from "../../utils/file";
import { DescriptionViewerDialogComponent, DescriptionViewerDialogData } from "../description-viewer-dialog/description-viewer-dialog.component";
import { SourceFileViewerDialogComponent, SourceFileViewerDialogData } from "../source-file-viewer-dialog/source-file-viewer-dialog.component";

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

  public displayedColumns: string[] = [];

  public columns: TableColumn[] = [];

  public dataSource = signal(new MatTableDataSource<SystemTask>([]));

  public totalCount = signal(0);

  public rowExpandable: (row: SystemTask) => boolean = (systemTask) => (systemTask?.childSystemTasks?.length || 0) > 0;

  // Descriptions shorter than this render inline without a "View More"
  // button; longer values get truncated and opened in a dialog on demand.
  public readonly descriptionInlineMaxLength = 120;

  constructor(
    @Inject(TABLE_SERVICE_INJECTION_TOKEN) public tableService: BaseTableService,
    private systemTaskService: SystemTaskService,
    private dialog: MatDialog,
  ) {}

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

  public ngOnInit(): void {
    this.getTableData();
  }

  public getTableData(): void {
    const pagedCommand = this.tableService.getPagedRequestCommand();
    const getSystemTaskCommand: GetSystemTaskCommand = {
      page: pagedCommand.page,
      pageSize: pagedCommand.pageSize,
      orderBy: pagedCommand.orderBy,
      sortDirection: pagedCommand.sortDirection,
      associatedEntityId: this.associatedEntityId(),
      associatedEntityType: this.associatedEntityType()
    };

    this.systemTaskService.getPagedSystemTasks(getSystemTaskCommand)
      .pipe(
        take(1),
        tap((pagedData) => {
          this.dataSource.set(new MatTableDataSource<SystemTask>((pagedData.data as any[]) as SystemTask[]));
          this.totalCount.set(pagedData.totalCount);
        })
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
