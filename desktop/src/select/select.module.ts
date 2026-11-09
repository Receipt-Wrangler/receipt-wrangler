import { CommonModule } from "@angular/common";
import { NgModule } from "@angular/core";
import { ReactiveFormsModule } from "@angular/forms";
import { MatInput } from "@angular/material/input";
import { MatSelectModule } from "@angular/material/select";
import { BadgeComponent } from "src/shared-ui/badge/badge.component";
import { OptionDisplayPipe } from "./pipes/option-display.pipe";
import { ReadonlyValuePipe } from "./pipes/readonly-value.pipe";
import { SelectComponent } from "./select/select.component";

@NgModule({
  imports: [
    CommonModule,
    MatSelectModule,
    ReactiveFormsModule,
    MatInput,
    // Imported directly rather than through SharedUiModule: that module imports
    // SelectModule, so going the other way round would be a circular import.
    BadgeComponent,
    SelectComponent, OptionDisplayPipe, ReadonlyValuePipe,
  ],
  exports: [SelectComponent],
  providers: [OptionDisplayPipe]
})
export class SelectModule {}
