import {
  AfterViewInit,
  Component,
  DestroyRef,
  TemplateRef,
  computed,
  inject,
  signal,
  viewChild,
} from "@angular/core";
import { takeUntilDestroyed } from "@angular/core/rxjs-interop";
import { FormControl } from "@angular/forms";
import { MatDialog } from "@angular/material/dialog";
import { MatTableDataSource } from "@angular/material/table";
import { Router } from "@angular/router";
import { EMPTY, catchError, take, tap } from "rxjs";
import {
  PermissionScope,
  Role,
  RoleService,
  PermissionService,
} from "../../open-api";
import { SnackbarService } from "../../services";
import { TableColumn } from "../../table/table-column.interface";
import { BreadcrumbItem } from "../../shared-ui/breadcrumb/breadcrumb-item.interface";
import { ConfirmationDialogComponent } from "../../shared-ui/confirmation-dialog/confirmation-dialog.component";
import { FilterTab } from "../../shared-ui/filter-bar/filter-tab.interface";
import {
  RoleListFilter,
  RoleListItem,
  RoleScope,
} from "./role-list-item.interface";

// A role carries no icon of its own; pick a sensible default per scope. These
// are inline style values, so the violet "group" accent is a literal (matching
// the documented out-of-palette accent used elsewhere on this page).
const SCOPE_ICON: Record<RoleScope, { icon: string; color: string; tint: string }> = {
  app: { icon: "apps", color: "#27b1ff", tint: "#ccecff" },
  group: { icon: "workspaces", color: "#6d28d9", tint: "#ede9fe" },
};

@Component({
  selector: "app-role-list",
  templateUrl: "./role-list.component.html",
  styleUrl: "./role-list.component.scss",
  standalone: false,
})
export class RoleListComponent implements AfterViewInit {
  private readonly permissionService = inject(PermissionService);
  private readonly roleService = inject(RoleService);
  private readonly router = inject(Router);
  private readonly matDialog = inject(MatDialog);
  private readonly snackbarService = inject(SnackbarService);
  private readonly destroyRef = inject(DestroyRef);

  private readonly roleCellTemplate = viewChild.required<TemplateRef<any>>("roleCell");
  private readonly typeCellTemplate = viewChild.required<TemplateRef<any>>("typeCell");
  private readonly permissionsCellTemplate =
    viewChild.required<TemplateRef<any>>("permissionsCell");
  private readonly membersCellTemplate = viewChild.required<TemplateRef<any>>("membersCell");
  private readonly actionsCellTemplate = viewChild.required<TemplateRef<any>>("actionsCell");

  public readonly columns = signal<TableColumn[]>([]);
  public readonly displayedColumns = signal<string[]>([
    "role",
    "type",
    "permissions",
    "members",
    "actions",
  ]);

  public readonly crumbs: BreadcrumbItem[] = [
    { label: "Admin" },
    { label: "Roles" },
  ];

  public readonly roles = signal<RoleListItem[]>([]);
  public readonly roleCount = computed(() => this.roles().length);

  public readonly filter = signal<RoleListFilter>("all");
  public readonly filteredRoles = computed(() => {
    const filter = this.filter();
    if (filter === "all") {
      return this.roles();
    }
    return this.roles().filter((role) => role.scope === filter);
  });

  public readonly counts = computed(() => {
    const roles = this.roles();
    return {
      all: roles.length,
      app: roles.filter((r) => r.scope === "app").length,
      group: roles.filter((r) => r.scope === "group").length,
    };
  });

  public readonly filterTabs = computed<FilterTab[]>(() => {
    const counts = this.counts();
    return [
      { value: "all", label: "All roles", icon: "list", count: counts.all },
      { value: "app", label: "Application", icon: "apps", count: counts.app },
      { value: "group", label: "Group", icon: "workspaces", count: counts.group },
    ];
  });

  public readonly dataSource = computed(
    () => new MatTableDataSource(this.filteredRoles()),
  );

