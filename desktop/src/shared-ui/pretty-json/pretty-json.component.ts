import { Component, input } from "@angular/core";
import { NgClass } from "@angular/common";
import { PrettyJsonPipe } from "../task-table/pretty-json.pipe";

@Component({
    selector: "app-pretty-json",
    templateUrl: "./pretty-json.component.html",
    styleUrl: "./pretty-json.component.scss",
    imports: [NgClass, PrettyJsonPipe]
})
export class PrettyJsonComponent {
  public readonly json = input<string | undefined>("");

  public readonly verticalJson = input(true);
}
