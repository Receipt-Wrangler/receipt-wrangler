import { AfterViewInit, Component, OnInit, TemplateRef, ViewChild } from "@angular/core";
import { MatDialog } from "@angular/material/dialog";
import { UntilDestroy, untilDestroyed } from "@ngneat/until-destroy";
import { Select, Store } from "@ngxs/store";
import { Observable, take, tap } from "rxjs";
import { ConfirmationDialogComponent } from "src/shared-ui/confirmation-dialog/confirmation-dialog.component";
import { TableComponent } from "src/table/table/table.component";
import { DEFAULT_DIALOG_CONFIG, DEFAULT_HOST_CLASS } from "../../constants";
import { AssociatedGroup, Group, GroupRole, GroupsService, UserRole } from "../../open-api";
import { SnackbarService } from "../../services";
import { BaseTableService } from "../../services/base-table.service";
import { BaseTableComponent } from "../../shared-ui/base-table/base-table.component";
import { AuthState, GroupState, RemoveGroup } from "../../store";
import { GroupTableState } from "../../store/group-table.state";
import { GroupTableFilterComponent } from "../group-table-filter/group-table-filter.component";
import { GroupTableService } from "./group-table.service";


@UntilDestroy()
@Component({
    selector: "app-group-table",
    templateUrl: "./group-table.component.html",
    styleUrls: ["./group-table.component.scss"],
    host: DEFAULT_HOST_CLASS,
    providers: [
        {
            provide: BaseTableService,
            useClass: GroupTableService
        }
    ],
    standalone: false
})
export class GroupTableComponent extends BaseTableComponent<Group> implements OnInit, AfterViewInit {
  @Select(GroupState.groups) public groups!: Observable<Group[]>;

  @ViewChild("nameCell") private nameCell!: TemplateRef<any>;

  @ViewChild("numberOfMembersCell")
  private numberOfMembersCell!: TemplateRef<any>;

  @ViewChild("createdAtCell") private createdAtCell!: TemplateRef<any>;

  @ViewChild("updatedAtCell") private updatedAtCell!: TemplateRef<any>;

  @ViewChild("actionsCell") private actionsCell!: TemplateRef<any>;

  @ViewChild(TableComponent) private table!: TableComponent;

  public groupRole = GroupRole;

  public isAdmin = false;

  public tableHeaderText = "My Groups";

  constructor(
    public override baseTableService: BaseTableService,
    private groupsService: GroupsService,
    private store: Store,
    private snackbarService: SnackbarService,
    private matDialog: MatDialog
  ) {
    super(baseTableService);
  }

  public ngOnInit(): void {
    this.isAdmin = this.store.selectSnapshot(AuthState.hasRole(UserRole.Admin));
    this.listenForFilterChanges();
  }

  public ngAfterViewInit(): void {
    this.initTable();
  }

  private listenForFilterChanges(): void {
    this.store.select(GroupTableState.filter)
      .pipe(
        untilDestroyed(this),
        tap((filter) => {
            this.getTableData();
            if (filter.associatedGroup === AssociatedGroup.All) {
              this.tableHeaderText = "All Groups";
            } else {
              this.tableHeaderText = "My Groups";
            }
          }
        )
      )
      .subscribe();
  }

  private initTable(): void {
    this.setColumns();
  }

  private setColumns(): void {
    this.columns = [
      {
        columnHeader: "Name",
        matColumnDef: "name",
        template: this.nameCell,
        sortable: true,
      },
      {
        columnHeader: "Number of Members",
        matColumnDef: "number_of_members",
        template: this.numberOfMembersCell,
        sortable: false,
      },
      {
        columnHeader: "Created At",
        matColumnDef: "created_at",
        template: this.createdAtCell,
        sortable: true,
      },
      {
        columnHeader: "Updated At",
        matColumnDef: "updated_at",
        template: this.updatedAtCell,
        sortable: true,
      },
      {
        columnHeader: "Actions",
        matColumnDef: "actions",
        template: this.actionsCell,
        sortable: false,
      },
    ];
    this.displayedColumns = [
      "name",
      "number_of_members",
      "created_at",
      "updated_at",
      "actions",
    ];
  }

  public deleteGroup(index: number): void {
    const groups = this.dataSource.data;
    if (groups.length > 1) {
      const group = groups[index];
      const dialogRef = this.matDialog.open(
        ConfirmationDialogComponent,
        DEFAULT_DIALOG_CONFIG
      );

      dialogRef.componentInstance.headerText = "Delete Group";
      dialogRef.componentInstance.dialogContent = `Are you sure you would like to the group '${group.name}'? All receipts will be deleted as a result.`;

      dialogRef.afterClosed().subscribe((r) => {
        if (r) {
          this.groupsService
            .deleteGroup(group.id)
            .pipe(
              take(1),
              tap(() => {
                this.snackbarService.success("Group successfully deleted");
                this.store.dispatch(new RemoveGroup(group.id.toString()));
                this.dataSource.data = this.store.selectSnapshot(
                  GroupState.groupsWithoutAll
                );
              })
            )
            .subscribe();
        }
      });
    }
  }

  public openFilterDialog(): void {
    const ref = this.matDialog.open(GroupTableFilterComponent, DEFAULT_DIALOG_CONFIG);
  }
}
