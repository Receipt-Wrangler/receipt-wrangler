import { Component, input, model } from "@angular/core";
import { FilterTab } from "./filter-tab.interface";

/**
 * Generic, reusable segmented filter bar — the standard pattern for simple
 * filters (e.g. "All / Application / Group").
 *
 * Presentational and fully input-driven: the consumer owns the data and passes
 * the {@link FilterTab}s (with optional per-tab counts); the selected value is a
 * two-way bound `string` id.
 *
 * @example
 * <app-filter-bar
 *   [tabs]="[{ value: 'all', label: 'All', count: 7 }, { value: 'app', label: 'App', count: 3 }]"
 *   [value]="filter()"
 *   (valueChange)="setFilter($event)"
 * ></app-filter-bar>
 */
@Component({
  selector: "app-filter-bar",
  templateUrl: "./filter-bar.component.html",
  styleUrl: "./filter-bar.component.scss",
  standalone: false,
})
export class FilterBarComponent {
  public readonly tabs = input.required<FilterTab[]>();
  public readonly value = model<string>("");

  public select(value: string): void {
    this.value.set(value);
  }
}
