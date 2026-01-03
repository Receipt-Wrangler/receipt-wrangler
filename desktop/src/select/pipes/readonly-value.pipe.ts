import { Pipe, PipeTransform } from "@angular/core";
import { OptionDisplayPipe } from "./option-display.pipe";

@Pipe({
  name: "readonlyValue",
  standalone: false,
})
export class ReadonlyValuePipe implements PipeTransform {
  constructor(private optionDisplayPipe: OptionDisplayPipe) {}

  public transform(option: any, options: any[], optionDisplayKey: string, optionValueKey?: string, optionsDisplayArray?: any[]): any {
    const transformedOptions = options.map(o => optionValueKey ? o[optionValueKey] : o);
    const optionIndex = transformedOptions.findIndex(o => o?.toString() === option?.toString());

    if (optionIndex === -1) {
      return null;
    }

    return this.optionDisplayPipe.transform(optionIndex, options[optionIndex], optionDisplayKey, optionsDisplayArray);
  }

}
