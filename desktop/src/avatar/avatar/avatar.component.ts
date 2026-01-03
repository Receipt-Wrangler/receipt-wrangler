import { Component, Input } from "@angular/core";
import { Group, User } from "../../open-api";

@Component({
    selector: "app-avatar",
    templateUrl: "./avatar.component.html",
    styleUrls: ["./avatar.component.scss"],
    standalone: false
})
export class AvatarComponent {
  @Input() public user?: User;
  @Input() public group?: Group;
}
