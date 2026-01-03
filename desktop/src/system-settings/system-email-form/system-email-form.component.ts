import { Component, OnInit, ViewChild } from "@angular/core";
import { FormBuilder, FormGroup, Validators } from "@angular/forms";
import { MatDialog } from "@angular/material/dialog";
import { ActivatedRoute, Router } from "@angular/router";
import { take, tap } from "rxjs";
import { FormMode } from "../../enums/form-mode.enum";
import { FormConfig } from "../../interfaces";
import {
  AssociatedEntityType,
  CheckEmailConnectivityCommand,
  ReceiptProcessingSettings,
  SystemEmail,
  SystemEmailService,
  SystemTaskStatus
} from "../../open-api";
import { SnackbarService } from "../../services";
import { TABLE_SERVICE_INJECTION_TOKEN } from "../../services/injection-tokens/table-service";
import { SystemEmailTaskTableService } from "../../services/system-email-task-table.service";
import { ConfirmationDialogComponent } from "../../shared-ui/confirmation-dialog/confirmation-dialog.component";
import { TaskTableComponent } from "../../shared-ui/task-table/task-table.component";

@Component({
  selector: "app-system-email-form",
  templateUrl: "./system-email-form.component.html",
  styleUrl: "./system-email-form.component.scss",
  providers: [{
    provide: TABLE_SERVICE_INJECTION_TOKEN,
    useClass: SystemEmailTaskTableService
  }],
  standalone: false
})
export class SystemEmailFormComponent implements OnInit {
  @ViewChild(TaskTableComponent) public taskTableComponent!: TaskTableComponent;

  protected readonly AssociatedEntityType = AssociatedEntityType;

  public formConfig!: FormConfig;

  public form!: FormGroup;

  public originalSystemEmail!: SystemEmail;

  protected readonly FormMode = FormMode;

  public allReceiptProcessingSettings: ReceiptProcessingSettings[] = [];

  constructor(
    private activatedRoute: ActivatedRoute,
    private formBuilder: FormBuilder,
    private systemEmailService: SystemEmailService,
    private snackbarService: SnackbarService,
    private router: Router,
    public matDialog: MatDialog
  ) {
  }

  public ngOnInit() {
    this.formConfig = this.activatedRoute.snapshot.data["formConfig"];
    this.originalSystemEmail = this.activatedRoute.snapshot.data["systemEmail"];
    this.allReceiptProcessingSettings = this.activatedRoute.snapshot.data["allReceiptProcessingSettings"];
    this.initForm();
  }

  private initForm(): void {
    this.form = this.formBuilder.group({
      host: [this.originalSystemEmail?.host, [Validators.required]],
      port: [this.originalSystemEmail?.port, [Validators.required]],
      username: [this.originalSystemEmail?.username, [Validators.required]],
      password: [null, this.formConfig.mode === FormMode.add ? [Validators.required] : []],
      useStartTLS: [this.originalSystemEmail?.useStartTLS],
    });
  }

  public submit(): void {
    if (this.formConfig.mode === FormMode.add) {
      this.createSystemEmail();
    } else if (this.formConfig.mode === FormMode.edit) {
      this.updateSystemEmail();
    }
  }

  private createSystemEmail(): void {
    this.systemEmailService.createSystemEmail(this.form.value)
      .pipe(
        take(1),
        tap((systemEmail) => {
          this.snackbarService.success("System Email created successfully");
          this.router.navigateByUrl(`/system-settings/system-emails/${systemEmail.id}/view`);
        })
      )
      .subscribe();
  }

  private updateSystemEmail(): void {
    if (this.form.get("password")?.dirty) {
      const dialogRef = this.matDialog.open(ConfirmationDialogComponent);
      dialogRef.componentInstance.headerText = "Update email password";
      dialogRef.componentInstance.dialogContent = "Are you sure you want to update the email password? This will replace the previous password.";

      dialogRef.afterClosed()
        .pipe(
          take(1),
          tap((result) => {
            this.callUpdateEndpoint(result);
          })
        )
        .subscribe();

    } else {
      this.callUpdateEndpoint(false);
    }
  }

  private callUpdateEndpoint(updatePassword: boolean): void {
    this.systemEmailService.updateSystemEmailById(this.originalSystemEmail.id, updatePassword, this.form.value)
      .pipe(
        take(1),
        tap(() => {
          this.snackbarService.success("System Email updated successfully");
          this.router.navigateByUrl(`/system-settings/system-emails/${this.originalSystemEmail.id}/view`);
        })
      )
      .subscribe();
  }

  public checkEmailConnectivity(): void {
    if (this.formConfig.mode === FormMode.edit && this.form.valid && this.form.dirty) {
      const dialogRef = this.matDialog.open(ConfirmationDialogComponent);
      dialogRef.componentInstance.headerText = "Check email connectivity";
      dialogRef.componentInstance.dialogContent = `You have made changes to the email settings. Would you like to check connectivity with
      the unsaved changes? Otherwise, the existing email settings will be used.`;

      dialogRef.afterClosed()
        .pipe(
          take(1),
          tap((result) => {
            if (result) {
              const command = {
                id: this.originalSystemEmail.id,
                ...this.form.value
              } as CheckEmailConnectivityCommand;
              this.checkConnectivitySettings(command);
            } else {
              this.checkConnectivitySettingsWithExistingSettings();
            }
          })
        )
        .subscribe();
    } else if (this.formConfig.mode === FormMode.add && this.form.valid) {
      const command: CheckEmailConnectivityCommand = this.form.value;
      this.checkConnectivitySettings(command);
    } else {
      this.checkConnectivitySettingsWithExistingSettings();
    }
  }

  private checkConnectivitySettingsWithExistingSettings(): void {
    const command: CheckEmailConnectivityCommand = { id: this.originalSystemEmail.id };
    this.checkConnectivitySettings(command);

  }

  private checkConnectivitySettings(command: CheckEmailConnectivityCommand): void {
    this.systemEmailService.checkSystemEmailConnectivity(command)
      .pipe(
        take(1),
        tap((systemTask) => {
          if (systemTask.status === SystemTaskStatus.Succeeded) {
            this.snackbarService.success("Successfully connected to email server");
          } else {
            this.snackbarService.error("Failed to connect to email server");
          }

          if (this.formConfig.mode !== FormMode.add) {
            this.taskTableComponent.getTableData();
          }
        })
      )
      .subscribe();
  }
}
