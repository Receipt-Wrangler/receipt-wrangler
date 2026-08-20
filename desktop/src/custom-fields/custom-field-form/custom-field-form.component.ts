import { Component, Input, OnInit } from "@angular/core";
import { FormArray, FormBuilder, FormGroup, Validators } from "@angular/forms";
import { MatDialogRef } from "@angular/material/dialog";
import { UntilDestroy, untilDestroyed } from "@ngneat/until-destroy";
import { Observable, take, tap } from "rxjs";
import { FormMode } from "../../enums/form-mode.enum";
import { FormOption } from "../../interfaces/form-option.interface";
import {
  CustomField,
  CustomFieldOption,
  CustomFieldService,
  CustomFieldType,
  UpsertCustomFieldCommand,
} from "../../open-api/index";
import { SnackbarService } from "../../services/index";

@UntilDestroy()
@Component({
  selector: "app-custom-field-form",
  standalone: false,
  templateUrl: "./custom-field-form.component.html",
  styleUrl: "./custom-field-form.component.scss"
})
export class CustomFieldFormComponent implements OnInit {
  @Input() public headerText: string = "";

  @Input() public customField?: CustomField;

  @Input() public mode: FormMode = FormMode.add;

  public typeOptions: FormOption[] = Object.keys(CustomFieldType).map((key) => {
    return {
      value: (CustomFieldType as any)[key],
      displayValue: key,
    };
  });

  public form!: FormGroup;

  protected readonly CustomFieldType = CustomFieldType;

  protected readonly FormMode = FormMode;

  constructor(
    private customFieldService: CustomFieldService,
    private formBuilder: FormBuilder,
    private matDialogRef: MatDialogRef<CustomFieldFormComponent>,
    private snackbarService: SnackbarService,
  ) {}

  public get options(): CustomFieldOption[] {
    return (this.form.get("options") as FormArray).value;
  }

  // The type is immutable once a custom field exists: a CustomFieldValue is
  // stored in a type-specific column, so re-typing would mis-column every value
  // already recorded against the field. The server rejects it too.
  public get typeReadonly(): boolean {
    return this.mode !== FormMode.add;
  }

  public ngOnInit(): void {
    this.initForm();
    this.listenForTypeChanges();
  }

  public submit(): void {
    if (!this.form.valid) {
      return;
    }

    const command = this.buildCommand();

    const request$: Observable<CustomField> =
      this.mode === FormMode.edit && this.customField
        ? this.customFieldService.updateCustomField(this.customField.id, command)
        : this.customFieldService.createCustomField(command);

    const successMessage =
      this.mode === FormMode.edit ? "Custom field updated" : "Custom field created";

    request$
      .pipe(
        take(1),
        tap(() => {
          this.snackbarService.success(successMessage);
          this.matDialogRef.close(true);
        })
      ).subscribe();
  }

  public closeDialog(): void {
    this.matDialogRef.close(false);
  }

  public addOption(): void {
    (this.form.get("options") as FormArray).push(this.buildOption());
  }

  public deleteOption(index: number): void {
    (this.form.get("options") as FormArray).removeAt(index);
  }

  // A saved option cannot be removed, only renamed: CustomFieldValue.SelectValue
  // holds an option id, so deleting one would orphan every receipt that picked
  // it. An option added in this session has no id yet and is safe to drop.
  public canDeleteOption(index: number): boolean {
    if (this.mode === FormMode.view || this.options.length <= 1) {
      return false;
    }

    return !this.options[index]?.id;
  }

  private buildCommand(): UpsertCustomFieldCommand {
    const { name, type, description, options } = this.form.value;

    return {
      name,
      type,
      description,
      options: (options ?? []).map((option: { id?: number; value: string }) => ({
        // A new option carries no id; the server reads that as "append".
        id: option.id ?? undefined,
        value: option.value,
        customFieldId: this.customField?.id ?? 0,
      })),
    };
  }

  private initForm(): void {
    this.form = this.formBuilder.group({
      name: [this.customField?.name, [Validators.required]],
      type: [this.customField?.type, [Validators.required]],
      description: [this.customField?.description],
      options: this.formBuilder.array(this.customField?.options?.map(option => this.buildOption(option)) ?? []),
    });
  }

  private listenForTypeChanges(): void {
    this.form.get("type")?.valueChanges.pipe(
      untilDestroyed(this),
      tap((type) => {
        if (type === CustomFieldType.Select) {
          (this.form.get("options") as FormArray).push(this.buildOption());
          this.form.get("options")?.addValidators(Validators.required);
        } else {
          (this.form.get("options") as FormArray).clear();
          this.form.get("options")?.removeValidators(Validators.required);
          this.form.setErrors(null);
        }
      })
    )
      .subscribe();
  }

  private buildOption(option?: CustomFieldOption): FormGroup {
    return this.formBuilder.group({
      // The real server id, or null for an option that does not exist yet. The
      // template tracks by index, so this never needs to be a synthetic key.
      id: option?.id ?? null,
      value: option?.value,
    });
  }
}
