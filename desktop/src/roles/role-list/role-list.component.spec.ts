import { CommonModule } from "@angular/common";
import { provideHttpClient, withInterceptorsFromDi } from "@angular/common/http";
import { provideHttpClientTesting } from "@angular/common/http/testing";
import { CUSTOM_ELEMENTS_SCHEMA, provideZonelessChangeDetection } from "@angular/core";
import { ComponentFixture, TestBed } from "@angular/core/testing";
import { MatDialog } from "@angular/material/dialog";
import { Router, RouterModule } from "@angular/router";
import { of, throwError } from "rxjs";
import {
  ApiModule,
  PermissionDescriptor,
  PermissionScope,
  PermissionService,
  Role,
  RoleService,
} from "../../open-api";
import { SnackbarService } from "../../services";
import { RoleListComponent } from "./role-list.component";
import { RoleListItem, RoleScope } from "./role-list-item.interface";

function listItem(overrides: Partial<RoleListItem> = {}): RoleListItem {
  return {
    id: "1",
    name: "Some Role",
    description: "",
    scope: "app" as RoleScope,
    permissionCount: 0,
    userCount: 0,
    isDefault: false,
    isSystem: false,
    icon: "apps",
    iconColor: "",
    iconTint: "",
    ...overrides,
  };
}

const DESCRIPTORS: PermissionDescriptor[] = [
  { key: "app.users.create", label: "", description: "", category: "Users", scope: "APP" },
  { key: "app.users.read", label: "", description: "", category: "Users", scope: "APP" },
  { key: "group.view", label: "", description: "", category: "Group", scope: "GROUP" },
  { key: "group.receipts.read", label: "", description: "", category: "Receipts", scope: "GROUP" },
  { key: "group.receipts.create", label: "", description: "", category: "Receipts", scope: "GROUP" },
];

const ROLES: Role[] = [
  {
    id: 1,
    name: "Administrator",
    description: "Full access",
    scope: "APP",
    isDefault: true,
    isSystem: true,
    permissions: ["app.users.create", "app.users.read"],
  },
  {
    id: 2,
    name: "Group Manager",
    description: "Runs groups",
    scope: "GROUP",
    isDefault: true,
    isSystem: false,
    permissions: ["group.view", "group.receipts.read"],
  },
];

async function setup(
  roles: Role[] = ROLES,
  permissions: PermissionDescriptor[] | "error" = DESCRIPTORS,
  confirmResult = true,
): Promise<{
  component: RoleListComponent;
  fixture: ComponentFixture<RoleListComponent>;
  navigateSpy: jest.SpyInstance;
  roleService: RoleService;
  snackbar: SnackbarService;
  dialogRef: { componentInstance: any; afterClosed: jest.Mock };
}> {
  const dialogRef = {
    componentInstance: {} as any,
    afterClosed: jest.fn().mockReturnValue(of(confirmResult)),
  };
  const matDialogMock = { open: jest.fn().mockReturnValue(dialogRef) };
  const snackbarMock = { info: jest.fn(), success: jest.fn(), error: jest.fn() };

  TestBed.resetTestingModule();
  await TestBed.configureTestingModule({
    declarations: [RoleListComponent],
    schemas: [CUSTOM_ELEMENTS_SCHEMA],
    imports: [ApiModule, CommonModule, RouterModule.forRoot([])],
    providers: [
      provideZonelessChangeDetection(),
      provideHttpClient(withInterceptorsFromDi()),
      provideHttpClientTesting(),
      { provide: MatDialog, useValue: matDialogMock },
      { provide: SnackbarService, useValue: snackbarMock },
    ],
  }).compileComponents();

  const roleService = TestBed.inject(RoleService);
  jest.spyOn(roleService, "getRoles").mockReturnValue(of(roles) as any);
  jest.spyOn(roleService, "deleteRole").mockReturnValue(of({}) as any);
  jest.spyOn(roleService, "setDefaultRole").mockReturnValue(of({}) as any);
  jest
    .spyOn(TestBed.inject(PermissionService), "getPermissions")
    .mockReturnValue(
      permissions === "error"
        ? (throwError(() => new Error("boom")) as any)
        : (of(permissions) as any),
    );

  const router = TestBed.inject(Router);
  const navigateSpy = jest.spyOn(router, "navigate").mockResolvedValue(true);

  const fixture = TestBed.createComponent(RoleListComponent);
  const component = fixture.componentInstance;
  await fixture.whenStable();
  return {
    component,
    fixture,
    navigateSpy,
    roleService,
    snackbar: TestBed.inject(SnackbarService),
    dialogRef,
  };
}

