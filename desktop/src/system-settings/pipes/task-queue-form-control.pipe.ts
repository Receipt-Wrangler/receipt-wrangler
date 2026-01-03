import { Pipe, PipeTransform } from "@angular/core";
import { FormArray, FormControl, FormGroup } from "@angular/forms";

@Pipe({
    name: "taskQueueFormControl",
    standalone: false
})
export class TaskQueueFormControlPipe implements PipeTransform {
  public transform(form: FormGroup, queueName: string, formKey: string): FormControl {
    const formControl = (form.get("taskQueueConfigurations") as FormArray)?.controls?.find(c => c?.value?.["name"] === queueName);

    if (formControl) {
      return formControl.get(formKey) as FormControl;
    }

    console.error(new Error(`Form control with name ${queueName} not found`));
    return new FormControl();
  }
}
