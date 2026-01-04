import { Component, OnInit, ViewChild } from "@angular/core";
import { FormBuilder, Validators } from "@angular/forms";
import { MatDialogRef } from "@angular/material/dialog";
import { DomSanitizer, SafeHtml } from "@angular/platform-browser";
import { Store } from "@ngxs/store";
import { prettyPrintJson } from "pretty-print-json";
import { switchMap, take, tap } from "rxjs";
import { fadeInOut } from "../../animations";
import { BaseFormComponent } from "../../form";
import { ReceiptFileUploadCommand } from "../../interfaces";
import { FormOption } from "../../interfaces/form-option.interface";
import { FeatureConfigService, ImportService, ImportType } from "../../open-api";
import { DEFAULT_PRETTY_JSON_OPTIONS } from "../../receipt-processing-settings/constants/pretty-json";
import { UploadImageComponent } from "../../receipts/upload-image/upload-image.component";
import { SnackbarService } from "../../services";
import { SetFeatureConfig } from "../../store";

@Component({
    selector: "app-import-form",
    templateUrl: "./import-form.component.html",
    styleUrl: "./import-form.component.scss",
    animations: [fadeInOut],
    standalone: false
})
export class ImportFormComponent extends BaseFormComponent implements OnInit {
  @ViewChild(UploadImageComponent)
  public uploadImageComponent!: UploadImageComponent;

  public readonly acceptedFileTypes = [
    "application/json",
  ];

  public file: File | null = null;

  public fileContents: SafeHtml = "";

  public test: string = "";

  constructor(
    private dialogRef: MatDialogRef<ImportFormComponent>,
    private formBuilder: FormBuilder,
    private importService: ImportService,
    private sanitizer: DomSanitizer,
    private snackbarService: SnackbarService,
    private featureConfigService: FeatureConfigService,
    private store: Store,
  ) {
    super();
  }

  public importTypeOptions: FormOption[] = [
    {
      value: ImportType.ImportConfig,
      displayValue: "Import Config"
    }
  ];

  public ngOnInit(): void {
    this.initForm();
  }

  private initForm(): void {
    this.form = this.formBuilder.group({
      importType: [ImportType.ImportConfig, Validators.required],
    });
  }

  public openFileUploadDialog(): void {
    this.uploadImageComponent.clickInput();
  }

  public fileLoaded(fileData: ReceiptFileUploadCommand): void {
    this.file = fileData.file;

    const reader = new FileReader();
    reader.onload = () => {
      let rawResult = reader.result as string;
      rawResult = JSON.parse(rawResult);
      const dirtyHTML = prettyPrintJson.toHtml(
        rawResult, DEFAULT_PRETTY_JSON_OPTIONS);
      this.fileContents = this.sanitizer.bypassSecurityTrustHtml(dirtyHTML);
    };
    reader.readAsText(this.file);
  }

  public clearFileContents(): void {
    this.file = null;
    this.fileContents = "";
  }

  public closeDialog(): void {
    this.dialogRef.close();
  }

  public submit(): void {
    const importType = this.form.get("importType")?.value;
    switch (importType) {
      case ImportType.ImportConfig:
        this.handleImportConfigSubmit();
        break;
      default:
        this.snackbarService.error("Invalid import type");
    }
  }

  private handleImportConfigSubmit(): void {
    if (!this.file) {
      this.snackbarService.error("Please select a config file to import");
      return;
    }
    this.importService.importConfigJson(this.file)
      .pipe(
        take(1),
        tap(() => {
          this.snackbarService.success("Config imported successfully");
          this.dialogRef.close();
        }),
        switchMap(() => this.featureConfigService.getFeatureConfig()),
        tap((featureConfig) => this.store.dispatch(new SetFeatureConfig(featureConfig)))
      )
      .subscribe();
  }

}
