import { provideHttpClient, withInterceptorsFromDi } from "@angular/common/http";
import { provideHttpClientTesting } from "@angular/common/http/testing";
import {
  Component,
  CUSTOM_ELEMENTS_SCHEMA,
  Input,
  provideZonelessChangeDetection,
} from "@angular/core";
import { ComponentFixture, TestBed } from "@angular/core/testing";
import { ReactiveFormsModule } from "@angular/forms";
import {
  ActivatedRoute,
  convertToParamMap,
  Router,
  RouterModule,
} from "@angular/router";
import { of, throwError } from "rxjs";
import {
  ApiModule,
  PermissionDescriptor,
  PermissionService,
  Role,
  RoleService,
} from "../../open-api";
import { PipesModule } from "../../pipes";
import { SnackbarService } from "../../services/snackbar.service";
import { RoleFormComponent } from "./role-form.component";

const APP_DESCRIPTORS: PermissionDescriptor[] = [
  { key: "app.users.create", label: "Create users", description: "", category: "Users", scope: "APP" },
  { key: "app.users.read", label: "Read users", description: "", category: "Users", scope: "APP" },
  { key: "app.users.update", label: "Update users", description: "", category: "Users", scope: "APP" },
  { key: "app.users.delete", label: "Delete users", description: "", category: "Users", scope: "APP" },
  { key: "app.api-keys.create", label: "Create API keys", description: "", category: "API Keys", scope: "APP" },
  { key: "app.api-keys.read", label: "Read API keys", description: "", category: "API Keys", scope: "APP" },
  { key: "app.system-emails.read", label: "Read system emails", description: "", category: "System Emails", scope: "APP" },
  { key: "app.categories.read", label: "Read categories", description: "", category: "Categories", scope: "APP" },
  { key: "app.system-settings.restart-task-server", label: "Restart task server", description: "", category: "System Settings", scope: "APP" },
];

const GROUP_DESCRIPTORS: PermissionDescriptor[] = [
  { key: "group.view", label: "View group", description: "", category: "Group", scope: "GROUP" },
  { key: "group.update", label: "Update group", description: "", category: "Group", scope: "GROUP" },
  { key: "group.receipts.create", label: "Create receipts", description: "", category: "Receipts", scope: "GROUP" },
  { key: "group.receipts.read", label: "Read receipts", description: "", category: "Receipts", scope: "GROUP" },
  { key: "group.dashboards.read", label: "Read dashboards", description: "", category: "Dashboards", scope: "GROUP" },
  { key: "group.widgets.read", label: "Read widgets", description: "", category: "Widgets", scope: "GROUP" },
];

// Stubbed because its `onlyIcon` input trips Angular's on*-event security check
// when treated as an unknown element under CUSTOM_ELEMENTS_SCHEMA.
@Component({ selector: "app-submit-button", template: "", standalone: false })
class StubSubmitButtonComponent {
  @Input() onlyIcon = false;
  @Input() buttonText = "";
}

const ALL_DESCRIPTORS: PermissionDescriptor[] = [...APP_DESCRIPTORS, ...GROUP_DESCRIPTORS];

const APP_KEYS = APP_DESCRIPTORS.map((d) => d.key);

interface SetupResult {
  component: RoleFormComponent;
  fixture: ComponentFixture<RoleFormComponent>;
  createRoleSpy: jest.SpyInstance;
  updateRoleSpy: jest.SpyInstance;
  navigateSpy: jest.SpyInstance;
  snackbar: { success: jest.Mock; error: jest.Mock };
}

interface SetupOptions {
  /** When set, the form opens in edit/view mode for this route id. */
  routeId?: string;
  routeScope?: string;
  /** Roles returned by getRoles() when loading an existing role. */
  roles?: Role[];
}

