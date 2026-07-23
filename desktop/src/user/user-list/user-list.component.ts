import { AfterViewInit, Component, inject, TemplateRef, viewChild } from "@angular/core";
import { takeUntilDestroyed, toSignal } from "@angular/core/rxjs-interop";
import { MatDialog } from "@angular/material/dialog";
import { Store } from "@ngxs/store";
import { take, tap } from "rxjs";
import { DEFAULT_HOST_CLASS } from "src/constants";
import { DEFAULT_DIALOG_CONFIG } from "src/constants/dialog.constant";
import { PagedTableInterface } from "src/interfaces/paged-table.interface";
import { ConfirmationDialogComponent } from "src/shared-ui/confirmation-dialog/confirmation-dialog.component";
import { TableColumn } from "src/table/table-column.interface";
import { TableComponent } from "src/table/table/table.component";
import { BulkUserDeleteCommand, Permission, PermissionScope, Role, RoleService, User, UserService } from "../../open-api";
import { loadAssignableRoles } from "../../roles/role-loading.util";
import { SnackbarService } from "../../services";
import { BaseTableService } from "../../services/base-table.service";
import { BaseTableComponent } from "../../shared-ui/base-table/base-table.component";
import { AuthState, RemoveUser, RemoveUsers } from "../../store";
import { DummyUserConversionDialogComponent } from "../dummy-user-conversion-dialog/dummy-user-conversion-dialog.component";
import { ResetPasswordComponent } from "../reset-password/reset-password.component";
import { UserFormComponent } from "../user-form/user-form.component";
import { UserTableService } from "./user-table.service";

@Component({
    selector: "app-user-list",
    templateUrl: "./user-list.component.html",
    styleUrls: ["./user-list.component.scss"],
    host: DEFAULT_HOST_CLASS,
    providers: [
        {
            provide: BaseTableService,
            useClass: UserTableService
        }
    ],
    standalone: false
})
export class UserListComponent extends BaseTableComponent<User> implements AfterViewInit {
  protected readonly Permission = Permission;
  protected readonly PermissionScope = PermissionScope;

  userId = this.store.selectSignal(AuthState.userId);

  public readonly usernameCell = viewChild.required<TemplateRef<any>>("usernameCell");

  public readonly displaynameCell = viewChild.required<TemplateRef<any>>("displayNameCell");

  public readonly appRoleCell = viewChild.required<TemplateRef<any>>("appRoleCell");

  public readonly createdAtCell = viewChild.required<TemplateRef<any>>("createdAtCell");

  public readonly updatedAtCell = viewChild.required<TemplateRef<any>>("updatedAtCell");

  public readonly actionsCell = viewChild.required<TemplateRef<any>>("actionsCell");

  public readonly table = viewChild.required(TableComponent);

  // Roles resolve each user's appRoleId to a role name. The request is skipped
  // unless the caller holds app.roles.read (otherwise the 403 would log them
  // out — see loadAssignableRoles); the cell then renders blank.
  public readonly roles = toSignal(
    loadAssignableRoles(inject(Store), inject(RoleService)),
    { initialValue: [] as Role[] }
  );

  public hasSelectedUsers: boolean = false;

  constructor(
    public override baseTableService: BaseTableService,
    private matDialog: MatDialog,
    private snackbarService: SnackbarService,
    private store: Store,
    private userService: UserService
  ) {
    super(baseTableService);
  }

  public ngAfterViewInit(): void {
    this.setColumns();
    this.getTableData();
    this.setupSelectionListener();
  }

  private setColumns(): void {
    this.columns = [
      {
        columnHeader: "Username",
        matColumnDef: "username",
        template: this.usernameCell(),
        sortable: true,
      },

      {
        columnHeader: "Displayname",
        matColumnDef: "display_name",
        template: this.displaynameCell(),
        sortable: true,
      },
      {
        columnHeader: "Role",
        matColumnDef: "app_role",
        template: this.appRoleCell(),
        sortable: false,
      },
      {
        columnHeader: "Created At",
        matColumnDef: "created_at",
        template: this.createdAtCell(),
        sortable: true,
      },
      {
        columnHeader: "Updated At",
        matColumnDef: "updated_at",
        template: this.updatedAtCell(),
        sortable: true,
      },
      {
        columnHeader: "Actions",
        matColumnDef: "actions",
        template: this.actionsCell(),
        sortable: false,
      },
    ];

    this.displayedColumns = [
      "select",
      "username",
      "display_name",
      "app_role",
      "created_at",
      "updated_at",
      "actions",
    ];

    this.setInitialSortedColumn(
      this.baseTableService.getPagedRequestCommand() as PagedTableInterface,
      this.columns
    );
  }

