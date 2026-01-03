import { AbstractControl, FormArray, FormGroup } from '@angular/forms';

export interface FormCommand {
  path?: string;
  payload?: any;
  command: keyof AbstractControl | keyof FormArray | keyof FormGroup;
}
