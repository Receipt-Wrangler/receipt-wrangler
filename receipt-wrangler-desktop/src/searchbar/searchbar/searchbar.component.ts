import { Component, ViewEncapsulation } from "@angular/core";
import { FormControl } from "@angular/forms";
import { Router } from "@angular/router";
import { UntilDestroy, untilDestroyed } from "@ngneat/until-destroy";
import { of, switchMap, take, tap } from "rxjs";
import { SearchResult, SearchService } from "../../open-api";

@UntilDestroy()
@Component({
    selector: "app-searchbar",
    templateUrl: "./searchbar.component.html",
    styleUrls: ["./searchbar.component.scss"],
    encapsulation: ViewEncapsulation.None,
    standalone: false
})
export class SearchbarComponent {
  public results: SearchResult[] = [];

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
          this.results = results;
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
