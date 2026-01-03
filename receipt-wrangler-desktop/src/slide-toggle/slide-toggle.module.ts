import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { MatSlideToggleModule } from '@angular/material/slide-toggle';

import { SlideToggleComponent } from './slide-toggle/slide-toggle.component';

@NgModule({
  declarations: [SlideToggleComponent],
  imports: [CommonModule, MatSlideToggleModule],
  exports: [SlideToggleComponent],
})
export class SlideToggleModule {}
