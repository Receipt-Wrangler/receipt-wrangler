import { Component, input } from "@angular/core";
import { BreadcrumbItem } from "./breadcrumb-item.interface";
import { RouterLink } from "@angular/router";
import { MatIcon } from "@angular/material/icon";

/**
 * Generic, reusable breadcrumb trail.
 *
 * Fully input-driven so it can be dropped into any page — the consumer passes
 * an ordered list of {@link BreadcrumbItem}s. The final item is rendered as the
 * current page (non-link, emphasised) regardless of whether it has a
 * `routerLink`.
 *
 * @example
 * <app-breadcrumb [items]="[{ label: 'Admin' }, { label: 'Roles' }]"></app-breadcrumb>
 */
@Component({
  selector: "app-breadcrumb",
  templateUrl: "./breadcrumb.component.html",
  styleUrl: "./breadcrumb.component.scss",
  imports: [RouterLink, MatIcon],
})
export class BreadcrumbComponent {
  public readonly items = input.required<BreadcrumbItem[]>();
}
