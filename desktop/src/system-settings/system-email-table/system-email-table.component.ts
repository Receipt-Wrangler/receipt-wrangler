import {AfterViewInit, Component, OnInit, signal, TemplateRef, viewChild} from "@angular/core";
import {MatDialog} from "@angular/material/dialog";
import {PageEvent} from "@angular/material/paginator";
import {Sort} from "@angular/material/sort";
import {MatTableDataSource} from "@angular/material/table";
import {ActivatedRoute} from "@angular/router";
import {Store} from "@ngxs/store";
import {take, tap} from "rxjs";
import {CheckEmailConnectivityCommand, Group, Permission, SystemEmail, SystemEmailService, SystemTaskStatus} from "../../open-api";
import {SnackbarService} from "../../services";
import {ConfirmationDialogComponent} from "../../shared-ui/confirmation-dialog/confirmation-dialog.component";
import {SystemEmailTableState} from "../../store/system-email-table.state";
import {SetOrderBy, SetPage, SetPageSize, SetSortDirection} from "../../store/system-email-table.state.actions";
import {TableColumn} from "../../table/table-column.interface";

@Component({
    selector: "app-system-email-table",
    templateUrl: "./system-email-table.component.html",
    styleUrl: "./system-email-table.component.scss",
    standalone: false
})
export class SystemEmailTableComponent implements OnInit, AfterViewInit {
  protected readonly Permission = Permission;

  public readonly usernameCell = viewChild.required<TemplateRef<any>>("usernameCell");

  public readonly hostCell = viewChild.required<TemplateRef<any>>("hostCell");

  public readonly createdAtCell = viewChild.required<TemplateRef<any>>("createdAtCell");

  public readonly updatedAtCell = viewChild.required<TemplateRef<any>>("updatedAtCell");

  public readonly actionsCell = viewChild.required<TemplateRef<any>>("actionsCell");

  public tableState = this.store.selectSignal(SystemEmailTableState.state);

  public displayedColumns: string[] = [];

  public columns: TableColumn[] = [];

  public dataSource = signal(new MatTableDataSource<SystemEmail>([]));

  public totalCount = signal(0);

  public allGroups: Group[] = [];

  public relatedSystemEmailMap: Map<number, Group[]> = new Map<number, Group[]>();

  constructor(
    private activatedRoute: ActivatedRoute,
    private snackbarService: SnackbarService,
    private store: Store,
    private systemEmailService: SystemEmailService,
    public matDialog: MatDialog,
  ) {
  }

  public ngOnInit(): void {
    this.allGroups = this.activatedRoute.snapshot.data?.["allGroups"];
    this.getTableData();
  }

  public ngAfterViewInit(): void {
    this.initTable();
  }

  private initTable(): void {
    this.setColumns();
  }

  private setColumns(): void {

    const columns = [
      {
        columnHeader: "Username",
        matColumnDef: "username",
        template: this.usernameCell(),
        sortable: true
      },
      {
        columnHeader: "Host",
        matColumnDef: "host",
        template: this.hostCell(),
        sortable: true
      },
      {
        columnHeader: "Created At",
        matColumnDef: "created_at",
        template: this.createdAtCell(),
        sortable: true
      },
      {
        columnHeader: "Updated At",
        matColumnDef: "updated_at",
        template: this.updatedAtCell(),
        sortable: true
      },
      {
        columnHeader: "Actions",
        matColumnDef: "actions",
        template: this.actionsCell(),
        sortable: false
      },

    ] as TableColumn[];

    const state = this.store.selectSnapshot(SystemEmailTableState.state);
    if (state.orderBy) {
      const column = columns.find((c) => c.matColumnDef === state.orderBy);
      if (column) {
        column.defaultSortDirection = state.sortDirection;
      }
    }

    this.columns = columns;
    this.displayedColumns = ["username", "host", "created_at", "updated_at", "actions"];
  }

  private getTableData(): void {
    const pagedRequestCommand = this.store.selectSnapshot(SystemEmailTableState.state);
    this.systemEmailService.getPagedSystemEmails(pagedRequestCommand)
      .pipe(
        take(1),
        tap((pagedData) => {
          this.dataSource.set(new MatTableDataSource(pagedData.data as SystemEmail[]));
          this.totalCount.set(pagedData.totalCount);
          this.setRelatedSystemEmailMap(pagedData.data as SystemEmail[]);
        })
      )
      .subscribe();
  }

  private setRelatedSystemEmailMap(systemEmails: SystemEmail[]): void {
    const map = new Map<number, Group[]>();
    systemEmails.forEach((systemEmail) => {
      const groups = this.allGroups.filter((group) => group.groupSettings?.systemEmailId === systemEmail.id);
      map.set(systemEmail.id, groups);
    });

    this.relatedSystemEmailMap = map;
  }

  public sorted(sort: Sort): void {
    this.store.dispatch(new SetOrderBy(sort.active));
    this.store.dispatch(new SetSortDirection(sort.direction));

    this.getTableData();
  }

  public pageChanged(pageEvent: PageEvent): void {
    const newPage = pageEvent.pageIndex + 1;

    this.store.dispatch(new SetPage(newPage));
    this.store.dispatch(new SetPageSize(pageEvent.pageSize));

    this.getTableData();
  }

  public deleteButtonClicked(systemEmail: SystemEmail): void {
    const dialogRef = this.matDialog.open(ConfirmationDialogComponent);

    dialogRef.componentInstance.headerText = "Delete System Email";
    dialogRef.componentInstance.dialogContent = `Are you sure you want to delete the email: ${systemEmail.username}?`;

    const index = this.dataSource().data.findIndex((se) => se.id === systemEmail.id);

    dialogRef.afterClosed()
      .pipe(
        take(1),
        tap((result) => {
          if (result) {
            this.callDeleteApi(systemEmail.id, index);
          }
        })
      )
      .subscribe();
  }

  private callDeleteApi(id: number, index: number): void {
    this.systemEmailService.deleteSystemEmailById(id)
      .pipe(
        take(1),
        tap(() => {
          this.getTableData();
          const data = Array.from(this.dataSource().data);
          data.splice(index, 1);
          this.dataSource.set(new MatTableDataSource(data));
          this.snackbarService.success("System email deleted successfully");
        })
      )
      .subscribe();
  }

  public checkEmailConnectivity(id: number): void {
    const command: CheckEmailConnectivityCommand = {id: id};
    this.systemEmailService.checkSystemEmailConnectivity(command)
      .pipe(
        take(1),
        tap(((systemTask) => {
          if (systemTask.status === SystemTaskStatus.Succeeded) {
            this.snackbarService.success("Successfully connected to email server");
          } else {
            this.snackbarService.error("Failed to connect to email server");
          }
        }))
      )
      .subscribe();
  }

  public disabledDeleteButtonClicked(email: SystemEmail): void {
    const mapData = this.relatedSystemEmailMap.get(email.id);
    const disabled = mapData && mapData.length > 0;

    if (disabled) {
      this.snackbarService.info(`Cannot delete ${email.username} because it is currently associated with the following groups: ${mapData.map((g) => g.name).join(", ")} `);
    }
  }
}
