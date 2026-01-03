import { CommonModule } from "@angular/common";
import { NgModule } from "@angular/core";
import { ReactiveFormsModule } from "@angular/forms";
import { MatAutocompleteModule } from "@angular/material/autocomplete";
import { MatChipsModule } from "@angular/material/chips";
import { MatIconModule } from "@angular/material/icon";
import { MatInputModule } from "@angular/material/input";
import { ButtonModule } from "../button";
import { AutocomleteComponent } from "./autocomlete/autocomlete.component";
import { OptionDisplayPipe } from "./autocomlete/option-display.pipe";

@NgModule({
  declarations: [AutocomleteComponent, OptionDisplayPipe],
  imports: [
    CommonModule,
    MatAutocompleteModule,
    ReactiveFormsModule,
    MatInputModule,
    MatChipsModule,
    MatIconModule,
    ButtonModule,
  ],
  exports: [AutocomleteComponent],
})
export class AutocompleteModule {}
