import { Component, Input } from '@angular/core';
import { FormControl } from '@angular/forms';
import { RadioButtonData } from '../models';

@Component({
    selector: 'app-radio-group',
    templateUrl: './radio-group.component.html',
    styleUrls: ['./radio-group.component.scss'],
    standalone: false
})
export class RadioGroupComponent {
  @Input() public radioButtonData: RadioButtonData[] = [];

  @Input() public inputFormControl: FormControl = new FormControl('');
}
