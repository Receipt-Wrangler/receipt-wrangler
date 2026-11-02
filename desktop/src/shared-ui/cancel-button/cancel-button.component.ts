import { Component, input, output } from '@angular/core';
import { FormMode } from 'src/enums/form-mode.enum';

@Component({
    selector: 'app-cancel-button',
    templateUrl: './cancel-button.component.html',
    styleUrls: ['./cancel-button.component.scss'],
    standalone: false
})
export class CancelButtonComponent {
  public readonly mode = input<FormMode>();

  public readonly disabled = input<boolean>(false);

  public readonly tooltip = input<string>("Cancel");

  public readonly clicked = output<void>();

  public formMode = FormMode;
}
