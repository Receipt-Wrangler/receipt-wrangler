import { provideHttpClient, withInterceptorsFromDi } from "@angular/common/http";
import { provideHttpClientTesting } from "@angular/common/http/testing";
import { CUSTOM_ELEMENTS_SCHEMA, provideZonelessChangeDetection } from "@angular/core";
import { ComponentFixture, TestBed } from "@angular/core/testing";
import { ReactiveFormsModule } from "@angular/forms";
import { MatDialog, MatDialogModule, MatDialogRef } from "@angular/material/dialog";
import { NgxsModule, Store } from "@ngxs/store";
import { NEVER, of, throwError } from "rxjs";
import {
  Category,
  CategoryService,
  Group,
  GroupMember,
  GroupsService,
  PermissionScope,
  Role,
  RoleService,
  Tag,
  TagService,
} from "../../open-api";
import { PipesModule } from "../../pipes";
import { AuthState } from "../../store";
import { GroupMemberFormComponent } from "./group-member-form.component";

describe("GroupMemberFormComponent", () => {
  const defaultGroupRole: Role = {
    id: 10,
    name: "Legacy Owner",
    description: "Group owner",
    scope: PermissionScope.Group,
    isDefault: true,
    isSystem: true,
    permissions: [],
  };

  // An app role is included to confirm the selector filters to group roles only.
  const appRole: Role = {
    id: 1,
    name: "Legacy Admin",
    description: "App admin",
    scope: PermissionScope.App,
    isDefault: true,
    isSystem: true,
    permissions: [],
  };

  let component: GroupMemberFormComponent;
  let fixture: ComponentFixture<GroupMemberFormComponent>;
  let getRolesMock: jest.Mock;
  let dialogRefMock: { close: jest.Mock };
  let getGroupByIdMock: jest.Mock;
  let updateGrantsMock: jest.Mock;

  const categoryPool = [
    { id: 1, name: "Child A" },
    { id: 2, name: "Child B" },
  ] as Category[];

  const cappedRole: Role = {
    id: 10,
    name: "Foster Parent",
    description: "",
    scope: PermissionScope.Group,
    isDefault: true,
    isSystem: false,
    permissions: [],
    categoryGrants: [1, 2],
    tagGrants: [],
  } as unknown as Role;

  /** A saved membership as the API returns it, carrying its grant ids. */
  function savedMember(categoryGrants: number[]): GroupMember {
    return {
      userId: 2,
      groupId: 5,
      groupRoleId: 10,
      categoryGrants,
      tagGrants: [],
    } as unknown as GroupMember;
  }

  beforeEach(async () => {
    getRolesMock = jest.fn().mockReturnValue(of([] as Role[]));
    dialogRefMock = { close: jest.fn() };
    getGroupByIdMock = jest.fn().mockReturnValue(of({ id: 5, groupMembers: [] } as unknown as Group));
    updateGrantsMock = jest.fn().mockReturnValue(of({}));
    await TestBed.configureTestingModule({
      declarations: [GroupMemberFormComponent],
      imports: [
        NgxsModule.forRoot([AuthState]),
        PipesModule,
        ReactiveFormsModule,
        MatDialogModule,
      ],
      schemas: [CUSTOM_ELEMENTS_SCHEMA],
      providers: [
        provideZonelessChangeDetection(),
        { provide: MatDialogRef, useValue: dialogRefMock },
        { provide: RoleService, useValue: { getRoles: getRolesMock } },
        {
          provide: GroupsService,
          useValue: {
            getGroupById: (...args: unknown[]) => getGroupByIdMock(...args),
            updateGroupMemberGrants: (...args: unknown[]) => updateGrantsMock(...args),
          },
        },
        { provide: CategoryService, useValue: { getAllCategories: () => of(categoryPool) } },
        { provide: TagService, useValue: { getAllTags: () => of([] as Tag[]) } },
        provideHttpClient(withInterceptorsFromDi()),
        provideHttpClientTesting(),
      ],
    }).compileComponents();

    TestBed.inject(Store).reset({ auth: { userId: "1" } });
  });

  // Creates and initialises a fresh component after the role mock is configured;
  // the roles toSignal subscribes at construction, so the mock must be set first.
  async function createComponent(): Promise<void> {
    fixture = TestBed.createComponent(GroupMemberFormComponent);
    component = fixture.componentInstance;
    await fixture.whenStable();
  }

  it("should create", async () => {
    await createComponent();
    expect(component).toBeTruthy();
  });

  it("defaults to the configured default group role on add", async () => {
    getRolesMock.mockReturnValue(of([appRole, defaultGroupRole]));
    await createComponent();

    expect(component.groupRoleOptions().map((role) => role.id)).toEqual([10]);
    expect(component.form.get("groupRoleId")?.value).toBe(10);
    expect(component.selectedRole()).toEqual(defaultGroupRole);
  });

  it("resolves the selected group role when an app role shares its id", async () => {
    // App and group role ids are independent and can collide. An app role
    // listed first with the same id must not be returned for the group selection.
    const collidingAppRole: Role = {
      id: 10,
      name: "Some App Role",
      description: "App role",
      scope: PermissionScope.App,
      isDefault: false,
      isSystem: false,
      permissions: [],
    };
    getRolesMock.mockReturnValue(of([collidingAppRole, defaultGroupRole]));
    await createComponent();

    expect(component.form.get("groupRoleId")?.value).toBe(10);
    expect(component.selectedRole()).toEqual(defaultGroupRole);
  });

  it("leaves the selector empty when roles fail to load", async () => {
    getRolesMock.mockReturnValue(throwError(() => new Error("forbidden")));
    await createComponent();

    expect(component).toBeTruthy();
    expect(component.groupRoleOptions()).toEqual([]);
    expect(component.form.get("groupRoleId")?.value).toBeNull();
  });

  it("opens the preview dialog for the selected role", async () => {
    getRolesMock.mockReturnValue(of([defaultGroupRole]));
    await createComponent();

    const openSpy = jest
      .spyOn(TestBed.inject(MatDialog), "open")
      .mockReturnValue({ afterClosed: () => of(undefined) } as any);

    component.previewRole();

    expect(openSpy).toHaveBeenCalledTimes(1);
    expect(openSpy.mock.calls[0][1]?.data?.role).toEqual(defaultGroupRole);
  });

  describe("per-member category/tag assignment", () => {
    /** Mounts the dialog for an already-saved membership carrying [grants]. */
    async function createForSavedMember(grants: number[] = [1]): Promise<void> {
      getRolesMock.mockReturnValue(of([cappedRole]));
      getGroupByIdMock.mockReturnValue(
        of({ id: 5, groupMembers: [savedMember(grants)] } as unknown as Group),
      );
      fixture = TestBed.createComponent(GroupMemberFormComponent);
      component = fixture.componentInstance;
      component.groupMember = savedMember(grants);
      await fixture.whenStable();
      component.ngOnInit();
      await fixture.whenStable();
    }

    it("loads the membership from the API, not from the store", async () => {
      // Grants are written through their own endpoint, which never refreshes
      // GroupState — so a store read can be stale by exactly the edit the admin
      // just made in the user form.
      await createForSavedMember([2]);

      expect(getGroupByIdMock).toHaveBeenCalledWith(5);
      expect(component.grantRow()?.current.categoryIds).toEqual([2]);
    });

    it("derives the ceiling from the member's group role", async () => {
      await createForSavedMember();

      expect(component.grantRow()?.ceiling?.categoryIds).toEqual([1, 2]);
      expect(component.grantRow()?.ceiling?.label).toContain("Foster Parent");
    });

    it("offers no assignment section for an unsaved member", async () => {
      // A member being added is only persisted when the parent group form saves,
      // so there is no membership to write grants against yet.
      getRolesMock.mockReturnValue(of([cappedRole]));
      await createComponent();

      expect(component.grantsReady()).toBe(false);
      expect(component.grantRow()).toBeNull();
      expect(getGroupByIdMock).not.toHaveBeenCalled();
    });

    it("withholds the section until the roles have also landed", async () => {
      // The section depends on four async inputs. Rendering before the roles
      // arrive flashes an unconstrained picker at the admin, because the ceiling
      // comes from the role — so readiness must include it.
      getRolesMock.mockReturnValue(NEVER);
      getGroupByIdMock.mockReturnValue(
        of({ id: 5, groupMembers: [savedMember([1])] } as unknown as Group),
      );
      fixture = TestBed.createComponent(GroupMemberFormComponent);
      component = fixture.componentInstance;
      component.groupMember = savedMember([1]);
      await fixture.whenStable();
      component.ngOnInit();
      await fixture.whenStable();

      expect(component.grantsReady()).toBe(false);
    });

    it("writes a changed assignment on submit", async () => {
      await createForSavedMember([1]);

      component.onGrantsChange(5, { categoryIds: [1, 2], tagIds: [] });
      component.submit();
      await fixture.whenStable();

      expect(updateGrantsMock).toHaveBeenCalledWith(5, 2, {
        categoryGrants: [1, 2],
        tagGrants: [],
      });
      expect(dialogRefMock.close).toHaveBeenCalledWith(component.form);
    });

    it("writes nothing when the assignment was not changed", async () => {
      // Opening and closing the dialog must not rewrite grants — a needless write
      // could also trip the endpoint's ceiling check on a since-narrowed role.
      await createForSavedMember([1]);

      component.submit();
      await fixture.whenStable();

      expect(updateGrantsMock).not.toHaveBeenCalled();
      expect(dialogRefMock.close).toHaveBeenCalledWith(component.form);
    });

    it("keeps the dialog open when the grants save fails", async () => {
      // The interceptor surfaces the error; the admin must keep their selection
      // rather than have it silently discarded by a closing dialog.
      await createForSavedMember([1]);
      updateGrantsMock.mockReturnValue(throwError(() => new Error("400")));

      component.onGrantsChange(5, { categoryIds: [2], tagIds: [] });
      component.submit();
      await fixture.whenStable();

      expect(updateGrantsMock).toHaveBeenCalled();
      expect(dialogRefMock.close).not.toHaveBeenCalled();
    });
  });

  it("returns the form with the selected modern role on submit", async () => {
    getRolesMock.mockReturnValue(of([defaultGroupRole]));
    await createComponent();

    component.form.get("userId")?.setValue("2");
    component.form.get("groupRoleId")?.setValue(10);

    component.submit();

    expect(component.form.get("groupRoleId")?.value).toBe(10);
    expect(dialogRefMock.close).toHaveBeenCalledWith(component.form);
  });
});
