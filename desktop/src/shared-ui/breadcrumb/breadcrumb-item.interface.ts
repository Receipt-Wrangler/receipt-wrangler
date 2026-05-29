/**
 * A single entry in a breadcrumb trail.
 *
 * Provide a `routerLink` for crumbs that should navigate; omit it for the
 * current/terminal crumb. `icon` is an optional leading Material icon.
 */
export interface BreadcrumbItem {
  label: string;
  routerLink?: string | any[];
  icon?: string;
}
