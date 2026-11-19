import { Component } from '@angular/core';
import { FormButtonComponent } from '../form-button/form-button.component';
import { ButtonComponent } from '../../button/button/button.component';

@Component({
    selector: 'app-back-button',
    templateUrl: './back-button.component.html',
    styleUrls: ['./back-button.component.scss'],
    imports: [ButtonComponent]
})
export class BackButtonComponent extends FormButtonComponent {}
