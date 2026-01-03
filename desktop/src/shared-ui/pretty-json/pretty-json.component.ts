import { Component, Input } from "@angular/core";

@Component({
    selector: "app-pretty-json",
    templateUrl: "./pretty-json.component.html",
    styleUrl: "./pretty-json.component.scss",
    standalone: false
})
export class PrettyJsonComponent {
  @Input() public json?: string = "";

  @Input() public verticalJson = true;
}
