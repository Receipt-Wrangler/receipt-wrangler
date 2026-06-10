import { CurrencyPipe } from "@angular/common";
import { provideHttpClient, withInterceptorsFromDi } from "@angular/common/http";
import { provideHttpClientTesting } from "@angular/common/http/testing";
import { CUSTOM_ELEMENTS_SCHEMA, QueryList, SimpleChange } from "@angular/core";
import { ComponentFixture, TestBed } from "@angular/core/testing";
import { FormArray, FormControl, FormGroup } from "@angular/forms";
import { MatExpansionPanel } from "@angular/material/expansion";
import { ActivatedRoute } from "@angular/router";
import { NgxsModule, Store } from "@ngxs/store";
import { FormMode } from "src/enums/form-mode.enum";
import { PipesModule } from "src/pipes/pipes.module";
import { InputComponent } from "../../input";
import { Category, Group, GroupRole, Item, ItemStatus, Receipt, Tag, User } from "../../open-api";
import { AuthState, UserState } from "../../store/index";
import { SystemSettingsState } from "../../store/system-settings.state";
import { UserTotalWithPercentagePipe } from "../user-total-with-percentage.pipe";
import { buildItemForm } from "../utils/form.utils";

import { ShareListComponent } from "./share-list.component";

describe("ShareListComponent", () => {
  let component: ShareListComponent;
  let fixture: ComponentFixture<ShareListComponent>;
  let store: Store;

  const mockUsers: User[] = [
    { id: 1, username: "user1", displayName: "User One" } as User,
    { id: 2, username: "user2", displayName: "User Two" } as User,
    { id: 3, username: "user3", displayName: "User Three" } as User,
  ];

  const mockItems: Item[] = [
    { id: 1, name: "Item 1", amount: "10.50", chargedToUserId: 1, status: ItemStatus.Open, receiptId: 1 } as Item,
    { id: 2, name: "Item 2", amount: "15.75", chargedToUserId: 2, status: ItemStatus.Open, receiptId: 1 } as Item,
    { id: 3, name: "Item 3", amount: "8.25", chargedToUserId: 1, status: ItemStatus.Resolved, receiptId: 1 } as Item,
    { id: 4, name: "Item 4", amount: "12.00", chargedToUserId: 3, status: ItemStatus.Open, receiptId: 1 } as Item,
  ];

  const mockCategories: Category[] = [
    { id: 1, name: "Food", description: "Food items" } as Category,
    { id: 2, name: "Entertainment", description: "Entertainment items" } as Category,
  ];

  const mockTags: Tag[] = [
    { id: 1, name: "Urgent", description: "Urgent items" } as Tag,
    { id: 2, name: "Business", description: "Business items" } as Tag,
  ];

  const mockReceipt: Receipt = {
    id: 1,
    name: "Test Receipt",
    amount: "46.50",
    date: "2023-01-01",
  } as any as Receipt;

  const mockGroup: Group = {
    id: 1,
    name: "Test Group",
    groupRole: GroupRole.Owner,
  } as any as Group;

  const mockActivatedRoute = {
    snapshot: {
      data: {
        receipt: mockReceipt,
        mode: FormMode.edit,
      },
    },
  };

  function createFormWithItems(items: Item[]): FormGroup {
    const receiptItems = new FormArray(
      items.map(item => buildItemForm(item, mockReceipt.id?.toString(), true, false))
    );

    return new FormGroup({
      receiptItems: receiptItems,
      amount: new FormControl("46.50"),
    });
  }

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ShareListComponent, UserTotalWithPercentagePipe],
      schemas: [CUSTOM_ELEMENTS_SCHEMA],
      imports: [CurrencyPipe, PipesModule, NgxsModule.forRoot([AuthState, UserState, SystemSettingsState])],
      providers: [
        {
          provide: ActivatedRoute,
          useValue: mockActivatedRoute,
        },
        CurrencyPipe,
        provideHttpClient(withInterceptorsFromDi()),
        provideHttpClientTesting(),
      ]
    }).compileComponents();

    fixture = TestBed.createComponent(ShareListComponent);
    component = fixture.componentInstance;
    store = TestBed.inject(Store);

    // Reset store with proper user data structure
    store.reset({
      auth: {},
      users: {
        users: mockUsers
      },
      systemSettings: {}
    });

    // Setup default component state
    fixture.componentRef.setInput('form', createFormWithItems(mockItems));
    fixture.componentRef.setInput('categories', mockCategories);
    fixture.componentRef.setInput('tags', mockTags);
    fixture.componentRef.setInput('selectedGroup', mockGroup);
    component.originalReceipt = mockReceipt;

    fixture.detectChanges();
  });

  it("should create", () => {
    expect(component).toBeTruthy();
  });

  describe("Component Structure & Initialization", () => {
    it("should have all required inputs", () => {
      expect(component.form).toBeDefined();
      expect(component.originalReceipt).toBeDefined();
      expect(component.categories).toBeDefined();
      expect(component.tags).toBeDefined();
      expect(component.selectedGroup).toBeDefined();
      expect(component.triggerAddMode).toBeDefined();
    });

    it("should have all required outputs", () => {
      expect(component.itemAdded).toBeDefined();
      expect(component.itemRemoved).toBeDefined();
      expect(component.allItemsResolved).toBeDefined();
    });

    it("should have all required ViewChildren", () => {
      expect(component.userExpansionPanels).toBeDefined();
      expect(component.nameFields).toBeDefined();
    });

    it("should initialize with route data on ngOnInit", () => {
      jest.spyOn(component, "setUserItemMap");
      component.ngOnInit();

      expect(component.originalReceipt).toEqual(mockReceipt);
      expect(component.mode).toBe(FormMode.edit);
      expect(component.setUserItemMap).toHaveBeenCalled();
    });

    it("should handle missing route data in ngOnInit", () => {
      const activatedRoute = TestBed.inject(ActivatedRoute);
      const originalData = activatedRoute.snapshot.data;
      activatedRoute.snapshot.data = {};

      component.ngOnInit();

      expect(component.originalReceipt).toBeUndefined();
      expect(component.mode).toBeUndefined();

      // Restore route data for subsequent tests
      activatedRoute.snapshot.data = originalData;
    });

    it("should handle ngOnChanges when triggerAddMode changes to true", () => {
      jest.spyOn(component, "initAddMode");
      const changes = {
        triggerAddMode: new SimpleChange(false, true, false)
      };

      component.ngOnChanges(changes);

      expect(component.initAddMode).toHaveBeenCalled();
    });

    it("should not call initAddMode when triggerAddMode is false", () => {
      jest.spyOn(component, "initAddMode");
      const changes = {
        triggerAddMode: new SimpleChange(true, false, false)
      };

      component.ngOnChanges(changes);

      expect(component.initAddMode).not.toHaveBeenCalled();
    });

    it("should not call initAddMode when triggerAddMode is not in changes", () => {
      jest.spyOn(component, "initAddMode");
      const changes = {
        someOtherProperty: new SimpleChange("old", "new", false)
      };

      component.ngOnChanges(changes);

      expect(component.initAddMode).not.toHaveBeenCalled();
    });

    it("should get receiptItems from form", () => {
      const receiptItems = component.receiptItems;

      expect(receiptItems).toBeInstanceOf(FormArray);
      expect(receiptItems.length).toBe(4);
    });

    it("should handle form without receiptItems", () => {
      fixture.componentRef.setInput('form', new FormGroup({}));

      const receiptItems = component.receiptItems;

      expect(receiptItems).toBeNull();
    });
  });

  describe("User Item Map Management (setUserItemMap)", () => {
    it("should correctly group items by user ID", () => {
      component.setUserItemMap();

      expect(component.userItemMap().size).toBe(3);
      expect(component.userItemMap().get("1")).toEqual([
        {
          item: expect.objectContaining({
            name: mockItems[0].name,
            amount: mockItems[0].amount,
            chargedToUserId: mockItems[0].chargedToUserId,
            status: mockItems[0].status,
            receiptId: mockItems[0].receiptId
          }), arrayIndex: 0
        },
        {
          item: expect.objectContaining({
            name: mockItems[2].name,
            amount: mockItems[2].amount,
            chargedToUserId: mockItems[2].chargedToUserId,
            status: mockItems[2].status,
            receiptId: mockItems[2].receiptId
          }), arrayIndex: 2
        }
      ]);
      expect(component.userItemMap().get("2")).toEqual([
        {
          item: expect.objectContaining({
            name: mockItems[1].name,
            amount: mockItems[1].amount,
            chargedToUserId: mockItems[1].chargedToUserId,
            status: mockItems[1].status,
            receiptId: mockItems[1].receiptId
          }), arrayIndex: 1
        }
      ]);
      expect(component.userItemMap().get("3")).toEqual([
        {
          item: expect.objectContaining({
            name: mockItems[3].name,
            amount: mockItems[3].amount,
            chargedToUserId: mockItems[3].chargedToUserId,
            status: mockItems[3].status,
            receiptId: mockItems[3].receiptId
          }), arrayIndex: 3
        }
      ]);
    });

    it("should handle items without chargedToUserId (null)", () => {
      const itemsWithNullUserId = [
        { id: 1, name: "Item 1", amount: "10.50", chargedToUserId: 1, status: ItemStatus.Open, receiptId: 1 } as Item,
        { id: 2, name: "Item 2", amount: "15.75", chargedToUserId: null, status: ItemStatus.Open, receiptId: 1 } as any as Item,
      ];
      fixture.componentRef.setInput('form', createFormWithItems(itemsWithNullUserId));

      component.setUserItemMap();

      expect(component.userItemMap().size).toBe(1);
      expect(component.userItemMap().get("1")).toEqual([
        {
          item: expect.objectContaining({
            name: "Item 1",
            amount: "10.50",
            chargedToUserId: 1,
            status: ItemStatus.Open,
            receiptId: 1
          }), arrayIndex: 0
        }
      ]);
    });

    it("should handle items without chargedToUserId (undefined)", () => {
      const itemsWithUndefinedUserId = [
        { id: 1, name: "Item 1", amount: "10.50", chargedToUserId: 1, status: ItemStatus.Open, receiptId: 1 } as Item,
        { id: 2, name: "Item 2", amount: "15.75", chargedToUserId: undefined, status: ItemStatus.Open, receiptId: 1 } as Item,
      ];
      fixture.componentRef.setInput('form', createFormWithItems(itemsWithUndefinedUserId));

      component.setUserItemMap();

      expect(component.userItemMap().size).toBe(1);
      expect(component.userItemMap().get("1")).toEqual([
        {
          item: expect.objectContaining({
            name: itemsWithUndefinedUserId[0].name,
            amount: itemsWithUndefinedUserId[0].amount,
            chargedToUserId: itemsWithUndefinedUserId[0].chargedToUserId,
            status: itemsWithUndefinedUserId[0].status,
            receiptId: itemsWithUndefinedUserId[0].receiptId
          }), arrayIndex: 0
        }
      ]);
    });

    it("should handle empty items array", () => {
      fixture.componentRef.setInput('form', createFormWithItems([]));

      component.setUserItemMap();

      expect(component.userItemMap().size).toBe(0);
    });

    it("should convert string user IDs to strings", () => {
      const itemsWithNumberUserId = [
        { id: 1, name: "Item 1", amount: "10.50", chargedToUserId: 123, status: ItemStatus.Open, receiptId: 1 } as Item,
      ];
      fixture.componentRef.setInput('form', createFormWithItems(itemsWithNumberUserId));

      component.setUserItemMap();

      expect(component.userItemMap().has("123")).toBe(true);
      expect(component.userItemMap().get("123")).toEqual([
        {
          item: expect.objectContaining({
            name: itemsWithNumberUserId[0].name,
            amount: itemsWithNumberUserId[0].amount,
            chargedToUserId: itemsWithNumberUserId[0].chargedToUserId,
            status: itemsWithNumberUserId[0].status,
            receiptId: itemsWithNumberUserId[0].receiptId
          }), arrayIndex: 0
        }
      ]);
    });

    it("should handle multiple items for same user", () => {
      const multipleItemsSameUser = [
        { id: 1, name: "Item 1", amount: "10.50", chargedToUserId: 1, status: ItemStatus.Open, receiptId: 1 } as Item,
        { id: 2, name: "Item 2", amount: "15.75", chargedToUserId: 1, status: ItemStatus.Open, receiptId: 1 } as Item,
        { id: 3, name: "Item 3", amount: "8.25", chargedToUserId: 1, status: ItemStatus.Resolved, receiptId: 1 } as Item,
      ];
      fixture.componentRef.setInput('form', createFormWithItems(multipleItemsSameUser));

      component.setUserItemMap();

      expect(component.userItemMap().size).toBe(1);
      expect(component.userItemMap().get("1")?.length).toBe(3);
    });

    it("should handle form without receiptItems control", () => {
      fixture.componentRef.setInput('form', new FormGroup({}));

      // The component doesn't clear the map when no receiptItems control, so we expect it to stay unchanged
      const originalMapSize = component.userItemMap().size;
      component.setUserItemMap();

      expect(component.userItemMap().size).toBe(originalMapSize);
    });

    it("should handle null items value", () => {
      // Create a form where receiptItems value would be null (empty FormArray gets null value)
      fixture.componentRef.setInput('form', new FormGroup({
        receiptItems: new FormArray([])
      }));

      // When form is empty, setUserItemMap should handle gracefully
      const originalMapSize = component.userItemMap().size;
      component.setUserItemMap();

      // Since the form has an empty receiptItems array, it should create an empty map
      expect(component.userItemMap().size).toBe(0);
    });

    it("should handle undefined items value", () => {
      fixture.componentRef.setInput('form', new FormGroup({
        receiptItems: new FormControl(undefined)
      }));

      component.setUserItemMap();

      expect(component.userItemMap().size).toBe(0);
    });

    it("should replace the map each call (no stale entries from prior state)", () => {
      // First populate with default mock items.
      component.setUserItemMap();
      expect(component.userItemMap().size).toBe(3);

      // Swap in a form with just one item, then rebuild — users that no longer
      // appear in the FormArray must not linger in the map.
      fixture.componentRef.setInput('form', createFormWithItems([
        { id: 99, name: "Solo", amount: "1.00", chargedToUserId: 2, status: ItemStatus.Open, receiptId: 1 } as Item,
      ]));
      component.setUserItemMap();

      expect(component.userItemMap().size).toBe(1);
      expect(component.userItemMap().has("1")).toBe(false);
      expect(component.userItemMap().has("2")).toBe(true);
      expect(component.userItemMap().has("3")).toBe(false);
    });
  });

  describe("Linked Items Handling (setUserItemMap)", () => {
    it("should expose linkedItems keyed by chargedToUserId with parent reference", () => {
      const parent = {
        id: 10, name: "Shared Appetizer", amount: "20.00",
        chargedToUserId: null, status: ItemStatus.Open, receiptId: 1,
        linkedItems: [
          { id: 11, name: "Half A", amount: "10.00", chargedToUserId: 2, status: ItemStatus.Open, receiptId: 1 },
          { id: 12, name: "Half B", amount: "10.00", chargedToUserId: 3, status: ItemStatus.Open, receiptId: 1 },
        ],
      } as any as Item;
      fixture.componentRef.setInput('form', createFormWithItems([parent]));

      component.setUserItemMap();

      const user2 = component.userItemMap().get("2");
      expect(user2?.length).toBe(1);
      expect(user2?.[0].isLinkedItem).toBe(true);
      expect(user2?.[0].linkedItemIndex).toBe(0);
      expect(user2?.[0].arrayIndex).toBe(0);
      expect(user2?.[0].parentItem).toBeDefined();
      expect(user2?.[0].parentItem?.name).toBe("Shared Appetizer");
      expect(user2?.[0].item.name).toBe("Half A");

      const user3 = component.userItemMap().get("3");
      expect(user3?.[0].linkedItemIndex).toBe(1);
      expect(user3?.[0].item.name).toBe("Half B");
    });

    it("should include both parent (if chargedToUserId is set) and its linkedItems for same user", () => {
      const parent = {
        id: 10, name: "Parent", amount: "20.00",
        chargedToUserId: 1, status: ItemStatus.Open, receiptId: 1,
        linkedItems: [
          { id: 11, name: "Linked Child", amount: "5.00", chargedToUserId: 1, status: ItemStatus.Open, receiptId: 1 },
        ],
      } as any as Item;
      fixture.componentRef.setInput('form', createFormWithItems([parent]));

      component.setUserItemMap();

      const user1 = component.userItemMap().get("1");
      expect(user1?.length).toBe(2);
      // Parent comes first, linked child second.
      expect(user1?.[0].isLinkedItem).toBeUndefined();
      expect(user1?.[0].item.name).toBe("Parent");
      expect(user1?.[1].isLinkedItem).toBe(true);
      expect(user1?.[1].item.name).toBe("Linked Child");
    });

    it("should skip linkedItems without a chargedToUserId", () => {
      const parent = {
        id: 10, name: "Parent", amount: "20.00",
        chargedToUserId: 1, status: ItemStatus.Open, receiptId: 1,
        linkedItems: [
          { id: 11, name: "Orphan", amount: "5.00", chargedToUserId: null, status: ItemStatus.Open, receiptId: 1 },
          { id: 12, name: "Has Owner", amount: "5.00", chargedToUserId: 2, status: ItemStatus.Open, receiptId: 1 },
        ],
      } as any as Item;
      fixture.componentRef.setInput('form', createFormWithItems([parent]));

      component.setUserItemMap();

      // Orphan linked item should not appear anywhere.
      expect(component.userItemMap().get("1")?.length).toBe(1); // just the parent
      expect(component.userItemMap().get("2")?.length).toBe(1); // only "Has Owner"
    });

    it("should handle items with an empty linkedItems array", () => {
      const item = {
        id: 1, name: "No Splits", amount: "5.00",
        chargedToUserId: 1, status: ItemStatus.Open, receiptId: 1,
        linkedItems: [],
      } as any as Item;
      fixture.componentRef.setInput('form', createFormWithItems([item]));

      expect(() => component.setUserItemMap()).not.toThrow();
      expect(component.userItemMap().get("1")?.length).toBe(1);
    });

    it("should preserve parent arrayIndex on linkedItems across multiple items", () => {
      const items = [
        { id: 1, name: "First", amount: "5.00", chargedToUserId: 1, status: ItemStatus.Open, receiptId: 1 } as Item,
        {
          id: 2, name: "Split Parent", amount: "10.00", chargedToUserId: null,
          status: ItemStatus.Open, receiptId: 1,
          linkedItems: [
            { id: 3, name: "L", amount: "5.00", chargedToUserId: 2, status: ItemStatus.Open, receiptId: 1 },
          ],
        } as any as Item,
      ];
      fixture.componentRef.setInput('form', createFormWithItems(items));

      component.setUserItemMap();

      const user2 = component.userItemMap().get("2");
      // arrayIndex points to the parent (index 1 in the FormArray), not to any
      // flattened position — that's how getFormControlPath reconstructs paths.
      expect(user2?.[0].arrayIndex).toBe(1);
      expect(user2?.[0].linkedItemIndex).toBe(0);
    });
  });

  describe("Add Mode Functionality", () => {
    it("should create form with correct structure in initAddMode", () => {
      component.initAddMode();

      expect(component.isAdding).toBe(true);
      expect(component.newItemFormGroup).toBeDefined();
      expect(component.newItemFormGroup.get("name")).toBeDefined();
      expect(component.newItemFormGroup.get("chargedToUserId")).toBeDefined();
      expect(component.newItemFormGroup.get("amount")).toBeDefined();
      expect(component.newItemFormGroup.get("status")).toBeDefined();
      expect(component.newItemFormGroup.get("receiptId")?.value).toBe(1);
    });

    it("should handle undefined originalReceipt in initAddMode", () => {
      component.originalReceipt = undefined;

      component.initAddMode();

      expect(component.isAdding).toBe(true);
      expect(component.newItemFormGroup.get("receiptId")?.value).toBeNaN();
    });

    it("should reset state in exitAddMode", () => {
      component.isAdding = true;
      component.newItemFormGroup = new FormGroup({
        name: new FormControl("test"),
        amount: new FormControl(10)
      });

      component.exitAddMode();

      expect(component.isAdding).toBe(false);
      expect(Object.keys(component.newItemFormGroup.controls).length).toBe(0);
    });

    it("should submit valid form in submitNewItemFormGroup", () => {
      jest.spyOn(component.itemAdded, "emit");
      jest.spyOn(component, "exitAddMode");

      component.initAddMode();
      component.newItemFormGroup.patchValue({
        name: "New Item",
        chargedToUserId: 1,
        amount: "25.50",
        status: ItemStatus.Open
      });

      component.submitNewItemFormGroup();

      expect(component.itemAdded.emit).toHaveBeenCalledWith(expect.objectContaining({
        name: "New Item",
        chargedToUserId: 1,
        amount: "25.50",
        status: ItemStatus.Open
      }));
      expect(component.exitAddMode).toHaveBeenCalled();
    });

    it("should not submit invalid form in submitNewItemFormGroup", () => {
      jest.spyOn(component.itemAdded, "emit");
      jest.spyOn(component, "exitAddMode");

      component.initAddMode();
      component.newItemFormGroup.patchValue({
        name: "",
        chargedToUserId: null,
        amount: null
      });
      component.newItemFormGroup.get("name")?.setErrors({ required: true });

      component.submitNewItemFormGroup();

      expect(component.itemAdded.emit).not.toHaveBeenCalled();
      expect(component.exitAddMode).not.toHaveBeenCalled();
    });
  });

  describe("Item Management", () => {
    it("should emit correct event in removeItem", () => {
      jest.spyOn(component.itemRemoved, "emit");
      const itemData = { item: mockItems[0], arrayIndex: 0 };

      component.removeItem(itemData);

      expect(component.itemRemoved.emit).toHaveBeenCalledWith({
        item: mockItems[0],
        arrayIndex: 0,
        isLinkedItem: undefined,
        linkedItemIndex: undefined
      });
    });

    it("should forward linkedItem metadata through removeItem", () => {
      jest.spyOn(component.itemRemoved, "emit");
      const linkedItemData = {
        item: { id: 99, name: "Linked", amount: "1.00", chargedToUserId: 2 } as Item,
        arrayIndex: 3,
        parentItem: { id: 10, name: "Parent" } as Item,
        isLinkedItem: true,
        linkedItemIndex: 1,
      };

      component.removeItem(linkedItemData);

      expect(component.itemRemoved.emit).toHaveBeenCalledWith({
        item: linkedItemData.item,
        arrayIndex: 3,
        isLinkedItem: true,
        linkedItemIndex: 1,
      });
    });

    it("should add inline item in edit mode", () => {
      jest.spyOn(component.itemAdded, "emit");
      component.mode = FormMode.edit;

      component.addInlineItem("2");

      expect(component.itemAdded.emit).toHaveBeenCalledWith(expect.objectContaining({
        name: "",
        chargedToUserId: 2
      }));
    });

    it("should add inline item in add mode", () => {
      jest.spyOn(component.itemAdded, "emit");
      component.mode = FormMode.add;

      component.addInlineItem("3");

      expect(component.itemAdded.emit).toHaveBeenCalledWith(expect.objectContaining({
        name: "",
        chargedToUserId: 3
      }));
    });

    it("should not add inline item in view mode", () => {
      jest.spyOn(component.itemAdded, "emit");
      component.mode = FormMode.view;

      component.addInlineItem("1");

      expect(component.itemAdded.emit).not.toHaveBeenCalled();
    });

    it("should stop event propagation in addInlineItem", () => {
      jest.spyOn(component.itemAdded, "emit");
      component.mode = FormMode.edit;
      const mockEvent = { stopImmediatePropagation: jest.fn() } as any;

      component.addInlineItem("1", mockEvent);

      expect(mockEvent.stopImmediatePropagation).toHaveBeenCalled();
      expect(component.itemAdded.emit).toHaveBeenCalled();
    });

    it("should handle undefined event in addInlineItem", () => {
      jest.spyOn(component.itemAdded, "emit");
      component.mode = FormMode.edit;

      expect(() => component.addInlineItem("1", undefined)).not.toThrow();
      expect(component.itemAdded.emit).toHaveBeenCalled();
    });

    // Simulates the parent ReceiptFormComponent.onItemAdded behavior so the
    // component under test can observe the resulting FormArray state.
    function wireParentOnItemAdded(): jest.Mock {
      const emitSpy = jest.fn();
      component.itemAdded.subscribe((item: Item) => {
        emitSpy(item);
        const newFormGroup = buildItemForm(item, mockReceipt.id?.toString(), true, false);
        component.receiptItems.push(newFormGroup);
        component.setUserItemMap();
      });
      return emitSpy;
    }

    it("should not add item on blur when editing an already-populated existing share (status change)", () => {
      // Regression: changing the status of the last existing share used to
      // spawn a blank placeholder because addInlineItemOnBlur fired for any
      // valid last item, not just pending inline placeholders.
      const emitSpy = wireParentOnItemAdded();
      component.mode = FormMode.edit;
      component.setUserItemMap();

      const userItems = component.userItemMap().get("1");
      const lastIndex = userItems!.length - 1;
      const lastItem = userItems![lastIndex];
      const formGroup = component.receiptItems.at(lastItem.arrayIndex);
      formGroup.patchValue({ status: ItemStatus.Resolved });

      component.addInlineItemOnBlur("1", lastIndex);

      expect(emitSpy).not.toHaveBeenCalled();
    });

    it("should not add item on blur when editing the name of an existing last share", () => {
      const emitSpy = wireParentOnItemAdded();
      component.mode = FormMode.edit;
      component.setUserItemMap();

      const userItems = component.userItemMap().get("1");
      const lastIndex = userItems!.length - 1;
      const formGroup = component.receiptItems.at(userItems![lastIndex].arrayIndex);
      formGroup.patchValue({ name: "Renamed Item" });

      component.addInlineItemOnBlur("1", lastIndex);

      expect(emitSpy).not.toHaveBeenCalled();
    });

    it("should not add item on blur when editing the amount of an existing last share", () => {
      const emitSpy = wireParentOnItemAdded();
      component.mode = FormMode.edit;
      component.setUserItemMap();

      const userItems = component.userItemMap().get("1");
      const lastIndex = userItems!.length - 1;
      const formGroup = component.receiptItems.at(userItems![lastIndex].arrayIndex);
      formGroup.patchValue({ amount: "7.25" });

      component.addInlineItemOnBlur("1", lastIndex);

      expect(emitSpy).not.toHaveBeenCalled();
    });

    it("should chain-add a placeholder after an inline add is filled in and blurred", () => {
      const emitSpy = wireParentOnItemAdded();
      component.mode = FormMode.edit;
      // Raise the receipt total so the added share fits within itemTotalValidator.
      component.form().get("amount")?.setValue("100.00");
      component.setUserItemMap();

      // User clicks the inline Add Share (+) for user "2".
      component.addInlineItem("2");
      expect(emitSpy).toHaveBeenCalledTimes(1);

      // User fills in the placeholder to make it valid.
      const userItems = component.userItemMap().get("2");
      const lastIndex = userItems!.length - 1;
      const lastItem = userItems![lastIndex];
      const formGroup = component.receiptItems.at(lastItem.arrayIndex);
      formGroup.patchValue({
        name: "Filled Share",
        amount: "5.00",
        chargedToUserId: 2,
        status: ItemStatus.Open,
      });
      formGroup.updateValueAndValidity();

      component.addInlineItemOnBlur("2", lastIndex);

      // The fill triggers the next cascading placeholder.
      expect(emitSpy).toHaveBeenCalledTimes(2);
    });

    it("should chain through multiple back-to-back inline adds for the same user", () => {
      const emitSpy = wireParentOnItemAdded();
      component.mode = FormMode.edit;
      component.form().get("amount")?.setValue("200.00");
      component.setUserItemMap();

      // Round 1: click (+), fill, blur → cascade adds placeholder #2.
      component.addInlineItem("2");
      let userItems = component.userItemMap().get("2");
      let currentIndex = userItems!.length - 1;
      let placeholder = component.receiptItems.at(userItems![currentIndex].arrayIndex);
      placeholder.patchValue({ name: "A", amount: "5.00", chargedToUserId: 2, status: ItemStatus.Open });
      placeholder.updateValueAndValidity();
      component.addInlineItemOnBlur("2", currentIndex);
      expect(emitSpy).toHaveBeenCalledTimes(2);

      // Round 2: fill that cascaded placeholder, blur → cascade adds placeholder #3.
      userItems = component.userItemMap().get("2");
      currentIndex = userItems!.length - 1;
      placeholder = component.receiptItems.at(userItems![currentIndex].arrayIndex);
      placeholder.patchValue({ name: "B", amount: "5.00", chargedToUserId: 2, status: ItemStatus.Open });
      placeholder.updateValueAndValidity();
      component.addInlineItemOnBlur("2", currentIndex);
      expect(emitSpy).toHaveBeenCalledTimes(3);
    });

    it("should keep pending placeholder state scoped per user", () => {
      // Adding an inline placeholder for user "2" must not cause blurs from
      // user "1"'s existing last share to cascade.
      const emitSpy = wireParentOnItemAdded();
      component.mode = FormMode.edit;
      component.form().get("amount")?.setValue("100.00");
      component.setUserItemMap();

      component.addInlineItem("2");
      expect(emitSpy).toHaveBeenCalledTimes(1);

      // Blur an existing share of user "1" — that share is NOT in the pending
      // set, so nothing further should emit.
      const user1Items = component.userItemMap().get("1");
      const lastUser1 = user1Items!.length - 1;
      component.addInlineItemOnBlur("1", lastUser1);

      expect(emitSpy).toHaveBeenCalledTimes(1);
    });

    it("should only chain-add once per placeholder (second blur is a no-op)", () => {
      const emitSpy = wireParentOnItemAdded();
      component.mode = FormMode.edit;
      component.form().get("amount")?.setValue("100.00");
      component.setUserItemMap();

      component.addInlineItem("2");
      const userItems = component.userItemMap().get("2");
      const placeholderIndex = userItems!.length - 1;
      const placeholder = component.receiptItems.at(userItems![placeholderIndex].arrayIndex);
      placeholder.patchValue({
        name: "Filled Share",
        amount: "5.00",
        chargedToUserId: 2,
        status: ItemStatus.Open,
      });
      placeholder.updateValueAndValidity();

      component.addInlineItemOnBlur("2", placeholderIndex);
      expect(emitSpy).toHaveBeenCalledTimes(2); // inline add + cascade

      // A second blur on the (now-filled, no-longer-pending) same control
      // must not re-trigger the cascade.
      component.addInlineItemOnBlur("2", placeholderIndex);
      expect(emitSpy).toHaveBeenCalledTimes(2);
    });

    it("should not add item on blur for non-last item", () => {
      jest.spyOn(component, "addInlineItem");
      component.mode = FormMode.edit;
      component.setUserItemMap();

      component.addInlineItemOnBlur("1", 0);

      expect(component.addInlineItem).not.toHaveBeenCalled();
    });

    it("should not chain-add on blur when the inline placeholder is still invalid", () => {
      const emitSpy = wireParentOnItemAdded();
      component.mode = FormMode.edit;
      component.setUserItemMap();

      component.addInlineItem("2");
      expect(emitSpy).toHaveBeenCalledTimes(1);

      const userItems = component.userItemMap().get("2");
      const placeholderIndex = userItems!.length - 1;
      const placeholder = component.receiptItems.at(userItems![placeholderIndex].arrayIndex);
      // Leave amount blank → placeholder stays invalid.
      placeholder.patchValue({ name: "Partial Fill" });

      component.addInlineItemOnBlur("2", placeholderIndex);

      // Still only the original inline add, no cascade.
      expect(emitSpy).toHaveBeenCalledTimes(1);
    });

    it("should handle undefined user items in addInlineItemOnBlur", () => {
      jest.spyOn(component, "addInlineItem");
      component.userItemMap.set(new Map());

      component.addInlineItemOnBlur("999", 0);

      expect(component.addInlineItem).not.toHaveBeenCalled();
    });

    it("should remove empty pristine items in checkLastInlineItem", () => {
      jest.spyOn(component.itemRemoved, "emit");
      component.mode = FormMode.edit;

      const multipleItemsUser = [
        { id: 1, name: "Item 1", amount: "10.50", chargedToUserId: 1, status: ItemStatus.Open, receiptId: 1 } as Item,
        { id: 2, name: "", amount: "0", chargedToUserId: 1, status: ItemStatus.Open, receiptId: 1 } as Item,
      ];
      fixture.componentRef.setInput('form', createFormWithItems(multipleItemsUser));
      component.setUserItemMap();

      const lastItemData = component.userItemMap().get("1")![1];
      const lastFormGroup = component.receiptItems.at(lastItemData.arrayIndex);
      lastFormGroup.markAsPristine();

      // Ensure the form has the expected values
      lastFormGroup.patchValue({ name: "", amount: 0 });
      lastFormGroup.markAsPristine();

      component.checkLastInlineItem("1");

      expect(component.itemRemoved.emit).toHaveBeenCalledWith({
        item: expect.objectContaining({
          name: "",
          amount: "0",
          chargedToUserId: 1,
          status: ItemStatus.Open
        }),
        arrayIndex: 1,
        isLinkedItem: undefined,
        linkedItemIndex: undefined
      });
    });

    it("should remove empty pristine items with whitespace name", () => {
      jest.spyOn(component.itemRemoved, "emit");
      component.mode = FormMode.edit;

      const multipleItemsUser = [
        { id: 1, name: "Item 1", amount: "10.50", chargedToUserId: 1, status: ItemStatus.Open, receiptId: 1 } as Item,
        { id: 2, name: "   ", amount: "0", chargedToUserId: 1, status: ItemStatus.Open, receiptId: 1 } as Item,
      ];
      fixture.componentRef.setInput('form', createFormWithItems(multipleItemsUser));
      component.setUserItemMap();

      const lastItemData = component.userItemMap().get("1")![1];
      const lastFormGroup = component.receiptItems.at(lastItemData.arrayIndex);
      lastFormGroup.patchValue({ name: "   ", amount: 0 });
      lastFormGroup.markAsPristine();

      component.checkLastInlineItem("1");

      expect(component.itemRemoved.emit).toHaveBeenCalledWith({
        item: expect.objectContaining({
          name: "   ",
          amount: "0",
          chargedToUserId: 1,
          status: ItemStatus.Open
        }),
        arrayIndex: 1,
        isLinkedItem: undefined,
        linkedItemIndex: undefined
      });
    });

    it("should not remove items that are not pristine", () => {
      jest.spyOn(component.itemRemoved, "emit");
      component.mode = FormMode.edit;

      const multipleItemsUser = [
        { id: 1, name: "Item 1", amount: "10.50", chargedToUserId: 1, status: ItemStatus.Open, receiptId: 1 } as Item,
        { id: 2, name: "", amount: "0", chargedToUserId: 1, status: ItemStatus.Open, receiptId: 1 } as Item,
      ];
      fixture.componentRef.setInput('form', createFormWithItems(multipleItemsUser));
      component.setUserItemMap();

      const lastFormGroup = component.receiptItems.at(1);
      lastFormGroup.markAsDirty();

      component.checkLastInlineItem("1");

      expect(component.itemRemoved.emit).not.toHaveBeenCalled();
    });

    it("should not remove items with valid name", () => {
      jest.spyOn(component.itemRemoved, "emit");
      component.mode = FormMode.edit;

      const multipleItemsUser = [
        { id: 1, name: "Item 1", amount: "10.50", chargedToUserId: 1, status: ItemStatus.Open, receiptId: 1 } as Item,
        { id: 2, name: "Valid Item", amount: "0", chargedToUserId: 1, status: ItemStatus.Open, receiptId: 1 } as Item,
      ];
      fixture.componentRef.setInput('form', createFormWithItems(multipleItemsUser));
      component.setUserItemMap();

      const lastFormGroup = component.receiptItems.at(1);
      lastFormGroup.markAsPristine();

      component.checkLastInlineItem("1");

      expect(component.itemRemoved.emit).not.toHaveBeenCalled();
    });

    it("should not remove items with valid amount", () => {
      jest.spyOn(component.itemRemoved, "emit");
      component.mode = FormMode.edit;

      const multipleItemsUser = [
        { id: 1, name: "Item 1", amount: "10.50", chargedToUserId: 1, status: ItemStatus.Open, receiptId: 1 } as Item,
        { id: 2, name: "", amount: "5.00", chargedToUserId: 1, status: ItemStatus.Open, receiptId: 1 } as Item,
      ];
      fixture.componentRef.setInput('form', createFormWithItems(multipleItemsUser));
      component.setUserItemMap();

      const lastFormGroup = component.receiptItems.at(1);
      lastFormGroup.markAsPristine();

      component.checkLastInlineItem("1");

      expect(component.itemRemoved.emit).not.toHaveBeenCalled();
    });

    it("should not remove when user has only one item", () => {
      jest.spyOn(component.itemRemoved, "emit");
      component.mode = FormMode.edit;

      const singleItemUser = [
        { id: 1, name: "", amount: "0", chargedToUserId: 1, status: ItemStatus.Open, receiptId: 1 } as Item,
      ];
      fixture.componentRef.setInput('form', createFormWithItems(singleItemUser));
      component.setUserItemMap();

      component.checkLastInlineItem("1");

      expect(component.itemRemoved.emit).not.toHaveBeenCalled();
    });

    it("should do nothing in view mode for checkLastInlineItem", () => {
      jest.spyOn(component.itemRemoved, "emit");
      component.mode = FormMode.view;

      component.checkLastInlineItem("1");

      expect(component.itemRemoved.emit).not.toHaveBeenCalled();
    });

    it("should handle undefined user items in checkLastInlineItem", () => {
      jest.spyOn(component.itemRemoved, "emit");
      component.mode = FormMode.edit;
      component.userItemMap.set(new Map());

      component.checkLastInlineItem("999");

      expect(component.itemRemoved.emit).not.toHaveBeenCalled();
    });
  });

  describe("Resolution Features", () => {
    it("should update all items to resolved in resolveAllItemsClicked", () => {
      jest.spyOn(component.allItemsResolved, "emit");
      const mockEvent = { stopImmediatePropagation: jest.fn() } as any;

      component.resolveAllItemsClicked(mockEvent, "1");

      expect(mockEvent.stopImmediatePropagation).toHaveBeenCalled();
      expect(component.allItemsResolved.emit).toHaveBeenCalledWith("1");

      const userItems = component.receiptItems.controls.filter(
        control => control.get("chargedToUserId")?.value?.toString() === "1"
      );
      expect(userItems.length).toBe(2);
      userItems.forEach(item => {
        expect(item.get("status")?.value).toBe(ItemStatus.Resolved);
      });
    });

    it("should return true when all items resolved in allUserItemsResolved", () => {
      const resolvedItems = [
        { id: 1, name: "Item 1", amount: "10.50", chargedToUserId: 1, status: ItemStatus.Resolved, receiptId: 1 } as Item,
        { id: 2, name: "Item 2", amount: "8.25", chargedToUserId: 1, status: ItemStatus.Resolved, receiptId: 1 } as Item,
      ];
      fixture.componentRef.setInput('form', createFormWithItems(resolvedItems));

      const result = component.allUserItemsResolved("1");

      expect(result).toBe(true);
    });

    it("should return false when some items open in allUserItemsResolved", () => {
      const mixedItems = [
        { id: 1, name: "Item 1", amount: "10.50", chargedToUserId: 1, status: ItemStatus.Resolved, receiptId: 1 } as Item,
        { id: 2, name: "Item 2", amount: "8.25", chargedToUserId: 1, status: ItemStatus.Open, receiptId: 1 } as Item,
      ];
      fixture.componentRef.setInput('form', createFormWithItems(mixedItems));

      const result = component.allUserItemsResolved("1");

      expect(result).toBe(false);
    });

    it("should return true for user with no items", () => {
      const result = component.allUserItemsResolved("999");

      expect(result).toBe(true);
    });

    it("should filter correctly in getItemsForUser", () => {
      const userItems = (component as any).getItemsForUser("1");

      expect(userItems.length).toBe(2);
      expect(userItems[0].get("chargedToUserId")?.value).toBe(1);
      expect(userItems[1].get("chargedToUserId")?.value).toBe(1);
    });

    it("should return empty array for non-existent user in getItemsForUser", () => {
      const userItems = (component as any).getItemsForUser("999");

      expect(userItems.length).toBe(0);
    });

    it("should handle string comparison in getItemsForUser", () => {
      const itemsWithStringUserId = [
        { id: 1, name: "Item 1", amount: "10.50", chargedToUserId: "123", status: ItemStatus.Open, receiptId: 1 } as any,
      ];
      fixture.componentRef.setInput('form', createFormWithItems(itemsWithStringUserId));

      const userItems = (component as any).getItemsForUser("123");

      expect(userItems.length).toBe(1);
    });

    it("should convert a mix of Open and Resolved items to all Resolved", () => {
      const mixedItems = [
        { id: 1, name: "I1", amount: "5.00", chargedToUserId: 1, status: ItemStatus.Open, receiptId: 1 } as Item,
        { id: 2, name: "I2", amount: "5.00", chargedToUserId: 1, status: ItemStatus.Resolved, receiptId: 1 } as Item,
        { id: 3, name: "I3", amount: "5.00", chargedToUserId: 2, status: ItemStatus.Open, receiptId: 1 } as Item,
      ];
      fixture.componentRef.setInput('form', createFormWithItems(mixedItems));
      jest.spyOn(component.allItemsResolved, "emit");
      const mockEvent = { stopImmediatePropagation: jest.fn() } as any;

      component.resolveAllItemsClicked(mockEvent, "1");

      // User 1's items are both Resolved; user 2's Open status is untouched.
      expect(component.receiptItems.at(0).get("status")?.value).toBe(ItemStatus.Resolved);
      expect(component.receiptItems.at(1).get("status")?.value).toBe(ItemStatus.Resolved);
      expect(component.receiptItems.at(2).get("status")?.value).toBe(ItemStatus.Open);
      expect(component.allItemsResolved.emit).toHaveBeenCalledWith("1");
    });

    it("should still emit allItemsResolved for a user with no matching items", () => {
      jest.spyOn(component.allItemsResolved, "emit");
      const mockEvent = { stopImmediatePropagation: jest.fn() } as any;

      expect(() => component.resolveAllItemsClicked(mockEvent, "999")).not.toThrow();
      expect(component.allItemsResolved.emit).toHaveBeenCalledWith("999");
    });
  });

  describe("getFormControlPath", () => {
    it("should build a path for a regular share", () => {
      const itemData = { item: mockItems[0], arrayIndex: 2 };

      expect(component.getFormControlPath(itemData, "status"))
        .toBe("receiptItems.2.status");
    });

    it("should build a path into linkedItems for a linked share", () => {
      const itemData = {
        item: { id: 11, name: "Linked" } as Item,
        arrayIndex: 4,
        isLinkedItem: true,
        linkedItemIndex: 1,
      };

      expect(component.getFormControlPath(itemData, "amount"))
        .toBe("receiptItems.4.linkedItems.1.amount");
    });

    it("should return empty string for missing itemData", () => {
      expect(component.getFormControlPath(undefined as any, "status")).toBe("");
    });

    it("should return empty string for missing fieldName", () => {
      const itemData = { item: mockItems[0], arrayIndex: 0 };
      expect(component.getFormControlPath(itemData, "")).toBe("");
    });

    it("should fall back to the regular path when isLinkedItem is true but linkedItemIndex is undefined", () => {
      // Defensive branch: the component treats partially-populated linked metadata
      // as a regular share rather than producing an invalid "undefined" path.
      const itemData = {
        item: mockItems[0],
        arrayIndex: 3,
        isLinkedItem: true,
        linkedItemIndex: undefined,
      };

      expect(component.getFormControlPath(itemData, "name"))
        .toBe("receiptItems.3.name");
    });

    it("should produce correct paths for every share field in the form", () => {
      const itemData = { item: mockItems[0], arrayIndex: 0 };

      for (const field of ["name", "amount", "status", "categories", "tags", "chargedToUserId"]) {
        expect(component.getFormControlPath(itemData, field))
          .toBe(`receiptItems.0.${field}`);
      }
    });
  });

  describe("Edge Cases", () => {
    it("should handle missing route data", () => {
      const activatedRoute = TestBed.inject(ActivatedRoute);
      const originalData = activatedRoute.snapshot.data;
      (activatedRoute.snapshot as any).data = null;

      // Component tries to access data["receipt"] which throws when data is null
      expect(() => component.ngOnInit()).toThrow();

      // Restore route data for subsequent tests
      activatedRoute.snapshot.data = originalData;
    });

    it("should handle form without receiptItems", () => {
      fixture.componentRef.setInput('form', new FormGroup({
        otherField: new FormControl("value")
      }));

      // The component doesn't clear the map when no receiptItems, so we expect it to stay unchanged
      const originalMapSize = component.userItemMap().size;
      component.setUserItemMap();

      expect(component.userItemMap().size).toBe(originalMapSize);
    });

    it("should handle invalid user IDs", () => {
      const itemsWithInvalidUserId = [
        { id: 1, name: "Item 1", amount: "10.50", chargedToUserId: 0, status: ItemStatus.Open, receiptId: 1 } as Item,
        { id: 2, name: "Item 2", amount: "15.75", chargedToUserId: -1, status: ItemStatus.Open, receiptId: 1 } as Item,
      ];
      fixture.componentRef.setInput('form', createFormWithItems(itemsWithInvalidUserId));

      component.setUserItemMap();

      expect(component.userItemMap().has("0")).toBe(true);
      expect(component.userItemMap().has("-1")).toBe(true);
    });

    it("should handle empty user maps", () => {
      component.userItemMap.set(new Map());
      // Also clear the form controls to match the empty userMap
      fixture.componentRef.setInput('form', new FormGroup({
        receiptItems: new FormArray([]),
        amount: new FormControl("0")
      }));

      expect(() => component.checkLastInlineItem("1")).not.toThrow();
      expect(() => component.addInlineItemOnBlur("1", 0)).not.toThrow();
      expect(component.allUserItemsResolved("1")).toBe(true);
    });

    it("should handle form mode transitions", () => {
      component.mode = FormMode.view;
      jest.spyOn(component.itemAdded, "emit");

      component.addInlineItem("1");
      expect(component.itemAdded.emit).not.toHaveBeenCalled();

      component.mode = FormMode.edit;
      component.addInlineItem("1");
      expect(component.itemAdded.emit).toHaveBeenCalled();

      component.mode = FormMode.add;
      component.addInlineItem("1");
      expect(component.itemAdded.emit).toHaveBeenCalledTimes(2);
    });

    it("should handle receipt without ID", () => {
      component.originalReceipt = { name: "Receipt without ID" } as Receipt;

      component.initAddMode();

      expect(component.newItemFormGroup.get("receiptId")?.value).toBeNaN();
    });

    it("should handle concurrent map updates", () => {
      component.setUserItemMap();
      const originalSize = component.userItemMap().size;

      // Add new item to form
      const newItem = buildItemForm(
        { id: 5, name: "New Item", amount: "20.00", chargedToUserId: 4, status: ItemStatus.Open, receiptId: 1 } as Item,
        "1",
        true,
        false
      );
      (component.receiptItems as FormArray).push(newItem);

      // Update map again
      component.setUserItemMap();

      expect(component.userItemMap().size).toBe(originalSize + 1);
      expect(component.userItemMap().has("4")).toBe(true);
    });
  });
});
