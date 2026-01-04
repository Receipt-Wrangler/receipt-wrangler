import { Component, EventEmitter, Input, Output } from "@angular/core";
import { ThemePalette } from "@angular/material/core";

@Component({
  selector: "app-base-button",
  standalone: false,
  templateUrl: "./base-button.component.html",
  styleUrl: "./base-button.component.scss"
})
export class BaseButtonComponent {
  @Input() public buttonClass: string = "";

  @Input() public color: string = "primary";

  @Input() public buttonText: string = "";

  @Input() public type: "button" | "menu" | "submit" | "reset" = "button";

  @Input() public matButtonType: "matRaisedButton" | "iconButton" | "basic" =
    "matRaisedButton";

  @Input() public icon: string = "";

  @Input() public customIcon: string = "";

  @Input() public disabled: boolean = false;

  @Input() public buttonRouterLink?: string[];

  @Input() public buttonQueryParams: any = {};

  @Input() public tooltip: string = "";

  @Input() public matBadgeContent?: any;

  @Input() public matBadgeColor: ThemePalette = "primary";

  @Output() public clicked: EventEmitter<MouseEvent> =
    new EventEmitter<MouseEvent>();
}
