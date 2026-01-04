import { provideHttpClient, withInterceptorsFromDi } from "@angular/common/http";
import { provideHttpClientTesting } from "@angular/common/http/testing";
import { CUSTOM_ELEMENTS_SCHEMA } from "@angular/core";
import { ComponentFixture, TestBed } from "@angular/core/testing";
import { ReactiveFormsModule } from "@angular/forms";
import { MatSnackBarModule } from "@angular/material/snack-bar";
import { ActivatedRoute } from "@angular/router";
import { Store } from "@ngxs/store";
import { of } from "rxjs";
import { InputReadonlyPipe } from "src/pipes/input-readonly.pipe";
import { SharedUiModule } from "src/shared-ui/shared-ui.module";
import { UserPreferences, UserPreferencesService } from "../../open-api";
import { PipesModule } from "../../pipes";
import { StoreModule } from "../../store/store.module";

import { UserPreferencesComponent } from "./user-preferences.component";

describe("UserPreferencesComponent", () => {
  let component: UserPreferencesComponent;
  let fixture: ComponentFixture<UserPreferencesComponent>;

  beforeEach(() => {
    TestBed.configureTestingModule({
      declarations: [UserPreferencesComponent, InputReadonlyPipe],
      schemas: [CUSTOM_ELEMENTS_SCHEMA],
      imports: [
        ReactiveFormsModule,
        StoreModule,
        MatSnackBarModule,
        PipesModule,
        SharedUiModule
      ],
      providers: [
        UserPreferencesService,
        {
          provide: ActivatedRoute,
          useValue: { snapshot: { data: { formConfig: {} } } },
        },
        provideHttpClient(withInterceptorsFromDi()),
        provideHttpClientTesting(),
      ]
    });
    fixture = TestBed.createComponent(UserPreferencesComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it("should create", () => {
    expect(component).toBeTruthy();
  });

  it("should init form correctly without data", () => {
    component.ngOnInit();

    expect(component.form.value).toEqual({
      quickScanDefaultPaidById: "",
      quickScanDefaultGroupId: "",
      quickScanDefaultStatus: "",
      showLargeImagePreviews: false,
      userShortcuts: []
    });
  });

  it("should init form with user preference data", () => {
    const store = TestBed.inject(Store);
    store.reset({
      ...store.snapshot(),
      auth: {
        ...store.snapshot().auth,
        userPreferences: {
          quickScanDefaultPaidById: "1",
          quickScanDefaultGroupId: "2",
          quickScanDefaultStatus: "OPEN",
          showLargeImagePreviews: true,
          userShortcuts: [{ id: 1, name: "Test", url: "test", icon: "icon" }],
        },
      },
    });
    component.ngOnInit();

    expect(component.form.value).toEqual({
      quickScanDefaultPaidById: "1",
      quickScanDefaultGroupId: "2",
      quickScanDefaultStatus: "OPEN",
      showLargeImagePreviews: true,
      userShortcuts: [{ name: "Test", url: "test", icon: "icon", trackby: 0 }],
    });
  });

  it("should attempt to call update endpoint", () => {
    const userPreference: UserPreferences = {
      quickScanDefaultPaidById: 1,
      quickScanDefaultGroupId: 2,
      quickScanDefaultStatus: "OPEN",
    } as UserPreferences;
    const serviceSpy = jest.spyOn(
      TestBed.inject(UserPreferencesService),
      "updateUserPreferences"
    ).mockReturnValue(of(userPreference) as any);
    const storeSpy = jest.spyOn(TestBed.inject(Store), "dispatch");

    component.ngOnInit();
    component.form.patchValue(userPreference);
    component.submit();

    expect(serviceSpy).toHaveBeenCalledWith(component.form.value);
    // TODO: Get this call covered expect(storeSpy).toHaveBeenCalled();
  });

  it("should attempt to call update endpoint with nulls", () => {
    // Reset store state to ensure clean test
    const store = TestBed.inject(Store);
    store.reset({
      ...store.snapshot(),
      auth: {
        ...store.snapshot().auth,
        userPreferences: undefined,
      },
    });

    const serviceSpy = jest.spyOn(
      TestBed.inject(UserPreferencesService),
      "updateUserPreferences"
    ).mockReturnValue(of(undefined as any));

    component.ngOnInit();
    component.submit();

    expect(serviceSpy).toHaveBeenCalledWith({
      quickScanDefaultPaidById: null,
      quickScanDefaultGroupId: null,
      quickScanDefaultStatus: "",
      showLargeImagePreviews: false,
      userShortcuts: [],
    } as any);
  });
});
