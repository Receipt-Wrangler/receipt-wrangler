import { Component, ViewEncapsulation, } from "@angular/core";
import { BaseButtonComponent } from "../base-button/base-button.component";

@Component({
  selector: "app-button",
  templateUrl: "./button.component.html",
  styleUrls: ["./button.component.scss"],
  encapsulation: ViewEncapsulation.None,
  standalone: false
})
export class ButtonComponent extends BaseButtonComponent {

  public emitClicked(event: MouseEvent): void {
    this.clicked.emit(event);
  }
}
