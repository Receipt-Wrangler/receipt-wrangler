import { CommonModule } from "@angular/common";
import { NgModule } from "@angular/core";
import { ReactiveFormsModule } from "@angular/forms";
import { MatButtonModule } from "@angular/material/button";
import { MatFormFieldModule } from "@angular/material/form-field";
import { MatIconModule } from "@angular/material/icon";
import { MatInputModule } from "@angular/material/input";
import { MatTooltipModule } from "@angular/material/tooltip";
import { NgxMaskDirective, provideNgxMask } from "ngx-mask";
import { InputComponent } from "./input/input.component";

@NgModule({
  declarations: [InputComponent],
  imports: [
    CommonModule,
    MatButtonModule,
    MatFormFieldModule,
    MatIconModule,
    MatInputModule,
    MatTooltipModule,
    NgxMaskDirective,
    ReactiveFormsModule,
  ],
  exports: [InputComponent],
  providers: [provideNgxMask()],
})
export class InputModule {}
