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
import { Router, RouterModule } from "@angular/router";
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
  navigateSpy: jest.SpyInstance;
  snackbar: { success: jest.Mock; error: jest.Mock };
}

async function setup(
  descriptors: PermissionDescriptor[] | "error" = ALL_DESCRIPTORS,
): Promise<SetupResult> {
  TestBed.resetTestingModule();
  const snackbar = { success: jest.fn(), error: jest.fn() };
  await TestBed.configureTestingModule({
    declarations: [RoleFormComponent, StubSubmitButtonComponent],
    schemas: [CUSTOM_ELEMENTS_SCHEMA],
    imports: [ApiModule, PipesModule, ReactiveFormsModule, RouterModule.forRoot([])],
    providers: [
      provideZonelessChangeDetection(),
      provideHttpClient(withInterceptorsFromDi()),
      provideHttpClientTesting(),
      { provide: SnackbarService, useValue: snackbar },
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

  const router = TestBed.inject(Router);
  const navigateSpy = jest.spyOn(router, "navigate").mockResolvedValue(true);

  const fixture = TestBed.createComponent(RoleFormComponent);
  const component = fixture.componentInstance;
  await fixture.whenStable();
  return { component, fixture, createRoleSpy, navigateSpy, snackbar };
}

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
});
