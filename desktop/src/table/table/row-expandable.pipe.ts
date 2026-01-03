import { Pipe, PipeTransform } from "@angular/core";

@Pipe({
    name: "isRowExpandable",
    standalone: false
})
export class RowExpandablePipe implements PipeTransform {
  public transform(row: any, isExpandableFunc: (row: any) => boolean): boolean {
    return isExpandableFunc(row);
  }
}
