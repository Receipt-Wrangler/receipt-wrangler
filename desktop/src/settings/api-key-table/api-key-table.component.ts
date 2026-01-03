import { AfterViewInit, Component, OnInit, TemplateRef, ViewChild } from "@angular/core";
import { MatDialog } from "@angular/material/dialog";
import { UntilDestroy, untilDestroyed } from "@ngneat/until-destroy";
import { Store } from "@ngxs/store";
import { take, tap } from "rxjs";
import { TableComponent } from "src/table/table/table.component";
import { DEFAULT_DIALOG_CONFIG, DEFAULT_HOST_CLASS } from "../../constants";
import { ApiKeyService, ApiKeyView, AssociatedApiKeys, UserRole } from "../../open-api";
import { BaseTableService } from "../../services/base-table.service";
import { SnackbarService } from "../../services/snackbar.service";
import { BaseTableComponent } from "../../shared-ui/base-table/base-table.component";
import { ConfirmationDialogComponent } from "../../shared-ui/confirmation-dialog/confirmation-dialog.component";
import { AuthState } from "../../store";
import { ApiKeyTableState } from "../../store/api-key-table.state";
import { ApiKeyFormDialogComponent } from "../api-key-form-dialog/api-key-form-dialog.component";
import { ApiKeyTableFilterComponent } from "../api-key-table-filter/api-key-table-filter.component";
import { ApiKeyTableService } from "./api-key-table.service";


@UntilDestroy()
@Component({
  selector: "app-api-key-table",
  templateUrl: "./api-key-table.component.html",
  styleUrls: ["./api-key-table.component.scss"],
  host: DEFAULT_HOST_CLASS,
  providers: [
    {
      provide: BaseTableService,
      useClass: ApiKeyTableService
    }
  ],
  standalone: false
})
export class ApiKeyTableComponent extends BaseTableComponent<ApiKeyView> implements OnInit, AfterViewInit {
  @ViewChild("nameCell") private nameCell!: TemplateRef<any>;

  @ViewChild("descriptionCell") private descriptionCell!: TemplateRef<any>;

  @ViewChild("createdByCell") private createdByCell!: TemplateRef<any>;

  @ViewChild("createdAtCell") private createdAtCell!: TemplateRef<any>;

  @ViewChild("lastUsedAtCell") private lastUsedAtCell!: TemplateRef<any>;


  @ViewChild("actionsCell") private actionsCell!: TemplateRef<any>;

  @ViewChild(TableComponent) private table!: TableComponent;

  public isAdmin = false;

  public tableHeaderText = "My API Keys";

  private currentUserId = "";

  constructor(
    public override baseTableService: BaseTableService,
    private store: Store,
    private matDialog: MatDialog,
    private apiKeyService: ApiKeyService,
    private snackbarService: SnackbarService
  ) {
    super(baseTableService);
  }

  public ngOnInit(): void {
    this.isAdmin = this.store.selectSnapshot(AuthState.hasRole(UserRole.Admin));
    this.currentUserId = this.store.selectSnapshot(AuthState.userId);
    this.listenForFilterChanges();
  }

  public ngAfterViewInit(): void {
    this.initTable();
  }

  private listenForFilterChanges(): void {
    this.store.select(ApiKeyTableState.filter)
      .pipe(
        untilDestroyed(this),
        tap((filter) => {
            this.getTableData();
            if (filter.associatedApiKeys === AssociatedApiKeys.All) {
              this.tableHeaderText = "All API Keys";
            } else {
              this.tableHeaderText = "My API Keys";
            }
          }
        )
      )
      .subscribe();
  }

  private initTable(): void {
    this.setColumns();
  }

  private setColumns(): void {
    this.columns = [
      {
        columnHeader: "Name",
        matColumnDef: "name",
        template: this.nameCell,
        sortable: true,
      },
      {
        columnHeader: "Description",
        matColumnDef: "description",
        template: this.descriptionCell,
        sortable: true,
      },
      {
        columnHeader: "Created By",
        matColumnDef: "created_by",
        template: this.createdByCell,
        sortable: true,
      },
      {
        columnHeader: "Created At",
        matColumnDef: "created_at",
        template: this.createdAtCell,
        sortable: true,
      },
      {
        columnHeader: "Last Used At",
        matColumnDef: "last_used_at",
        template: this.lastUsedAtCell,
        sortable: true,
      },
      {
        columnHeader: "Actions",
        matColumnDef: "actions",
        template: this.actionsCell,
        sortable: false,
      },
    ];
    this.displayedColumns = [
      "name",
      "description",
      "created_by",
      "created_at",
      "last_used_at",
      "actions",
    ];
  }

  public isOwner(apiKey: ApiKeyView): boolean {
    return apiKey.userId?.toString() === this.currentUserId;
  }

  public openFilterDialog(): void {
    const ref = this.matDialog.open(ApiKeyTableFilterComponent, DEFAULT_DIALOG_CONFIG);
  }

  public openCreateApiKeyDialog(): void {
    const ref = this.matDialog.open(ApiKeyFormDialogComponent, DEFAULT_DIALOG_CONFIG);

    ref.componentInstance.headerText = "Create API Key";

    ref.afterClosed().subscribe((result) => {
      if (result) {
        this.getTableData();
      }
    });
  }

  public openEditDialog(apiKeyView: ApiKeyView): void {
    const ref = this.matDialog.open(ApiKeyFormDialogComponent, DEFAULT_DIALOG_CONFIG);

    ref.componentInstance.apiKey = apiKeyView;
    ref.componentInstance.headerText = `Edit ${apiKeyView.name}`;

    ref.afterClosed().subscribe((result) => {
      if (result) {
        this.getTableData();
      }
    });
  }

  public deleteApiKey(apiKeyView: ApiKeyView): void {
    if (!apiKeyView.id) {
      this.snackbarService.error("Cannot delete API key: ID is missing");
      return;
    }

    const dialogRef = this.matDialog.open(
      ConfirmationDialogComponent,
      DEFAULT_DIALOG_CONFIG
    );

    dialogRef.componentInstance.headerText = "Delete API Key";
    dialogRef.componentInstance.dialogContent = `Are you sure you would like to delete the API key "${apiKeyView.name}"? This action is irreversible and the API key will no longer be usable.`;

    dialogRef.afterClosed().subscribe((confirmed) => {
      if (confirmed && apiKeyView.id) {
        this.apiKeyService
          .deleteApiKey(apiKeyView.id)
          .pipe(
            take(1),
            tap(() => {
              this.snackbarService.success(`API key "${apiKeyView.name}" successfully deleted`);
              this.getTableData();
            })
          )
          .subscribe();
      }
    });
  }
}
