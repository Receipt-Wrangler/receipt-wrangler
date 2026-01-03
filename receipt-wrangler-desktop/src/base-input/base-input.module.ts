import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { BaseInputComponent } from './base-input/base-input.component';

@NgModule({
  declarations: [BaseInputComponent],
  imports: [CommonModule],
  exports: [BaseInputComponent],
})
export class BaseInputModule {}