  // Options for the two "default role" selectors: roles of each scope, valued by
  // id and displayed by name.
  public readonly appRoleOptions = computed(() =>
    this.roles()
      .filter((role) => role.scope === "app")
      .map((role) => ({ id: Number(role.id), name: role.name })),
  );
  public readonly groupRoleOptions = computed(() =>
    this.roles()
      .filter((role) => role.scope === "group")
      .map((role) => ({ id: Number(role.id), name: role.name })),
  );

  // The two default-role selectors. Exactly one role per scope is the default,
  // so there is no empty option. Values are patched from the server with
  // emitEvent:false; a genuine user selection drives setDefaultRole.
  public readonly defaultAppRoleControl = new FormControl<number | null>(null);
  public readonly defaultGroupRoleControl = new FormControl<number | null>(null);

  // Per-scope permission totals from the registry — the denominator for each
  // role's meter (a role's bar is relative to its own scope's total).
  private readonly scopeTotals = signal<Record<RoleScope, number>>({ app: 0, group: 0 });

  constructor() {
    this.loadScopeTotals();
    this.loadRoles();

    // A user picking a different role in either selector makes it the new default
    // for that scope. Programmatic sync from the server uses emitEvent:false, so
    // those updates do not re-trigger the API call.
    this.defaultAppRoleControl.valueChanges
      .pipe(takeUntilDestroyed())
      .subscribe((roleId) => this.onDefaultChange("app", roleId));
    this.defaultGroupRoleControl.valueChanges
      .pipe(takeUntilDestroyed())
      .subscribe((roleId) => this.onDefaultChange("group", roleId));
  }

  public ngAfterViewInit(): void {
    this.columns.set([
      {
        columnHeader: "Role",
        matColumnDef: "role",
        sortable: false,
        template: this.roleCellTemplate(),
      },
      {
        columnHeader: "Type",
        matColumnDef: "type",
        sortable: false,
        template: this.typeCellTemplate(),
      },
      {
        columnHeader: "Permissions",
        matColumnDef: "permissions",
        sortable: false,
        template: this.permissionsCellTemplate(),
      },
      {
        columnHeader: "Members",
        matColumnDef: "members",
        sortable: false,
        template: this.membersCellTemplate(),
      },
      {
        columnHeader: "",
        matColumnDef: "actions",
        sortable: false,
        template: this.actionsCellTemplate(),
      },
    ]);
  }

  public setFilter(filter: string): void {
    this.filter.set(filter as RoleListFilter);
  }

  /** Ten meter segments; filled ones tagged with the role's scope for color. */
  public meter(role: RoleListItem): (string | null)[] {
    const total = this.scopeTotals()[role.scope] || 1;
    const filled = Math.round((role.permissionCount / total) * 10);
    return Array.from({ length: 10 }, (_, i) => (i < filled ? role.scope : null));
  }

  public scopeTotal(role: RoleListItem): number {
    return this.scopeTotals()[role.scope];
  }

  public addRole(): void {
    this.router.navigate(["/roles/new"]);
  }

  // Opens the role for editing (or, for system roles, a read-only view). The
  // scope disambiguates app/group roles, which have independent id sequences.
  public editRole(role: RoleListItem): void {
    this.router.navigate(["/roles", role.id, "edit"], {
      queryParams: { scope: role.scope },
    });
  }

  // A role can be deleted only when it is not a system role and nobody is
  // assigned to it. Assigned roles must be unassigned before deletion.
  public canDelete(role: RoleListItem): boolean {
    return !role.isSystem && role.userCount === 0;
  }

  public deleteRole(role: RoleListItem): void {
    // Defensive: the button gates this via [disabled], but guard programmatic
    // calls so a non-deletable role never opens the confirmation or hits the API.
    if (!this.canDelete(role)) {
      this.disabledDeleteClicked(role);
      return;
    }

    const dialogRef = this.matDialog.open(ConfirmationDialogComponent);
    dialogRef.componentInstance.headerText = "Delete Role";
    dialogRef.componentInstance.dialogContent = `Are you sure you want to delete the role: ${role.name}?`;

    dialogRef
      .afterClosed()
      .pipe(
        take(1),
        tap((result) => {
          if (result) {
            this.callDeleteApi(role);
          }
        }),
      )
      .subscribe();
  }

