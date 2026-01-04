import { Component, EventEmitter, Input, Output, TemplateRef, } from "@angular/core";
import { FormGroup } from "@angular/forms";
import { FormConfig } from "src/interfaces";

@Component({
    selector: "app-form",
    templateUrl: "./form.component.html",
    styleUrls: ["./form.component.scss"],
    standalone: false
})
export class FormComponent {
  @Input() public formConfig!: FormConfig;

  @Input() public form!: FormGroup;

  @Input() public formTemplate!: TemplateRef<any>;

  @Input() public editButtonRouterLink: string[] = [];

  @Input() public editButtonQueryParams: any = {};

  @Input() public canEdit = true;

  @Input() public bottomSpacing = false;

  @Input() public submitButtonDisabled = false;

  @Output() public submitted: EventEmitter<void> = new EventEmitter<void>();
}
