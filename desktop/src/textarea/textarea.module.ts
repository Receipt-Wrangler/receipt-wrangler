import { CommonModule } from "@angular/common";
import { NgModule } from "@angular/core";
import { ReactiveFormsModule } from "@angular/forms";
import { MatAutocomplete, MatAutocompleteTrigger, MatOption } from "@angular/material/autocomplete";
import { MatFormFieldModule } from "@angular/material/form-field";
import { MatInputModule } from "@angular/material/input";
import { MatMenu, MatMenuItem, MatMenuTrigger } from "@angular/material/menu";
import { TextareaComponent } from "./textarea/textarea.component";

@NgModule({
  declarations: [TextareaComponent],
  imports: [
    CommonModule,
    MatFormFieldModule,
    MatInputModule,
    ReactiveFormsModule,
    MatMenu,
    MatMenuItem,
    MatMenuTrigger,
    MatAutocomplete,
    MatOption,
    MatAutocompleteTrigger,
  ],
  exports: [TextareaComponent],
})
export class TextareaModule {}
