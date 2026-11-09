import { Component, Input, TemplateRef, input } from "@angular/core";
import { NgTemplateOutlet } from "@angular/common";
import { CdkScrollable } from "@angular/cdk/scrolling";
import { MatDialogContent } from "@angular/material/dialog";

@Component({
    selector: "app-dialog",
    templateUrl: "./dialog.component.html",
    styleUrls: ["./dialog.component.scss"],
    imports: [NgTemplateOutlet, CdkScrollable, MatDialogContent]
})
export class DialogComponent {
  public readonly headerText = input<string>("");

  @Input() public actionsTemplate?: TemplateRef<any>;
}
