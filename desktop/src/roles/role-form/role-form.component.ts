import { Component, computed, signal } from "@angular/core";
import { takeUntilDestroyed } from "@angular/core/rxjs-interop";
import { FormControl, FormGroup, Validators } from "@angular/forms";
import { Router } from "@angular/router";
import { catchError, EMPTY, take, tap } from "rxjs";
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
  public readonly breadcrumbs: BreadcrumbItem[] = [
    { label: "Admin" },
    { label: "Roles", routerLink: "/roles" },
    { label: "New Role" },
  ];

  public readonly roleTypes = ROLE_TYPES;

  public readonly form = new FormGroup({
    name: new FormControl<string>("", { nonNullable: true, validators: [Validators.required] }),
    description: new FormControl<string>("", { nonNullable: true }),
  });

  // ----- State signals -----
  public readonly type = signal<RoleType>("app");
  public readonly granted = signal<Set<string>>(new Set<string>());
  public readonly selectedPreset = signal<string>(CUSTOM_PRESET_ID);
  public readonly openPanels = signal<Record<string, boolean>>({});

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
    private readonly snackbar: SnackbarService,
  ) {
    this.permissionService
      .getPermissions()
      .pipe(
        take(1),
        catchError(() => EMPTY),
        takeUntilDestroyed(),
      )
      .subscribe((descriptors) => this.registry.set(descriptors));
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
    if (type === this.type()) {
      return;
    }
    this.type.set(type);
    this.selectedPreset.set(CUSTOM_PRESET_ID);
    this.granted.set(new Set<string>());
    this.openPanels.set({});
  }

  public pickPreset(preset: RolePreset): void {
    this.selectedPreset.set(preset.id);
    this.granted.set(preset.resolve(this.scopeKeys()));
  }

  public toggle(key: string): void {
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

    this.roleService
      .createRole(payload)
      .pipe(
        take(1),
        tap((role) => {
          this.snackbar.success(`Role "${role.name}" created`);
          this.router.navigate(["/roles"]);
        }),
      )
      .subscribe();
  }

  private markCustom(): void {
    this.selectedPreset.set(CUSTOM_PRESET_ID);
  }
}
