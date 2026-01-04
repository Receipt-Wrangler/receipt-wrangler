import { Pipe, PipeTransform } from "@angular/core";
import { OcrEngine } from "../../open-api";
import { ocrEngineOptions } from "../constants/ocr-engine-options";

@Pipe({
  name: "ocrEngine",
  standalone: true
})
export class OcrEnginePipe implements PipeTransform {
  public transform(value: OcrEngine): string {
    return ocrEngineOptions.find(option => option.value === value)?.displayValue ?? "";
  }
}
