import { Component } from '@angular/core';
import { FormButtonComponent } from '../form-button/form-button.component';

@Component({
    selector: 'app-edit-button',
    templateUrl: './edit-button.component.html',
    styleUrls: ['./edit-button.component.scss'],
    standalone: false
})
export class EditButtonComponent extends FormButtonComponent {}
