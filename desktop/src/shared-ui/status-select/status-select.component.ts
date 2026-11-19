import { Component, OnChanges, SimpleChanges, input } from "@angular/core";
import { FormControl } from "@angular/forms";
import { RECEIPT_STATUS_OPTIONS } from "src/constants";
import { SelectComponent } from "../../select/select/select.component";

@Component({
    selector: "app-status-select",
    templateUrl: "./status-select.component.html",
    styleUrls: ["./status-select.component.scss"],
    imports: [SelectComponent]
})
export class StatusSelectComponent implements OnChanges {
  public readonly inputFormControl = input.required<FormControl>();

  public readonly readonly = input<boolean>(false);

  public readonly addBlankOption = input<boolean>(false);

  public readonly label = input("Status");

  public receiptStatusOptions = [...RECEIPT_STATUS_OPTIONS];

  public ngOnChanges(changes: SimpleChanges): void {
    if (changes["addBlankOption"]?.currentValue) {
      this.receiptStatusOptions.unshift({
        value: null,
        displayValue: "",
      });
    } else {
      this.receiptStatusOptions = [...RECEIPT_STATUS_OPTIONS];
    }
  }
}
