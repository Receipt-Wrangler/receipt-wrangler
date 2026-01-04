import { DatePipe } from "@angular/common";
import { Pipe, PipeTransform } from "@angular/core";
import { Store } from "@ngxs/store";
import { SearchResult } from "../../open-api";
import { GroupState } from "../../store";

@Pipe({
    name: "searchResult",
    standalone: false
})
export class SearchResultPipe implements PipeTransform {
  constructor(private datePipe: DatePipe, private store: Store) {}

  public transform(searchResult: SearchResult): string {
    const date = this.datePipe.transform(searchResult.date);
    const group = this.store.selectSnapshot(
      GroupState.getGroupById(searchResult?.groupId?.toString())
    );

    return `${date} - ${searchResult.name} (${group?.name})`;
  }
}
