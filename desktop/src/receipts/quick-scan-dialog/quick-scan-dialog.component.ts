import { Component, DestroyRef, OnInit, ViewEncapsulation, inject, viewChild } from "@angular/core";
import { takeUntilDestroyed } from "@angular/core/rxjs-interop";
import { FormArray, FormBuilder, FormControl, FormGroup, Validators } from "@angular/forms";
import { MatDialogRef } from "@angular/material/dialog";
import { Store } from "@ngxs/store";
import { take, tap } from "rxjs";
import { ReceiptFileUploadCommand } from "../../interfaces";
import { Category, GroupReceiptSettings, ReceiptService, ReceiptStatus, Tag } from "../../open-api";
import { SnackbarService } from "../../services";
import { AuthState, GroupState } from "../../store";
import { UploadImageComponent } from "../upload-image/upload-image.component";

@Component({
    selector: "app-quick-scan-dialog",
    templateUrl: "./quick-scan-dialog.component.html",
    styleUrls: ["./quick-scan-dialog.component.scss"],
    encapsulation: ViewEncapsulation.None,
    standalone: false
})
export class QuickScanDialogComponent implements OnInit {
  public readonly uploadImageComponent = viewChild.required(UploadImageComponent);

  public form: FormGroup = new FormGroup({});

  public images: ReceiptFileUploadCommand[] = [];

  public currentlySelectedIndex: number = 0;

  private readonly destroyRef = inject(DestroyRef);

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

  public get categories(): FormArray {
    return this.form.get("categories") as FormArray;
  }

  public get tags(): FormArray {
    return this.form.get("tags") as FormArray;
  }

  public ngOnInit(): void {
    this.initForm();
  }

  private initForm(): void {
    this.form = this.formBuilder.group({
      paidByUserIds: this.formBuilder.array<number>([]),
      statuses: this.formBuilder.array<ReceiptStatus>([]),
      groupIds: this.formBuilder.array<number>([]),
      categories: this.formBuilder.array<Category[]>([]),
      tags: this.formBuilder.array<Tag[]>([]),
    });

    // Re-resolve each image's field config whenever anything changes (a group selection can flip
    // which fields are shown/required for that image).
    this.form.valueChanges
      .pipe(takeUntilDestroyed(this.destroyRef))
      .subscribe(() => this.configureImages());
  }

  public fileLoaded(fileData: ReceiptFileUploadCommand): void {
    if (fileData.file && !fileData.encodedImage) {
      fileData.url = URL.createObjectURL(fileData.file);
    }
    this.images.push(fileData);
    const userPreferences = this.store.selectSnapshot(AuthState.userPreferences);

    // Group is always required; paid-by/status validators are applied per-image by configureImages.
    this.paidByUserIds.push(new FormControl(userPreferences?.quickScanDefaultPaidById ?? ""));
    this.statuses.push(new FormControl(userPreferences?.quickScanDefaultStatus ?? ""));
    this.groupIds.push(new FormControl(userPreferences?.quickScanDefaultGroupId ?? "", Validators.required));
    this.categories.push(new FormControl([]));
    this.tags.push(new FormControl([]));

    this.configureImages();
  }

  // Returns the receipt settings for the group currently selected on the given image, if any.
  public settingsForIndex(index: number): GroupReceiptSettings | undefined {
    const groupId = this.groupIds.at(index)?.value;
    if (!groupId) {
      return undefined;
    }

    return this.store.selectSnapshot(GroupState.getGroupById(groupId.toString()))?.groupReceiptSettings;
  }

  public showPaidBy(index: number): boolean {
    return this.settingsForIndex(index)?.quickScanPaidByEnabled ?? true;
  }

  public showStatus(index: number): boolean {
    return this.settingsForIndex(index)?.quickScanStatusEnabled ?? true;
  }

  public showCategories(index: number): boolean {
    return this.settingsForIndex(index)?.quickScanCategoriesEnabled ?? false;
  }

  public showTags(index: number): boolean {
    return this.settingsForIndex(index)?.quickScanTagsEnabled ?? false;
  }

  public categoriesForIndex(index: number): Category[] {
    const groupId = this.groupIds.at(index)?.value;
    return groupId ? this.store.selectSnapshot(AuthState.groupCategories(Number(groupId))) : [];
  }

  public tagsForIndex(index: number): Tag[] {
    const groupId = this.groupIds.at(index)?.value;
    return groupId ? this.store.selectSnapshot(AuthState.groupTags(Number(groupId))) : [];
  }

  // Applies the per-image required validators and clears hidden fields so the server backfills their
  // group-configured defaults (a hidden paid-by/status must be sent empty, not with a stale value).
  private configureImages(): void {
    for (let i = 0; i < this.images.length; i++) {
      const settings = this.settingsForIndex(i);
      const paidByShown = settings?.quickScanPaidByEnabled ?? true;
      const statusShown = settings?.quickScanStatusEnabled ?? true;
      const categoriesShown = settings?.quickScanCategoriesEnabled ?? false;
      const tagsShown = settings?.quickScanTagsEnabled ?? false;

      this.setRequired("paidByUserIds." + i, paidByShown && (settings?.quickScanPaidByRequired ?? true));
      this.setRequired("statuses." + i, statusShown && (settings?.quickScanStatusRequired ?? true));
      this.setRequired("categories." + i, categoriesShown && (settings?.quickScanCategoriesRequired ?? false));
      this.setRequired("tags." + i, tagsShown && (settings?.quickScanTagsRequired ?? false));

      if (!paidByShown) {
        this.form.get("paidByUserIds." + i)?.setValue("", { emitEvent: false });
      }
      if (!statusShown) {
        this.form.get("statuses." + i)?.setValue("", { emitEvent: false });
      }
      if (!categoriesShown) {
        this.form.get("categories." + i)?.setValue([], { emitEvent: false });
      }
      if (!tagsShown) {
        this.form.get("tags." + i)?.setValue([], { emitEvent: false });
      }
    }
  }

  private setRequired(path: string, required: boolean): void {
    const control = this.form.get(path);
    if (!control) {
      return;
    }

    if (required) {
      control.addValidators(Validators.required);
    } else {
      control.removeValidators(Validators.required);
    }
    control.updateValueAndValidity({ emitEvent: false });
  }

  public openImageUploadComponent(): void {
    this.uploadImageComponent().clickInput();
  }

  public removeImage(index: number): void {
    this.paidByUserIds.removeAt(index);
    this.statuses.removeAt(index);
    this.groupIds.removeAt(index);
    this.categories.removeAt(index);
    this.tags.removeAt(index);
    this.images.splice(index, 1);
  }

  public submitButtonClicked(): void {
    if (this.form.valid && this.images.length > 0) {
      this.receiptService
        .quickScanReceipt(
          this.images.map((i) => i.file),
          this.groupIds.value,
          this.paidByUserIds.value,
          this.statuses.value,
          this.joinIds(this.categories),
          this.joinIds(this.tags)
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

  // Serializes each image's selected category/tag objects into a comma-joined id string (one entry
  // per image), matching the multipart shape the API expects.
  private joinIds(array: FormArray): string[] {
    return array.controls.map((control) =>
      ((control.value ?? []) as { id: number }[]).map((entity) => entity.id).join(",")
    );
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
