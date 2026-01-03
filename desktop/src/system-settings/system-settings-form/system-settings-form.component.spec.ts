import { provideHttpClientTesting } from "@angular/common/http/testing";
import { CUSTOM_ELEMENTS_SCHEMA } from "@angular/core";
import { ComponentFixture, TestBed } from "@angular/core/testing";
import { FormArray, FormControl, FormGroup, ReactiveFormsModule } from "@angular/forms";
import { NoopAnimationsModule } from "@angular/platform-browser/animations";
import { ActivatedRoute, Router } from "@angular/router";
import { NgxsModule, Store } from "@ngxs/store";
import { of } from "rxjs";
import { AutocompleteModule } from "../../autocomplete/autocomplete.module";
import { CheckboxModule } from "../../checkbox/checkbox.module";
import { InputModule } from "../../input/index";
import { CurrencySeparator, CurrencySymbolPosition, QueueName, SystemSettingsService } from "../../open-api";
import { PipesModule } from "../../pipes";
import { CustomCurrencyPipe } from "../../pipes/custom-currency.pipe";
import { SnackbarService } from "../../services";
import { SharedUiModule } from "../../shared-ui/shared-ui.module";
import { SystemSettingsState } from "../../store/system-settings.state";
import { TaskQueueFormControlPipe } from "../pipes/task-queue-form-control.pipe";

import { SystemSettingsFormComponent } from "./system-settings-form.component";
import { provideHttpClient, withInterceptorsFromDi } from "@angular/common/http";

describe("SystemSettingsFormComponent", () => {
  let component: SystemSettingsFormComponent;
  let fixture: ComponentFixture<SystemSettingsFormComponent>;
  let store: Store;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
    declarations: [SystemSettingsFormComponent, CustomCurrencyPipe, TaskQueueFormControlPipe],
    schemas: [CUSTOM_ELEMENTS_SCHEMA],
    imports: [AutocompleteModule,
        CheckboxModule,
        InputModule,
        NgxsModule.forRoot([SystemSettingsState]),
        PipesModule,
        ReactiveFormsModule,
        SharedUiModule,
        NoopAnimationsModule],
    providers: [
        CustomCurrencyPipe,
        {
            provide: ActivatedRoute,
            useValue: {
                snapshot: {
                    data: {
                        allReceiptProcessingSettings: [],
                        systemSettings: {
                            taskQueueConfigurations: [
                                { name: "email_polling", priority: 1 },
                                { name: "email_receipt_processing", priority: 1 },
                                { name: "email_receipt_image_cleanup", priority: 1 },
                                { name: "quick_scan", priority: 1 }
                            ]
                        },
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

    store = TestBed.inject(Store);
    fixture = TestBed.createComponent(SystemSettingsFormComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it("should create", () => {
    expect(component).toBeTruthy();
  });

  it("init form with no data", () => {
    component.ngOnInit();

    expect(component.form.value).toEqual({
      enableLocalSignUp: null,
      debugOcr: null,
      currencyDisplay: null,
      emailPollingInterval: null,
      receiptProcessingSettingsId: null,
      fallbackReceiptProcessingSettingsId: null,
      currencyThousandthsSeparator: null,
      currencyDecimalSeparator: null,
      currencySymbolPosition: null,
      currencyHideDecimalPlaces: null,
      taskConcurrency: null,
      taskQueueConfigurations: [
        { name: "email_polling", priority: 1 },
        { name: "email_receipt_processing", priority: 1 },
        { name: "email_receipt_image_cleanup", priority: 1 },
        { name: "quick_scan", priority: 1 }
      ]
    });
  });

  it("init form with data", () => {
    const activatedRoute = TestBed.inject(ActivatedRoute);
    activatedRoute.snapshot.data["systemSettings"] = {
      enableLocalSignUp: true,
      debugOcr: true,
      currencyDisplay: "USD",
      emailPollingInterval: 5,
      receiptProcessingSettingsId: 1,
      fallbackReceiptProcessingSettingsId: 2,
      currencyThousandthsSeparator: CurrencySeparator.Comma,
      currencyDecimalSeparator: CurrencySeparator.Period,
      currencySymbolPosition: CurrencySymbolPosition.Start,
      currencyHideDecimalPlaces: true,
      taskConcurrency: 12,
      taskQueueConfigurations: [{
        name: QueueName.QuickScan,
        priority: 1,
      }]
    };

    component.ngOnInit();

    expect(component.form.getRawValue()).toEqual({
      enableLocalSignUp: true,
      debugOcr: true,
      currencyDisplay: "USD",
      emailPollingInterval: 5,
      receiptProcessingSettingsId: 1,
      fallbackReceiptProcessingSettingsId: 2,
      currencyThousandthsSeparator: CurrencySeparator.Comma,
      currencyDecimalSeparator: CurrencySeparator.Period,
      currencySymbolPosition: CurrencySymbolPosition.Start,
      currencyHideDecimalPlaces: true,
      taskConcurrency: 12,
      taskQueueConfigurations: [{
        name: QueueName.QuickScan,
        priority: 1,
      }]
    });
  });

  it("should submit form", () => {
    const systemSettingsService = TestBed.inject(SystemSettingsService);
    const snackbarService = TestBed.inject(SnackbarService);
    const router = TestBed.inject(Router);

    const updateSystemSettingsSpy = jest.spyOn(systemSettingsService, "updateSystemSettings").mockReturnValue(of(null as any));
    const snackbarServiceSpy = jest.spyOn(snackbarService, "success");
    const routerSpy = jest.spyOn(router, "navigate");

    component.originalSystemSettings.taskQueueConfigurations = [{
      name: QueueName.QuickScan,
      priority: "1",
    } as any];

    component.form.patchValue({
      enableLocalSignUp: true,
      debugOcr: true,
      currencyDisplay: "USD",
      emailPollingInterval: "5",
      receiptProcessingSettingsId: 1,
      fallbackReceiptProcessingSettingsId: 2,
      currencyThousandthsSeparator: CurrencySeparator.Comma,
      currencyDecimalSeparator: CurrencySeparator.Period,
      currencySymbolPosition: CurrencySymbolPosition.Start,
      currencyHideDecimalPlaces: false,
      taskConcurrency: "12"
    });

    // Update the quick_scan queue priority specifically
    const queueArray = component.form.get("taskQueueConfigurations") as FormArray;
    const quickScanIndex = queueArray.controls.findIndex(control => 
      control.get('name')?.value === 'quick_scan'
    );
    if (quickScanIndex >= 0) {
      queueArray.at(quickScanIndex).get('priority')?.setValue("1");
    }

    component.submit();

    expect(updateSystemSettingsSpy).toHaveBeenCalledWith({
      enableLocalSignUp: true,
      debugOcr: true,
      currencyDisplay: "USD",
      emailPollingInterval: 5,
      receiptProcessingSettingsId: 1,
      fallbackReceiptProcessingSettingsId: 2,
      currencyThousandthsSeparator: CurrencySeparator.Comma,
      currencyDecimalSeparator: CurrencySeparator.Period,
      currencySymbolPosition: CurrencySymbolPosition.Start,
      currencyHideDecimalPlaces: false,
      taskConcurrency: 12,
      taskQueueConfigurations: [
        { name: 'email_polling', priority: 1 },
        { name: 'email_receipt_processing', priority: 1 },
        { name: 'email_receipt_image_cleanup', priority: 1 },
        { name: 'quick_scan', priority: 1 }
      ]
    });

    expect(snackbarServiceSpy).toHaveBeenCalled();
    expect(routerSpy).toHaveBeenCalled();
  });
});
