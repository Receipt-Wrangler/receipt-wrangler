import { provideZonelessChangeDetection } from "@angular/core";
import { ComponentFixture, TestBed } from "@angular/core/testing";
import { MatIconModule } from "@angular/material/icon";
import { FilterBarComponent } from "./filter-bar.component";
import { FilterTab } from "./filter-tab.interface";

describe("FilterBarComponent", () => {
  let component: FilterBarComponent;
  let fixture: ComponentFixture<FilterBarComponent>;

  const TABS: FilterTab[] = [
    { value: "all", label: "All roles", icon: "list", count: 7 },
    { value: "app", label: "Application", icon: "apps", count: 3 },
    { value: "group", label: "Group", count: 4 },
  ];

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [FilterBarComponent],
      imports: [MatIconModule],
      providers: [provideZonelessChangeDetection()],
    }).compileComponents();

    fixture = TestBed.createComponent(FilterBarComponent);
    component = fixture.componentInstance;
  });

  it("creates", async () => {
    fixture.componentRef.setInput("tabs", TABS);
    await fixture.whenStable();
    expect(component).toBeTruthy();
  });

  it("renders one button per tab", async () => {
    fixture.componentRef.setInput("tabs", TABS);
    await fixture.whenStable();

    expect(fixture.nativeElement.querySelectorAll(".filter-tab").length).toBe(3);
  });

  it("marks the tab matching value as active", async () => {
    fixture.componentRef.setInput("tabs", TABS);
    fixture.componentRef.setInput("value", "app");
    await fixture.whenStable();

    const active = fixture.nativeElement.querySelectorAll(".filter-tab.active");
    expect(active.length).toBe(1);
    expect(active[0].textContent).toContain("Application");
    expect(active[0].getAttribute("aria-selected")).toBe("true");
  });

  it("updates the bound value when a tab is clicked", async () => {
    fixture.componentRef.setInput("tabs", TABS);
    fixture.componentRef.setInput("value", "all");
    await fixture.whenStable();

    const buttons = fixture.nativeElement.querySelectorAll(".filter-tab");
    buttons[1].click();
    await fixture.whenStable();

    expect(component.value()).toBe("app");
  });

  it("renders the count pill only when a count is provided", async () => {
    fixture.componentRef.setInput("tabs", [
      { value: "a", label: "A", count: 5 },
      { value: "b", label: "B" },
    ] satisfies FilterTab[]);
    await fixture.whenStable();

    const pills = fixture.nativeElement.querySelectorAll(".filter-tab .pill");
    expect(pills.length).toBe(1);
    expect(pills[0].textContent.trim()).toBe("5");
  });

  it("renders an icon only when provided", async () => {
    fixture.componentRef.setInput("tabs", [
      { value: "a", label: "A", icon: "apps" },
      { value: "b", label: "B" },
    ] satisfies FilterTab[]);
    await fixture.whenStable();

    expect(fixture.nativeElement.querySelectorAll(".filter-tab mat-icon").length).toBe(1);
  });
});
