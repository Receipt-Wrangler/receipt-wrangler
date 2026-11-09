import { Component, input } from "@angular/core";
import { SystemTaskStatus } from "../../open-api";
import { MatIcon } from "@angular/material/icon";

@Component({
    selector: "app-status-icon",
    templateUrl: "./status-icon.component.html",
    styleUrl: "./status-icon.component.scss",
    imports: [MatIcon]
})
export class StatusIconComponent {
  public readonly taskStatus = input.required<SystemTaskStatus>();
  protected readonly SystemTaskStatus = SystemTaskStatus;
}
