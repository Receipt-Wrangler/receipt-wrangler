import { CommonModule } from "@angular/common";
import { Component, Input } from "@angular/core";
import { PipesModule } from "../../../pipes/index";

@Component({
    selector: "app-date-block",
    imports: [CommonModule, PipesModule],
    templateUrl: "./date-block.component.html",
    styleUrl: "./date-block.component.scss"
})
export class DateBlockComponent {
  @Input() public date!: Date | string;
}
