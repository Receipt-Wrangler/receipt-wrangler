import { Inject, LOCALE_ID, Pipe, PipeTransform } from "@angular/core";
import { Store } from "@ngxs/store";
import { SearchResult } from "../../open-api";
import { formatAppDate } from "../../pipes/app-date.pipe";
import { GroupState } from "../../store";
import { SystemSettingsState } from "../../store/system-settings.state";

@Pipe({
    name: "searchResult",
    standalone: false
})
export class SearchResultPipe implements PipeTransform {
  constructor(private store: Store, @Inject(LOCALE_ID) private locale: string) {}

  public transform(searchResult: SearchResult): string {
    // A receipt's date is a calendar day: its stored UTC day, whatever the zone.
    const date = formatAppDate(
      searchResult.date,
      "mediumDate",
      "calendar",
      this.store.selectSnapshot(SystemSettingsState.timeZone),
      this.locale
    );
    const group = this.store.selectSnapshot(
      GroupState.getGroupById(searchResult?.groupId?.toString())
    );

    return `${date} - ${searchResult.name} (${group?.name})`;
  }
}
