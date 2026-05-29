import { provideZonelessChangeDetection } from "@angular/core";
import { ComponentFixture, TestBed } from "@angular/core/testing";
import { MatIconModule } from "@angular/material/icon";
import { RouterModule } from "@angular/router";
import { BreadcrumbComponent } from "./breadcrumb.component";

describe("BreadcrumbComponent", () => {
  let component: BreadcrumbComponent;
  let fixture: ComponentFixture<BreadcrumbComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [BreadcrumbComponent],
      imports: [MatIconModule, RouterModule.forRoot([])],
      providers: [provideZonelessChangeDetection()],
    }).compileComponents();

    fixture = TestBed.createComponent(BreadcrumbComponent);
    component = fixture.componentInstance;
  });

  it("creates", async () => {
    fixture.componentRef.setInput("items", [{ label: "Admin" }]);
    await fixture.whenStable();
    expect(component).toBeTruthy();
  });

  it("renders one crumb per item with a separator between them", async () => {
    fixture.componentRef.setInput("items", [
      { label: "Admin", routerLink: ["/admin"] },
      { label: "Roles" },
    ]);
    await fixture.whenStable();

    const crumbs = fixture.nativeElement.querySelectorAll(".crumb");
    const separators = fixture.nativeElement.querySelectorAll(".separator");
    expect(crumbs.length).toBe(2);
    expect(separators.length).toBe(1);
  });

  it("links crumbs that have a routerLink and marks the last as the current page", async () => {
    fixture.componentRef.setInput("items", [
      { label: "Admin", routerLink: ["/admin"] },
      { label: "Roles" },
    ]);
    await fixture.whenStable();

    const link = fixture.nativeElement.querySelector("a.crumb");
    const current = fixture.nativeElement.querySelector(".crumb.current");
    expect(link?.textContent.trim()).toContain("Admin");
    expect(current?.getAttribute("aria-current")).toBe("page");
    expect(current?.textContent.trim()).toContain("Roles");
  });

  it("renders the last crumb as non-link even when it has a routerLink", async () => {
    fixture.componentRef.setInput("items", [
      { label: "Admin", routerLink: ["/admin"] },
      { label: "Roles", routerLink: ["/roles"] },
    ]);
    await fixture.whenStable();

    const links = fixture.nativeElement.querySelectorAll("a.crumb");
    expect(links.length).toBe(1);
    expect(links[0].textContent.trim()).toContain("Admin");
  });
});
