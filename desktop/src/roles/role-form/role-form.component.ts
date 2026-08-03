import { Component, DestroyRef, Signal, computed, effect, signal } from "@angular/core";
import { takeUntilDestroyed } from "@angular/core/rxjs-interop";
import { FormArray, FormControl, FormGroup, Validators } from "@angular/forms";
import { ActivatedRoute, Router } from "@angular/router";
import { Store } from "@ngxs/store";
import { catchError, EMPTY, take, tap } from "rxjs";
import { FormMode } from "../../enums/form-mode.enum";
import { SnackbarService } from "../../services/snackbar.service";
import { BreadcrumbItem } from "../../shared-ui/breadcrumb/breadcrumb-item.interface";
import { UserState } from "../../store";
import {
  Category,
  CategoryService,
  Permission,
  PermissionDescriptor,
  PermissionScope,
  PermissionService,
  ReportService,
  ReportTemplateOption,
  RoleService,
  Tag,
  TagService,
  UpsertRoleCommand,
  User,
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

// An option in the paid-by visibility picker: either a real user (id > 0) or the
// pinned "their own receipts" sentinel (id === OWN_PAID_RECEIPTS_OPTION_ID).
export interface PaidByGrantOption {
  id: number;
  displayName: string;
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
    // Group-role-only flag: holders can see, and be seen by, all members of an
    // isolated group. Ignored for app roles (serialized only when showGrants()).
    seesAllMembers: new FormControl<boolean>(false, { nonNullable: true }),
    // App-role-only flag: users created with this role skip the personal
    // "My Receipts" group. Ignored for group roles (serialized only when
    // showAppOptions()).
    skipDefaultGroupCreation: new FormControl<boolean>(false, { nonNullable: true }),
  });

  // ----- State signals -----
  public readonly mode = signal<FormMode>(FormMode.add);
  private readonly roleId = signal<number | null>(null);
  public readonly type = signal<RoleType>("app");
  public readonly granted = signal<Set<string>>(new Set<string>());
  public readonly selectedPreset = signal<string>(CUSTOM_PRESET_ID);
  public readonly openPanels = signal<Record<string, boolean>>({});

  // ----- Category/tag grants (group roles only) -----
  // The full pool an admin picks grants from, and the selected grants. The
  // autocomplete drives these as FormArrays of the selected category/tag objects.
  public readonly categoryPool = signal<Category[]>([]);
  public readonly tagPool = signal<Tag[]>([]);
  public readonly grantedCategories = new FormArray<FormControl<Category>>([]);
  public readonly grantedTags = new FormArray<FormControl<Tag>>([]);

  // Grant ids loaded for an existing role, resolved to pool objects by an effect
  // once the pool arrives (pool and role load independently/asynchronously).
  private readonly pendingCategoryGrantIds = signal<number[]>([]);
  private readonly pendingTagGrantIds = signal<number[]>([]);

  // ----- Paid-by visibility grants (group roles only) -----
  // The sentinel id of the pinned "their own receipts" option. Real user ids are
  // positive, so a negative sentinel never collides and is filtered out of the
  // serialized paidByUserGrants (it maps to includeOwnPaidReceipts instead).
  public static readonly OWN_PAID_RECEIPTS_OPTION_ID = -1;
  private readonly ownPaidReceiptsOption: PaidByGrantOption = {
    id: RoleFormComponent.OWN_PAID_RECEIPTS_OPTION_ID,
    displayName: "Their own receipts",
  };

  // The full user pool, read reactively so the picker repopulates if AppData lands
  // after the form mounts (a snapshot inside the computed would cache once and drop
  // granted users on edit). Assigned in the constructor — a field initializer can't
  // reference this.store before it is set.
  private usersSignal!: Signal<User[]>;

  // The picker options: the pinned "their own" sentinel first, then every user.
  // The role editor is admin-only and global (no group context), so it lists all
  // users — sourced from AppData like every other user picker.
  public readonly paidByOptions = computed<PaidByGrantOption[]>(() => [
    this.ownPaidReceiptsOption,
    ...this.usersSignal().map((user) => ({
      id: user.id,
      displayName: user.displayName?.length ? user.displayName : user.username,
    })),
  ]);
  public readonly grantedPaidByUsers = new FormArray<FormControl<PaidByGrantOption>>([]);
  private readonly pendingPaidByUserGrantIds = signal<number[]>([]);
  private readonly pendingIncludeOwnPaidReceipts = signal<boolean>(false);

  // ----- Report template access matrix (group roles only) -----
  // The scopable per-template actions (create is excluded — it makes a new
  // template, so there is nothing to scope). The label is the admin-facing verb.
  public readonly reportActions: { key: string; label: string }[] = [
    { key: "read", label: "View" },
    { key: "generate", label: "Generate" },
    { key: "update", label: "Edit" },
    { key: "delete", label: "Delete" },
    { key: "duplicate", label: "Duplicate" },
  ];
  // The full template pool (rows of the matrix) and the granted actions per
  // template. An empty grant map means unrestricted (every template the role's
  // group access reaches); a template maps to the subset of actions the role may
  // perform on it. The map is replaced immutably so zoneless CD picks up changes.
  public readonly reportTemplateOptions = signal<ReportTemplateOption[]>([]);
  public readonly reportTemplateGrants = signal<Map<number, Set<string>>>(new Map());

  /** Grants are a group-role concept; hidden for app roles. */
  public readonly showGrants = computed<boolean>(() => this.type() === "group");

  /** Personal-group creation is an app-role concept; hidden for group roles. */
  public readonly showAppOptions = computed<boolean>(() => this.type() === "app");

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
    private readonly categoryService: CategoryService,
    private readonly tagService: TagService,
    private readonly reportService: ReportService,
    private readonly router: Router,
    private readonly route: ActivatedRoute,
    private readonly snackbar: SnackbarService,
    private readonly destroyRef: DestroyRef,
    private readonly store: Store,
  ) {
    // Reactive user pool for the paid-by picker (see usersSignal above).
    this.usersSignal = this.store.selectSignal(UserState.users);

    this.permissionService
      .getPermissions()
      .pipe(
        take(1),
        catchError(() => EMPTY),
        takeUntilDestroyed(),
      )
      .subscribe((descriptors) => this.registry.set(descriptors));

    // The role editor is admin-only (app.roles.*), so it may read the full
    // category/tag pool to choose grants from.
    this.categoryService
      .getAllCategories()
      .pipe(take(1), catchError(() => EMPTY), takeUntilDestroyed())
      .subscribe((categories) => this.categoryPool.set(categories));
    this.tagService
      .getAllTags()
      .pipe(take(1), catchError(() => EMPTY), takeUntilDestroyed())
      .subscribe((tags) => this.tagPool.set(tags));

    // Report templates for the access matrix. Gated on app.roles.read (the role
    // editor's own gate), so the admin may list them even without report access.
    this.reportService
      .getReportTemplateOptions()
      .pipe(take(1), catchError(() => EMPTY), takeUntilDestroyed())
      .subscribe((options) => this.reportTemplateOptions.set(options));

    // A view-mode role is read-only; keep the reactive form in sync with that.
    effect(() => {
      if (this.isViewMode()) {
        this.form.disable({ emitEvent: false });
      } else {
        this.form.enable({ emitEvent: false });
      }
    });

    // Resolve a loaded role's grant ids to pool objects once the pool is
    // available. Pool and grant ids are stable after load, so this does not
    // clobber subsequent manual edits.
    effect(() => {
      const pool = this.categoryPool();
      const ids = this.pendingCategoryGrantIds();
      this.setGrantArray(this.grantedCategories, pool.filter((category) => category.id !== undefined && ids.includes(category.id)));
    });
    effect(() => {
      const pool = this.tagPool();
      const ids = this.pendingTagGrantIds();
      this.setGrantArray(this.grantedTags, pool.filter((tag) => tag.id !== undefined && ids.includes(tag.id)));
    });

    // Resolve a loaded role's paid-by config to picker options, prepending the
    // "their own receipts" sentinel when include-own is set. Filtering the shared
    // paidByOptions list keeps object references stable so the autocomplete
    // correctly excludes already-selected options.
    effect(() => {
      const options = this.paidByOptions();
      const ids = this.pendingPaidByUserGrantIds();
      const includeOwn = this.pendingIncludeOwnPaidReceipts();
      const selected = options.filter(
        (option) =>
          (option.id === RoleFormComponent.OWN_PAID_RECEIPTS_OPTION_ID && includeOwn) ||
          (option.id !== RoleFormComponent.OWN_PAID_RECEIPTS_OPTION_ID && ids.includes(option.id)),
      );
      this.setGrantArray(this.grantedPaidByUsers, selected);
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
        this.form.patchValue({
          name: role.name,
          description: role.description ?? "",
          seesAllMembers: role.seesAllMembers ?? false,
          skipDefaultGroupCreation: role.skipDefaultGroupCreation ?? false,
        });
        this.type.set(role.scope === PermissionScope.Group ? "group" : "app");
        this.granted.set(new Set(role.permissions));
        this.pendingCategoryGrantIds.set(role.categoryGrants ?? []);
        this.pendingTagGrantIds.set(role.tagGrants ?? []);
        this.pendingPaidByUserGrantIds.set(role.paidByUserGrants ?? []);
        this.pendingIncludeOwnPaidReceipts.set(role.includeOwnPaidReceipts ?? false);

        // The matrix grants are self-contained ({templateId, permissions}), so they
        // hydrate directly — no pool resolution needed (the options only supply names).
        const grants = new Map<number, Set<string>>();
        for (const grant of role.reportTemplateGrants ?? []) {
          grants.set(grant.reportTemplateId, new Set(grant.permissions));
        }
        this.reportTemplateGrants.set(grants);
      });
  }

  // Replace a grant FormArray's contents with the given option objects, without
  // emitting (the autocomplete reads the value, not change events, on load).
  private setGrantArray<T>(array: FormArray, values: T[]): void {
    array.clear({ emitEvent: false });
    for (const value of values) {
      array.push(new FormControl(value), { emitEvent: false });
    }
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
    // Grants only apply to group roles — reset them on a type switch.
    this.pendingCategoryGrantIds.set([]);
    this.pendingTagGrantIds.set([]);
    this.pendingPaidByUserGrantIds.set([]);
    this.pendingIncludeOwnPaidReceipts.set(false);
    this.setGrantArray(this.grantedCategories, []);
    this.setGrantArray(this.grantedTags, []);
    this.setGrantArray(this.grantedPaidByUsers, []);
    this.reportTemplateGrants.set(new Map());
    this.form.controls.seesAllMembers.setValue(false);
    // Personal-group creation only applies to app roles — reset it too.
    this.form.controls.skipDefaultGroupCreation.setValue(false);
  }

  // ----- Report template matrix helpers -----
  public isTemplateActionGranted(templateId: number, action: string): boolean {
    return this.reportTemplateGrants().get(templateId)?.has(action) ?? false;
  }

  public toggleTemplateAction(templateId: number, action: string): void {
    if (this.isViewMode()) {
      return;
    }
    const next = new Map(this.reportTemplateGrants());
    const actions = new Set(next.get(templateId) ?? []);
    if (actions.has(action)) {
      actions.delete(action);
    } else {
      actions.add(action);
    }
    // Drop the template entirely once it grants nothing, so an all-empty matrix
    // serializes as "unrestricted".
    if (actions.size === 0) {
      next.delete(templateId);
    } else {
      next.set(templateId, actions);
    }
    this.reportTemplateGrants.set(next);
  }

  public isTemplateAllGranted(templateId: number): boolean {
    const actions = this.reportTemplateGrants().get(templateId);
    return !!actions && actions.size === this.reportActions.length;
  }

  public toggleTemplateAll(templateId: number, on: boolean): void {
    if (this.isViewMode()) {
      return;
    }
    const next = new Map(this.reportTemplateGrants());
    if (on) {
      next.set(templateId, new Set(this.reportActions.map((action) => action.key)));
    } else {
      next.delete(templateId);
    }
    this.reportTemplateGrants.set(next);
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

    // Category/tag/paid-by grants are only valid on group roles.
    if (this.showGrants()) {
      payload.categoryGrants = (this.grantedCategories.value as Category[])
        .map((category) => category.id)
        .filter((id): id is number => id !== undefined);
      payload.tagGrants = (this.grantedTags.value as Tag[])
        .map((tag) => tag.id)
        .filter((id): id is number => id !== undefined);

      // Split the picker's selections into the relative "their own" toggle and
      // the absolute user-id grants (the sentinel never goes to the backend).
      const paidBySelections = this.grantedPaidByUsers.value as PaidByGrantOption[];
      payload.includeOwnPaidReceipts = paidBySelections.some(
        (option) => option.id === RoleFormComponent.OWN_PAID_RECEIPTS_OPTION_ID,
      );
      payload.paidByUserGrants = paidBySelections
        .map((option) => option.id)
        .filter((id) => id !== RoleFormComponent.OWN_PAID_RECEIPTS_OPTION_ID);

      // An empty map serializes as [] → unrestricted; otherwise one entry per
      // template carrying its granted actions.
      payload.reportTemplateGrants = [...this.reportTemplateGrants().entries()].map(
        ([reportTemplateId, actions]) => ({ reportTemplateId, permissions: [...actions] }),
      );

      // Supervisor exemption: holders see, and are seen by, all members of an
      // isolated group. Only meaningful on group roles.
      payload.seesAllMembers = this.form.controls.seesAllMembers.value;
    } else {
      // Personal-group creation is an app-role concept: users created with this
      // role skip the automatic "My Receipts" group.
      payload.skipDefaultGroupCreation =
        this.form.controls.skipDefaultGroupCreation.value;
    }

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
