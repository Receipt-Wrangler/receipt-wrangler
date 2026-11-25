import { provideHttpClient, withInterceptorsFromDi } from "@angular/common/http";
import { provideHttpClientTesting } from "@angular/common/http/testing";
import { CUSTOM_ELEMENTS_SCHEMA } from "@angular/core";
import { ComponentFixture, TestBed } from "@angular/core/testing";
import { FormGroup, ReactiveFormsModule } from "@angular/forms";
import { MatDialog, MatDialogModule } from "@angular/material/dialog";
import { MatSnackBarModule } from "@angular/material/snack-bar";
import { ActivatedRoute, Router } from "@angular/router";
import { NgxsModule } from "@ngxs/store";
import { of } from "rxjs";
import { FormMode } from "../../enums/form-mode.enum";
import { ApiModule, SystemEmail, SystemEmailService } from "../../open-api";
import { TABLE_SERVICE_INJECTION_TOKEN } from "../../services/injection-tokens/table-service";
import { SystemEmailTaskTableService } from "../../services/system-email-task-table.service";
import { ConfirmationDialogComponent } from "../../shared-ui/confirmation-dialog/confirmation-dialog.component";
import { SystemEmailTaskTableState } from "../../store/system-email-task-table.state";

import { SystemEmailFormComponent } from "./system-email-form.component";
import { provideNgxMask } from "ngx-mask";
import { mockDialogRef } from "../../testing/dialog-ref.testing";

describe("SystemEmailFormComponent", () => {
  let component: SystemEmailFormComponent;
  let fixture: ComponentFixture<SystemEmailFormComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      schemas: [CUSTOM_ELEMENTS_SCHEMA],
      imports: [ApiModule,
        NgxsModule.forRoot([SystemEmailTaskTableState]),
        MatDialogModule,
        MatSnackBarModule,
        ReactiveFormsModule,
        SystemEmailFormComponent],
      providers: [
        provideNgxMask(),
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
        { provide: Router, useValue: { navigate: vi.fn().mockResolvedValue(true), navigateByUrl: vi.fn().mockResolvedValue(true) } },
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
    const serviceSpy = vi.spyOn(TestBed.inject(SystemEmailService), "createSystemEmail").mockReturnValue(of({ id: 7 } as any));
    component.formConfig = { mode: FormMode.add } as any;
    component.submit();

    expect(serviceSpy).toHaveBeenCalled();
    expect(TestBed.inject(Router).navigateByUrl).toHaveBeenCalledWith("/system-settings/system-emails/7/view");
  });

  it("should call update", () => {
    component.originalSystemEmail = { id: 1 } as any;
    const serviceSpy = vi.spyOn(TestBed.inject(SystemEmailService), "updateSystemEmailById").mockReturnValue(of({} as any));
    component.formConfig = { mode: FormMode.edit } as any;
    component.submit();

    expect(serviceSpy).toHaveBeenCalled();
    expect(TestBed.inject(Router).navigateByUrl).toHaveBeenCalledWith("/system-settings/system-emails/1/view");
  });

  describe("changing the password", () => {
    // Submits an edit with a dirty password and closes the confirmation with [result].
    const submitPasswordChange = (result: boolean | undefined) => {
      component.ngOnInit();
      component.originalSystemEmail = { id: 1 } as any;
      component.form.get("password")?.markAsDirty();
      const matDialogSpy = vi.spyOn(TestBed.inject(MatDialog), "open").mockReturnValue(mockDialogRef({ afterClosed: of(result) }) as any);
      const updateSpy = vi.spyOn(TestBed.inject(SystemEmailService), "updateSystemEmailById").mockReturnValue(of({} as any));
      component.formConfig = { mode: FormMode.edit } as any;
      component.submit();
      return { matDialogSpy, updateSpy };
    };

    it("should pop dialog", () => {
      const { matDialogSpy } = submitPasswordChange(false);

      expect(matDialogSpy).toHaveBeenCalledWith(ConfirmationDialogComponent);
    });

    it("replaces the password when the dialog is confirmed", () => {
      const { updateSpy } = submitPasswordChange(true);

      expect(updateSpy).toHaveBeenCalledWith(1, true, component.form.value);
    });

    it("keeps the password when the dialog is cancelled", () => {
      const { updateSpy } = submitPasswordChange(false);

      expect(updateSpy).toHaveBeenCalledWith(1, false, component.form.value);
    });

    // Esc and a backdrop click close the dialog with undefined, which the
    // generated client rejects outright, so the save used to do nothing.
    it("treats a dismissed dialog like Cancel", () => {
      const { updateSpy } = submitPasswordChange(undefined);

      expect(updateSpy).toHaveBeenCalledWith(1, false, component.form.value);
    });
  });
});
