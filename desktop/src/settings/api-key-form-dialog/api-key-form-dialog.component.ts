import { Component, OnInit, signal, input } from "@angular/core";
import { FormBuilder, FormGroup, Validators, ReactiveFormsModule } from "@angular/forms";
import { MatDialogRef } from "@angular/material/dialog";
import { take, tap } from "rxjs";
import { API_KEY_SCOPE_OPTIONS } from "../../constants";
import { ApiKeyResult, ApiKeyScope, ApiKeyService, ApiKeyView, UpsertApiKeyCommand } from "../../open-api";
import { SnackbarService } from "../../services";
import { DialogComponent } from "../../shared-ui/dialog/dialog.component";
import { InputComponent } from "../../input/input/input.component";
import { SelectComponent } from "../../select/select/select.component";
import { DialogFooterComponent } from "../../shared-ui/dialog-footer/dialog-footer.component";
import { MatIcon } from "@angular/material/icon";
import { ButtonComponent } from "../../button/button/button.component";
import { FormGetPipe } from "../../pipes/form-get.pipe";

@Component({
  selector: "app-api-key-form-dialog",
  templateUrl: "./api-key-form-dialog.component.html",
  styleUrls: ["./api-key-form-dialog.component.scss"],
  imports: [DialogComponent, ReactiveFormsModule, InputComponent, SelectComponent, DialogFooterComponent, MatIcon, ButtonComponent, FormGetPipe]
})
export class ApiKeyFormDialogComponent implements OnInit {
  public readonly headerText = input<string>("");
  public readonly apiKey = input<ApiKeyView>();

  public form!: FormGroup;
  public apiKeyResult = signal<ApiKeyResult | undefined>(undefined);
  public apiKeyScopes = API_KEY_SCOPE_OPTIONS;

  constructor(
    private dialogRef: MatDialogRef<ApiKeyFormDialogComponent>,
    private formBuilder: FormBuilder,
    private apiKeyService: ApiKeyService,
    private snackbarService: SnackbarService
  ) {}

  public ngOnInit(): void {
    this.initForm();
  }

  private initForm(): void {
    this.form = this.formBuilder.group({
      name: [this.apiKey()?.name ?? "", Validators.required],
      description: [this.apiKey()?.description ?? ""],
      scope: [this.apiKey()?.scope ?? ApiKeyScope.R, Validators.required]
    });
  }

  public submitButtonClicked(): void {
    if (this.form.valid) {
      const command: UpsertApiKeyCommand = this.form.value;
      const apiKey = this.apiKey();

      if (apiKey && apiKey.id) {
        this.apiKeyService.updateApiKey(apiKey.id, command)
          .pipe(
            take(1),
            tap(() => {
              this.snackbarService.success("API key updated successfully");
              this.dialogRef.close(true);
            })
          )
          .subscribe();
      } else {
        this.apiKeyService.createApiKey(command)
          .pipe(
            take(1),
            tap((result: ApiKeyResult) => {
              this.apiKeyResult.set(result);
              this.snackbarService.success("API key created successfully");
            })
          )
          .subscribe();
      }
    } else {
      this.snackbarService.error("Please fill in all required fields");
    }
  }

  public cancelButtonClicked(): void {
    this.dialogRef.close();
  }

  public copyToClipboard(): void {
    const key = this.apiKeyResult()?.key;
    if (key) {
      navigator.clipboard.writeText(key).then(() => {
        this.snackbarService.success("API key copied to clipboard");
      }).catch(() => {
        this.snackbarService.error("Failed to copy API key to clipboard");
      });
    }
  }

  public closeDialog(): void {
    this.dialogRef.close(this.apiKey() ? true : this.apiKeyResult());
  }
}
