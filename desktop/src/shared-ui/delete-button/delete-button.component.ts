import { Component } from '@angular/core';
import { FormButtonComponent } from '../form-button/form-button.component';

@Component({
    selector: 'app-delete-button',
    templateUrl: './delete-button.component.html',
    styleUrls: ['./delete-button.component.scss'],
    standalone: false
})
export class DeleteButtonComponent extends FormButtonComponent {
  constructor() {
    super();
    this.color = 'warn';
  }
}
