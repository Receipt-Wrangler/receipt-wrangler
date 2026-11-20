import { Component, input, output } from "@angular/core";
import { FormMode } from "src/enums/form-mode.enum";

@Component({
  selector: "app-form-button",
  templateUrl: "./form-button.component.html",
  styleUrls: ["./form-button.component.scss"]
})
export class FormButtonComponent {
  public readonly mode = input<FormMode>(FormMode.view);

  public readonly tooltip = input<string>();

  public readonly disabled = input<boolean>(false);

  public readonly color = input<string>("primary");

  public readonly buttonRouterLink = input<string[]>();

  public readonly buttonQueryParams = input<any>({});

  public readonly buttonText = input<string>();

  public readonly type = input<"button" | "submit">("button");

  public readonly clicked = output<MouseEvent | void>();
}
