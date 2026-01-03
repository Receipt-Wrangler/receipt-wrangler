import { Component, Input } from '@angular/core';

@Component({
    selector: 'app-slide-toggle',
    templateUrl: './slide-toggle.component.html',
    styleUrls: ['./slide-toggle.component.scss'],
    standalone: false
})
export class SlideToggleComponent {
  @Input() public color: string = '';

  @Input() public checked: boolean = false;

  @Input() public disabled: boolean = false;
}
