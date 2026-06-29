import { provideHttpClient, withInterceptorsFromDi } from "@angular/common/http";
import { provideHttpClientTesting } from "@angular/common/http/testing";
import { Component, CUSTOM_ELEMENTS_SCHEMA, input } from "@angular/core";
import { ComponentFixture, TestBed } from "@angular/core/testing";
import { FormArray, FormGroup, ReactiveFormsModule } from "@angular/forms";
import { MatSnackBarModule } from "@angular/material/snack-bar";
import { ActivatedRoute, Router } from "@angular/router";
import { Store } from "@ngxs/store";
import { of } from "rxjs";
import { FormConfig } from "src/interfaces/form-config.interface";
import { InputReadonlyPipe } from "src/pipes/input-readonly.pipe";
import { EditableListComponent } from "src/shared-ui/editable-list/editable-list.component";
import { SharedUiModule } from "src/shared-ui/shared-ui.module";
import { Permission, UserPreferences, UserPreferencesService, UserShortcut } from "../../open-api";
import { PipesModule } from "../../pipes";
import { SetPermissions } from "../../store/auth.state.actions";
import { StoreModule } from "../../store/store.module";
import { UserShortcutComponent } from "../user-shortcut/user-shortcut.component";

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
        {
          provide: Router,
          useValue: { navigate: jest.fn().mockResolvedValue(true) },
        },
        provideHttpClient(withInterceptorsFromDi()),
        provideHttpClientTesting(),
      ]
    });
    fixture = TestBed.createComponent(UserPreferencesComponent);
    component = fixture.componentInstance;
    Object.defineProperty(component, 'userShortcutComponent', {
      value: () => ({
        editableListComponent: () => ({
          openLastRow: jest.fn(),
          closeRow: jest.fn(),
          getCurrentRowOpen: jest.fn(),
        }),
        isAddingShortcut: false,
      }),
      configurable: true,
    });
    fixture.detectChanges();
  });

  it("should create", () => {
    expect(component).toBeTruthy();
  });

  it("gates the edit button on app.user-preferences.update", () => {
    const store = TestBed.inject(Store);

    store.dispatch(new SetPermissions([], {}));
    expect((component as any).canEdit()).toEqual(false);

    store.dispatch(
      new SetPermissions([Permission.AppUserPreferencesUpdate], {})
    );
    expect((component as any).canEdit()).toEqual(true);
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

  describe("when userShortcutComponent is not yet available", () => {
    beforeEach(() => {
      Object.defineProperty(component, 'userShortcutComponent', {
        value: () => undefined,
        configurable: true,
      });
    });

    it("should not throw on addNewShortcut", () => {
      component.ngOnInit();
      expect(() => component.addNewShortcut()).not.toThrow();
    });

    it("should not throw on shortcutDoneClicked", () => {
      component.ngOnInit();
      expect(() => component.shortcutDoneClicked()).not.toThrow();
    });

    it("should not throw on shortcutCancelClicked", () => {
      component.ngOnInit();
      expect(() => component.shortcutCancelClicked()).not.toThrow();
    });
  });
});
