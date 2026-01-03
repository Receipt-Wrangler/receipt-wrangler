import { Component, OnInit, ViewChild, ViewEncapsulation, } from "@angular/core";
import { FormArray, FormBuilder, FormControl, FormGroup, Validators } from "@angular/forms";
import { MatDialogRef } from "@angular/material/dialog";
import { Store } from "@ngxs/store";
import { take, tap } from "rxjs";
import { ReceiptFileUploadCommand } from "../../interfaces";
import { ReceiptService, ReceiptStatus } from "../../open-api";
import { SnackbarService } from "../../services";
import { AuthState } from "../../store";
import { UploadImageComponent } from "../upload-image/upload-image.component";

@Component({
    selector: "app-quick-scan-dialog",
    templateUrl: "./quick-scan-dialog.component.html",
    styleUrls: ["./quick-scan-dialog.component.scss"],
    encapsulation: ViewEncapsulation.None,
    standalone: false
})
export class QuickScanDialogComponent implements OnInit {
  @ViewChild(UploadImageComponent)
  public uploadImageComponent!: UploadImageComponent;

  public form: FormGroup = new FormGroup({});

  public images: ReceiptFileUploadCommand[] = [];

  public currentlySelectedIndex: number = 0;

  constructor(
    private dialogRef: MatDialogRef<QuickScanDialogComponent>,
    private formBuilder: FormBuilder,
    private receiptService: ReceiptService,
    private snackbarService: SnackbarService,
    private store: Store
  ) {}

  public get paidByUserIds(): FormArray {
    return this.form.get("paidByUserIds") as FormArray;
  }

  public get statuses(): FormArray {
    return this.form.get("statuses") as FormArray;
  }

  public get groupIds(): FormArray {
    return this.form.get("groupIds") as FormArray;
  }

  public ngOnInit(): void {
    this.initForm();
  }

  private initForm(): void {
    this.form = this.formBuilder.group({
      paidByUserIds: this.formBuilder.array<number>([]),
      statuses: this.formBuilder.array<ReceiptStatus>([]),
      groupIds: this.formBuilder.array<number>([]),
    });
  }

  public fileLoaded(fileData: ReceiptFileUploadCommand): void {
    if (fileData.file && !fileData.encodedImage) {
      fileData.url = URL.createObjectURL(fileData.file);
    }
    this.images.push(fileData);
    const userPreferences = this.store.selectSnapshot(AuthState.userPreferences);

    this.paidByUserIds.push(new FormControl(userPreferences?.quickScanDefaultPaidById ?? "", Validators.required));
    this.statuses.push(new FormControl(userPreferences?.quickScanDefaultStatus ?? "", Validators.required));
    this.groupIds.push(new FormControl(userPreferences?.quickScanDefaultGroupId ?? "", Validators.required));
  }


  public openImageUploadComponent(): void {
    this.uploadImageComponent.clickInput();
  }

  public removeImage(index: number): void {
    this.paidByUserIds.removeAt(index);
    this.statuses.removeAt(index);
    this.groupIds.removeAt(index);
    this.images.splice(index, 1);
  }

  public submitButtonClicked(): void {
    if (this.form.valid && this.images.length > 0) {
      this.receiptService
        .quickScanReceipt(
          this.images.map((i) => i.file),
          this.groupIds.value,
          this.paidByUserIds.value,
          this.statuses.value
        )
        .pipe(
          take(1),
          tap(() => {
            const imageWord = this.images.length === 1 ? "image" : "images";
            this.snackbarService.success(`Successfully queued ${imageWord} for processing`);
            this.dialogRef.close();
          }),
        )
        .subscribe();
    }
    if (this.images.length === 0) {
      this.snackbarService.error("Please select images to upload");
    }
    if (this.form.invalid) {
      this.snackbarService.error("Please fill in all required fields. Some images are missing required fields.");
    }
  }

  public cancelButtonClicked(): void {
    this.dialogRef.close();
  }

  public navigateImages(delta: number): void {
    const newValue = this.currentlySelectedIndex + delta;
    if (newValue === -1) {
      this.currentlySelectedIndex = this.images.length - 1;
      return;
    }

    if (newValue > this.images.length - 1) {
      this.currentlySelectedIndex = 0;
      return;
    }

    this.currentlySelectedIndex = newValue;
  }
}
