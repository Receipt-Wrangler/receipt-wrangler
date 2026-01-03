import { provideHttpClient, withInterceptorsFromDi } from "@angular/common/http";
import { provideHttpClientTesting } from "@angular/common/http/testing";
import { CUSTOM_ELEMENTS_SCHEMA } from "@angular/core";
import { ComponentFixture, TestBed } from "@angular/core/testing";
import { MatDialogRef } from "@angular/material/dialog";
import { MatSnackBarModule } from "@angular/material/snack-bar";
import { Store } from "@ngxs/store";
import { of } from "rxjs";
import { ApiModule, UserService } from "../../open-api";
import { PipesModule } from "../../pipes";
import { SnackbarService } from "../../services";
import { UpdateUser } from "../../store";
import { StoreModule } from "../../store/store.module";
import { DummyUserConversionDialogComponent } from "./dummy-user-conversion-dialog.component";

describe("DummyUserConversionDialogComponent", () => {
  let component: DummyUserConversionDialogComponent;
  let fixture: ComponentFixture<DummyUserConversionDialogComponent>;

  beforeEach(() => {
    TestBed.configureTestingModule({
      declarations: [DummyUserConversionDialogComponent],
      schemas: [CUSTOM_ELEMENTS_SCHEMA],
      imports: [
        ApiModule,
        MatSnackBarModule,
        StoreModule,
        PipesModule,
      ],
      providers: [
        SnackbarService,
        {
          provide: MatDialogRef,
          useValue: {
            close: () => { },
          },
        },
        provideHttpClient(withInterceptorsFromDi()),
        provideHttpClientTesting(),
      ]
    });
    fixture = TestBed.createComponent(DummyUserConversionDialogComponent);
    component = fixture.componentInstance;
    component.user = {} as any;
    fixture.detectChanges();
  });

  it("should create", () => {
    expect(component).toBeTruthy();
  });

  it("should init form correctly", () => {
    component.ngOnInit();

    expect(component.form.value).toEqual({
      password: "",
    });
  });

  it("should call close when cancel button is clicked", () => {
    const dialogRefSpy = jest.spyOn(component.matDialogRef, "close");
    component.cancelButtonClicked();

    expect(dialogRefSpy).toHaveBeenCalledTimes(1);
  });

  it("should call service and call update in state", () => {
    const usersService = TestBed.inject(UserService);

    const usersSpy = jest.spyOn(usersService, "convertDummyUserById");
    usersSpy.mockReturnValue(of(undefined as any));

    const store = TestBed.inject(Store);
    const storeSpy = jest.spyOn(store, "dispatch");
    storeSpy.mockReturnValue(of(undefined));

    const dialogSpy = jest.spyOn(component.matDialogRef, "close");
    dialogSpy.mockReturnValue(undefined);

    jest.spyOn(TestBed.inject(SnackbarService), "success").mockReturnValue(
      undefined
    );

    component.user = {
      id: 1,
    } as any;

    component.ngOnInit();

    component.form.patchValue({
      password: "hello world",
    });

    component.submitButtonClicked();

    expect(usersSpy).toHaveBeenCalledWith(
      1,
      {
        password: "hello world",
      },
    );
    expect(storeSpy).toHaveBeenCalledWith(
      new UpdateUser("1", { id: 1, isDummyUser: false } as any)
    );
    expect(dialogSpy).toHaveBeenCalledWith(true);
  });
});
