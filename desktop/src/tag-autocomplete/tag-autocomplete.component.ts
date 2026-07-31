import { Component, input } from "@angular/core";
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
  public readonly tags = input<Tag[]>([]);

  public readonly inputFormControl = input.required<FormControl>();

  public readonly readonly = input(false);

  // Whether the user may create a brand-new tag inline. Defaults to true to
  // preserve existing call sites; gated callers bind it to app.tags.create.
  public readonly creatable = input(true);

  // See CategoryAutocompleteComponent.inputId — overrides the underlying input's
  // DOM id so multiple instances on one page don't collide.
  public readonly inputId = input("");
}
