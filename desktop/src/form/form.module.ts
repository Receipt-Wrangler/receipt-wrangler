import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { BaseFormComponent } from './base-form/base-form.component';



@NgModule({
  imports: [
    CommonModule,
    BaseFormComponent
  ],
  exports: [
    BaseFormComponent
  ]
})
export class FormModule { }
