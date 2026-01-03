import { Pipe, PipeTransform } from "@angular/core";

@Pipe({
  name: "optionDisplay",
  standalone: false,
})
export class OptionDisplayPipe implements PipeTransform {

  public transform(index: number, option: any, optionDisplayKey: string, optionsDisplayArray?: any[]): string {
    if (optionsDisplayArray && optionsDisplayArray.length > 0) {
      return optionsDisplayArray[index];
    } else {
      return optionDisplayKey ? option[optionDisplayKey] : option;
    }
  }
}
