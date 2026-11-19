import { Component, ViewEncapsulation, } from "@angular/core";
import { BaseButtonComponent } from "../base-button/base-button.component";
import { MatButton, MatIconButton } from "@angular/material/button";
import { MatTooltip } from "@angular/material/tooltip";
import { RouterLink } from "@angular/router";
import { MatBadge } from "@angular/material/badge";
import { MatIcon } from "@angular/material/icon";

@Component({
  selector: "app-button",
  templateUrl: "./button.component.html",
  styleUrls: ["./button.component.scss"],
  encapsulation: ViewEncapsulation.None,
  imports: [MatButton, MatTooltip, RouterLink, MatBadge, MatIcon, MatIconButton]
})
export class ButtonComponent extends BaseButtonComponent {

  public emitClicked(event: MouseEvent): void {
    this.clicked.emit(event);
  }
}
