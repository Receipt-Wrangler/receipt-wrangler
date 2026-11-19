import { Component, OnInit } from "@angular/core";
import { FormBuilder, ReactiveFormsModule } from "@angular/forms";
import { MatDialogRef } from "@angular/material/dialog";
import { Store } from "@ngxs/store";
import { BaseFormComponent } from "../../form";
import { GroupTableState } from "../../store/group-table.state";
import { SetFilter } from "../../store/group-table.state.actions";
import { associatedGroupOptions } from "./associated-group-options";
import { DialogComponent } from "../../shared-ui/dialog/dialog.component";
import { SelectComponent } from "../../select/select/select.component";
import { DialogFooterComponent } from "../../shared-ui/dialog-footer/dialog-footer.component";
import { FormGetPipe } from "../../pipes/form-get.pipe";

@Component({
    selector: "app-group-table-filter",
    templateUrl: "./group-table-filter.component.html",
    styleUrl: "./group-table-filter.component.scss",
    imports: [DialogComponent, ReactiveFormsModule, SelectComponent, DialogFooterComponent, FormGetPipe]
})
export class GroupTableFilterComponent extends BaseFormComponent implements OnInit {
  constructor(
    private dialogRef: MatDialogRef<GroupTableFilterComponent>,
    private formBuilder: FormBuilder,
    private store: Store
  ) {
    super();
  }

  protected readonly associatedGroupOptions = associatedGroupOptions;

  public ngOnInit(): void {
    this.initForm();
  }

  private initForm(): void {
    const filter = this.store.selectSnapshot(GroupTableState.filter);
    this.form = this.formBuilder.group({
        associatedGroup: [filter.associatedGroup],
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
