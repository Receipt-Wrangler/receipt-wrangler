import { Pipe, PipeTransform } from "@angular/core";
import { FormControl, FormGroup } from "@angular/forms";

@Pipe({
    name: "formGet",
    standalone: false
})
export class FormGetPipe implements PipeTransform {
  transform(form: FormGroup, path: string): FormControl {
    const result = form.get(path);
    if (result) {
      return result as FormControl;
    } else {
      return new FormControl();
    }
  }
}
