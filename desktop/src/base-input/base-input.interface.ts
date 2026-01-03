import { FormControl } from '@angular/forms';

export interface BaseInputInterface {
  inputFormControl: FormControl;
  label: string;
  additionalErrorMessages?: { [key: string]: string };
  readonly: boolean;
  placeholder?: string;
  helpText?: string;
}
