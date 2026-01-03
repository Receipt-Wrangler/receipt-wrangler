import { AfterViewInit, Component, OnInit, TemplateRef, ViewChild } from "@angular/core";
import { MatDialog } from "@angular/material/dialog";
import { ActivatedRoute } from "@angular/router";
import { Store } from "@ngxs/store";
import { take, tap } from "rxjs";
import { DEFAULT_DIALOG_CONFIG } from "../../constants";
import {
  CheckReceiptProcessingSettingsConnectivityCommand,
  ReceiptProcessingSettings,
  ReceiptProcessingSettingsService,
  SystemSettings,
  SystemTaskStatus
} from "../../open-api";
import { SnackbarService } from "../../services";
import { BaseTableService } from "../../services/base-table.service";
import { BaseTableComponent } from "../../shared-ui/base-table/base-table.component";
import { ConfirmationDialogComponent } from "../../shared-ui/confirmation-dialog/confirmation-dialog.component";
import { ReceiptProcessingSettingsTableState } from "../../store/receipt-processing-settings-table.state";
import { TableColumn } from "../../table/table-column.interface";
import { ReceiptProcessingSettingsTableService } from "./receipt-processing-settings-table.service";

@Component({
    selector: "app-receipt-processing-settings-table",
    templateUrl: "./receipt-processing-settings-table.component.html",
    styleUrl: "./receipt-processing-settings-table.component.scss",
    providers: [
        {
            provide: BaseTableService,
            useClass: ReceiptProcessingSettingsTableService
        }
    ],
    standalone: false
})
export class ReceiptProcessingSettingsTableComponent extends BaseTableComponent<ReceiptProcessingSettings> implements OnInit, AfterViewInit {
  @ViewChild("nameCell") public nameCell!: TemplateRef<any>;

  @ViewChild("descriptionCell") public descriptionCell!: TemplateRef<any>;

  @ViewChild("promptCell") public promptCell!: TemplateRef<any>;

  @ViewChild("aiTypeCell") public aiTypeCell!: TemplateRef<any>;

  @ViewChild("ocrEngineCell") public ocrEngineCell!: TemplateRef<any>;

  @ViewChild("createdAtCell") public createdAtCell!: TemplateRef<any>;

  @ViewChild("updatedAtCell") public updatedAtCell!: TemplateRef<any>;

  @ViewChild("actionsCell") public actionsCell!: TemplateRef<any>;

  public systemSettings!: SystemSettings;

  constructor(
    public override baseTableService: BaseTableService,
    private activatedRoute: ActivatedRoute,
    private dialog: MatDialog,
    private receiptProcessingSettingsService: ReceiptProcessingSettingsService,
    private snackbarService: SnackbarService,
    private store: Store,
  ) {
    super(baseTableService);
  }

  public ngOnInit(): void {
    this.getTableData();
    this.systemSettings = this.activatedRoute.snapshot.data?.["systemSettings"];
  }

  public ngAfterViewInit(): void {
    this.initTable();
  }

  private initTable(): void {
    this.setColumns();
  }

  private setColumns(): void {
    const columns: TableColumn[] = [
      {
        columnHeader: "Name",
        matColumnDef: "name",
        template: this.nameCell,
        sortable: true
      },
      {
        columnHeader: "Description",
        matColumnDef: "description",
        template: this.descriptionCell,
        sortable: true
      },
      {
        columnHeader: "Prompt",
        matColumnDef: "prompt",
        template: this.promptCell,
        sortable: false
      },
      {
        columnHeader: "AI Type",
        matColumnDef: "ai_type",
        template: this.aiTypeCell,
        sortable: true
      },
      {
        columnHeader: "OCR Engine",
        matColumnDef: "ocr_engine",
        template: this.ocrEngineCell,
        sortable: true
      },
      {
        columnHeader: "Created At",
        matColumnDef: "created_at",
        template: this.createdAtCell,
        sortable: true
      },
      {
        columnHeader: "Updated At",
        matColumnDef: "updated_at",
        template: this.updatedAtCell,
        sortable: true
      },
      {
        columnHeader: "Actions",
        matColumnDef: "actions",
        template: this.actionsCell,
        sortable: false
      }
    ] as TableColumn[];
    this.setInitialSortedColumn(this.store.selectSnapshot(ReceiptProcessingSettingsTableState.state), columns);

    this.columns = columns;
    this.displayedColumns = columns.map((c) => c.matColumnDef);
  }

  public deleteReceiptProcessingSettings(receiptProcessingSettings: ReceiptProcessingSettings): void {
    const ref = this.dialog.open(ConfirmationDialogComponent, DEFAULT_DIALOG_CONFIG);

    ref.componentInstance.headerText = "Delete Receipt Processing Settings";
    ref.componentInstance.dialogContent = `Are you sure you want to delete the receipt processing settings: ${receiptProcessingSettings.name}?`;

    ref.afterClosed()
      .pipe(
        take(1),
        tap((confirmed) => {
          if (confirmed) {
            this.callDeleteApi(receiptProcessingSettings);
          }
        })
      )
      .subscribe();

  }

  private callDeleteApi(receiptProcessingSettings: ReceiptProcessingSettings): void {
    this.receiptProcessingSettingsService.deleteReceiptProcessingSettingsById(receiptProcessingSettings.id)
      .pipe(
        take(1),
        tap(() => {
          this.snackbarService.success("Receipt Processing Settings deleted successfully");
          this.getTableData();
        })
      )
      .subscribe();
  }

  public checkConnectivity(id: number): void {
    const command: CheckReceiptProcessingSettingsConnectivityCommand = {
      id: id
    };

    this.receiptProcessingSettingsService.checkReceiptProcessingSettingsConnectivity(command)
      .pipe(
        take(1),
        tap((systemTask) => {
          if (systemTask.status === SystemTaskStatus.Succeeded) {
            this.snackbarService.success("Successfully connected to the server");
          } else {
            this.snackbarService.error("Failed to connect to the server");
          }
        })
      )
      .subscribe();
  }

  public disabledDeleteButtonClicked(settings: ReceiptProcessingSettings, disabled: boolean): void {
    if (disabled) {
      this.snackbarService.info(`Cannot delete ${settings.name} because it is currently active in system settings.`);
    }
  }
}
