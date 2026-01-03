import { Component, OnInit } from "@angular/core";
import { FormBuilder } from "@angular/forms";
import { MatDialogRef } from "@angular/material/dialog";
import { Store } from "@ngxs/store";
import { BaseFormComponent } from "../../form";
import { ApiKeyTableState } from "../../store/api-key-table.state";
import { SetFilter } from "../../store/api-key-table.state.actions";
import { associatedApiKeyOptions } from "./associated-api-key-options";

@Component({
    selector: "app-api-key-table-filter",
    templateUrl: "./api-key-table-filter.component.html",
    styleUrl: "./api-key-table-filter.component.scss",
    standalone: false
})
export class ApiKeyTableFilterComponent extends BaseFormComponent implements OnInit {
  constructor(
    private dialogRef: MatDialogRef<ApiKeyTableFilterComponent>,
    private formBuilder: FormBuilder,
    private store: Store
  ) {
    super();
  }

  protected readonly associatedApiKeyOptions = associatedApiKeyOptions;

  public ngOnInit(): void {
    this.initForm();
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