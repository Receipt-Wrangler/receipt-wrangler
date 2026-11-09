import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { ReactiveFormsModule } from '@angular/forms';
import { MatCheckboxModule } from '@angular/material/checkbox';
import { CheckboxComponent } from './checkbox/checkbox.component';
import { MatIconModule } from '@angular/material/icon';
import { MatTooltipModule } from '@angular/material/tooltip';
import { SharedUiModule } from 'src/shared-ui/shared-ui.module';

@NgModule({
  imports: [
    CommonModule,
    ReactiveFormsModule,
    MatCheckboxModule,
    MatIconModule,
    MatTooltipModule,
    SharedUiModule,
    CheckboxComponent,
  ],
  exports: [CheckboxComponent],
})
export class CheckboxModule {}
