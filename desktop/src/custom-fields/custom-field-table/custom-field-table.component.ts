import { AfterViewInit, Component, OnInit, signal, TemplateRef, viewChild } from "@angular/core";
import { MatDialog } from "@angular/material/dialog";
import { PageEvent } from "@angular/material/paginator";
import { Sort } from "@angular/material/sort";
import { MatTableDataSource } from "@angular/material/table";
import { Store } from "@ngxs/store";
import { of, switchMap, take, tap } from "rxjs";
import { CustomField, CustomFieldService, PagedDataDataInner, PagedRequestCommand, Permission } from "src/open-api";
import { ConfirmationDialogComponent } from "src/shared-ui/confirmation-dialog/confirmation-dialog.component";
import { TableComponent } from "src/table/table/table.component";
import { DEFAULT_DIALOG_CONFIG } from "../../constants/index";
import { FormMode } from "../../enums/form-mode.enum";
import { SnackbarService } from "../../services/index";
import { CustomFieldTableState } from "../../store/custom-field-table.state";
import { SetOrderBy, SetPage, SetPageSize, SetSortDirection } from "../../store/custom-field-table.state.actions";
import { AuthState } from "../../store/index";
import { TableColumn } from "../../table/table-column.interface";
import { CustomFieldFormComponent } from "../custom-field-form/custom-field-form.component";

@Component({
  selector: "app-custom-field-table",
  templateUrl: "./custom-field-table.component.html",
  styleUrl: "./custom-field-table.component.scss",
  standalone: false
})
export class CustomFieldTableComponent implements OnInit, AfterViewInit {
  public readonly nameCell = viewChild.required<TemplateRef<any>>("nameCell");

  public readonly typeCell = viewChild.required<TemplateRef<any>>("typeCell");

  public readonly descriptionCell = viewChild.required<TemplateRef<any>>("descriptionCell");

  public readonly actionsCell = viewChild.required<TemplateRef<any>>("actionsCell");

  public readonly table = viewChild.required(TableComponent);

  public state = this.store.selectSignal(CustomFieldTableState.state);

  public dataSource = signal(new MatTableDataSource<PagedDataDataInner>([]));

  public displayedColumns: string[] = [];

  public columns: TableColumn[] = [];

  public totalCount = signal(0);

  protected readonly Permission = Permission;

  protected readonly FormMode = FormMode;

  private readonly canUpdate = this.store.selectSignal(
    AuthState.hasAppPermission(Permission.AppCustomFieldsUpdate)
  );

  constructor(
    private customFieldService: CustomFieldService,
    private matDialog: MatDialog,
    private snackbarService: SnackbarService,
    private store: Store,
  ) {}

  public ngOnInit(): void {
    this.initTableData();
  }

  public ngAfterViewInit(): void {
    this.initTable();
  }

  public updatePageData(pageEvent: PageEvent) {
    const newPage = pageEvent.pageIndex + 1;

    this.store.dispatch(new SetPage(newPage));
    this.store.dispatch(new SetPageSize(pageEvent.pageSize));

    this.getCustomFields();
  }

  public sorted({ sortState }: { sortState: Sort }): void {
    this.store.dispatch(new SetOrderBy(sortState.active));
    this.store.dispatch(new SetSortDirection(sortState.direction));

    this.getCustomFields();
  }

  public openCustomFieldDialog(customField?: CustomField, mode?: FormMode): void {
    const dialogRef = this.matDialog.open(CustomFieldFormComponent, DEFAULT_DIALOG_CONFIG);
    const resolvedMode = mode ?? this.resolveDialogMode(customField);

    dialogRef.componentInstance.headerText = this.buildDialogHeaderText(customField, resolvedMode);
    dialogRef.componentInstance.customField = customField;
    dialogRef.componentInstance.mode = resolvedMode;

    dialogRef
      .afterClosed()
      .pipe(
        take(1),
        tap((refreshData) => {
          if (refreshData) {
            this.getCustomFields();
          }
        })
      )
      .subscribe();
  }

  // Clicking a custom field's name opens the editor for anyone who may edit it,
  // and the read-only view for everyone else.
  private resolveDialogMode(customField?: CustomField): FormMode {
    if (!customField) {
      return FormMode.add;
    }

    return this.canUpdate() ? FormMode.edit : FormMode.view;
  }

  private buildDialogHeaderText(customField: CustomField | undefined, mode: FormMode): string {
    if (!customField) {
      return "Add Custom Field";
    }

    return mode === FormMode.edit ? `Edit ${customField.name}` : "View Custom Field";
  }

  private initTableData(): void {
    this.getCustomFields();
  }

  private getCustomFields(): void {
    const command: PagedRequestCommand = this.store.selectSnapshot(
      CustomFieldTableState.state
    );

    this.customFieldService
      .getPagedCustomFields(command)
      .pipe(
        take(1),
        tap((pagedData) => {
          this.dataSource.set(new MatTableDataSource<PagedDataDataInner>(
            pagedData.data
          ));
          this.totalCount.set(pagedData.totalCount);
        })
      )
      .subscribe();
  }

  private initTable(): void {
    this.setColumns();
  }

  private setColumns(): void {
    const columns = [
      {
        columnHeader: "Name",
        matColumnDef: "name",
        template: this.nameCell(),
        sortable: true,
      },
      {
        columnHeader: "Type",
        matColumnDef: "type",
        template: this.typeCell(),
        sortable: true,
      },
      {
        columnHeader: "Description",
        matColumnDef: "description",
        template: this.descriptionCell(),
        sortable: true,
      },
      {
        columnHeader: "Actions",
        matColumnDef: "actions",
        template: this.actionsCell(),
        sortable: false,
      }
    ] as TableColumn[];

    const tableState = this.store.selectSnapshot(CustomFieldTableState.state);
    if (tableState.orderBy) {
      const column = columns.find(c => c.matColumnDef === tableState.orderBy);
      if (column) {
        column.defaultSortDirection = tableState.sortDirection;
      }
    }


    this.columns = columns;
    this.displayedColumns = [
      "name",
      "type",
      "description",
    ];

    const canUseRowActions = this.store.selectSnapshot(
      AuthState.hasAnyAppPermission([
        Permission.AppCustomFieldsUpdate,
        Permission.AppCustomFieldsDelete,
      ])
    );
    if (canUseRowActions) {
      this.displayedColumns.push("actions");
    }
  }

  public openDeleteConfirmationDialog(customField: CustomField) {
    const dialogRef = this.matDialog.open(
      ConfirmationDialogComponent,
      DEFAULT_DIALOG_CONFIG
    );

    dialogRef.componentInstance.headerText = `Delete ${customField.name}`;
    dialogRef.componentInstance.dialogContent = `Are you sure you want to delete ${customField.name}? This action is irreversible and will remove this custom field from the receipts it is associated with.`;

    dialogRef
      .afterClosed()
      .pipe(
        take(1),
        switchMap((confirmed) => {
          if (confirmed) {
            return this.customFieldService.deleteCustomField(customField.id).pipe(
              tap(() => {
                this.snackbarService.success("Custom field successfully deleted");
                this.getCustomFields();
              })
            );
          } else {
            return of(undefined);
          }
        })
      )
      .subscribe();
  }
}
