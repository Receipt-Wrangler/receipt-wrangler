import { Component, Input, ViewEncapsulation } from "@angular/core";

@Component({
  selector: "app-card",
  templateUrl: "./card.component.html",
  styleUrls: ["./card.component.scss"],
  standalone: false,
  encapsulation: ViewEncapsulation.None,
})
export class CardComponent {
  @Input() public cardStyle: string = "";
}
