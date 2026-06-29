import { Component, input } from "@angular/core";
import { BreadcrumbItem } from "./breadcrumb-item.interface";

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
  standalone: false,
})
export class BreadcrumbComponent {
  public readonly items = input.required<BreadcrumbItem[]>();
}
