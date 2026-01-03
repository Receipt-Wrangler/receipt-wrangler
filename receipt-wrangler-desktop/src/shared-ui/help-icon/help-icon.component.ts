import { Component, Input } from '@angular/core';

@Component({
    selector: 'app-help-icon',
    templateUrl: './help-icon.component.html',
    styleUrls: ['./help-icon.component.scss'],
    standalone: false
})
export class HelpIconComponent {
  @Input() public helpText?: string;
}
