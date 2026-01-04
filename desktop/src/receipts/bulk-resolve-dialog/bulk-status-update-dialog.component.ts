import { Component, OnInit } from "@angular/core";
import { FormBuilder, FormGroup, Validators } from "@angular/forms";
import { MatDialogRef } from "@angular/material/dialog";
import { RECEIPT_STATUS_OPTIONS } from "src/constants/receipt-status-options";
import { ReceiptStatus } from "../../open-api";

@Component({
    selector: "app-bulk-status-update-dialog",
    templateUrl: "./bulk-status-update-dialog.component.html",
    styleUrls: ["./bulk-status-update-dialog.component.scss"],
    standalone: false
})
export class BulkStatusUpdateComponent implements OnInit {
  public form: FormGroup = new FormGroup({});

  public receiptStatusOptions = RECEIPT_STATUS_OPTIONS;

  constructor(
    private formBuilder: FormBuilder,
    public matDialogRef: MatDialogRef<BulkStatusUpdateComponent>
  ) {}

  public ngOnInit(): void {
    this.initForm();
  }

  private initForm(): void {
    this.form = this.formBuilder.group({
      status: [ReceiptStatus.Resolved, Validators.required],
      comment: "",
    });
  }

  public cancelButtonClicked(): void {
    this.matDialogRef.close(undefined);
  }

  public submitButtonClicked(): void {
    if (this.form.valid) {
      this.matDialogRef.close(this.form.value);
    }
  }
}
