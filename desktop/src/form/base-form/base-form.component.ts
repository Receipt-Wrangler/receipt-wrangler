import { Component, EventEmitter, Output } from "@angular/core";
import { FormGroup } from "@angular/forms";
import { ActivatedRoute } from "@angular/router";
import { FormMode } from "src/enums/form-mode.enum";
import { FormConfig } from "src/interfaces";
import { applyFormCommand } from "../../utils/index";
import { FormCommand } from "../interfaces/form-command";

@Component({
    selector: "app-base-form",
    templateUrl: "./base-form.component.html",
    styleUrls: ["./base-form.component.scss"],
    standalone: false
})
export class BaseFormComponent {
  protected readonly FormMode = FormMode;

  public formConfig!: FormConfig;

  public form: FormGroup = new FormGroup({});

  @Output() public formCommand: EventEmitter<FormCommand> =
    new EventEmitter<FormCommand>();

  constructor() {}

  public setFormConfigFromRoute(activatedRoute: ActivatedRoute): void {
    this.formConfig = activatedRoute?.snapshot?.data?.["formConfig"];
  }

  public handleFormCommand(formCommand: FormCommand): void {
    applyFormCommand(this.form, formCommand);
  }

  public emitFormCommand(formCommand: FormCommand): void {
    this.formCommand.emit(formCommand);
  }
}
