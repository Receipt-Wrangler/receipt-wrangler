import { provideHttpClientTesting } from "@angular/common/http/testing";
import { CUSTOM_ELEMENTS_SCHEMA } from "@angular/core";
import { ComponentFixture, TestBed } from "@angular/core/testing";
import { ReactiveFormsModule, Validators } from "@angular/forms";
import { MatDialogModule, MatDialogRef } from "@angular/material/dialog";
import { MatSnackBarModule } from "@angular/material/snack-bar";
import { NgxsModule, Store } from "@ngxs/store";
import { of } from "rxjs";
import { ApiModule, AuthService, User, UserRole, UserService } from "../../open-api";
import { PipesModule } from "../../pipes";
import { SnackbarService } from "../../services";
import { AddUser, AuthState, UpdateUser, UserState } from "../../store";
import { UserFormComponent } from "./user-form.component";
import { provideHttpClient, withInterceptorsFromDi } from "@angular/common/http";

describe("UserFormComponent", () => {
  let component: UserFormComponent;
  let fixture: ComponentFixture<UserFormComponent>;
  let store: Store;

  beforeEach(async () => {
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
        {
            provide: MatDialogRef,
            useValue: {
                close: () => { },
            },
        },
        SnackbarService,
        provideHttpClient(withInterceptorsFromDi()),
        provideHttpClientTesting(),
    ]
}).compileComponents();

    store = TestBed.inject(Store);
    fixture = TestBed.createComponent(UserFormComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it("should create", () => {
    expect(component).toBeTruthy();
  });

  it("should init empty form correctly", () => {
    component.ngOnInit();

    expect(component.form.value).toEqual({
      displayName: "",
      username: "",
      userRole: "",
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
      userRole: UserRole.Admin,
    } as User;

    component.user = user;
    component.ngOnInit();

    expect(component.form.value).toEqual({
      displayName: "Pizza man",
      username: "Waffle guy",
      userRole: UserRole.Admin,
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
      TestBed.inject(AuthService),
      "getNewRefreshToken"
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
      userRole: UserRole.Admin,
    } as User;

    component.user = user;
    component.ngOnInit();

    component.submit();

    expect(userServiceSpy).toHaveBeenCalledWith(
      1,
      {
        displayName: "Pizza man",
        username: "Waffle guy",
        userRole: UserRole.Admin,
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
      TestBed.inject(AuthService),
      "getNewRefreshToken"
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
      userRole: UserRole.Admin,
    } as User;

    component.user = user;
    component.ngOnInit();

    component.submit();

    expect(userServiceSpy).toHaveBeenCalledWith(
      1,
      {
        displayName: "Pizza man",
        username: "Waffle guy",
        userRole: UserRole.Admin,
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
      userRole: UserRole.Admin,
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
      userRole: UserRole.Admin,
    });

    component.submit();

    expect(userServiceSpy).toHaveBeenCalledWith({
      displayName: "Pizza man",
      username: "Waffle guy",
      isDummyUser: false,
      userRole: UserRole.Admin,
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
});