describe("RoleListComponent", () => {
  it("creates", async () => {
    const { component } = await setup();
    expect(component).toBeTruthy();
  });

  it("renders the breadcrumb and the page title", async () => {
    const { fixture } = await setup();
    expect(fixture.nativeElement.querySelector("app-breadcrumb")).toBeTruthy();
    const header = fixture.nativeElement.querySelector("app-table-header");
    expect(header).toBeTruthy();
    expect(header?.getAttribute("headerText")).toBe("Roles");
  });

  it("shows the empty state and no table when there are no roles", async () => {
    const { component, fixture } = await setup([]);
    expect(component.roleCount()).toBe(0);
    expect(fixture.nativeElement.querySelector(".empty-state")).toBeTruthy();
    expect(fixture.nativeElement.querySelector("app-table")).toBeNull();
  });

  it("renders the table fed with a row per role when roles are present", async () => {
    const { component, fixture } = await setup();
    expect(component.roleCount()).toBe(2);
    expect(fixture.nativeElement.querySelector(".empty-state")).toBeNull();
    expect(fixture.nativeElement.querySelector("app-table")).toBeTruthy();
    expect(component.dataSource().data.length).toBe(2);
  });

  it("builds five table columns once the view initializes", async () => {
    const { component } = await setup();
    expect(component.columns().map((c) => c.matColumnDef)).toEqual([
      "role",
      "type",
      "permissions",
      "members",
      "actions",
    ]);
  });

  it("computes per-type counts for the filter tabs", async () => {
    const { component } = await setup();
    expect(component.counts()).toEqual({ all: 2, app: 1, group: 1 });
  });

  it("filters rows by the selected type", async () => {
    const { component } = await setup();
    component.setFilter("group");
    expect(component.filteredRoles().map((r) => r.name)).toEqual(["Group Manager"]);
    expect(component.dataSource().data.length).toBe(1);

    component.setFilter("app");
    expect(component.filteredRoles().map((r) => r.name)).toEqual(["Administrator"]);

    component.setFilter("all");
    expect(component.filteredRoles().length).toBe(2);
    expect(component.dataSource().data.length).toBe(2);
  });

  it("uses the role's own scope total as the meter denominator", async () => {
    const { component } = await setup();
    const appRole = component.roles().find((r) => r.scope === "app")!;
    const groupRole = component.roles().find((r) => r.scope === "group")!;

    expect(component.scopeTotal(appRole)).toBe(2); // two APP descriptors
    expect(component.scopeTotal(groupRole)).toBe(3); // three GROUP descriptors
    expect(component.meter(appRole).length).toBe(10);
    expect(component.meter(appRole).filter((s) => s === "app").length).toBe(10); // 2 of 2
  });

  it("navigates to the create page from Add Role", async () => {
    const { component, navigateSpy } = await setup();
    component.addRole();
    expect(navigateSpy).toHaveBeenCalledWith(["/roles/new"]);
  });

  it("navigates to the edit page with the role scope", async () => {
    const { component, navigateSpy } = await setup();
    const groupRole = component.roles().find((r) => r.scope === "group")!;

    component.editRole(groupRole);

    expect(navigateSpy).toHaveBeenCalledWith(["/roles", groupRole.id, "edit"], {
      queryParams: { scope: "group" },
    });
  });

  it("routes system roles through the same edit/view path", async () => {
    const { component, navigateSpy } = await setup();
    const systemRole = component.roles().find((r) => r.isSystem)!;

    component.editRole(systemRole);

    expect(navigateSpy).toHaveBeenCalledWith(["/roles", systemRole.id, "edit"], {
      queryParams: { scope: "app" },
    });
  });

  it("does not crash when the permission registry errors", async () => {
    const { component } = await setup(ROLES, "error");
    expect(component).toBeTruthy();
    expect(component.roleCount()).toBe(2);
  });

  it("populates the member count from the role's assigned count", async () => {
    const { component } = await setup([
      { ...ROLES[1], id: 5, assignedCount: 3 },
    ]);
    expect(component.roles()[0].userCount).toBe(3);
  });

  it("allows deletion only for unassigned, non-system roles", async () => {
    const { component } = await setup();
    expect(component.canDelete(listItem())).toBe(true);
    expect(component.canDelete(listItem({ isSystem: true }))).toBe(false);
    expect(component.canDelete(listItem({ userCount: 2 }))).toBe(false);
  });

  it("deletes a role with its scope after confirmation, then reloads", async () => {
    const { component, roleService } = await setup();
    const reloadSpy = jest.spyOn(roleService, "getRoles");
    const callsBefore = reloadSpy.mock.calls.length;

    component.deleteRole(listItem({ id: "5", scope: "group" }));

    expect(roleService.deleteRole).toHaveBeenCalledWith(PermissionScope.Group, 5);
    expect(reloadSpy.mock.calls.length).toBe(callsBefore + 1);
  });

  it("does not delete when the confirmation is cancelled", async () => {
    const { component, roleService } = await setup(ROLES, DESCRIPTORS, false);

    component.deleteRole(listItem({ id: "5", scope: "app" }));

    expect(roleService.deleteRole).not.toHaveBeenCalled();
  });

  it("explains why a disabled delete is blocked without revealing assignees", async () => {
    const { component, snackbar } = await setup();

    component.disabledDeleteClicked(listItem({ name: "Sys", isSystem: true }));
    expect(snackbar.info).toHaveBeenCalledWith(
      "Cannot delete Sys because it is a system role.",
    );

    component.disabledDeleteClicked(listItem({ name: "Busy", userCount: 4 }));
    expect(snackbar.info).toHaveBeenCalledWith(
      "Cannot delete Busy because it is currently assigned.",
    );
  });

  it("shows an error and does not reload when delete fails", async () => {
    const { component, roleService, snackbar } = await setup();
    const reloadSpy = jest.spyOn(roleService, "getRoles");
    const callsBefore = reloadSpy.mock.calls.length;
    (roleService.deleteRole as jest.Mock).mockReturnValueOnce(
      throwError(() => new Error("boom")) as any,
    );

    component.deleteRole(listItem({ id: "5", scope: "group" }));

    expect(snackbar.error).toHaveBeenCalledWith("Failed to delete role");
    expect(reloadSpy.mock.calls.length).toBe(callsBefore);
  });

  it("preselects each default selector from the role marked default for its scope", async () => {
    const { component } = await setup();
    // ROLES marks the APP role (id 1) and the GROUP role (id 2) as default.
    expect(component.defaultAppRoleControl.value).toBe(1);
    expect(component.defaultGroupRoleControl.value).toBe(2);
  });

  it("offers each scope's roles as options for its selector", async () => {
    const { component } = await setup();
    expect(component.appRoleOptions()).toEqual([{ id: 1, name: "Administrator" }]);
    expect(component.groupRoleOptions()).toEqual([{ id: 2, name: "Group Manager" }]);
  });

  it("sets the new default for the scope and reloads when a selector changes", async () => {
    const roles: Role[] = [
      { ...ROLES[0], id: 1, isDefault: true },
      { ...ROLES[0], id: 9, name: "Restricted", isDefault: false, isSystem: false },
    ];
    const { component, roleService, snackbar } = await setup(roles);
    const reloadSpy = jest.spyOn(roleService, "getRoles");
    const callsBefore = reloadSpy.mock.calls.length;

    component.defaultAppRoleControl.setValue(9);

    expect(roleService.setDefaultRole).toHaveBeenCalledWith(PermissionScope.App, 9);
    expect(snackbar.success).toHaveBeenCalledWith("Default role updated");
    expect(reloadSpy.mock.calls.length).toBe(callsBefore + 1);
  });

  it("does not call the API when syncing selectors from the server", async () => {
    const { roleService } = await setup();
    // The constructor loads roles and patches the selectors with emitEvent:false,
    // so no set-default request is made just from rendering.
    expect(roleService.setDefaultRole).not.toHaveBeenCalled();
  });

  it("shows an error and reverts the selector when setting the default fails", async () => {
    const roles: Role[] = [
      { ...ROLES[0], id: 1, isDefault: true },
      { ...ROLES[0], id: 9, name: "Restricted", isDefault: false, isSystem: false },
    ];
    const { component, roleService, snackbar } = await setup(roles);
    const reloadSpy = jest.spyOn(roleService, "getRoles");
    const callsBefore = reloadSpy.mock.calls.length;
    (roleService.setDefaultRole as jest.Mock).mockReturnValueOnce(
      throwError(() => new Error("boom")) as any,
    );

    component.defaultAppRoleControl.setValue(9);

    expect(snackbar.error).toHaveBeenCalledWith("Failed to update default role");
    // Reloaded to revert the selector to the server's current default.
    expect(reloadSpy.mock.calls.length).toBe(callsBefore + 1);
  });
});
