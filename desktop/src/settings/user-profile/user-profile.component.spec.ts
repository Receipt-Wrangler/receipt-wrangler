import { provideHttpClient, withInterceptorsFromDi } from "@angular/common/http";
import { CUSTOM_ELEMENTS_SCHEMA } from "@angular/core";
import { ComponentFixture, TestBed } from "@angular/core/testing";
import { ReactiveFormsModule } from "@angular/forms";
import { MatSnackBarModule } from "@angular/material/snack-bar";
import { ActivatedRoute } from "@angular/router";
import { NgxsModule, Store } from "@ngxs/store";
import { of } from "rxjs";
import { PipesModule } from "src/pipes/pipes.module";
import { ApiModule, AuthService, UserService } from "../../open-api";
import { AuthState, UserState } from "../../store";
import { UserProfileComponent } from "./user-profile.component";

describe("UserProfileComponent", () => {
  let component: UserProfileComponent;
  let fixture: ComponentFixture<UserProfileComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
    declarations: [UserProfileComponent],
    schemas: [CUSTOM_ELEMENTS_SCHEMA],
    imports: [ApiModule,
        PipesModule,
        MatSnackBarModule,
        NgxsModule.forRoot([AuthState, UserState]),
        PipesModule,
        ReactiveFormsModule],
    providers: [
        {
            provide: ActivatedRoute,
            useValue: { snapshot: { data: { formConfig: {} } } },
        },
        provideHttpClient(withInterceptorsFromDi()),
    ]
}).compileComponents();

    fixture = TestBed.createComponent(UserProfileComponent);
    component = fixture.componentInstance;
  });

  it("should create", () => {
    expect(component).toBeTruthy();
  });

  it("should init form correctly", () => {
    const store = TestBed.inject(Store);
    store.reset({
      auth: {
        username: "cheetos",
        displayname: "burger",
        defaultAvatarColor: "#CD5C5C",
      },
    });

    component.ngOnInit();

    expect(component.form.value).toEqual({
      username: "cheetos",
      displayName: "burger",
      defaultAvatarColor: "#CD5C5C",
    });
  });

  it("should submit form and update state correctly", () => {
    const store = TestBed.inject(Store);
    const serviceSpy = jest.spyOn(TestBed.inject(UserService), "updateUserProfile");
    const authSpy = jest.spyOn(TestBed.inject(AuthService), "getNewRefreshToken");

    jest.spyOn(TestBed.inject(UserService), "getUserClaims").mockReturnValue(
      of({
        userId: "1",
        displayname: "store",
        username: "general",
      } as any)
    );

    serviceSpy.mockReturnValue(of(undefined) as any);
    authSpy.mockReturnValue(of(undefined as any));

    store.reset({
      users: { users: [{ id: 1, displayName: "cheetos", username: "burger" }] },
      auth: {
        userId: "1",
        username: "cheetos",
        displayname: "burger",
        defaultAvatarColor: "#CD5C5C",
      },
    });

    component.ngOnInit();
    component.form.patchValue({
      username: "general",
      displayName: "store",
    });

    component.submit();

    const updatedUser = store.selectSnapshot(AuthState.loggedInUser);
    const updatedUsers = store.selectSnapshot(UserState.users);

    expect(serviceSpy).toHaveBeenCalledWith({
      username: "general",
      displayName: "store",
      defaultAvatarColor: "#CD5C5C",
    } as any);
    expect(authSpy).toHaveBeenCalled();

    // expect(updatedUser).toEqual({
    //   id: 1,
    //   username: 'general',
    //   displayName: 'store',
    // } as any);

    // expect(updatedUsers[0]).toEqual({
    //   id: 1,
    //   username: 'general',
    //   displayName: 'store',
    // } as any);
  });
});
