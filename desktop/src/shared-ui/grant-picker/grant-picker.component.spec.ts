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

  it("keeps a user's edit when the pool changes afterwards", async () => {
    // Regression: the seeding effect re-runs whenever its inputs change, and the
    // pool arrives asynchronously AFTER the host has mounted the picker. Before
    // the guard, that late run reset the FormArray and silently discarded the
    // admin's pick — the host then saved the original value.
    fixture.componentRef.setInput("categoryPool", categories);
    fixture.componentRef.setInput("selectedCategoryIds", [1]);
    await fixture.whenStable();

    component.grantedCategories.push(new FormControl(categories[1] as any));
    await fixture.whenStable();
    expect(component.selection().categoryIds).toEqual([1, 2]);

    // A late pool refresh must not undo it.
    fixture.componentRef.setInput("categoryPool", [...categories, { id: 4, name: "Child D" } as Category]);
    await fixture.whenStable();

    expect(component.selection().categoryIds).toEqual([1, 2]);
  });

  it("keeps a user's edit when the ceiling arrives afterwards", async () => {
    // The host's role list resolves after the picker mounts, so the ceiling flips
    // from null to a real set — the other late input that used to clobber an edit.
    fixture.componentRef.setInput("categoryPool", categories);
    await fixture.whenStable();

    component.grantedCategories.push(new FormControl(categories[0] as any));
    await fixture.whenStable();

    fixture.componentRef.setInput("ceiling", {
      categoryIds: [1, 2],
      tagIds: [],
      label: "role Foster Parent",
    });
    await fixture.whenStable();

    expect(component.selection().categoryIds).toEqual([1]);
  });

  it("keeps a tag edit independently of the category pool changing", async () => {
    fixture.componentRef.setInput("tagPool", tags);
    await fixture.whenStable();

    component.grantedTags.push(new FormControl(tags[0] as any));
    await fixture.whenStable();

    fixture.componentRef.setInput("categoryPool", categories);
    await fixture.whenStable();

    expect(component.selection().tagIds).toEqual([10]);
  });

  it("gives each instance a distinct input id", async () => {
    // The user form renders one picker per group. The base autocomplete otherwise
    // derives the input's DOM id from its label, so every instance would render
    // id="categories" — which breaks <label for> association for every field after
    // the first and misdirects the base component's getElementById filter clear.
    const second = TestBed.createComponent(GrantPickerComponent);

    expect(component.categoryInputId).not.toBe(second.componentInstance.categoryInputId);
    expect(component.tagInputId).not.toBe(second.componentInstance.tagInputId);
    expect(component.categoryInputId).not.toBe(component.tagInputId);
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
