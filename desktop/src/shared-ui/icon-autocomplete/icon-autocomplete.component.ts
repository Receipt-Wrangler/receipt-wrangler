import { Component, Input } from "@angular/core";
import { FormControl } from "@angular/forms";
import { Select, Store } from "@ngxs/store";
import { Observable } from "rxjs";
import { Icon } from "../../open-api/index";
import { AuthState } from "../../store/index";

@Component({
    selector: "app-icon-autocomplete",
    templateUrl: "./icon-autocomplete.component.html",
    styleUrl: "./icon-autocomplete.component.scss",
    standalone: false
})
export class IconAutocompleteComponent {
  @Input() public inputFormControl!: FormControl;

  @Input() public label = "";

  @Select(AuthState.icons) public icons!: Observable<Icon[]>;

  constructor(private store: Store) {}

  public displayWith(value: string): string {
    return this.store.selectSnapshot(AuthState.icons).find((icon) => icon.value === value)?.displayValue ?? "";
  }
}
