import { Pipe, PipeTransform } from "@angular/core";
import { AiType } from "../../open-api/index";

@Pipe({
  name: "urlLabel"
})
export class UrlLabelPipe implements PipeTransform {
  public transform(aiType: AiType): string {
    if (aiType === AiType.OpenAiCustom) {
      return "Base Url";
    }

    return "Url";
  }
}
