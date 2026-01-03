import { Component, Input, TemplateRef } from "@angular/core";

@Component({
    selector: "app-form-section",
    templateUrl: "./form-section.component.html",
    styleUrls: ["./form-section.component.scss"],
    standalone: false
})
export class FormSectionComponent {
  @Input() public headerText: string = "";

  @Input() public headerButtonsTemplate?: TemplateRef<any>;

  @Input() public indent: boolean = true;

  @Input() public subtitle: string = "";
}
