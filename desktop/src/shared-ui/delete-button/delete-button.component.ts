import { Component } from '@angular/core';
import { FormButtonComponent } from '../form-button/form-button.component';
import { ButtonComponent } from '../../button/button/button.component';

@Component({
    selector: 'app-delete-button',
    templateUrl: './delete-button.component.html',
    styleUrls: ['./delete-button.component.scss'],
    imports: [ButtonComponent]
})
export class DeleteButtonComponent extends FormButtonComponent {
  constructor() {
    super();
    this.color = 'warn';
  }
}
