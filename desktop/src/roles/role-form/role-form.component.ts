import { Component, DestroyRef, computed, effect, signal } from "@angular/core";
import { takeUntilDestroyed } from "@angular/core/rxjs-interop";
import { FormControl, FormGroup, Validators } from "@angular/forms";
import { ActivatedRoute, Router } from "@angular/router";
import { catchError, EMPTY, take, tap } from "rxjs";
import { FormMode } from "../../enums/form-mode.enum";
import { SnackbarService } from "../../services/snackbar.service";
import { BreadcrumbItem } from "../../shared-ui/breadcrumb/breadcrumb-item.interface";
import {
  Permission,
  PermissionDescriptor,
  PermissionScope,
  PermissionService,
  RoleService,
  UpsertRoleCommand,
} from "../../open-api";
import {
  CUSTOM_PRESET_ID,
  friendlyActionLabel,
  friendlyResourceName,
  iconForResource,
  presetsForType,
  resourceKeyOf,
  RolePreset,
  RoleType,
  ROLE_TYPES,
  RoleTypeMeta,
} from "../role-presets";

export interface RolePermissionRow {
  key: string;
  label: string;
  description: string;
  actionLabel: string;
}

export interface RoleResourceGroup {
  resourceKey: string;
  name: string;
  icon: string;
  rows: RolePermissionRow[];
}


@Component({
  selector: "app-role-form",
  templateUrl: "./role-form.component.html",
  styleUrls: ["./role-form.component.scss"],
  standalone: false,
})
export class RoleFormComponent {
  public readonly roleTypes = ROLE_TYPES;
  public readonly formMode = FormMode;

  public readonly form = new FormGroup({
    name: new FormControl<string>("", { nonNullable: true, validators: [Validators.required] }),
    description: new FormControl<string>("", { nonNullable: true }),
  });

  // ----- State signals -----
  public readonly mode = signal<FormMode>(FormMode.add);
  private readonly roleId = signal<number | null>(null);
  public readonly type = signal<RoleType>("app");
  public readonly granted = signal<Set<string>>(new Set<string>());
  public readonly selectedPreset = signal<string>(CUSTOM_PRESET_ID);
  public readonly openPanels = signal<Record<string, boolean>>({});

  // ----- Mode-derived state -----
  public readonly isEditMode = computed<boolean>(() => this.mode() === FormMode.edit);
  public readonly isViewMode = computed<boolean>(() => this.mode() === FormMode.view);
  /** The role type is fixed once a role exists — locked in both edit and view. */
  public readonly isTypeLocked = computed<boolean>(() => this.mode() !== FormMode.add);

  public readonly title = computed<string>(() => {
    if (this.isViewMode()) {
      return "View a Role";
    }
    return this.isEditMode() ? "Edit a Role" : "Create a Role";
  });

  public readonly breadcrumbs = computed<BreadcrumbItem[]>(() => {
    const last = this.isViewMode() ? "View Role" : this.isEditMode() ? "Edit Role" : "New Role";
    return [
      { label: "Admin" },
      { label: "Roles", routerLink: "/roles" },
      { label: last },
    ];
  });

  /** Last assembled payload — exposed for tests. */
  public lastPayload: UpsertRoleCommand | null = null;

  private readonly registry = signal<PermissionDescriptor[]>([]);

  // ----- Derived state -----
  public readonly typeMeta = computed<RoleTypeMeta>(
    () => this.roleTypes.find((t) => t.id === this.type()) ?? this.roleTypes[0],
  );

  public readonly presets = computed<RolePreset[]>(() => presetsForType(this.type()));

  public readonly activePreset = computed<RolePreset>(() => {
    const presets = this.presets();
    return (
      presets.find((p) => p.id === this.selectedPreset()) ??
      presets.find((p) => p.id === CUSTOM_PRESET_ID) ??
      presets[0]
    );
  });

  /** Registry descriptors for the active scope. */
  public readonly scopeDescriptors = computed<PermissionDescriptor[]>(() => {
    const scope = this.typeMeta().scope;
    return this.registry().filter((d) => d.scope === scope);
  });

  /** Full permission keys available in the active scope. */
  public readonly scopeKeys = computed<string[]>(() =>
    this.scopeDescriptors().map((d) => d.key),
  );

  public readonly scopeTotal = computed<number>(() => this.scopeKeys().length);

  /** Accordions grouped by resource (key minus last segment). */
  public readonly resourceGroups = computed<RoleResourceGroup[]>(() => {
    const groups = new Map<string, RoleResourceGroup>();

    for (const descriptor of this.scopeDescriptors()) {
      const resourceKey = resourceKeyOf(descriptor.key);
      let group = groups.get(resourceKey);
      if (!group) {
        group = {
          resourceKey,
          name: friendlyResourceName(resourceKey),
          icon: iconForResource(resourceKey),
          rows: [],
        };
        groups.set(resourceKey, group);
      }

      group.rows.push({
        key: descriptor.key,
        label: descriptor.label,
        description: descriptor.description,
        actionLabel: friendlyActionLabel(descriptor.key),
      });
    }

    return [...groups.values()];
  });

  public readonly scopeColor = computed<string>(() =>
    this.type() === "app" ? "#27b1ff" : "#8b5cf6",
  );

