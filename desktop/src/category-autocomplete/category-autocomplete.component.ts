import { Component, input } from "@angular/core";
import { FormControl } from "@angular/forms";
import { Category } from "../open-api/index";
import { AutocomleteComponent } from "../autocomplete/autocomlete/autocomlete.component";

@Component({
  selector: "app-category-autocomplete",
  standalone: true,
  imports: [
    AutocomleteComponent
  ],
  templateUrl: "./category-autocomplete.component.html",
  styleUrl: "./category-autocomplete.component.scss"
})
export class CategoryAutocompleteComponent {
  public readonly categories = input<Category[]>([]);

  public readonly inputFormControl = input.required<FormControl>();

  public readonly readonly = input(false);

  // Whether the user may create a brand-new category inline. Defaults to true to
  // preserve existing call sites; gated callers (e.g. the receipt form) bind it
  // to the app.categories.create permission.
  public readonly creatable = input(true);

  // Overrides the underlying input's DOM id. The base autocomplete otherwise
  // derives it from the label, so rendering this component more than once on a
  // page produces duplicate ids — which breaks <label for> association (only the
  // first field gets an accessible name) and makes the base component's
  // getElementById-based filter clear target the wrong input.
  public readonly inputId = input("");
}
