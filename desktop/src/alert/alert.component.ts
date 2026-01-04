import { CommonModule } from "@angular/common";
import { Component, Input } from "@angular/core";
import { MatIconModule } from "@angular/material/icon";

@Component({
    selector: "app-alert",
    standalone: true,
    imports: [
        CommonModule,
        MatIconModule
    ],
    templateUrl: "./alert.component.html",
    styleUrl: "./alert.component.scss"
})
export class AlertComponent {
  @Input() public type: "warning" = "warning";
  @Input() public message: string = "";
}
