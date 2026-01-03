import { Component, EventEmitter, Input, Output, TemplateRef, } from "@angular/core";
import { Select } from "@ngxs/store";
import { Observable } from "rxjs";
import { LayoutState } from "src/store/layout.state";

@Component({
    selector: "app-dialog-footer",
    templateUrl: "./dialog-footer.component.html",
    styleUrls: ["./dialog-footer.component.scss"],
    standalone: false
})
export class DialogFooterComponent {
  @Input() public additionalButtonsTemplate?: TemplateRef<any>;
  @Input() public submitButtonTooltip: string = "Save";
  @Input() public submitButtonType: "button" | "submit" = "submit";
  @Input() public disableWhenProgressBarIsShown: boolean = false;
  @Output() public cancelClicked: EventEmitter<void> = new EventEmitter<void>();
  @Output() public submitClicked: EventEmitter<void> = new EventEmitter<void>();

  @Select(LayoutState.showProgressBar)
  public showProgressBar!: Observable<boolean>;
}
