import { Component, Input, OnInit } from "@angular/core";
import { FormArray, FormBuilder, FormGroup, Validators } from "@angular/forms";
import { MatDialogRef } from "@angular/material/dialog";
import { UntilDestroy, untilDestroyed } from "@ngneat/until-destroy";
import { take, tap } from "rxjs";
import { CategoryForm } from "../../categories/category-form/category-form.component";
import { FormOption } from "../../interfaces/form-option.interface";
import { CustomField, CustomFieldOption, CustomFieldService, CustomFieldType } from "../../open-api/index";
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

  public typeOptions: FormOption[] = Object.keys(CustomFieldType).map((key) => {
    return {
      value: (CustomFieldType as any)[key],
      displayValue: key,
    };
  });

  public form!: FormGroup;

  public readonly = false;

  protected readonly CustomFieldType = CustomFieldType;

  constructor(
    private customFieldService: CustomFieldService,
    private formBuilder: FormBuilder,
    private matDialogRef: MatDialogRef<CategoryForm>,
    private snackbarService: SnackbarService,
  ) {}

  public get options(): CustomFieldOption[] {
    return (this.form.get("options") as FormArray).value;
  }

  public ngOnInit(): void {
    this.initForm();
    this.listenForTypeChanges();
    this.readonly = !!this.customField;
  }

  public submit(): void {
    if (this.form.valid) {
      const command = this.form.value;
      this.customFieldService.createCustomField(command)
        .pipe(
          take(1),
          tap(() => {
            this.snackbarService.success("Custom field created");
            this.matDialogRef.close(true);
          })
        ).subscribe();
    }
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
      id: Math.random(),
      value: option?.value,
      customFieldId: option?.customFieldId ?? 0,
    });
  }
}
