import { Component, OnInit } from "@angular/core";
import { FormBuilder, ReactiveFormsModule } from "@angular/forms";
import { MatDialogRef } from "@angular/material/dialog";
import { Store } from "@ngxs/store";
import { BaseFormComponent } from "../../form";
import { FormOption } from "../../interfaces/form-option.interface";
import { AssociatedApiKeys, Permission } from "../../open-api";
import { AuthState } from "../../store";
import { ApiKeyTableState } from "../../store/api-key-table.state";
import { SetFilter } from "../../store/api-key-table.state.actions";
import { associatedApiKeyOptions } from "./associated-api-key-options";
import { DialogComponent } from "../../shared-ui/dialog/dialog.component";
import { SelectComponent } from "../../select/select/select.component";
import { DialogFooterComponent } from "../../shared-ui/dialog-footer/dialog-footer.component";
import { FormGetPipe } from "../../pipes/form-get.pipe";

@Component({
    selector: "app-api-key-table-filter",
    templateUrl: "./api-key-table-filter.component.html",
    styleUrl: "./api-key-table-filter.component.scss",
    imports: [DialogComponent, ReactiveFormsModule, SelectComponent, DialogFooterComponent, FormGetPipe]
})
export class ApiKeyTableFilterComponent extends BaseFormComponent implements OnInit {
  constructor(
    private dialogRef: MatDialogRef<ApiKeyTableFilterComponent>,
    private formBuilder: FormBuilder,
    private store: Store
  ) {
    super();
  }

  protected associatedApiKeyOptions: FormOption[] = associatedApiKeyOptions;

  public ngOnInit(): void {
    this.initOptions();
    this.initForm();
  }

  private initOptions(): void {
    const canViewAll = this.store.selectSnapshot(
      AuthState.hasAppPermission(Permission.AppApiKeysReadAny)
    );

    this.associatedApiKeyOptions = canViewAll
      ? associatedApiKeyOptions
      : associatedApiKeyOptions.filter(
        (option) => option.value !== AssociatedApiKeys.All
      );
  }

  private initForm(): void {
    const filter = this.store.selectSnapshot(ApiKeyTableState.filter);
    this.form = this.formBuilder.group({
        associatedApiKeys: [filter.associatedApiKeys],
      }
    );
  }

  public submit(): void {
    this.store.dispatch(new SetFilter(this.form.value));
    this.dialogRef.close();
  }

  cancel(): void {
    this.dialogRef.close();
  }
}