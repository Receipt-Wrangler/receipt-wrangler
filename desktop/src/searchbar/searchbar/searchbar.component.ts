import { Component, signal, ViewEncapsulation } from "@angular/core";
import { FormControl, ReactiveFormsModule } from "@angular/forms";
import { Router } from "@angular/router";
import { UntilDestroy, untilDestroyed } from "@ngneat/until-destroy";
import { of, switchMap, take, tap } from "rxjs";
import { SearchResult, SearchService } from "../../open-api";
import { MatFormField, MatPrefix, MatInput } from "@angular/material/input";
import { MatIcon } from "@angular/material/icon";
import { MatAutocompleteTrigger, MatAutocomplete, MatOption } from "@angular/material/autocomplete";
import { SearchResultPipe } from "../pipes/search-result.pipe";

@UntilDestroy()
@Component({
    selector: "app-searchbar",
    templateUrl: "./searchbar.component.html",
    styleUrls: ["./searchbar.component.scss"],
    encapsulation: ViewEncapsulation.None,
    imports: [ReactiveFormsModule, MatFormField, MatPrefix, MatIcon, MatInput, MatAutocompleteTrigger, MatAutocomplete, MatOption, SearchResultPipe]
})
export class SearchbarComponent {
  public results = signal<SearchResult[]>([]);

  public searchFormControl = new FormControl("");

  public displayWith = (searchResult: SearchResult) => {
    return searchResult?.name;
  };

  constructor(private searchService: SearchService, private router: Router) {}

  public ngOnInit(): void {
    this.searchFormControl.valueChanges
      .pipe(
        untilDestroyed(this),
        switchMap((value) =>
          value
            ? this.searchService.receiptSearch(value ?? "").pipe(take(1))
            : of([] as SearchResult[])
        ),
        tap((results) => {
          this.results.set(results);
        })
      )
      .subscribe();
  }

  public navigateToResult(result: SearchResult) {
    switch (result.type) {
      case "Receipt":
        this.router.navigateByUrl(`/receipts/${result.id}/view`);
        break;
    }
  }
}
