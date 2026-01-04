import { Pipe, PipeTransform } from "@angular/core";
import { AiType } from "../../open-api";
import { aiTypeOptions } from "../constants/ai-type-options";

@Pipe({
  name: "aiType",
  standalone: true,
})
export class AiTypePipe implements PipeTransform {
  public transform(value: AiType): string {
    return aiTypeOptions.find(option => option.value === value)?.displayValue ?? "";
  }
}
