import {
  AfterViewInit,
  Component,
  TemplateRef,
  computed,
  inject,
  signal,
  viewChild,
} from "@angular/core";
import { takeUntilDestroyed } from "@angular/core/rxjs-interop";
import { MatTableDataSource } from "@angular/material/table";
import { Router } from "@angular/router";
import { EMPTY, catchError, take, tap } from "rxjs";
import { Role, RoleService, PermissionService } from "../../open-api";
import { TableColumn } from "../../table/table-column.interface";
import { BreadcrumbItem } from "../../shared-ui/breadcrumb/breadcrumb-item.interface";
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

  // Per-scope permission totals from the registry — the denominator for each
  // role's meter (a role's bar is relative to its own scope's total).
  private readonly scopeTotals = signal<Record<RoleScope, number>>({ app: 0, group: 0 });

  constructor() {
    this.loadScopeTotals();
    this.loadRoles();
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

  /** Up-to-two-letter initials for a member avatar. */
  public initials(name: string): string {
    return (name || "?")
      .split(" ")
      .map((part) => part[0])
      .join("")
      .slice(0, 2)
      .toUpperCase();
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

  // TODO: open the per-role actions menu once those actions exist.
  public openRoleMenu(_role: RoleListItem): void {}

  private loadRoles(): void {
    this.roleService
      .getRoles()
      .pipe(
        take(1),
        tap((roles) => this.roles.set(roles.map((role) => this.toListItem(role)))),
        catchError(() => EMPTY),
        takeUntilDestroyed(),
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
      members: [],
      userCount: 0,
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