  // Explains why the delete button is disabled when a disabled button is
  // clicked. We intentionally do not reveal who the role is assigned to.
  public disabledDeleteClicked(role: RoleListItem): void {
    if (role.isSystem) {
      this.snackbarService.info(
        `Cannot delete ${role.name} because it is a system role.`,
      );
    } else if (role.userCount > 0) {
      this.snackbarService.info(
        `Cannot delete ${role.name} because it is currently assigned.`,
      );
    }
  }

  private callDeleteApi(role: RoleListItem): void {
    this.roleService
      .deleteRole(this.toScopeEnum(role.scope), Number(role.id))
      .pipe(
        take(1),
        tap(() => {
          this.snackbarService.success("Role deleted successfully");
          this.loadRoles();
        }),
        catchError(() => {
          this.snackbarService.error("Failed to delete role");
          return EMPTY;
        }),
      )
      .subscribe();
  }

  // The view-model scope ('app'/'group') maps to the API's PermissionScope
  // enum ('APP'/'GROUP') — the inverse of the mapping in toListItem.
  private toScopeEnum(scope: RoleScope): PermissionScope {
    return scope === "group" ? PermissionScope.Group : PermissionScope.App;
  }

  // Makes the chosen role the default for its scope, then reloads so every
  // selector and the Default badges reflect the server (setting one default
  // clears the previous one). A null value (no selection) is ignored.
  private onDefaultChange(scope: RoleScope, roleId: number | null): void {
    if (roleId == null) {
      return;
    }

    this.roleService
      .setDefaultRole(this.toScopeEnum(scope), roleId)
      .pipe(
        take(1),
        tap(() => {
          this.snackbarService.success("Default role updated");
          this.loadRoles();
        }),
        catchError(() => {
          this.snackbarService.error("Failed to update default role");
          // Revert the selector to the server's current default.
          this.loadRoles();
          return EMPTY;
        }),
        takeUntilDestroyed(this.destroyRef),
      )
      .subscribe();
  }

  private loadRoles(): void {
    this.roleService
      .getRoles()
      .pipe(
        take(1),
        tap((roles) => {
          this.roles.set(roles.map((role) => this.toListItem(role)));
          this.syncDefaultControls();
        }),
        catchError(() => EMPTY),
        takeUntilDestroyed(this.destroyRef),
      )
      .subscribe();
  }

  // Patches the two selectors to the server's current default per scope without
  // emitting (so it never re-triggers setDefaultRole).
  private syncDefaultControls(): void {
    const appDefault = this.roles().find((role) => role.scope === "app" && role.isDefault);
    const groupDefault = this.roles().find((role) => role.scope === "group" && role.isDefault);

    this.defaultAppRoleControl.setValue(appDefault ? Number(appDefault.id) : null, {
      emitEvent: false,
    });
    this.defaultGroupRoleControl.setValue(groupDefault ? Number(groupDefault.id) : null, {
      emitEvent: false,
    });
  }

  private toListItem(role: Role): RoleListItem {
    const scope: RoleScope = role.scope === "GROUP" ? "group" : "app";
    const visual = SCOPE_ICON[scope];
    return {
      id: String(role.id),
      name: role.name,
      description: role.description ?? "",
      scope,
      permissionCount: role.permissions?.length ?? 0,
      userCount: role.assignedCount ?? 0,
      isDefault: role.isDefault,
      isSystem: role.isSystem,
      icon: visual.icon,
      iconColor: visual.color,
      iconTint: visual.tint,
    };
  }

  private loadScopeTotals(): void {
    this.permissionService
      .getPermissions()
      .pipe(
        take(1),
        tap((permissions) => {
          this.scopeTotals.set({
            app: permissions.filter((p) => p.scope === "APP").length,
            group: permissions.filter((p) => p.scope === "GROUP").length,
          });
        }),
        catchError(() => {
          this.scopeTotals.set({ app: 0, group: 0 });
          return EMPTY;
        }),
        takeUntilDestroyed(),
      )
      .subscribe();
  }
}
