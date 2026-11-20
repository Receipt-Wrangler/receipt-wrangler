import { Component, input } from '@angular/core';
import { MatIcon } from '@angular/material/icon';
import { MatTooltip } from '@angular/material/tooltip';

@Component({
    selector: 'app-help-icon',
    templateUrl: './help-icon.component.html',
    styleUrls: ['./help-icon.component.scss'],
    imports: [MatIcon, MatTooltip]
})
export class HelpIconComponent {
  public readonly helpText = input<string>();
}
