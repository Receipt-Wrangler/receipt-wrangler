import { AfterViewInit, Component, Inject, OnInit, signal, TemplateRef, input, viewChild } from "@angular/core";
import { PageEvent } from "@angular/material/paginator";
import { Sort } from "@angular/material/sort";
import { MatTableDataSource } from "@angular/material/table";
import { take, tap } from "rxjs";
import { AssociatedEntityType, GetSystemTaskCommand, SystemTask, SystemTaskService, SystemTaskType } from "../../open-api";
import { BaseTableService } from "../../services/base-table.service";
import { TABLE_SERVICE_INJECTION_TOKEN } from "../../services/injection-tokens/table-service";
import { TableColumn } from "../../table/table-column.interface";

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

  public readonly associatedEntityType = input<AssociatedEntityType>();

  public readonly associatedEntityId = input<number>();

  public readonly expandedRowTemplate = input<TemplateRef<any>>();

  public displayedColumns: string[] = [];

  public columns: TableColumn[] = [];

  public dataSource = signal(new MatTableDataSource<SystemTask>([]));

  public totalCount = signal(0);

  public rowExpandable: (row: SystemTask) => boolean = (systemTask) => (systemTask?.childSystemTasks?.length || 0) > 0;

  constructor(
    @Inject(TABLE_SERVICE_INJECTION_TOKEN) public tableService: BaseTableService,
    private systemTaskService: SystemTaskService
  ) {}

  public ngOnInit(): void {
    this.getTableData();
  }

  public getTableData(): void {
    const pagedCommand = this.tableService.getPagedRequestCommand() as any;
    const getSystemTaskCommand: GetSystemTaskCommand = {
      page: pagedCommand.page,
      pageSize: pagedCommand.pageSize,
      orderBy: pagedCommand.orderBy,
      sortDirection: pagedCommand.sortDirection,
      associatedEntityId: this.associatedEntityId(),
      associatedEntityType: this.associatedEntityType(),
      filter: pagedCommand.filter,
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
