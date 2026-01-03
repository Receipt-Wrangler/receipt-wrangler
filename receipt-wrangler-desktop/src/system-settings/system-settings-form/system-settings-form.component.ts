import { AfterViewInit, Component, ElementRef, OnInit, ViewChild } from "@angular/core";
import { FormArray, FormBuilder, FormGroup, Validators } from "@angular/forms";
import { ActivatedRoute, Router } from "@angular/router";
import { UntilDestroy, untilDestroyed } from "@ngneat/until-destroy";
import { Store } from "@ngxs/store";
import { startWith, switchMap, take, tap } from "rxjs";
import { fadeInOut } from "../../animations/index";
import { AutocomleteComponent } from "../../autocomplete/autocomlete/autocomlete.component";
import { BaseFormComponent } from "../../form";
import { FormOption } from "../../interfaces/form-option.interface";
import {
  CurrencySeparator,
  CurrencySymbolPosition,
  FeatureConfigService,
  QueueName,
  ReceiptProcessingSettings,
  SystemSettings,
  SystemSettingsService
} from "../../open-api";
import { InputReadonlyPipe } from "../../pipes/input-readonly.pipe";
import { SnackbarService } from "../../services";
import { SetFeatureConfig } from "../../store";
import { SetCurrencyData, SetCurrencyDisplay } from "../../store/system-settings.state.actions";

interface QueueData extends FormOption {
  description: string;
}

@UntilDestroy()
@Component({
  selector: "app-system-settings-form",
  templateUrl: "./system-settings-form.component.html",
  styleUrl: "./system-settings-form.component.scss",
  providers: [InputReadonlyPipe],
  animations: [fadeInOut],
  standalone: false
})
export class SystemSettingsFormComponent extends BaseFormComponent implements OnInit, AfterViewInit {
  @ViewChild("fallbackReceiptProcessingSettings")
  public fallbackReceiptProcessingSettings!: AutocomleteComponent;

  @ViewChild("alert", { static: false, read: ElementRef })
  public alert!: ElementRef;

  public originalSystemSettings!: SystemSettings;

  public allReceiptProcessingSettings: ReceiptProcessingSettings[] = [];

  public filteredReceiptProcessingSettings: ReceiptProcessingSettings[] = [];

  public showRestartTaskServerAlert = false;

  public readonly queueData: QueueData[] = [
    {
      value: QueueName.EmailPolling,
      displayValue: "Email Polling",
      description: "Polls system emails for receipts to process"
    },
    {
      value: QueueName.EmailReceiptProcessing,
      displayValue: "Email Receipt Processing",
      description: "Processing of captured emails"
    },
    {
      value: QueueName.EmailReceiptImageCleanup,
      displayValue: "Email Receipt Image Cleanup",
      description: "Cleans up email receipt images after all processing is done"
    },
    {
      value: QueueName.QuickScan,
      displayValue: "Quick Scan",
      description: "Processes quick scan receipts"
    }
  ];

  public readonly symbolPositions: FormOption[] = [
    {
      displayValue: "Start",
      value: CurrencySymbolPosition.Start
    },
    {
      displayValue: "End",
      value: CurrencySymbolPosition.End,
    }
  ];

  public readonly decimalSeparators: FormOption[] = [
    {
      displayValue: ", (Comma)",
      value: CurrencySeparator.Comma
    },
    {
      displayValue: ". (Dot)",
      value: CurrencySeparator.Period
    }
  ];

  constructor(
    private activatedRoute: ActivatedRoute,
    private featureConfigService: FeatureConfigService,
    private formBuilder: FormBuilder,
    private router: Router,
    private snackbarService: SnackbarService,
    private store: Store,
    private systemSettingsService: SystemSettingsService,
    private inputReadonlyPipe: InputReadonlyPipe,
  ) {
    super();
  }

  public ngOnInit(): void {
    this.setFormConfigFromRoute(this.activatedRoute);
    this.allReceiptProcessingSettings = this.activatedRoute.snapshot.data?.["allReceiptProcessingSettings"];
    this.originalSystemSettings = this.activatedRoute.snapshot.data?.["systemSettings"];
    this.showRestartTaskServerAlert = this.activatedRoute.snapshot?.queryParams?.["restartTaskServer"] === "true";
    this.initForm();
  }

  public ngAfterViewInit(): void {
    setTimeout(() => {

      if (this.showRestartTaskServerAlert) {
        this.alert.nativeElement.scrollIntoView({ behavior: "smooth" });
      }

    }, 0);
  }

