import { Component, Input, OnInit } from "@angular/core";
import { FormBuilder, FormGroup, Validators } from "@angular/forms";
import { MatDialogRef } from "@angular/material/dialog";
import { take, tap } from "rxjs";
import { DuplicateValidator } from "src/validators/duplicate-validator";
import { Category, CategoryService, CategoryView } from "../../open-api";
import { SnackbarService } from "../../services";

@Component({
    selector: "app-category-form",
    templateUrl: "./category-form.component.html",
    styleUrls: ["./category-form.component.scss"],
    standalone: false
})
export class CategoryForm implements OnInit {
  @Input() public headerText: string = "";

  @Input() public category?: CategoryView;

  public form: FormGroup = new FormGroup({});

  constructor(
    private formBuilder: FormBuilder,
    private matDialogRef: MatDialogRef<CategoryForm>,
    private categoryService: CategoryService,
    private snackService: SnackbarService,
    private duplicateValidator: DuplicateValidator
  ) {}

  public ngOnInit(): void {
    this.initForm();
  }

  private initForm(): void {
    const name = this.category?.name ?? "";

    const nameValidator = this.duplicateValidator.isUnique("category", 0, name);
    this.form = this.formBuilder.group({
      name: [name, Validators.required, nameValidator],
      description: [this.category?.description ?? ""],
    });
  }

  public submit(): void {
    if (this.form.valid && this.category) {
      const category: Category = {
        id: this.category?.id,
        name: this.form.value.name,
        description: this.form.value.description,
      };
      this.categoryService
        .updateCategory(category.id as number, category)
        .pipe(
          take(1),
          tap(() => {
            this.snackService.success("Category updated successfully");
            this.matDialogRef.close(true);
          })
        )
        .subscribe();
    } else if (this.form.valid && !this.category) {
      this.categoryService
        .createCategory(this.form.value as Category)
        .pipe(
          take(1),
          tap(() => {
            this.snackService.success("Category created successfully");
            this.matDialogRef.close(true);
          })
        )
        .subscribe();
    }
  }

  public closeDialog(): void {
    this.matDialogRef.close(false);
  }
}