  private setupSelectionListener(): void {
    this.table()
      .selection.changed.pipe(takeUntilDestroyed())
      .subscribe(() => {
        this.updateSelectionState();
      });
  }

  private updateSelectionState(): void {
    this.hasSelectedUsers = this.table().selection.selected.length > 0;
  }

  public openUserFormDialog(user?: User): void {
    const dialogRef = this.matDialog.open(
      UserFormComponent,
      DEFAULT_DIALOG_CONFIG
    );

    dialogRef.componentInstance.user = user;

    dialogRef.afterClosed().subscribe((refresh) => {
      if (refresh) {
        this.getTableData();
      }
    });
  }

  public openResetPasswordDialog(user: User): void {
    const dialogRef = this.matDialog.open(
      ResetPasswordComponent,
      DEFAULT_DIALOG_CONFIG
    );

    dialogRef.componentInstance.user = user;
  }

  public openDummyUserConversionDialog(user: User): void {
    const dialogRef = this.matDialog.open(
      DummyUserConversionDialogComponent,
      DEFAULT_DIALOG_CONFIG
    );

    dialogRef.componentInstance.user = user;
  }

  public deleteUser(index: number) {
    const users = this.dataSource().data;
    const userId = this.store.selectSnapshot(AuthState.userId);
    const user = users[index];

    if (user.id.toString() !== userId) {
      const dialogRef = this.matDialog.open(
        ConfirmationDialogComponent,
        DEFAULT_DIALOG_CONFIG
      );

      dialogRef.componentInstance.headerText = "Delete User";
      dialogRef.componentInstance.dialogContent = `Are you sure you would like to delete the user '${user.username}'? This will remove the user from the user from groups, the user's receipt items, groups where this user is the only member, and receipts where the user paid. This action is irreversible.`;

      dialogRef.afterClosed().subscribe((r) => {
        if (r) {
          this.userService
            .deleteUserById(user.id)
            .pipe(
              take(1),
              tap(() => {
                this.snackbarService.success("User successfully deleted");
                this.store.dispatch(new RemoveUser(user.id.toString()));
                this.getTableData();
              })
            )
            .subscribe();
        }
      });
    }
  }

  public bulkDeleteUsers(): void {
    const selectedUsers = this.table().selection.selected;
    const currentUserId = this.store.selectSnapshot(AuthState.userId);

    const usersToDelete = selectedUsers.filter(user => user.id.toString() !== currentUserId);

    if (usersToDelete.length === 0) {
      this.snackbarService.error("Cannot delete current user or no valid users selected");
      return;
    }

    const dialogRef = this.matDialog.open(
      ConfirmationDialogComponent,
      DEFAULT_DIALOG_CONFIG
    );

    const usernames = usersToDelete.map(user => user.username).join(", ");
    dialogRef.componentInstance.headerText = "Delete Users";
    dialogRef.componentInstance.dialogContent = `Are you sure you would like to delete ${usersToDelete.length} user(s): ${usernames}? This will remove these users from groups, their receipt items, groups where they are the only member, and receipts where they paid. This action is irreversible.`;

    dialogRef.afterClosed().subscribe((confirmed) => {
      if (confirmed) {
        const bulkDeleteCommand: BulkUserDeleteCommand = {
          userIds: usersToDelete.map(user => user.id.toString())
        };

        this.userService
          .bulkDeleteUsers(bulkDeleteCommand)
          .pipe(
            take(1),
            tap(() => {
              this.snackbarService.success(`${usersToDelete.length} user(s) successfully deleted`);
              this.store.dispatch(new RemoveUsers(bulkDeleteCommand.userIds));
              this.table().selection.clear();
              this.getTableData();
            })
          )
          .subscribe();
      }
    });
  }
}
