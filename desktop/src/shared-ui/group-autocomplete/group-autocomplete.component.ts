import { Component, input } from "@angular/core";
import { FormControl } from "@angular/forms";
import { Store } from "@ngxs/store";
import { Group } from "../../open-api";
import { GroupState } from "../../store";
import { AutocomleteComponent } from "../../autocomplete/autocomlete/autocomlete.component";

@Component({
    selector: "app-group-autocomplete",
    templateUrl: "./group-autocomplete.component.html",
    styleUrls: ["./group-autocomplete.component.scss"],
    imports: [AutocomleteComponent]
})
export class GroupAutocompleteComponent {
  public readonly inputFormControl = input.required<FormControl>();

  public readonly readonly = input<boolean>(false);

  public groups = this.store.selectSignal(GroupState.groupsWithoutAll);

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
