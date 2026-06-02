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
import { MatDialog } from "@angular/material/dialog";
import { PageEvent } from "@angular/material/paginator";
import { Sort } from "@angular/material/sort";
import { MatTableDataSource } from "@angular/material/table";
import { Router } from "@angular/router";
import { Store } from "@ngxs/store";
import { EMPTY, catchError, take, tap } from "rxjs";
import {
  PagedRoleRequestCommand,
  PermissionScope,
  Role,
  RoleService,
  PermissionService,
} from "../../open-api";
import { SnackbarService } from "../../services";
import { RoleTableState } from "../../store/role-table.state";
import {
  SetOrderBy,
  SetPage,
  SetPageSize,
  SetScope,
  SetSortDirection,
} from "../../store/role-table.state.actions";
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
  private readonly store = inject(Store);
  private readonly destroyRef = inject(DestroyRef);

  private readonly roleCellTemplate = viewChild.required<TemplateRef<any>>("roleCell");
  private readonly typeCellTemplate = viewChild.required<TemplateRef<any>>("typeCell");
  private readonly permissionsCellTemplate =
    viewChild.required<TemplateRef<any>>("permissionsCell");
  private readonly membersCellTemplate = viewChild.required<TemplateRef<any>>("membersCell");
  private readonly actionsCellTemplate = viewChild.required<TemplateRef<any>>("actionsCell");

  // Both are populated together in ngAfterViewInit (the cell templates aren't
  // available earlier). They must stay empty until then: the paged table renders
  // on first paint, and MatTable errors if displayedColumns references a column
  // whose matColumnDef hasn't been registered yet.
  public readonly columns = signal<TableColumn[]>([]);
  public readonly displayedColumns = signal<string[]>([]);

  public readonly crumbs: BreadcrumbItem[] = [
    { label: "Admin" },
    { label: "Roles" },
  ];

  // Server-side paged table state. Paging fields come from the shared `state`
  // selector; the scope filter (which backs the filter bar) is exposed
  // separately and is mapped to the API scope when building the request.
  public readonly state = this.store.selectSignal(RoleTableState.state);
  public readonly filter = this.store.selectSignal(RoleTableState.scope);

  public readonly dataSource = signal(new MatTableDataSource<RoleListItem>([]));
  public readonly totalCount = signal(0);

  // Becomes true after the first response so the empty-state CTA doesn't flash
  // before the initial page has loaded.
  private readonly loaded = signal(false);
  public readonly showEmptyState = computed(
    () => this.loaded() && this.filter() === "all" && this.totalCount() === 0,
  );

  // The numeric count badges were dropped when the table moved to server-side
  // paging — a single page can't report per-scope totals.
  public readonly filterTabs: FilterTab[] = [
    { value: "all", label: "All roles", icon: "list" },
    { value: "app", label: "Application", icon: "apps" },
    { value: "group", label: "Group", icon: "workspaces" },
  ];

  // Per-scope permission totals from the registry — the denominator for each
  // role's meter (a role's bar is relative to its own scope's total).
  private readonly scopeTotals = signal<Record<RoleScope, number>>({ app: 0, group: 0 });

  constructor() {
    this.loadScopeTotals();
    this.loadRoles();
  }

  public ngAfterViewInit(): void {
    const state = this.store.selectSnapshot(RoleTableState.state);
    const columns: TableColumn[] = [
      {
        columnHeader: "Role",
        matColumnDef: "name",
        sortable: true,
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
    ];

    // Restore the persisted sort indicator on the sortable column.
    if (state.orderBy) {
      const column = columns.find((c) => c.matColumnDef === state.orderBy);
      if (column) {
        column.defaultSortDirection = state.sortDirection;
      }
    }

    this.columns.set(columns);
    this.displayedColumns.set(["name", "type", "permissions", "members", "actions"]);
  }

  public setFilter(filter: string): void {
    this.store.dispatch(new SetScope(filter as RoleListFilter));
    this.loadRoles();
  }

  public updatePageData(pageEvent: PageEvent): void {
    this.store.dispatch(new SetPage(pageEvent.pageIndex + 1));
    this.store.dispatch(new SetPageSize(pageEvent.pageSize));
    this.loadRoles();
  }

  public sorted(sortState: Sort): void {
    this.store.dispatch(new SetOrderBy(sortState.active));
    this.store.dispatch(new SetSortDirection(sortState.direction));
    this.loadRoles();
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

  // Maps the filter-bar scope ('all'/'app'/'group') to the API filter. "all"
  // omits the scope so the backend returns both scopes.
  private toScopeFilter(scope: RoleListFilter): PermissionScope | undefined {
    if (scope === "app") {
      return PermissionScope.App;
    }
    if (scope === "group") {
      return PermissionScope.Group;
    }
    return undefined;
  }

  private buildRequest(): PagedRoleRequestCommand {
    const state = this.store.selectSnapshot(RoleTableState.state);
    const scope = this.store.selectSnapshot(RoleTableState.scope);
    return {
      page: state.page,
      pageSize: state.pageSize,
      orderBy: state.orderBy,
      sortDirection: state.sortDirection,
      filter: { scope: this.toScopeFilter(scope) },
    };
  }

  private loadRoles(): void {
    this.roleService
      .getPagedRoles(this.buildRequest())
      .pipe(
        take(1),
        tap((pagedData) => {
          const items = (pagedData.data ?? []).map((role) =>
            this.toListItem(role as Role),
          );
          this.dataSource.set(new MatTableDataSource(items));
          this.totalCount.set(pagedData.totalCount);
          this.loaded.set(true);
        }),
        catchError(() => EMPTY),
        takeUntilDestroyed(this.destroyRef),
      )
      .subscribe();
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
