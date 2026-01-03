import { provideHttpClient, withInterceptorsFromDi } from "@angular/common/http";
import { provideHttpClientTesting } from "@angular/common/http/testing";
import { CUSTOM_ELEMENTS_SCHEMA } from "@angular/core";
import { ComponentFixture, TestBed } from "@angular/core/testing";
import { FormGroup, ReactiveFormsModule } from "@angular/forms";
import { MatDialog, MatDialogModule } from "@angular/material/dialog";
import { MatSnackBarModule } from "@angular/material/snack-bar";
import { NoopAnimationsModule } from "@angular/platform-browser/animations";
import { ActivatedRoute } from "@angular/router";
import { NgxsModule } from "@ngxs/store";
import { of } from "rxjs";
import { FormMode } from "../../enums/form-mode.enum";
import { ApiModule, SystemEmail, SystemEmailService } from "../../open-api";
import { PipesModule } from "../../pipes";
import { TABLE_SERVICE_INJECTION_TOKEN } from "../../services/injection-tokens/table-service";
import { SystemEmailTaskTableService } from "../../services/system-email-task-table.service";
import { ConfirmationDialogComponent } from "../../shared-ui/confirmation-dialog/confirmation-dialog.component";
import { SharedUiModule } from "../../shared-ui/shared-ui.module";
import { SystemEmailTaskTableState } from "../../store/system-email-task-table.state";

import { SystemEmailFormComponent } from "./system-email-form.component";

describe("SystemEmailFormComponent", () => {
  let component: SystemEmailFormComponent;
  let fixture: ComponentFixture<SystemEmailFormComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [SystemEmailFormComponent],
      schemas: [CUSTOM_ELEMENTS_SCHEMA],
      imports: [ApiModule,
        NgxsModule.forRoot([SystemEmailTaskTableState]),
        PipesModule,
        MatDialogModule,
        MatSnackBarModule,
        ReactiveFormsModule,
        SharedUiModule,
        NoopAnimationsModule],
      providers: [
        {
          provide: TABLE_SERVICE_INJECTION_TOKEN,
          useClass: SystemEmailTaskTableService
        },
        {
          provide: ActivatedRoute,
          useValue: {
            snapshot: {
              data: {
                formConfig: {}
              }
            }
          }
        },
        provideHttpClient(withInterceptorsFromDi()),
        provideHttpClientTesting()
      ]
    })
      .compileComponents();

    fixture = TestBed.createComponent(SystemEmailFormComponent);
    component = fixture.componentInstance;
    component.form = new FormGroup({});
  });

  it("should create", () => {
    expect(component).toBeTruthy();
  });

  it("init form when there is no original data", () => {
    component.ngOnInit();

    expect(component.form.value).toEqual({
      host: null,
      port: null,
      username: null,
      password: null,
      useStartTLS: null,
    });
  });

  it("init form when there is original data", () => {
    const activatedRoute = TestBed.inject(ActivatedRoute);
    activatedRoute.snapshot.data["systemEmail"] = {
      host: "host",
      port: "123",
      username: "username",
      password: "password",
      useStartTLS: true
    } as SystemEmail;
    component.ngOnInit();


    expect(component.form.value).toEqual({
      host: "host",
      port: "123",
      username: "username",
      password: null,
      useStartTLS: true
    });
  });

  it("should call create", () => {
    const serviceSpy = jest.spyOn(TestBed.inject(SystemEmailService), "createSystemEmail").mockReturnValue(of({} as any));
    component.formConfig = { mode: FormMode.add } as any;
    component.submit();

    expect(serviceSpy).toHaveBeenCalled();
  });

  it("should call update", () => {
    component.originalSystemEmail = { id: 1 } as any;
    const serviceSpy = jest.spyOn(TestBed.inject(SystemEmailService), "updateSystemEmailById").mockReturnValue(of({} as any));
    component.formConfig = { mode: FormMode.edit } as any;
    component.submit();

    expect(serviceSpy).toHaveBeenCalled();
  });

  it("should pop dialog", () => {
    component.ngOnInit();
    component.originalSystemEmail = { id: 1 } as any;
    component.form.get("password")?.markAsDirty();
    const matDialogSpy = jest.spyOn(TestBed.inject(MatDialog), "open").mockReturnValue({
      componentInstance: {},
      afterClosed: () => of(undefined)
    } as any);
    component.formConfig = { mode: FormMode.edit } as any;
    component.submit();

    expect(matDialogSpy).toHaveBeenCalledWith(ConfirmationDialogComponent);
  });
});
