import { Component } from '@angular/core';
import { FormButtonComponent } from '../form-button/form-button.component';
import { ButtonComponent } from '../../button/button/button.component';
import { RouterLink } from '@angular/router';

@Component({
    selector: 'app-add-button',
    templateUrl: './add-button.component.html',
    styleUrls: ['./add-button.component.scss'],
    imports: [ButtonComponent, RouterLink]
})
export class AddButtonComponent extends FormButtonComponent {}
