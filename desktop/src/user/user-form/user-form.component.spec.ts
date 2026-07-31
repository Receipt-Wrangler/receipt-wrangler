import { provideHttpClientTesting } from "@angular/common/http/testing";
import { CUSTOM_ELEMENTS_SCHEMA, provideZonelessChangeDetection } from "@angular/core";
import { ComponentFixture, TestBed } from "@angular/core/testing";
import { ReactiveFormsModule, Validators } from "@angular/forms";
import { MatDialog, MatDialogModule, MatDialogRef } from "@angular/material/dialog";
import { MatSnackBarModule } from "@angular/material/snack-bar";
import { NgxsModule, Store } from "@ngxs/store";
import { Observable, of, throwError } from "rxjs";
import {
  ApiModule,
  GroupsService,
  PermissionScope,
  Role,
  RoleService,
  User,
  UserService,
} from "../../open-api";
import { PipesModule } from "../../pipes";
import { SnackbarService, TokenRefreshService } from "../../services";
import { AddUser, AuthState, GroupState, SetGroups, UpdateUser, UserState } from "../../store";
import { UserFormComponent } from "./user-form.component";
import { provideHttpClient, withInterceptorsFromDi } from "@angular/common/http";

describe("UserFormComponent", () => {
  const defaultAppRole: Role = {
    id: 7,
    name: "Legacy User",
    description: "Standard user",
    scope: PermissionScope.App,
    isDefault: true,
    isSystem: true,
    permissions: [],
  };

  // A group role is included to confirm the selector filters to app roles only.
  const groupRole: Role = {
    id: 8,
    name: "Legacy Owner",
    description: "Group owner",
    scope: PermissionScope.Group,
    isDefault: true,
    isSystem: true,
    permissions: [],
  };

  let component: UserFormComponent;
  let fixture: ComponentFixture<UserFormComponent>;
  let store: Store;
  let getRolesMock: jest.Mock;

  beforeEach(async () => {
    getRolesMock = jest.fn().mockReturnValue(of([] as Role[]));
    await TestBed.configureTestingModule({
    declarations: [UserFormComponent],
    schemas: [CUSTOM_ELEMENTS_SCHEMA],
    imports: [NgxsModule.forRoot([AuthState, UserState, GroupState]),
        ReactiveFormsModule,
        PipesModule,
        MatDialogModule,
        MatSnackBarModule,
        ApiModule],
    providers: [
        provideZonelessChangeDetection(),
        {
            provide: MatDialogRef,
            useValue: {
                close: () => { },
            },
        },
        { provide: RoleService, useValue: { getRoles: getRolesMock } },
        SnackbarService,
        provideHttpClient(withInterceptorsFromDi()),
        provideHttpClientTesting(),
    ]
}).compileComponents();

    store = TestBed.inject(Store);
    fixture = TestBed.createComponent(UserFormComponent);
    component = fixture.componentInstance;
    await fixture.whenStable();
  });

  // Creates and initialises a fresh component after the role mock is configured;
  // the roles toSignal subscribes at construction, so the mock must be set first.
  async function createWithRoles(roles$: any): Promise<UserFormComponent> {
    getRolesMock.mockReturnValue(roles$);
    const freshFixture = TestBed.createComponent(UserFormComponent);
    const freshComponent = freshFixture.componentInstance;
    await freshFixture.whenStable();
    return freshComponent;
  }

  it("should create", () => {
    expect(component).toBeTruthy();
  });

  it("should init empty form correctly", () => {
    component.ngOnInit();

    expect(component.form.value).toEqual({
      displayName: "",
      username: "",
      appRoleId: null,
      password: "",
      isDummyUser: false,
    });
  });

  it("should init form with user data correctly", () => {
    const user: User = {
      id: 1,
      displayName: "Pizza man",
      username: "Waffle guy",
      isDummyUser: false,
      appRoleId: 5,
    } as User;

    component.user = user;
    component.ngOnInit();

    expect(component.form.value).toEqual({
      displayName: "Pizza man",
      username: "Waffle guy",
      appRoleId: 5,
    });
    expect(component.form.get("isDummyUser")?.value).toEqual(false);
  });

  it("should attempt to update user on submit and get refresh token if the user is updating his/her own record", () => {
    store.reset({
      auth: {
        userId: "1",
      },
    });

    jest.spyOn(TestBed.inject(SnackbarService), "success").mockReturnValue();
    jest.spyOn(TestBed.inject(UserService), "getUsernameCount").mockReturnValue(
      of(0 as any)
    );
    const userServiceSpy = jest.spyOn(TestBed.inject(UserService), "updateUserById");
    userServiceSpy.mockReturnValue(of(undefined as any));

    const authServiceSpy = jest.spyOn(
      TestBed.inject(TokenRefreshService),
      "refreshToken"
    );
    authServiceSpy.mockReturnValue(of(undefined as any));

    const storeSpy = jest.spyOn(TestBed.inject(Store), "dispatch");
    storeSpy.mockReturnValue(of(undefined));

    const dialogRefSpy = jest.spyOn(component.matDialogRef, "close");

    const user: User = {
      id: 1,
      displayName: "Pizza man",
      username: "Waffle guy",
      isDummyUser: false,
      appRoleId: 5,
    } as User;

    component.user = user;
    component.ngOnInit();

    component.submit();

    expect(userServiceSpy).toHaveBeenCalledWith(
      1,
      {
        displayName: "Pizza man",
        username: "Waffle guy",
        appRoleId: 5,
      } as User,
    );

    expect(storeSpy).toHaveBeenCalledWith(
      new UpdateUser("1", {
        ...component.user,
        ...component.form.value,
      })
    );

    expect(authServiceSpy).toHaveBeenCalled();

    expect(dialogRefSpy).toHaveBeenCalledWith(true);
  });

  it("should attempt to update user on submit and not get refresh token if the user is not updating his/her own record", () => {
    store.reset({
      auth: {
        userId: "2",
      },
    });

    jest.spyOn(TestBed.inject(SnackbarService), "success").mockReturnValue();
    jest.spyOn(TestBed.inject(UserService), "getUsernameCount").mockReturnValue(
      of(0 as any)
    );
    const userServiceSpy = jest.spyOn(TestBed.inject(UserService), "updateUserById");
    userServiceSpy.mockReturnValue(of(undefined as any));

    const authServiceSpy = jest.spyOn(
      TestBed.inject(TokenRefreshService),
      "refreshToken"
    );
    authServiceSpy.mockReturnValue(of(undefined as any));

    const storeSpy = jest.spyOn(TestBed.inject(Store), "dispatch");
    storeSpy.mockReturnValue(of(undefined));

    const dialogRefSpy = jest.spyOn(component.matDialogRef, "close");

    const user: User = {
      id: 1,
      displayName: "Pizza man",
      username: "Waffle guy",
      isDummyUser: false,
      appRoleId: 5,
    } as User;

    component.user = user;
    component.ngOnInit();

    component.submit();

    expect(userServiceSpy).toHaveBeenCalledWith(
      1,
      {
        displayName: "Pizza man",
        username: "Waffle guy",
        appRoleId: 5,
      } as User,
    );

    expect(storeSpy).toHaveBeenCalledWith(
      new UpdateUser("1", {
        ...component.user,
        ...component.form.value,
      })
    );
    expect(component.form.get("isDummyUser")?.value).toEqual(false);

    expect(authServiceSpy).toHaveBeenCalledTimes(0);

    expect(dialogRefSpy).toHaveBeenCalledWith(true);
  });

  it("should attempt to create user", () => {
    jest.spyOn(TestBed.inject(SnackbarService), "success").mockReturnValue();
    jest.spyOn(TestBed.inject(UserService), "getUsernameCount").mockReturnValue(
      of(0 as any)
    );
    const user: User = {
      id: 1,
      displayName: "Pizza man",
      username: "Waffle guy",
      isDummyUser: false,
      appRoleId: 5,
    } as User;

    const userServiceSpy = jest.spyOn(TestBed.inject(UserService), "createUser");
    userServiceSpy.mockReturnValue(of(user as any));

    const storeSpy = jest.spyOn(TestBed.inject(Store), "dispatch");
    storeSpy.mockReturnValue(of(undefined));

    const dialogRefSpy = jest.spyOn(component.matDialogRef, "close");
    component.ngOnInit();

    component.form.patchValue({
      displayName: "Pizza man",
      username: "Waffle guy",
      isDummyUser: false,
      password: "Dough boy",
      appRoleId: 5,
    });

    component.submit();

    expect(userServiceSpy).toHaveBeenCalledWith({
      displayName: "Pizza man",
      username: "Waffle guy",
      isDummyUser: false,
      appRoleId: 5,
      password: "Dough boy",
    } as any);

    expect(storeSpy).toHaveBeenCalledWith(new AddUser(user));

    expect(dialogRefSpy).toHaveBeenCalledWith(true);
  });

  it("keeps the dialog open when creating the user fails", () => {
    // Closing on failure would report a create that never happened and discard
    // everything the admin typed.
    const errorSpy = jest
      .spyOn(TestBed.inject(SnackbarService), "error")
      .mockReturnValue();
    const successSpy = jest
      .spyOn(TestBed.inject(SnackbarService), "success")
      .mockReturnValue();
    jest
      .spyOn(TestBed.inject(UserService), "getUsernameCount")
      .mockReturnValue(of(0 as any));
    jest
      .spyOn(TestBed.inject(UserService), "createUser")
      .mockReturnValue(
        throwError(() => ({ error: { username: "Username already exists" } })) as any
      );

    const dialogRefSpy = jest.spyOn(component.matDialogRef, "close");
    component.ngOnInit();

    component.form.patchValue({
      displayName: "Pizza man",
      username: "Waffle guy",
      isDummyUser: false,
      password: "Dough boy",
      appRoleId: 5,
    });

    component.submit();

    expect(errorSpy).toHaveBeenCalledWith("Username already exists");
    expect(successSpy).not.toHaveBeenCalled();
    expect(dialogRefSpy).not.toHaveBeenCalled();
  });

  it("should disable empty and disable password field if isDummyUser is true", () => {
    component.ngOnInit();
    component.form.patchValue({
      isDummyUser: true,
    });

    const passwordField = component.form.get("password");

    expect(passwordField?.disabled).toEqual(true);
    expect(passwordField?.value).toEqual("");
    expect(passwordField?.hasValidator(Validators.required)).toEqual(false);
  });

  it("should disable empty and disable password field if isDummyUser is false", () => {
    component.ngOnInit();
    component.form.patchValue({
      isDummyUser: true,
    });

    component.form.patchValue({
      isDummyUser: false,
    });

    const passwordField = component.form.get("password");

    expect(passwordField?.disabled).toEqual(false);
    expect(passwordField?.value).toEqual("");
    expect(passwordField?.hasValidator(Validators.required)).toEqual(true);
  });

  it("should disable isDummyUser if user is defined", () => {
    component.user = {} as User;

    component.ngOnInit();

    const isDummyUserField = component.form.get("isDummyUser");

    expect(isDummyUserField?.disabled).toEqual(true);
  });

  it("defaults to the configured default app role on add", async () => {
    const freshComponent = await createWithRoles(of([groupRole, defaultAppRole]));

    expect(freshComponent.appRoleOptions().map((role) => role.id)).toEqual([7]);
    expect(freshComponent.form.get("appRoleId")?.value).toBe(7);
    expect(freshComponent.selectedRole()).toEqual(defaultAppRole);
  });

  it("resolves the selected app role when a group role shares its id", async () => {
    // App and group role ids are independent and can collide. A group role
    // listed first with the same id must not be returned for the app selection.
    const collidingGroupRole: Role = {
      id: 7,
      name: "Some Group Role",
      scope: PermissionScope.Group,
      isDefault: false,
      isSystem: false,
      permissions: [],
    };
    const freshComponent = await createWithRoles(
      of([collidingGroupRole, defaultAppRole])
    );

    expect(freshComponent.form.get("appRoleId")?.value).toBe(7);
    expect(freshComponent.selectedRole()).toEqual(defaultAppRole);
  });

  it("leaves the selector empty when roles fail to load", async () => {
    const freshComponent = await createWithRoles(
      throwError(() => new Error("forbidden"))
    );

    expect(freshComponent).toBeTruthy();
    expect(freshComponent.appRoleOptions()).toEqual([]);
    expect(freshComponent.form.get("appRoleId")?.value).toBeNull();
  });

  it("opens the preview dialog for the selected role", async () => {
    const freshComponent = await createWithRoles(of([defaultAppRole]));

    const openSpy = jest
      .spyOn(TestBed.inject(MatDialog), "open")
      .mockReturnValue({ afterClosed: () => of(undefined) } as any);

    freshComponent.previewRole();

    expect(openSpy).toHaveBeenCalledTimes(1);
    expect(openSpy.mock.calls[0][1]?.data?.role).toEqual(defaultAppRole);
  });

  describe("per-member category/tag assignment", () => {
    const groups = [
      {
        id: 100,
        name: "Agency",
        groupMembers: [
          { userId: 1, groupId: 100, groupRoleId: 8, categoryGrants: [5], tagGrants: [] },
        ],
      },
    ] as any[];

    const user: User = {
      id: 1,
      displayName: "Foster Parent",
      username: "fparent",
      appRoleId: 7,
    } as User;

    function seedGroups(): void {
      store.dispatch(new SetGroups(groups));
    }

    /**
     * Builds an edit-mode component with everything submit() touches stubbed.
     *
     * Every test below needs the same six spies and differs only in what the two
     * writes return, so they are parameters rather than a copy of the block.
     */
    async function createEditForm(
      results: {
        grants?: Observable<unknown>;
        userUpdate?: Observable<unknown>;
      } = {}
    ) {
      // Seed groups BEFORE stubbing dispatch — the rows are read from GroupState.
      seedGroups();
      getRolesMock.mockReturnValue(of([defaultAppRole, groupRole]));

      const successSpy = jest
        .spyOn(TestBed.inject(SnackbarService), "success")
        .mockReturnValue();
      jest
        .spyOn(TestBed.inject(UserService), "getUsernameCount")
        .mockReturnValue(of(0 as any));
      jest
        .spyOn(TestBed.inject(UserService), "updateUserById")
        .mockReturnValue((results.userUpdate ?? of(undefined)) as any);
      jest
        .spyOn(TestBed.inject(TokenRefreshService), "refreshToken")
        .mockReturnValue(of(undefined as any));
      const updateGrantsSpy = jest
        .spyOn(TestBed.inject(GroupsService), "updateGroupMemberGrants")
        .mockReturnValue((results.grants ?? of({})) as any);

      jest.spyOn(TestBed.inject(Store), "dispatch").mockReturnValue(of(undefined));

      const freshFixture = TestBed.createComponent(UserFormComponent);
      const freshComponent = freshFixture.componentInstance;
      const closeSpy = jest.spyOn(freshComponent.matDialogRef, "close");
      freshComponent.user = user;
      freshComponent.ngOnInit();
      await freshFixture.whenStable();

      return { freshFixture, freshComponent, successSpy, updateGrantsSpy, closeSpy };
    }

    it("shows no assignment rows when creating a user", async () => {
      // Grants hang off a MEMBERSHIP, and a user being created has none yet.
      seedGroups();
      const freshComponent = await createWithRoles(of([defaultAppRole, groupRole]));
      freshComponent.ngOnInit();

      expect(freshComponent.grantRows()).toEqual([]);
    });

    it("builds one row per group the edited user belongs to", async () => {
      seedGroups();
      getRolesMock.mockReturnValue(of([defaultAppRole, groupRole]));
      const freshFixture = TestBed.createComponent(UserFormComponent);
      const freshComponent = freshFixture.componentInstance;
      freshComponent.user = user;
      freshComponent.ngOnInit();
      await freshFixture.whenStable();

      const rows = freshComponent.grantRows();
      expect(rows.length).toBe(1);
      expect(rows[0].groupName).toBe("Agency");
      expect(rows[0].roleName).toBe("Legacy Owner");
      expect(rows[0].current.categoryIds).toEqual([5]);
    });

    it("writes changed assignments through the grants endpoint on submit", async () => {
      const { freshFixture, freshComponent, updateGrantsSpy } = await createEditForm();

      freshComponent.onGrantsChange(100, { categoryIds: [5, 6], tagIds: [] });
      freshComponent.submit();
      await freshFixture.whenStable();

      expect(updateGrantsSpy).toHaveBeenCalledWith(100, 1, {
        categoryGrants: [5, 6],
        tagGrants: [],
      });
    });

    // The user record and the grants are two independent writes, so the success
    // message must speak for both — announcing it after the first one claims an
    // assignment that never saved.
    it("reports success only after the grants save succeeds", async () => {
      const { freshFixture, freshComponent, successSpy } = await createEditForm({
        grants: throwError(() => new Error("400 ceiling violation")),
      });

      freshComponent.onGrantsChange(100, { categoryIds: [5, 6], tagIds: [] });
      freshComponent.submit();
      await freshFixture.whenStable();

      expect(successSpy).not.toHaveBeenCalled();
    });

    it("keeps the dialog open when the grants save fails", async () => {
      // The other half of the contract: the interceptor reports the failure, and
      // the selection survives so the admin can correct it instead of retyping.
      const { freshFixture, freshComponent, closeSpy } = await createEditForm({
        grants: throwError(() => new Error("400 ceiling violation")),
      });

      freshComponent.onGrantsChange(100, { categoryIds: [5, 6], tagIds: [] });
      freshComponent.submit();
      await freshFixture.whenStable();

      expect(closeSpy).not.toHaveBeenCalled();
    });

    it("reports nothing and keeps the dialog open when the user update fails", async () => {
      const { freshFixture, freshComponent, successSpy, closeSpy, updateGrantsSpy } =
        await createEditForm({
          userUpdate: throwError(() => new Error("500")),
        });

      freshComponent.onGrantsChange(100, { categoryIds: [5, 6], tagIds: [] });
      freshComponent.submit();
      await freshFixture.whenStable();

      expect(successSpy).not.toHaveBeenCalled();
      expect(closeSpy).not.toHaveBeenCalled();
      // The grants write is downstream of the user update, so it never runs.
      expect(updateGrantsSpy).not.toHaveBeenCalled();
    });

    it("reports success once and closes the dialog when both writes land", async () => {
      const { freshFixture, freshComponent, successSpy, closeSpy } = await createEditForm();

      freshComponent.onGrantsChange(100, { categoryIds: [5, 6], tagIds: [] });
      freshComponent.submit();
      await freshFixture.whenStable();

      // Exactly once: a toast left behind at the old site would double up here.
      expect(successSpy).toHaveBeenCalledTimes(1);
      expect(successSpy).toHaveBeenCalledWith("User successfully updated");
      expect(closeSpy).toHaveBeenCalledWith(true);
    });

    it("does not write assignments that were never edited", async () => {
      const { freshFixture, freshComponent, updateGrantsSpy } = await createEditForm();

      // Renaming the user must not rewrite their assignment.
      freshComponent.submit();
      await freshFixture.whenStable();

      expect(updateGrantsSpy).not.toHaveBeenCalled();
    });
  });
});
