import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { MatSlideToggleModule } from '@angular/material/slide-toggle';

import { SlideToggleComponent } from './slide-toggle/slide-toggle.component';

@NgModule({
  imports: [CommonModule, MatSlideToggleModule, SlideToggleComponent],
  exports: [SlideToggleComponent],
})
export class SlideToggleModule {}
