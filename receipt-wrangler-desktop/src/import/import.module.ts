import { CommonModule } from "@angular/common";
import { NgModule } from "@angular/core";
import { ReactiveFormsModule } from "@angular/forms";
import { ButtonModule } from "../button";
import { PipesModule } from "../pipes";
import { ReceiptsModule } from "../receipts/receipts.module";
import { SelectModule } from "../select/select.module";
import { SharedUiModule } from "../shared-ui/shared-ui.module";
import { ImportFormComponent } from "./import-form/import-form.component";


@NgModule({
  declarations: [
    ImportFormComponent
  ],
  imports: [
    ButtonModule,
    CommonModule,
    PipesModule,
    ReactiveFormsModule,
    ReceiptsModule,
    SelectModule,
    SharedUiModule,
  ],
  exports: [
    ImportFormComponent
  ]
})
export class ImportModule {}