async function setup(
  descriptors: PermissionDescriptor[] | "error" = ALL_DESCRIPTORS,
  options: SetupOptions = {},
): Promise<SetupResult> {
  TestBed.resetTestingModule();
  const snackbar = { success: jest.fn(), error: jest.fn() };

  const params: Record<string, string> = {};
  if (options.routeId) {
    params["id"] = options.routeId;
  }
  const query: Record<string, string> = {};
  if (options.routeScope) {
    query["scope"] = options.routeScope;
  }
  const activatedRoute = {
    snapshot: {
      paramMap: convertToParamMap(params),
      queryParamMap: convertToParamMap(query),
    },
  };

  await TestBed.configureTestingModule({
    declarations: [RoleFormComponent, StubSubmitButtonComponent],
    schemas: [CUSTOM_ELEMENTS_SCHEMA],
    imports: [ApiModule, PipesModule, ReactiveFormsModule, RouterModule.forRoot([])],
    providers: [
      provideZonelessChangeDetection(),
      provideHttpClient(withInterceptorsFromDi()),
      provideHttpClientTesting(),
      { provide: SnackbarService, useValue: snackbar },
      { provide: ActivatedRoute, useValue: activatedRoute },
    ],
  }).compileComponents();

  const permissionService = TestBed.inject(PermissionService);
  jest
    .spyOn(permissionService, "getPermissions")
    .mockReturnValue(
      descriptors === "error"
        ? (throwError(() => new Error("boom")) as any)
        : (of(descriptors) as any),
    );

  const createdRole: Role = {
    id: 1,
    name: "Created",
    scope: "APP",
    isSystem: false,
    permissions: [],
  };
  const roleService = TestBed.inject(RoleService);
  const createRoleSpy = jest
    .spyOn(roleService, "createRole")
    .mockReturnValue(of(createdRole) as any);
  const updateRoleSpy = jest
    .spyOn(roleService, "updateRole")
    .mockImplementation((_id: number, command: any) =>
      of({ id: 1, isSystem: false, ...command } as Role) as any,
    );
  jest
    .spyOn(roleService, "getRoles")
    .mockReturnValue(of(options.roles ?? []) as any);

  const router = TestBed.inject(Router);
  const navigateSpy = jest.spyOn(router, "navigate").mockResolvedValue(true);

  const fixture = TestBed.createComponent(RoleFormComponent);
  const component = fixture.componentInstance;
  await fixture.whenStable();
  return { component, fixture, createRoleSpy, updateRoleSpy, navigateSpy, snackbar };
}

const EDITABLE_APP_ROLE: Role = {
  id: 7,
  name: "App Manager",
  description: "Manages the app",
  scope: "APP",
  isSystem: false,
  permissions: ["app.users.create", "app.users.read"],
};

const SYSTEM_APP_ROLE: Role = {
  id: 9,
  name: "Administrator",
  description: "Built in",
  scope: "APP",
  isSystem: true,
  permissions: ["app.users.create"],
};

