import { Component, Input } from "@angular/core";
import { FormControl } from "@angular/forms";
import { AutocompleteModule } from "../autocomplete/autocomplete.module";
import { Category } from "../open-api/index";
import { PipesModule } from "../pipes/index";

@Component({
  selector: "app-category-autocomplete",
  standalone: true,
  imports: [
    AutocompleteModule,
    PipesModule
  ],
  templateUrl: "./category-autocomplete.component.html",
  styleUrl: "./category-autocomplete.component.scss"
})
export class CategoryAutocompleteComponent {
  @Input() public categories: Category[] = [];

  @Input() public inputFormControl!: FormControl;

  @Input() public readonly = false;
}
