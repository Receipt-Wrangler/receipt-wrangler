import { CommonModule } from "@angular/common";
import { Component, input } from "@angular/core";
import { AppDatePipe } from "../../../pipes/app-date.pipe";

@Component({
    selector: "app-date-block",
    imports: [CommonModule, AppDatePipe],
    templateUrl: "./date-block.component.html",
    styleUrl: "./date-block.component.scss"
})
export class DateBlockComponent {
  public readonly date = input.required<Date | string>();
}
