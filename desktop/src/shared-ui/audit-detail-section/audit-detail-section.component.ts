import { Component, Input } from "@angular/core";
import { FormMode } from "../../enums/form-mode.enum";

@Component({
    selector: "app-audit-detail-section",
    templateUrl: "./audit-detail-section.component.html",
    styleUrl: "./audit-detail-section.component.scss",
    standalone: false
})
export class AuditDetailSectionComponent {
  @Input() data!: any;

  @Input() formMode!: FormMode;

  @Input() public indent: boolean = true;

  protected readonly FormMode = FormMode;
}
