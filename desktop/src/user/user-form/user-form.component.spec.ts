import { provideHttpClientTesting } from "@angular/common/http/testing";
import { CUSTOM_ELEMENTS_SCHEMA, provideZonelessChangeDetection } from "@angular/core";
import { ComponentFixture, TestBed } from "@angular/core/testing";
import { ReactiveFormsModule, Validators } from "@angular/forms";
import { MatDialog, MatDialogModule, MatDialogRef } from "@angular/material/dialog";
import { MatSnackBarModule } from "@angular/material/snack-bar";
import { NgxsModule, Store } from "@ngxs/store";
import { of, throwError } from "rxjs";
import { ApiModule, PermissionScope, Role, RoleService, User, UserService } from "../../open-api";
import { PipesModule } from "../../pipes";
import { SnackbarService, TokenRefreshService } from "../../services";
import { AddUser, AuthState, UpdateUser, UserState } from "../../store";
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
    imports: [NgxsModule.forRoot([AuthState, UserState]),
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
});
