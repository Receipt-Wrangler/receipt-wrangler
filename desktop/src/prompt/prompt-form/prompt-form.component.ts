import { Component, OnInit } from "@angular/core";
import { AbstractControl, FormBuilder, ValidationErrors, ValidatorFn, Validators } from "@angular/forms";
import { ActivatedRoute, Router } from "@angular/router";
import { take, tap } from "rxjs";
import { BaseFormComponent } from "../../form";
import { Prompt, PromptService } from "../../open-api";
import { SnackbarService } from "../../services";

@Component({
    selector: "app-prompt-form",
    templateUrl: "./prompt-form.component.html",
    styleUrl: "./prompt-form.component.scss",
    standalone: false
})
export class PromptFormComponent extends BaseFormComponent implements OnInit {
  public originalPrompt?: Prompt;

  public promptVariables = ["categories", "tags", "currentYear", "ocrText"];

  constructor(
    private activatedRoute: ActivatedRoute,
    private formBuilder: FormBuilder,
    private promptService: PromptService,
    private router: Router,
    private snackbarService: SnackbarService,
  ) {
    super();
  }


  public ngOnInit() {
    this.originalPrompt = this.activatedRoute.snapshot.data["prompt"];
    this.setFormConfigFromRoute(this.activatedRoute);
    this.initForm();
  }

  private initForm(): void {
    this.form = this.formBuilder.group({
      name: [this.originalPrompt?.name, Validators.required],
      description: [this.originalPrompt?.description],
      prompt: [this.originalPrompt?.prompt, [Validators.required, this.templateVariableValidator()]],
    });
  }

  private templateVariableValidator(): ValidatorFn {
    return (control: AbstractControl): ValidationErrors | null => {
      const regex = /@\w+/g;
      const matches = (control.value ?? "").match(regex);
      if (!matches) {
        return null;
      }

      const allValidVariables = matches?.every((match: string) => this.promptVariables.includes(match.slice(1)));

      if (allValidVariables) {
        return null;
      } else {
        return { invalidVariables: { value: true } };
      }
    };
  }

  public submit(): void {
    if (this.form.valid && this.originalPrompt) {
      this.updatePrompt();
    } else if (this.form.valid && !this.originalPrompt) {
      this.createPrompt();
    }
  }

  private createPrompt(): void {
    this.promptService
      .createPrompt(this.form.value)
      .pipe(
        take(1),
        tap((prompt) => {
          this.snackbarService.success("Prompt created successfully");
          this.router.navigate([`/system-settings/prompts/${prompt.id}/view`]);
        })
      )
      .subscribe();
  }

  private updatePrompt(): void {
    this.promptService
      .updatePromptById(this.originalPrompt?.id ?? 0, this.form.value)
      .pipe(
        take(1),
        tap(() => {
          this.snackbarService.success("Prompt updated successfully");
          this.router.navigate(["/system-settings/prompts", this.originalPrompt?.id, "view"]);
        })
      )
      .subscribe();
  }
}
