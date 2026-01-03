import { Component, EventEmitter, Input, Output } from '@angular/core';
import { FormMode } from 'src/enums/form-mode.enum';

@Component({
    selector: 'app-cancel-button',
    templateUrl: './cancel-button.component.html',
    styleUrls: ['./cancel-button.component.scss'],
    standalone: false
})
export class CancelButtonComponent {
  @Input() public mode?: FormMode;

  @Input() public disabled: boolean = false;

  @Output() public clicked: EventEmitter<void> = new EventEmitter<void>();

  public formMode = FormMode;
}
