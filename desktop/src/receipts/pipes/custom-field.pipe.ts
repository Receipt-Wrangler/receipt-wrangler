import { Pipe, PipeTransform } from "@angular/core";
import { CustomField } from "../../open-api/index";

@Pipe({
  name: "customField",
  standalone: false
})
export class CustomFieldPipe implements PipeTransform {

  public transform(customFieldId: number, customFields: CustomField[]): CustomField | null {
    return customFields.find(customField => customField.id === customFieldId) || null;
  }

}
