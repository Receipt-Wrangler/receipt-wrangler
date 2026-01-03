import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { AvatarComponent } from './avatar/avatar.component';
import { MatTooltipModule } from '@angular/material/tooltip';

@NgModule({
  declarations: [AvatarComponent],
  imports: [CommonModule, MatTooltipModule],
  exports: [AvatarComponent],
})
export class AvatarModule {}
