import { Component, Input } from "@angular/core";
import { FormControl } from "@angular/forms";
import { AutocompleteModule } from "../autocomplete/autocomplete.module";
import { Tag } from "../open-api/index";
import { PipesModule } from "../pipes/index";

@Component({
    selector: "app-tag-autocomplete",
    standalone: true,
    imports: [
        AutocompleteModule,
        PipesModule
    ],
    templateUrl: "./tag-autocomplete.component.html",
    styleUrl: "./tag-autocomplete.component.scss"
})
export class TagAutocompleteComponent {
  @Input() public tags: Tag[] = [];

  @Input() public inputFormControl!: FormControl;

  @Input() public readonly = false;
}
