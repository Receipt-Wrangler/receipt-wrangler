import { FormGroup } from "@angular/forms";
import { FormMode } from "src/enums/form-mode.enum";
import { FormConfig } from "src/interfaces";
import { FormCommand } from "../form/index";

export function applyApiErrors(form: FormGroup, errors: any): void {
  const keys = Array.from(Object.keys(errors?.error));

  keys?.forEach((k) => {
    const field = form.get(k);
    if (field) {
      field.setErrors({
        ...form.errors,
        error: errors.error[k],
      });
    }
  });
}

export function setEntityHeaderText(
  entity: any,
  nameKey: string,
  formConfig: FormConfig,
  entityType?: string
): string {
  if (entity && entity[nameKey] && formConfig) {
    let entityName = entity[nameKey];
    let headerText = "";

    if (formConfig.mode === FormMode.view) {
      headerText = `View ${entityName}`;
    }

    if (formConfig.mode === FormMode.edit) {
      headerText = `Edit ${entityName}`;
    }

    if (formConfig.mode === FormMode.add) {
      headerText = `Create ${entityName}`;
    }

    if (entityType) {
      headerText += ` ${entityType}`;
    }

    return headerText;
  }

  return "";
}

export function applyFormCommand(form: FormGroup, formCommand: FormCommand) {
  if (formCommand.path) {
    (form.get(formCommand.path) as any)[formCommand.command](
      formCommand.payload
    );
  } else {
    (form as any)[formCommand.command](formCommand.payload);
  }
}