describe("RoleFormComponent", () => {
  it("should create", async () => {
    const { component } = await setup();
    expect(component).toBeTruthy();
  });

  it("renders the two role-type cards", async () => {
    const { fixture } = await setup();
    const cards = fixture.nativeElement.querySelectorAll(".rw-typecard");
    expect(cards.length).toBe(2);
  });

  it("defaults the type to app", async () => {
    const { component } = await setup();
    expect(component.type()).toBe("app");
    expect(component.typeMeta().scope).toBe("APP");
  });

  it("resets granted and preset when switching to group", async () => {
    const { component } = await setup();
    component.pickPreset(component.presets()[0]); // Administrator -> grants
    expect(component.granted().size).toBeGreaterThan(0);

    component.pickType("group");

    expect(component.type()).toBe("group");
    expect(component.granted().size).toBe(0);
    expect(component.selectedPreset()).toBe("custom");
  });

  it("fills granted with all app keys when picking Administrator", async () => {
    const { component } = await setup();
    const admin = component.presets().find((p) => p.id === "admin")!;
    component.pickPreset(admin);

    expect([...component.granted()].sort()).toEqual([...APP_KEYS].sort());
    expect(component.selectedPreset()).toBe("admin");
  });

  it("limits the User Manager preset to its resources", async () => {
    const { component } = await setup();
    const userMgr = component.presets().find((p) => p.id === "user-mgr")!;
    component.pickPreset(userMgr);

    expect([...component.granted()].sort()).toEqual(
      [
        "app.api-keys.create",
        "app.api-keys.read",
        "app.system-emails.read",
        "app.users.create",
        "app.users.delete",
        "app.users.read",
        "app.users.update",
      ].sort(),
    );
  });

  it("flips selectedPreset to custom when toggling a permission", async () => {
    const { component } = await setup();
    const admin = component.presets().find((p) => p.id === "admin")!;
    component.pickPreset(admin);
    expect(component.selectedPreset()).toBe("admin");

    component.toggle("app.users.create");

    expect(component.selectedPreset()).toBe("custom");
    expect(component.granted().has("app.users.create")).toBe(false);
  });

  it("groups permissions by resource into accordions", async () => {
    const { component } = await setup();
    const resourceKeys = component.resourceGroups().map((g) => g.resourceKey);

    expect(resourceKeys).toEqual([
      "app.users",
      "app.api-keys",
      "app.system-emails",
      "app.categories",
      "app.system-settings",
    ]);
  });

  it("groups group-scope permissions including the bare group resource", async () => {
    const { component } = await setup();
    component.pickType("group");

    const resourceKeys = component.resourceGroups().map((g) => g.resourceKey);
    expect(resourceKeys).toContain("group");
    expect(resourceKeys).toContain("group.receipts");
    expect(resourceKeys).toContain("group.dashboards");
  });

  it("reflects granted.size in the summary count", async () => {
    const { component, fixture } = await setup();
    component.toggle("app.users.create");
    component.toggle("app.users.read");
    await fixture.whenStable();

    const big = fixture.nativeElement.querySelector(".rw-summary-stat .big");
    expect(big.textContent.trim()).toBe("2");
    expect(component.granted().size).toBe(2);
  });

  it("assembles the correct payload for an app role", async () => {
    const { component } = await setup();
    component.form.controls.name.setValue("App Admin");
    component.form.controls.description.setValue("Runs the app");
    component.toggle("app.users.create");

    component.submit();

    expect(component.lastPayload).toEqual({
      name: "App Admin",
      description: "Runs the app",
      scope: "APP",
      permissions: ["app.users.create"],
    });
  });

  it("assembles a GROUP scope payload after switching type", async () => {
    const { component } = await setup();
    component.pickType("group");
    component.form.controls.name.setValue("Group Manager");
    component.toggle("group.view");

    component.submit();

    expect(component.lastPayload?.scope).toBe("GROUP");
    expect(component.lastPayload?.permissions).toEqual(["group.view"]);
  });

  it("calls createRole, notifies, and navigates back to the list on submit", async () => {
    const { component, createRoleSpy, navigateSpy, snackbar } = await setup();
    component.form.controls.name.setValue("App Admin");
    component.toggle("app.users.create");

    component.submit();

    expect(createRoleSpy).toHaveBeenCalledWith({
      name: "App Admin",
      description: "",
      scope: "APP",
      permissions: ["app.users.create"],
    });
    expect(snackbar.success).toHaveBeenCalled();
    expect(navigateSpy).toHaveBeenCalledWith(["/roles"]);
  });

  it("does not submit when the form is invalid", async () => {
    const { component, createRoleSpy } = await setup();
    component.submit();
    expect(component.lastPayload).toBeNull();
    expect(createRoleSpy).not.toHaveBeenCalled();
  });

  it("does not crash when the permission registry errors", async () => {
    const { component } = await setup("error");
    expect(component).toBeTruthy();
    expect(component.resourceGroups()).toEqual([]);
    expect(component.scopeTotal()).toBe(0);
  });

  describe("edit mode", () => {
    it("loads the role into the form", async () => {
      const { component } = await setup(ALL_DESCRIPTORS, {
        routeId: "7",
        routeScope: "app",
        roles: [EDITABLE_APP_ROLE],
      });

      expect(component.isEditMode()).toBe(true);
      expect(component.form.controls.name.value).toBe("App Manager");
      expect(component.form.controls.description.value).toBe("Manages the app");
      expect(component.type()).toBe("app");
      expect([...component.granted()].sort()).toEqual(
        ["app.users.create", "app.users.read"].sort(),
      );
    });

    it("locks the role type so it cannot be switched", async () => {
      const { component } = await setup(ALL_DESCRIPTORS, {
        routeId: "7",
        routeScope: "app",
        roles: [EDITABLE_APP_ROLE],
      });

      expect(component.isTypeLocked()).toBe(true);
      component.pickType("group");
      expect(component.type()).toBe("app");
    });

    it("calls updateRole and navigates on submit", async () => {
      const { component, updateRoleSpy, createRoleSpy, navigateSpy, snackbar } =
        await setup(ALL_DESCRIPTORS, {
          routeId: "7",
          routeScope: "app",
          roles: [EDITABLE_APP_ROLE],
        });

      component.form.controls.name.setValue("Renamed");
      component.submit();

      expect(createRoleSpy).not.toHaveBeenCalled();
      expect(updateRoleSpy).toHaveBeenCalledWith(7, {
        name: "Renamed",
        description: "Manages the app",
        scope: "APP",
        permissions: ["app.users.create", "app.users.read"],
      });
      expect(snackbar.success).toHaveBeenCalled();
      expect(navigateSpy).toHaveBeenCalledWith(["/roles"]);
    });

    it("disambiguates overlapping ids using the scope query param", async () => {
      const groupRoleSameId: Role = {
        id: 7,
        name: "Group Role 7",
        scope: "GROUP",
        isSystem: false,
        permissions: ["group.view"],
      };
      const { component } = await setup(ALL_DESCRIPTORS, {
        routeId: "7",
        routeScope: "group",
        roles: [EDITABLE_APP_ROLE, groupRoleSameId],
      });

      expect(component.type()).toBe("group");
      expect(component.form.controls.name.value).toBe("Group Role 7");
    });

    it("redirects to the list when the role is not found", async () => {
      const { navigateSpy, snackbar } = await setup(ALL_DESCRIPTORS, {
        routeId: "404",
        routeScope: "app",
        roles: [EDITABLE_APP_ROLE],
      });

      expect(snackbar.error).toHaveBeenCalled();
      expect(navigateSpy).toHaveBeenCalledWith(["/roles"]);
    });
  });

  describe("view mode (system roles)", () => {
    it("opens system roles read-only", async () => {
      const { component } = await setup(ALL_DESCRIPTORS, {
        routeId: "9",
        routeScope: "app",
        roles: [SYSTEM_APP_ROLE],
      });

      expect(component.isViewMode()).toBe(true);
      expect(component.isTypeLocked()).toBe(true);
      expect(component.form.disabled).toBe(true);
    });

    it("ignores permission edits", async () => {
      const { component } = await setup(ALL_DESCRIPTORS, {
        routeId: "9",
        routeScope: "app",
        roles: [SYSTEM_APP_ROLE],
      });

      const before = component.granted().size;
      component.toggle("app.users.read");
      expect(component.granted().size).toBe(before);
    });

    it("does not call updateRole on submit", async () => {
      const { component, updateRoleSpy } = await setup(ALL_DESCRIPTORS, {
        routeId: "9",
        routeScope: "app",
        roles: [SYSTEM_APP_ROLE],
      });

      component.submit();
      expect(updateRoleSpy).not.toHaveBeenCalled();
    });
  });
});
