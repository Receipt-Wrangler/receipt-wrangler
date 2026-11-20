import { Component, signal, ViewEncapsulation, DestroyRef, inject } from "@angular/core";
import { takeUntilDestroyed } from "@angular/core/rxjs-interop";
import { FormControl, ReactiveFormsModule } from "@angular/forms";
import { Router } from "@angular/router";
import { of, switchMap, take, tap } from "rxjs";
import { SearchResult, SearchService } from "../../open-api";
import { MatFormField, MatPrefix, MatInput } from "@angular/material/input";
import { MatIcon } from "@angular/material/icon";
import { MatAutocompleteTrigger, MatAutocomplete, MatOption } from "@angular/material/autocomplete";
import { SearchResultPipe } from "../pipes/search-result.pipe";

@Component({
    selector: "app-searchbar",
    templateUrl: "./searchbar.component.html",
    styleUrls: ["./searchbar.component.scss"],
    encapsulation: ViewEncapsulation.None,
    imports: [ReactiveFormsModule, MatFormField, MatPrefix, MatIcon, MatInput, MatAutocompleteTrigger, MatAutocomplete, MatOption, SearchResultPipe]
})
export class SearchbarComponent {
  private readonly destroyRef = inject(DestroyRef);

  public results = signal<SearchResult[]>([]);

  public searchFormControl = new FormControl("");

  public displayWith = (searchResult: SearchResult) => {
    return searchResult?.name;
  };

  constructor(private searchService: SearchService, private router: Router) {}

  public ngOnInit(): void {
    this.searchFormControl.valueChanges
      .pipe(
        takeUntilDestroyed(this.destroyRef),
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
