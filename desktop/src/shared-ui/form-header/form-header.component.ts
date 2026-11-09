import { Location, NgClass, NgTemplateOutlet } from "@angular/common";
import { Component, Input, TemplateRef, input } from "@angular/core";
import { BackButtonComponent } from "../back-button/back-button.component";

@Component({
    selector: "app-form-header",
    templateUrl: "./form-header.component.html",
    styleUrls: ["./form-header.component.scss"],
    imports: [NgClass, BackButtonComponent, NgTemplateOutlet]
})
export class FormHeaderComponent {
  public readonly headerText = input<string>("");

  @Input() public headerButtonsTemplate?: TemplateRef<any>;

  public readonly bottomSpacing = input<boolean>(false);

  public readonly displayBackButton = input<boolean>(true);

  constructor(private location: Location) {}

  public navigateBack(): void {
    this.location.back();
  }
}
