import { Component, TemplateRef, input, output } from "@angular/core";
import { FormGroup, ReactiveFormsModule } from "@angular/forms";
import { FormConfig } from "src/interfaces";
import { FormHeaderComponent } from "../form-header/form-header.component";
import { NgTemplateOutlet } from "@angular/common";
import { FormButtonBarComponent } from "../form-button-bar/form-button-bar.component";
import { SubmitButtonComponent } from "../submit-button/submit-button.component";
import { EditButtonComponent } from "../edit-button/edit-button.component";
import { InputReadonlyPipe } from "../../pipes/input-readonly.pipe";

@Component({
    selector: "app-form",
    templateUrl: "./form.component.html",
    styleUrls: ["./form.component.scss"],
    imports: [FormHeaderComponent, ReactiveFormsModule, NgTemplateOutlet, FormButtonBarComponent, SubmitButtonComponent, EditButtonComponent, InputReadonlyPipe]
})
export class FormComponent {
  public readonly formConfig = input.required<FormConfig>();

  public readonly form = input.required<FormGroup>();

  public readonly formTemplate = input.required<TemplateRef<any>>();

  public readonly editButtonRouterLink = input<string[]>([]);

  public readonly editButtonQueryParams = input<any>({});

  public readonly canEdit = input(true);

  public readonly bottomSpacing = input(false);

  public readonly submitButtonDisabled = input(false);

  public readonly submitted = output<void>();
}
