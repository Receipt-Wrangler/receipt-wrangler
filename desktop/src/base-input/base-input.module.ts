import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { BaseInputComponent } from './base-input/base-input.component';

@NgModule({
  imports: [CommonModule, BaseInputComponent],
  exports: [BaseInputComponent],
})
export class BaseInputModule {}
