import { Component, OnInit, ViewChild } from "@angular/core";
import { FormBuilder, Validators } from "@angular/forms";
import { MatDialog } from "@angular/material/dialog";
import { ActivatedRoute, Router } from "@angular/router";
import { UntilDestroy, untilDestroyed } from "@ngneat/until-destroy";
import { startWith, take, tap } from "rxjs";
import { DEFAULT_DIALOG_CONFIG } from "../../constants";
import { FormMode } from "../../enums/form-mode.enum";
import { BaseFormComponent } from "../../form";
import { FormOption } from "../../interfaces/form-option.interface";
import {
  AiType,
  AssociatedEntityType,
  CheckReceiptProcessingSettingsConnectivityCommand,
  Prompt,
  ReceiptProcessingSettings,
  ReceiptProcessingSettingsService,
  SystemTaskStatus,
  SystemTaskType
} from "../../open-api";
import { SnackbarService } from "../../services";
import { TABLE_SERVICE_INJECTION_TOKEN } from "../../services/injection-tokens/table-service";
import { ReceiptProcessingSettingsTaskTableService } from "../../services/receipt-processing-settings-task-table.service";
import { ConfirmationDialogComponent } from "../../shared-ui/confirmation-dialog/confirmation-dialog.component";
import { TaskTableComponent } from "../../shared-ui/task-table/task-table.component";
import { aiTypeOptions } from "../constants/ai-type-options";
import { ocrEngineOptions } from "../constants/ocr-engine-options";

@UntilDestroy()
@Component({
    selector: "app-receipt-processing-settings-form",
    templateUrl: "./receipt-processing-settings-form.component.html",
    styleUrl: "./receipt-processing-settings-form.component.scss",
    providers: [{
            provide: TABLE_SERVICE_INJECTION_TOKEN,
            useClass: ReceiptProcessingSettingsTaskTableService
        }],
    standalone: false
})
export class ReceiptProcessingSettingsFormComponent extends BaseFormComponent implements OnInit {
  @ViewChild(TaskTableComponent) public taskTableComponent!: TaskTableComponent;

  public originalReceiptProcessingSettings?: ReceiptProcessingSettings;

  protected readonly AiType = AiType;

  protected readonly openAiGeminiSpecificFields: string[] = ["key", "model", "isVisionModel"];

  protected readonly openAiCustomSpecificFields: string[] = ["url", "model"];

  protected readonly ollamaSpecificFields = ["url", "model", "isVisionModel"];

  public readonly aiTypeOptions: FormOption[] = aiTypeOptions;

  public readonly ocrEngineOptions: FormOption[] = ocrEngineOptions;

  protected readonly AssociatedEntityType = AssociatedEntityType;

  public prompts: Prompt[] = [];

  constructor(
    private activatedRoute: ActivatedRoute,
    private dialog: MatDialog,
    private formBuilder: FormBuilder,
    private receiptProcessingSettingsService: ReceiptProcessingSettingsService,
    private router: Router,
    private snackbarService: SnackbarService,
  ) {
    super();
  }


  public ngOnInit(): void {
    this.originalReceiptProcessingSettings = this.activatedRoute?.snapshot?.data?.["receiptProcessingSettings"];
    this.prompts = this.activatedRoute.snapshot.data["prompts"];
    this.setFormConfigFromRoute(this.activatedRoute);
    this.initForm();
  }


  private initForm(): void {
    this.form = this.formBuilder.group({
      name: [this.originalReceiptProcessingSettings?.name, Validators.required],
      description: [this.originalReceiptProcessingSettings?.description],
      ocrEngine: [this.originalReceiptProcessingSettings?.ocrEngine, Validators.required],
      aiType: [this.originalReceiptProcessingSettings?.aiType, Validators.required],
      promptId: [this.originalReceiptProcessingSettings?.promptId, Validators.required],
      key: [this.originalReceiptProcessingSettings?.key],
      url: [this.originalReceiptProcessingSettings?.url],
      model: [this.originalReceiptProcessingSettings?.model],
      isVisionModel: [this.originalReceiptProcessingSettings?.isVisionModel],
    });

    this.listenForTypeChange();
    this.listenForIsVisionModelChange();

    if (this.formConfig.mode === FormMode.view) {
      this.form.get("ocrEngine")?.disable();
      this.form.get("aiType")?.disable();
      this.form.get("promptId")?.disable();
      this.form.get("isVisionModel")?.disable();
    }
  }

  private listenForIsVisionModelChange(): void {
    this.form.get("isVisionModel")?.valueChanges
      .pipe(
        untilDestroyed(this),
        startWith(this.form.get("isVisionModel")?.value),
        tap((isVisionModel: boolean) => {
          if (isVisionModel) {
            this.form.get("ocrEngine")?.disable();
          } else {
            this.form.get("ocrEngine")?.enable();
          }
        })).subscribe();
  }

  private listenForTypeChange(): void {
    this.form.get("aiType")?.valueChanges
      .pipe(
        untilDestroyed(this),
        tap((type: AiType) => {
          switch (type) {
            case AiType.Gemini:
            case AiType.OpenAi:
              this.updateOpenAiGeminiForm();
              break;
            case AiType.OpenAiCustom:
              this.updateOpenAiCustomForm();
              break;
            case AiType.Ollama:
              this.updateOllamaForm();
              break;
          }
        })
      )
      .subscribe();
  }

