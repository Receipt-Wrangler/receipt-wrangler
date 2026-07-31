import { provideZonelessChangeDetection } from "@angular/core";
import { ComponentFixture, TestBed } from "@angular/core/testing";
import { FormControl } from "@angular/forms";
import { Category, Tag } from "../../open-api";
import { GrantPickerComponent, GrantSelection } from "./grant-picker.component";

describe("GrantPickerComponent", () => {
  let fixture: ComponentFixture<GrantPickerComponent>;
  let component: GrantPickerComponent;

  const categories = [
    { id: 1, name: "Child A" },
    { id: 2, name: "Child B" },
    { id: 3, name: "Child C" },
  ] as Category[];
  const tags = [
    { id: 10, name: "Urgent" },
    { id: 11, name: "Archived" },
  ] as Tag[];

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [GrantPickerComponent],
      providers: [provideZonelessChangeDetection()],
    }).compileComponents();

    fixture = TestBed.createComponent(GrantPickerComponent);
    component = fixture.componentInstance;
  });

  it("offers the whole pool when no ceiling is given", async () => {
    fixture.componentRef.setInput("categoryPool", categories);
    fixture.componentRef.setInput("tagPool", tags);
    await fixture.whenStable();

    expect(component.availableCategories().length).toBe(3);
    expect(component.availableTags().length).toBe(2);
    expect(component.categoryCeilingHint()).toBe("");
  });

  it("narrows the offered pool to the ceiling and explains why", async () => {
    fixture.componentRef.setInput("categoryPool", categories);
    fixture.componentRef.setInput("ceiling", {
      categoryIds: [1, 2],
      tagIds: [],
      label: "role Foster Parent",
    });
    await fixture.whenStable();

    expect(component.availableCategories().map((c) => c.id)).toEqual([1, 2]);
    expect(component.categoryCeilingHint()).toContain("role Foster Parent");
  });

  it("treats an empty ceiling id list as unrestricted for that resource", async () => {
    // Mirrors the backend rule: a role with no grants for a resource restricts nothing.
    fixture.componentRef.setInput("tagPool", tags);
    fixture.componentRef.setInput("ceiling", {
      categoryIds: [1],
      tagIds: [],
      label: "role Foster Parent",
    });
    await fixture.whenStable();

    expect(component.availableTags().length).toBe(2);
    expect(component.tagCeilingHint()).toBe("");
  });

  it("resolves selected ids to pool objects once the pool arrives", async () => {
    fixture.componentRef.setInput("selectedCategoryIds", [2]);
    await fixture.whenStable();
    expect(component.grantedCategories.length).toBe(0);

    fixture.componentRef.setInput("categoryPool", categories);
    await fixture.whenStable();

    expect(component.grantedCategories.length).toBe(1);
    expect(component.selection().categoryIds).toEqual([2]);
  });

  it("drops a selected id that falls outside the ceiling", async () => {
    // A role narrowed after the member was assigned: the stale id must not be
    // offered back as though it were still granted.
    fixture.componentRef.setInput("categoryPool", categories);
    fixture.componentRef.setInput("selectedCategoryIds", [3]);
    fixture.componentRef.setInput("ceiling", {
      categoryIds: [1, 2],
      tagIds: [],
      label: "role Foster Parent",
    });
    await fixture.whenStable();

    expect(component.selection().categoryIds).toEqual([]);
  });

  it("does not emit while merely seeding itself from inputs", async () => {
    const emissions: GrantSelection[] = [];
    component.grantsChange.subscribe((selection) => emissions.push(selection));

    fixture.componentRef.setInput("categoryPool", categories);
    fixture.componentRef.setInput("selectedCategoryIds", [1]);
    await fixture.whenStable();

    // Seeding uses emitEvent: false — an emission here would look like a user edit
    // and could overwrite a host form's loaded values.
    expect(emissions).toEqual([]);
  });

  it("emits the selection when the user edits it", async () => {
    const emissions: GrantSelection[] = [];
    component.grantsChange.subscribe((selection) => emissions.push(selection));

    fixture.componentRef.setInput("categoryPool", categories);
    await fixture.whenStable();

    component.grantedCategories.push(new FormControl(categories[0] as any));
    await fixture.whenStable();

    expect(emissions.length).toBe(1);
    expect(emissions[0].categoryIds).toEqual([1]);
  });
});
