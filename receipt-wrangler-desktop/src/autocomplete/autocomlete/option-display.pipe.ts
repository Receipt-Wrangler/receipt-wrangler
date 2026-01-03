import { Pipe, PipeTransform } from "@angular/core";

@Pipe({
    name: "optionDisplay",
    standalone: false
})
export class OptionDisplayPipe implements PipeTransform {
  public transform(
    option: any,
    options: any[],
    optionDisplayKey?: string,
    optionValueKey?: string
  ): any {
    // If we have just the value
    if (optionValueKey && optionDisplayKey) {
      return options.find(
        (o) => o?.[optionValueKey] === (option?.[optionValueKey] ?? option)
      )?.[optionDisplayKey];
    }

    // If we have the whole option object
    if (optionDisplayKey) {
      return option?.[optionDisplayKey];
    }

    return option;
  }
}
