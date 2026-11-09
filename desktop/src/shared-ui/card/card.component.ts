import { Component, ViewEncapsulation, input } from "@angular/core";
import { MatCard, MatCardHeader, MatCardTitle, MatCardContent } from "@angular/material/card";
import { NgClass } from "@angular/common";

@Component({
  selector: "app-card",
  templateUrl: "./card.component.html",
  styleUrls: ["./card.component.scss"],
  encapsulation: ViewEncapsulation.None,
  imports: [
    MatCard,
    NgClass,
    MatCardHeader,
    MatCardTitle,
    MatCardContent,
  ],
})
export class CardComponent {
  public readonly cardStyle = input<string>("");
}
