import { Component, Input, input } from "@angular/core";
import { ReceiptStatus } from "../../open-api";
import { MatChip } from "@angular/material/chips";
import { NgClass } from "@angular/common";
import { StatusPipe } from "../../pipes/status.pipe";

@Component({
    selector: "app-status-chip",
    templateUrl: "./status-chip.component.html",
    styleUrls: ["./status-chip.component.scss"],
    imports: [MatChip, NgClass, StatusPipe]
})
export class StatusChipComponent {
  @Input() public status: string = "";

  @Input() public customStatus: string = "";

  public readonly customStatusColor = input<"red" | "green" | "gray" | "yellow">();

  public receiptStatus = ReceiptStatus;
}
