import { Component, Input, input } from "@angular/core";
import { FormMode } from "../../enums/form-mode.enum";
import { FormSectionComponent } from "../form-section/form-section.component";
import { AppDatePipe } from "../../pipes/app-date.pipe";
import { UserPipe } from "../../pipes/user.pipe";

@Component({
    selector: "app-audit-detail-section",
    templateUrl: "./audit-detail-section.component.html",
    styleUrl: "./audit-detail-section.component.scss",
    imports: [FormSectionComponent, AppDatePipe, UserPipe]
})
export class AuditDetailSectionComponent {
  @Input() data!: any;

  readonly formMode = input<FormMode>(FormMode.view);

  public readonly indent = input<boolean>(true);

  protected readonly FormMode = FormMode;
}
