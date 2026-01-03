import { Component, Input } from "@angular/core";
import { FormControl } from "@angular/forms";
import { Select, Store } from "@ngxs/store";
import { Observable } from "rxjs";
import { Group } from "../../open-api";
import { GroupState } from "../../store";

@Component({
    selector: "app-group-autocomplete",
    templateUrl: "./group-autocomplete.component.html",
    styleUrls: ["./group-autocomplete.component.scss"],
    standalone: false
})
export class GroupAutocompleteComponent {
  @Input() public inputFormControl!: FormControl;

  @Input() public readonly: boolean = false;

  @Select(GroupState.groupsWithoutAll) public groups!: Observable<Group[]>;

  constructor(private store: Store) {}

  public groupDisplayWith(id: number): string {
    const group = this.store.selectSnapshot(
      GroupState.getGroupById(id?.toString())
    );

    if (group) {
      return group.name;
    }
    return "";
  }
}