  private updateOpenAiGeminiForm(): void {
    this.openAiCustomSpecificFields.forEach((field: string) => {
      this.form.get(field)?.setValidators(null);
      this.form.get(field)?.setErrors(null);
    });

    const originalType = this.originalReceiptProcessingSettings?.aiType;
    if (this.formConfig.mode === FormMode.edit && (originalType !== AiType.OpenAi || originalType !== AiType.Gemini)) {
      this.form.get("key")?.setValidators(Validators.required);
    }

    if (this.form?.get("aiType")?.value === AiType.Gemini) {
      this.form.get("isVisionModel")?.setValue(false);
    }
  }

  private updateOpenAiCustomForm(): void {
    this.openAiGeminiSpecificFields.forEach((field: string) => {
      this.form.get(field)?.setValidators(null);
      this.form.get(field)?.setErrors(null);
    });

    this.form.get("url")?.setValidators(Validators.required);
  }

  private updateOllamaForm(): void {
    this.ollamaSpecificFields.forEach((field: string) => {
      this.form.get(field)?.setValidators(null);
      this.form.get(field)?.setErrors(null);
    });

    this.form.get("url")?.setValidators(Validators.required);
  }

  public promptDisplayWith(promptId: number): string {
    if (promptId) {
      const prompt = this.prompts.find((p) => p.id === promptId);
      return prompt?.name ?? "";
    }

    return "";

  }

  public submit(): void {
    if (this.form.valid && !this.originalReceiptProcessingSettings) {
      this.createSettings();
    } else if (this.form.valid) {
      this.updateSettings();
    }
  }

  private createSettings(): void {
    this.receiptProcessingSettingsService.createReceiptProcessingSettings(this.form.value)
      .pipe(
        take(1),
        tap((settings) => {
          this.router.navigate([`system-settings/receipt-processing-settings/${settings.id}/view`]);
          this.snackbarService.success("Receipt processing settings created successfully");
        })
      ).subscribe();

  }

  private updateSettings(): void {
    if (this.form.get("key")?.dirty) {
      this.showConfirmationDialog();
    } else {
      this.callUpdateApi(false);
    }
  }

  private showConfirmationDialog(): void {
    const ref = this.dialog.open(ConfirmationDialogComponent, DEFAULT_DIALOG_CONFIG);

    ref.componentInstance.headerText = "Update Key";
    ref.componentInstance.dialogContent = "Are you sure you want to update the key? This will replace the previous key.";

    ref.afterClosed()
      .pipe(
        take(1),
        tap((confirmed) => {
          this.callUpdateApi(confirmed);

        })
      )
      .subscribe();
  }

  private callUpdateApi(updateKey: boolean): void {
    this.receiptProcessingSettingsService.updateReceiptProcessingSettingsById(
      this.originalReceiptProcessingSettings?.id ?? 0,
      updateKey,
      this.form.value)
      .pipe(
        take(1),
        tap(() => {
          this.router.navigate([`system-settings/receipt-processing-settings/${this.originalReceiptProcessingSettings?.id}/view`]);
          this.snackbarService.success("Receipt processing settings updated successfully");
        })
      )
      .subscribe();
  }

  public checkConnectivity(): void {
    if (this.form.valid && this.formConfig.mode === FormMode.edit && this.form.dirty) {
      const ref = this.dialog.open(ConfirmationDialogComponent, DEFAULT_DIALOG_CONFIG);
      ref.componentInstance.headerText = "Check Connectivity";
      ref.componentInstance.dialogContent = `You have made changes to the receipt processing settings.
       Would you like to check connectivity with the unsaved changes? Otherwise the existing settings will be used.`;

      ref.afterClosed()
        .pipe(
          take(1),
          tap((useUnsavedChanges) => {
            if (useUnsavedChanges) {
              this.callCheckConnectivityApi({
                ...this.form.value,
                id: this.originalReceiptProcessingSettings?.id
              }, true);
            } else {
              this.callCheckConnectivityApi({
                id: this.originalReceiptProcessingSettings?.id
              }, true);
            }

          })
        )
        .subscribe();
      return;
    }

    if (this.form.valid && (this.formConfig.mode === FormMode.edit || this.formConfig.mode === FormMode.view)) {
      const command: CheckReceiptProcessingSettingsConnectivityCommand = {
        id: this.originalReceiptProcessingSettings?.id
      };
      this.callCheckConnectivityApi(command, true);
    } else if (this.form.valid) {
      this.callCheckConnectivityApi(this.form.value);
    }
  }

  private callCheckConnectivityApi(command: CheckReceiptProcessingSettingsConnectivityCommand, refreshSystemTaskTable = false): void {
    this.receiptProcessingSettingsService.checkReceiptProcessingSettingsConnectivity(command)
      .pipe(
        take(1),
        tap((systemTask) => {
          if (systemTask.status === SystemTaskStatus.Succeeded) {
            this.snackbarService.success("Successfully connected to the server");
          } else {
            this.snackbarService.error("Failed to connect to the server");
          }

          if (refreshSystemTaskTable) {
            this.taskTableComponent.getTableData();
          }
        })
      ).subscribe();
  }

  protected readonly SystemTaskType = SystemTaskType;
}
