import { Location } from "@angular/common";
import { Component, Input, TemplateRef } from "@angular/core";

@Component({
    selector: "app-form-header",
    templateUrl: "./form-header.component.html",
    styleUrls: ["./form-header.component.scss"],
    standalone: false
})
export class FormHeaderComponent {
  @Input() public headerText: string = "";

  @Input() public headerButtonsTemplate?: TemplateRef<any>;

  @Input() public bottomSpacing: boolean = false;

  @Input() public displayBackButton: boolean = true;

  constructor(private location: Location) {}

  public navigateBack(): void {
    this.location.back();
  }
}
