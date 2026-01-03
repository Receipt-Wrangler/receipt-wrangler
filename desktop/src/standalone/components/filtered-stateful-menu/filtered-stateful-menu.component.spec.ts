import { CUSTOM_ELEMENTS_SCHEMA } from "@angular/core";
import { ComponentFixture, TestBed } from "@angular/core/testing";
import { NoopAnimationsModule } from "@angular/platform-browser/animations";
import { ActivatedRoute } from "@angular/router";
import { ButtonModule } from "../../../button/index";
import { InputModule } from "../../../input/index";

import { FilteredStatefulMenuComponent } from "./filtered-stateful-menu.component";
import { StatefulMenuItem } from "./stateful-menu-item";

describe("FilteredStatefulMenuComponent", () => {
  let component: FilteredStatefulMenuComponent;
  let fixture: ComponentFixture<FilteredStatefulMenuComponent>;
  let mockItems: StatefulMenuItem[];

  beforeEach(async () => {
    mockItems = [
      { displayValue: "Item 1", value: "item1", selected: false },
      { displayValue: "Item 2", value: "item2", selected: true },
      { displayValue: "Another Item", value: "item3", selected: false }
    ];

    await TestBed.configureTestingModule({
      imports: [
        FilteredStatefulMenuComponent,
        ButtonModule,
        InputModule,
        NoopAnimationsModule
      ],
      providers: [
        {
          provide: ActivatedRoute,
          useValue: {}
        }
      ],
      schemas: [CUSTOM_ELEMENTS_SCHEMA]
    })
      .compileComponents();

    fixture = TestBed.createComponent(FilteredStatefulMenuComponent);
    component = fixture.componentInstance;
    component.items = mockItems;
    fixture.detectChanges();
  });

  it("should create", () => {
    expect(component).toBeTruthy();
  });

  it("should initialize filteredItems with all items", () => {
    component.ngOnChanges({ items: { currentValue: mockItems, previousValue: null, firstChange: true, isFirstChange: () => true } });
    expect(component.filteredItems).toEqual(mockItems);
  });

  it("should filter items based on filter string", () => {
    const filterString = "Another";
    component.filterFormControl.setValue(filterString);
    fixture.detectChanges();

    expect(component.filteredItems.length).toBe(1);
    expect(component.filteredItems[0].displayValue).toBe("Another Item");
  });

  it("should filter items using custom filter function if provided", () => {
    // Custom filter function that only returns items with selected = true
    component.filterFunc = (item: StatefulMenuItem) => item.selected;
    component.filterFormControl.setValue("dummy"); // Value doesn't matter with our custom filter
    fixture.detectChanges();

    expect(component.filteredItems.length).toBe(1);
    expect(component.filteredItems[0].value).toBe("item2");
  });

  it("should reset filter when resetFilter is called", () => {
    component.filterFormControl.setValue("some filter");
    component.resetFilter();
    expect(component.filterFormControl.value).toBe("");
    expect(component.filteredItems.length).toBe(mockItems.length);
  });

  it("should handle item selection and emit itemSelected event", () => {
    jest.spyOn(component.itemSelected, "emit");
    const mockEvent = new MouseEvent("click");
    const mockItem = mockItems[0];

    component.onItemSelected(mockItem, mockEvent);

    expect(component.itemSelected.emit).toHaveBeenCalledWith(mockItem);
  });

  it("should not emit itemSelected event when readonly is true", () => {
    jest.spyOn(component.itemSelected, "emit");
    component.readonly = true;
    const mockEvent = new MouseEvent("click");
    const mockItem = mockItems[0];

    component.onItemSelected(mockItem, mockEvent);

    expect(component.itemSelected.emit).not.toHaveBeenCalled();
  });

  it("should update filteredItems when items input changes", () => {
    const newItems: StatefulMenuItem[] = [
      { displayValue: "New Item", value: "new", selected: false }
    ];

    component.ngOnChanges({
      items: {
        currentValue: newItems,
        previousValue: mockItems,
        firstChange: false,
        isFirstChange: () => false
      }
    });

    expect(component.filteredItems).toEqual(newItems);
  });

  it("should show all items when filter is cleared", () => {
    // First apply a filter
    component.filterFormControl.setValue("Another");
    fixture.detectChanges();
    expect(component.filteredItems.length).toBe(1);

    // Then clear the filter
    component.filterFormControl.setValue("");
    fixture.detectChanges();
    expect(component.filteredItems.length).toBe(mockItems.length);
  });
});
