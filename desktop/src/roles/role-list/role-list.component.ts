import { Component, computed, inject, signal } from "@angular/core";
import { takeUntilDestroyed } from "@angular/core/rxjs-interop";
import { take, tap } from "rxjs";
import { PermissionService } from "../../open-api";
import { BreadcrumbItem } from "../../shared-ui/breadcrumb/breadcrumb-item.interface";
import { RoleListItem } from "./role-list-item.interface";

@Component({
  selector: "app-role-list",
  templateUrl: "./role-list.component.html",
  styleUrl: "./role-list.component.scss",
  standalone: false,
})
export class RoleListComponent {
  private readonly permissionService = inject(PermissionService);

  public readonly crumbs: BreadcrumbItem[] = [
    { label: "Admin" },
    { label: "Roles" },
  ];

  // Empty until the backend role-list endpoint exists. The table below is built
  // to render real rows the moment this signal is populated.
  public readonly roles = signal<RoleListItem[]>([]);
  public readonly roleCount = computed(() => this.roles().length);

  // Total number of permissions in the registry — the denominator for each
  // role's permission meter. Sourced from the one role/permission endpoint that
  // exists today.
  private readonly totalPermissions = signal(0);

  constructor() {
    this.loadPermissionTotal();
  }

  /** Number of filled segments (out of 10) for a role's permission meter. */
  public meter(role: RoleListItem): (string | null)[] {
    const total = this.totalPermissions() || 1;
    const granted = role.appCount + role.groupCount;
    const filled = Math.round((granted / total) * 10);
    const fillClass = role.appCount >= role.groupCount ? "app" : "group";
    return Array.from({ length: 10 }, (_, i) => (i < filled ? fillClass : null));
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

  // TODO: open the create-role flow once it exists (separate slice).
  public addRole(): void {}

  // TODO: open the edit/view-role flow once it exists (separate slice).
  public editRole(_role: RoleListItem): void {}

  // TODO: open the per-role actions menu once those actions exist.
  public openRoleMenu(_role: RoleListItem): void {}

  private loadPermissionTotal(): void {
    this.permissionService
      .getPermissions()
      .pipe(
        take(1),
        takeUntilDestroyed(),
        tap((permissions) => this.totalPermissions.set(permissions.length)),
      )
      .subscribe();
  }
}
