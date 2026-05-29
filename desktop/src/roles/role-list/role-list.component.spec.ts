import { CUSTOM_ELEMENTS_SCHEMA, provideZonelessChangeDetection } from "@angular/core";
import { ComponentFixture, TestBed } from "@angular/core/testing";
import { RouterModule } from "@angular/router";
import { of, throwError } from "rxjs";
import { ButtonModule } from "../../button";
import { PermissionService } from "../../open-api";
import { RoleListComponent } from "./role-list.component";
import { RoleListItem } from "./role-list-item.interface";

describe("RoleListComponent", () => {
  let component: RoleListComponent;
  let fixture: ComponentFixture<RoleListComponent>;
  let permissionService: { getPermissions: jest.Mock };

  const buildRole = (overrides: Partial<RoleListItem> = {}): RoleListItem => ({
    id: "role-1",
    name: "Administrator",
    description: "Full access",
    scopes: ["app", "group"],
    appCount: 10,
    groupCount: 5,
    members: [],
    userCount: 0,
    isSystem: true,
    icon: "shield_person",
    iconColor: "#27b1ff",
    iconTint: "#ccecff",
    ...overrides,
  });

  beforeEach(async () => {
    permissionService = { getPermissions: jest.fn().mockReturnValue(of([])) };

    await TestBed.configureTestingModule({
      declarations: [RoleListComponent],
      schemas: [CUSTOM_ELEMENTS_SCHEMA],
      imports: [ButtonModule, RouterModule.forRoot([])],
      providers: [
        provideZonelessChangeDetection(),
        { provide: PermissionService, useValue: permissionService },
      ],
    }).compileComponents();

    fixture = TestBed.createComponent(RoleListComponent);
    component = fixture.componentInstance;
    await fixture.whenStable();
  });

  it("creates", () => {
    expect(component).toBeTruthy();
  });

  it("renders the breadcrumb and the page title", () => {
    expect(fixture.nativeElement.querySelector("app-breadcrumb")).toBeTruthy();
    expect(fixture.nativeElement.querySelector("h1")?.textContent).toContain(
      "Roles",
    );
  });

  it("shows the empty state and no table when there are no roles", () => {
    expect(component.roleCount()).toBe(0);
    expect(fixture.nativeElement.querySelector(".empty-state")).toBeTruthy();
    expect(fixture.nativeElement.querySelector(".roles-table")).toBeNull();
  });

  it("loads the permission total on init for the meter denominator", () => {
    expect(permissionService.getPermissions).toHaveBeenCalledTimes(1);
  });

  it("does not throw and keeps the empty state when getPermissions fails", async () => {
    permissionService.getPermissions.mockReturnValue(
      throwError(() => new Error("boom")),
    );

    const errorFixture = TestBed.createComponent(RoleListComponent);
    await errorFixture.whenStable();

    expect(permissionService.getPermissions).toHaveBeenCalled();
    expect(errorFixture.componentInstance).toBeTruthy();
    expect(errorFixture.nativeElement.querySelector(".empty-state")).toBeTruthy();
  });

  it("renders the table with rows when roles are present", async () => {
    component.roles.set([buildRole()]);
    await fixture.whenStable();

    expect(fixture.nativeElement.querySelector(".empty-state")).toBeNull();
    expect(
      fixture.nativeElement.querySelectorAll(".roles-table tbody tr").length,
    ).toBe(1);
    expect(component.roleCount()).toBe(1);
  });

  describe("interactions", () => {
    let role: RoleListItem;

    beforeEach(async () => {
      role = buildRole();
      component.roles.set([role]);
      await fixture.whenStable();
    });

    it("invokes addRole when the Add Role button is clicked", () => {
      const addRoleSpy = jest.spyOn(component, "addRole");
      const addButton = Array.from(
        fixture.nativeElement.querySelectorAll("button"),
      ).find((button) =>
        (button as HTMLElement).textContent?.includes("Add Role"),
      ) as HTMLButtonElement;

      addButton.click();

      expect(addRoleSpy).toHaveBeenCalled();
    });

    it("invokes editRole with the role when the edit action is clicked", () => {
      const editRoleSpy = jest.spyOn(component, "editRole");
      const actionButtons =
        fixture.nativeElement.querySelectorAll(".row-actions button");

      (actionButtons[0] as HTMLButtonElement).click();

      expect(editRoleSpy).toHaveBeenCalledWith(role);
    });

    it("invokes openRoleMenu with the role when the more action is clicked", () => {
      const openRoleMenuSpy = jest.spyOn(component, "openRoleMenu");
      const actionButtons =
        fixture.nativeElement.querySelectorAll(".row-actions button");

      (actionButtons[1] as HTMLButtonElement).click();

      expect(openRoleMenuSpy).toHaveBeenCalledWith(role);
    });
  });

  describe("meter", () => {
    it("returns ten segments with the filled ones flagged app/group", async () => {
      permissionService.getPermissions.mockReturnValue(of(new Array(20).fill({})));
      // Re-create so the total is loaded from the new mock value.
      fixture = TestBed.createComponent(RoleListComponent);
      component = fixture.componentInstance;
      await fixture.whenStable();

      const segments = component.meter(buildRole({ appCount: 10, groupCount: 0 }));
      expect(segments.length).toBe(10);
      // 10 of 20 permissions granted -> ~5 filled segments, app-dominant.
      expect(segments.filter((s) => s === "app").length).toBe(5);
      expect(segments.filter((s) => s === null).length).toBe(5);
    });
  });

  describe("initials", () => {
    it("returns up to two uppercased initials", () => {
      expect(component.initials("Noah Hall")).toBe("NH");
      expect(component.initials("dana")).toBe("D");
    });
  });
});