  private initForm(): void {
    this.form = this.formBuilder.group({
      enableLocalSignUp: [this.originalSystemSettings?.enableLocalSignUp],
      debugOcr: [this.originalSystemSettings?.debugOcr],
      emailPollingInterval: [this.originalSystemSettings?.emailPollingInterval, [Validators.required, Validators.min(0)]],
      currencyDisplay: [this.originalSystemSettings?.currencyDisplay],
      currencyThousandthsSeparator: [this.originalSystemSettings.currencyThousandthsSeparator, [Validators.required]],
      currencyDecimalSeparator: [this.originalSystemSettings.currencyDecimalSeparator, [Validators.required]],
      currencySymbolPosition: [this.originalSystemSettings.currencySymbolPosition, [Validators.required]],
      currencyHideDecimalPlaces: [this.originalSystemSettings.currencyHideDecimalPlaces],
      receiptProcessingSettingsId: [this.originalSystemSettings?.receiptProcessingSettingsId],
      fallbackReceiptProcessingSettingsId: [this.originalSystemSettings?.fallbackReceiptProcessingSettingsId],
      taskConcurrency: [this.originalSystemSettings?.taskConcurrency, [Validators.min(0), Validators.required]],
      taskQueueConfigurations: this.formBuilder.array(this.buildAsynqQueueConfigurations())
    });

    if (this.inputReadonlyPipe.transform(this.formConfig.mode)) {
      this.form.get("debugOcr")?.disable();
      this.form.get("enableLocalSignUp")?.disable();
      this.form.get("currencyThousandthsSeparator")?.disable();
      this.form.get("currencyDecimalSeparator")?.disable();
      this.form.get("currencySymbolPosition")?.disable();
      this.form.get("currencyHideDecimalPlaces")?.disable();
    }

    this.listenForReceiptProcessingSettingsChanges();
    this.listenForHideDecimalPlacesChanges();
  }

  // TODO: finish implementing UI for taskQueueConfigurations
  private buildAsynqQueueConfigurations(): FormGroup[] {
    return (this.originalSystemSettings?.taskQueueConfigurations ?? []).map(config => {
      return this.formBuilder.group({
        name: [config.name],
        priority: [config.priority],
      });
    });
  }

  private listenForReceiptProcessingSettingsChanges(): void {
    this.form.get("receiptProcessingSettingsId")?.valueChanges
      .pipe(
        startWith(this.form.get("receiptProcessingSettingsId")?.value),
        untilDestroyed(this),
        tap((value: number) => {
          this.filteredReceiptProcessingSettings = this.allReceiptProcessingSettings.filter((rps) => rps.id !== value);

          if (!value) {
            this.fallbackReceiptProcessingSettings?.clearFilter();
          }
        })
      )
      .subscribe();
  }

  private listenForHideDecimalPlacesChanges(): void {
    this.form.get("currencyHideDecimalPlaces")?.valueChanges
      .pipe(
        startWith(this.form.get("currencyHideDecimalPlaces")?.value),
        untilDestroyed(this),
        tap((hide: boolean) => {
          if (hide) {
            this.form.get("currencyDecimalSeparator")?.disable();
          } else {
            this.form.get("currencyDecimalSeparator")?.enable();
          }
        })
      )
      .subscribe();
  }

  public displayWith(id: number): string {
    return this.allReceiptProcessingSettings.find((rps) => rps.id === id)?.name ?? "";
  }

  private doesTaskServerRequiresRestart(): boolean {
    let requiresRestart = this.originalSystemSettings.taskConcurrency !== this.form.get("taskConcurrency")?.value;
    for (let i = 0; i < this.originalSystemSettings.taskQueueConfigurations.length; i++) {
      const originalConfig = this.originalSystemSettings.taskQueueConfigurations[i];
      const formConfig = (this.form.get("taskQueueConfigurations") as FormArray).controls.find((control) => control.get("name")?.value === originalConfig.name);

      if (originalConfig.priority !== formConfig?.get("priority")?.value) {
        requiresRestart = true;
        break;
      }
    }

    return requiresRestart;
  }

  public submit(): void {
    const formValue = this.form.getRawValue();
    formValue["emailPollingInterval"] = Number.parseInt(formValue["emailPollingInterval"]);
    formValue["taskConcurrency"] = Number.parseInt(formValue["taskConcurrency"]);
    (formValue["taskQueueConfigurations"] as Array<any>).forEach(config => {
      config.priority = Number.parseInt(config.priority);
    });
    const restartTaskServer = this.doesTaskServerRequiresRestart();

    this.systemSettingsService.updateSystemSettings(formValue)
      .pipe(
        take(1),
        tap(() => {
          this.snackbarService.success("System settings updated successfully");
          this.router.navigate(["/system-settings/settings/view"], {
            queryParams: {
              restartTaskServer: restartTaskServer
            }
          });
        }),
        switchMap(() => this.featureConfigService.getFeatureConfig()),
        tap((featureConfig) => this.store.dispatch(new SetFeatureConfig(featureConfig))),
        switchMap(() => this.store.dispatch(new SetCurrencyDisplay(formValue["currencyDisplay"]?.toString()))),
        switchMap(() => this.store.dispatch(
          new SetCurrencyData(formValue["currencySymbolPosition"],
            formValue["currencyDecimalSeparator"],
            formValue["currencyThousandthsSeparator"],
            formValue["currencyHideDecimalPlaces"]
          ))),
      )
      .subscribe();
  }

  public restartTaskServer(): void {
    this.systemSettingsService.restartTaskServer().pipe(
      take(1),
      tap(() => {
          this.snackbarService.success("Task server restarted successfully");
          this.showRestartTaskServerAlert = false;
        },
      )).subscribe();
  }
}
