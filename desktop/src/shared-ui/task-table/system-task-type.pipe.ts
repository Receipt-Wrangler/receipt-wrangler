import { Pipe, PipeTransform } from "@angular/core";
import { SystemTaskType } from "../../open-api";

@Pipe({
  name: "systemTaskType",
  standalone: true
})
export class SystemTaskTypePipe implements PipeTransform {

  public transform(value: SystemTaskType): string {
    switch (value) {
      case "MAGIC_FILL":
        return "Magic Fill";
      case "QUICK_SCAN":
        return "Quick Scan";
      case "SYSTEM_EMAIL_CONNECTIVITY_CHECK":
        return "System Email Connectivity Check";
      case "RECEIPT_PROCESSING_SETTINGS_CONNECTIVITY_CHECK":
        return "Receipt Processing Settings Connectivity Check";
      case "EMAIL_READ":
        return "Email Read";
      case "EMAIL_UPLOAD":
        return "Email Upload";
      case "CHAT_COMPLETION":
        return "Chat Completion";
      case "OCR_PROCESSING":
        return "OCR Processing";
      case "RECEIPT_UPLOADED":
        return "Receipt Uploaded";
      case "PROMPT_GENERATED":
        return "Prompt Generated";
      case "RECEIPT_UPDATED":
        return "Updated Receipt";
    }
    return "";
  }
}
