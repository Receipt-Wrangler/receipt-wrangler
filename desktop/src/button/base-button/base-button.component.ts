import { Component, Input, input, output } from "@angular/core";
import { ThemePalette } from "@angular/material/core";

@Component({
  selector: "app-base-button",
  templateUrl: "./base-button.component.html",
  styleUrl: "./base-button.component.scss"
})
export class BaseButtonComponent {
  public readonly buttonClass = input<string>("");

  public readonly color = input<string>("primary");

  public readonly buttonText = input<string>("");

  public readonly type = input<"button" | "menu" | "submit" | "reset">("button");

  public readonly matButtonType = input<"matRaisedButton" | "iconButton" | "basic">("matRaisedButton");

  @Input() public icon: string = "";

  @Input() public customIcon: string = "";

  public readonly disabled = input<boolean>(false);

  public readonly buttonRouterLink = input<string[]>();

  public readonly buttonQueryParams = input<any>({});

  public readonly tooltip = input<string>("");

  /**
   * The accessible name of an icon-only button. Its mat-icon is aria-hidden and a
   * tooltip is only a description, so without this a screen reader has no name.
   */
  public readonly ariaLabel = input<string>();

  public readonly matBadgeContent = input<any>();

  public readonly matBadgeColor = input<ThemePalette>("primary");

  public readonly clicked = output<MouseEvent>();
}
