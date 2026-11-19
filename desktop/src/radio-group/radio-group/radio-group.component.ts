import { Component, input } from '@angular/core';
import { FormControl, ReactiveFormsModule } from '@angular/forms';
import { RadioButtonData } from '../models';
import { MatRadioGroup, MatRadioButton } from '@angular/material/radio';

@Component({
    selector: 'app-radio-group',
    templateUrl: './radio-group.component.html',
    styleUrls: ['./radio-group.component.scss'],
    imports: [MatRadioGroup, ReactiveFormsModule, MatRadioButton]
})
export class RadioGroupComponent {
  public readonly radioButtonData = input<RadioButtonData[]>([]);

  public readonly inputFormControl = input<FormControl>(new FormControl(''));
}
