import { Pipe, PipeTransform } from "@angular/core";

@Pipe({ name: "isRowExpandable" })
export class RowExpandablePipe implements PipeTransform {
  public transform(row: any, isExpandableFunc: (row: any) => boolean): boolean {
    return isExpandableFunc(row);
  }
}
