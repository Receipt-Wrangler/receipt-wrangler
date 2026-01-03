import { Pipe, PipeTransform } from '@angular/core';
import { FormMode } from 'src/enums/form-mode.enum';

@Pipe({
    name: 'inputReadonly',
    standalone: false
})
export class InputReadonlyPipe implements PipeTransform {
  public transform(mode: FormMode): boolean {
    return mode === FormMode.view;
  }
}
