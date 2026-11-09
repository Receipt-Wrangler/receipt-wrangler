import { Component, input } from '@angular/core';
import { MatSlideToggle } from '@angular/material/slide-toggle';

@Component({
    selector: 'app-slide-toggle',
    templateUrl: './slide-toggle.component.html',
    styleUrls: ['./slide-toggle.component.scss'],
    imports: [MatSlideToggle]
})
export class SlideToggleComponent {
  public readonly color = input<string>('');

  public readonly checked = input<boolean>(false);

  public readonly disabled = input<boolean>(false);
}
