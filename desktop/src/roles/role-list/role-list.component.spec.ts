import { CommonModule } from "@angular/common";
import { provideHttpClient, withInterceptorsFromDi } from "@angular/common/http";
import { provideHttpClientTesting } from "@angular/common/http/testing";
import { CUSTOM_ELEMENTS_SCHEMA, provideZonelessChangeDetection } from "@angular/core";
import { ComponentFixture, TestBed } from "@angular/core/testing";
import { Router, RouterModule } from "@angular/router";
import { of, throwError } from "rxjs";
import {
  ApiModule,
  PermissionDescriptor,
  PermissionService,
  Role,
  RoleService,
} from "../../open-api";
import { RoleListComponent } from "./role-list.component";

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
    isSystem: true,
    permissions: ["app.users.create", "app.users.read"],
  },
  {
    id: 2,
    name: "Group Manager",
    description: "Runs groups",
    scope: "GROUP",
    isSystem: false,
    permissions: ["group.view", "group.receipts.read"],
  },
];

async function setup(
  roles: Role[] = ROLES,
  permissions: PermissionDescriptor[] | "error" = DESCRIPTORS,
): Promise<{
  component: RoleListComponent;
  fixture: ComponentFixture<RoleListComponent>;
  navigateSpy: jest.SpyInstance;
}> {
  TestBed.resetTestingModule();
  await TestBed.configureTestingModule({
    declarations: [RoleListComponent],
    schemas: [CUSTOM_ELEMENTS_SCHEMA],
    imports: [ApiModule, CommonModule, RouterModule.forRoot([])],
    providers: [
      provideZonelessChangeDetection(),
      provideHttpClient(withInterceptorsFromDi()),
      provideHttpClientTesting(),
    ],
  }).compileComponents();

  jest.spyOn(TestBed.inject(RoleService), "getRoles").mockReturnValue(of(roles) as any);
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
  return { component, fixture, navigateSpy };
}

describe("RoleListComponent", () => {
  it("creates", async () => {
    const { component } = await setup();
    expect(component).toBeTruthy();
  });

  it("renders the breadcrumb and the page title", async () => {
    const { fixture } = await setup();
    expect(fixture.nativeElement.querySelector("app-breadcrumb")).toBeTruthy();
    expect(fixture.nativeElement.querySelector("h1")?.textContent).toContain("Roles");
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

  it("does not crash when the permission registry errors", async () => {
    const { component } = await setup(ROLES, "error");
    expect(component).toBeTruthy();
    expect(component.roleCount()).toBe(2);
  });

  it("returns up to two uppercased initials", async () => {
    const { component } = await setup();
    expect(component.initials("Noah Hall")).toBe("NH");
    expect(component.initials("dana")).toBe("D");
  });
});