  constructor(
    private readonly permissionService: PermissionService,
    private readonly roleService: RoleService,
    private readonly router: Router,
    private readonly route: ActivatedRoute,
    private readonly snackbar: SnackbarService,
    private readonly destroyRef: DestroyRef,
  ) {
    this.permissionService
      .getPermissions()
      .pipe(
        take(1),
        catchError(() => EMPTY),
        takeUntilDestroyed(),
      )
      .subscribe((descriptors) => this.registry.set(descriptors));

    // A view-mode role is read-only; keep the reactive form in sync with that.
    effect(() => {
      if (this.isViewMode()) {
        this.form.disable({ emitEvent: false });
      } else {
        this.form.enable({ emitEvent: false });
      }
    });

    const id = this.route.snapshot.paramMap.get("id");
    if (id) {
      // Scope is required to identify a role: app and group roles have
      // independent id sequences, so an id alone is ambiguous.
      const scope = this.route.snapshot.queryParamMap.get("scope");
      if (!scope) {
        this.snackbar.error("Role type is required to edit a role");
        this.router.navigate(["/roles"]);
        return;
      }
      this.mode.set(FormMode.edit);
      this.roleId.set(Number(id));
      this.loadRole(id, scope);
    }
  }

  private loadRole(id: string, scope: string | null): void {
    this.roleService
      .getRoles()
      .pipe(
        take(1),
        catchError(() => {
          this.snackbar.error("Unable to load role");
          this.router.navigate(["/roles"]);
          return EMPTY;
        }),
        takeUntilDestroyed(this.destroyRef),
      )
      .subscribe((roles) => {
        const role = roles.find(
          (candidate) =>
            String(candidate.id) === id && this.scopeMatches(candidate.scope, scope),
        );

        if (!role) {
          this.snackbar.error("Role not found");
          this.router.navigate(["/roles"]);
          return;
        }

        // System roles are immutable — open them read-only.
        this.mode.set(role.isSystem ? FormMode.view : FormMode.edit);
        this.form.patchValue({ name: role.name, description: role.description ?? "" });
        this.type.set(role.scope === PermissionScope.Group ? "group" : "app");
        this.granted.set(new Set(role.permissions));
      });
  }

  private scopeMatches(scope: PermissionScope, scopeParam: string | null): boolean {
    if (!scopeParam) {
      return false;
    }
    const normalized = scope === PermissionScope.Group ? "group" : "app";
    return normalized === scopeParam;
  }

  // ----- Per-resource helpers -----
  public onCount(group: RoleResourceGroup): number {
    const granted = this.granted();
    return group.rows.filter((row) => granted.has(row.key)).length;
  }

  public isAllOn(group: RoleResourceGroup): boolean {
    return group.rows.length > 0 && this.onCount(group) === group.rows.length;
  }

  public isGranted(key: string): boolean {
    return this.granted().has(key);
  }

  public isPanelOpen(resourceKey: string): boolean {
    return !!this.openPanels()[resourceKey];
  }

  // ----- Actions -----
  public pickType(type: RoleType): void {
    // The type is fixed once the role exists.
    if (this.isTypeLocked() || type === this.type()) {
      return;
    }
    this.type.set(type);
    this.selectedPreset.set(CUSTOM_PRESET_ID);
    this.granted.set(new Set<string>());
    this.openPanels.set({});
  }

  public pickPreset(preset: RolePreset): void {
    if (this.isViewMode()) {
      return;
    }
    this.selectedPreset.set(preset.id);
    this.granted.set(preset.resolve(this.scopeKeys()));
  }

  public toggle(key: string): void {
    if (this.isViewMode()) {
      return;
    }
    const next = new Set(this.granted());
    if (next.has(key)) {
      next.delete(key);
    } else {
      next.add(key);
    }
    this.granted.set(next);
    this.markCustom();
  }

  public toggleAll(group: RoleResourceGroup, on: boolean): void {
    if (this.isViewMode()) {
      return;
    }
    const next = new Set(this.granted());
    for (const row of group.rows) {
      if (on) {
        next.add(row.key);
      } else {
        next.delete(row.key);
      }
    }
    this.granted.set(next);
    this.markCustom();
  }

  public togglePanel(resourceKey: string): void {
    this.openPanels.update((panels) => ({
      ...panels,
      [resourceKey]: !panels[resourceKey],
    }));
  }

  public submit(): void {
    if (this.isViewMode()) {
      return;
    }

    this.form.markAllAsTouched();
    if (this.form.invalid) {
      return;
    }

    const payload: UpsertRoleCommand = {
      name: this.form.controls.name.value,
      description: this.form.controls.description.value,
      scope: this.typeMeta().scope,
      permissions: [...this.granted()] as Permission[],
    };

    this.lastPayload = payload;

    const roleId = this.roleId();
    const request$ =
      this.isEditMode() && roleId !== null
        ? this.roleService.updateRole(roleId, payload)
        : this.roleService.createRole(payload);

    request$
      .pipe(
        take(1),
        tap((role) => {
          this.snackbar.success(
            `Role "${role.name}" ${this.isEditMode() ? "updated" : "created"}`,
          );
          this.router.navigate(["/roles"]);
        }),
      )
      .subscribe();
  }

  private markCustom(): void {
    this.selectedPreset.set(CUSTOM_PRESET_ID);
  }
}
