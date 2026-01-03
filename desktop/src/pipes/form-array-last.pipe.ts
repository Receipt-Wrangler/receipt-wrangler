import { Pipe, PipeTransform } from '@angular/core';
import { FormArray, FormControl, FormGroup } from '@angular/forms';

@Pipe({
    name: 'formArrayLast',
    standalone: false
})
export class FormArrayLastPipe implements PipeTransform {
  public transform(array: FormArray): any {
    if (array.length === 0) {
      return null;
    }

    if (array.length > 0) {
      return array.at(array.length - 1);
    }
  }
}
