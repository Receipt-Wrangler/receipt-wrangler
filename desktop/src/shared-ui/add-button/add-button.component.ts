import { Component } from '@angular/core';
import { FormButtonComponent } from '../form-button/form-button.component';

@Component({
    selector: 'app-add-button',
    templateUrl: './add-button.component.html',
    styleUrls: ['./add-button.component.scss'],
    standalone: false
})
export class AddButtonComponent extends FormButtonComponent {}
