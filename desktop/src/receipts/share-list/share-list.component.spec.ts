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
import { UserState } from "../../store/index";
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
      imports: [CurrencyPipe, PipesModule, NgxsModule.forRoot([UserState, SystemSettingsState])],
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
      users: {
        users: mockUsers
      },
      systemSettings: {}
    });

    // Setup default component state
    component.form = createFormWithItems(mockItems);
    component.categories = mockCategories;
    component.tags = mockTags;
    component.selectedGroup = mockGroup;
    component.originalReceipt = mockReceipt;

    // Mock ViewChildren
    component.userExpansionPanels = new QueryList<MatExpansionPanel>();
    component.nameFields = new QueryList<InputComponent>();

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
      const emptyRouteData = {
        snapshot: {
          data: {}
        }
      };

      const newComponent = new ShareListComponent(emptyRouteData as any);
      newComponent.form = createFormWithItems(mockItems);

      newComponent.ngOnInit();

      expect(newComponent.originalReceipt).toBeUndefined();
      expect(newComponent.mode).toBeUndefined();
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
      component.form = new FormGroup({});

      const receiptItems = component.receiptItems;

      expect(receiptItems).toBeNull();
    });
  });

  describe("User Item Map Management (setUserItemMap)", () => {
    it("should correctly group items by user ID", () => {
      component.setUserItemMap();

      expect(component.userItemMap.size).toBe(3);
      expect(component.userItemMap.get("1")).toEqual([
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
      expect(component.userItemMap.get("2")).toEqual([
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
      expect(component.userItemMap.get("3")).toEqual([
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
      component.form = createFormWithItems(itemsWithNullUserId);

      component.setUserItemMap();

      expect(component.userItemMap.size).toBe(1);
      expect(component.userItemMap.get("1")).toEqual([
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
      component.form = createFormWithItems(itemsWithUndefinedUserId);

      component.setUserItemMap();

      expect(component.userItemMap.size).toBe(1);
      expect(component.userItemMap.get("1")).toEqual([
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
      component.form = createFormWithItems([]);

      component.setUserItemMap();

      expect(component.userItemMap.size).toBe(0);
    });

    it("should convert string user IDs to strings", () => {
      const itemsWithNumberUserId = [
        { id: 1, name: "Item 1", amount: "10.50", chargedToUserId: 123, status: ItemStatus.Open, receiptId: 1 } as Item,
      ];
      component.form = createFormWithItems(itemsWithNumberUserId);

      component.setUserItemMap();

      expect(component.userItemMap.has("123")).toBe(true);
      expect(component.userItemMap.get("123")).toEqual([
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
      component.form = createFormWithItems(multipleItemsSameUser);

      component.setUserItemMap();

      expect(component.userItemMap.size).toBe(1);
      expect(component.userItemMap.get("1")?.length).toBe(3);
    });

    it("should handle form without receiptItems control", () => {
      component.form = new FormGroup({});

      // The component doesn't clear the map when no receiptItems control, so we expect it to stay unchanged
      const originalMapSize = component.userItemMap.size;
      component.setUserItemMap();

      expect(component.userItemMap.size).toBe(originalMapSize);
    });

    it("should handle null items value", () => {
      // Create a form where receiptItems value would be null (empty FormArray gets null value)
      component.form = new FormGroup({
        receiptItems: new FormArray([])
      });

      // When form is empty, setUserItemMap should handle gracefully
      const originalMapSize = component.userItemMap.size;
      component.setUserItemMap();

      // Since the form has an empty receiptItems array, it should create an empty map
      expect(component.userItemMap.size).toBe(0);
    });

    it("should handle undefined items value", () => {
      component.form = new FormGroup({
        receiptItems: new FormControl(undefined)
      });

      component.setUserItemMap();

      expect(component.userItemMap.size).toBe(0);
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

    it("should add item on blur for last item with valid form", () => {
      jest.spyOn(component, "addInlineItem");
      component.mode = FormMode.edit;
      component.setUserItemMap();

      const userItems = component.userItemMap.get("1");
      const lastIndex = userItems!.length - 1;
      const lastItem = userItems![lastIndex];
      const formGroup = component.receiptItems.at(lastItem.arrayIndex);
      formGroup.patchValue({ name: "Valid Item", amount: "10" });

      // Ensure the form is valid by clearing any validation errors
      formGroup.get("name")?.setErrors(null);
      formGroup.get("amount")?.setErrors(null);
      formGroup.updateValueAndValidity();

      component.addInlineItemOnBlur("1", lastIndex);

      expect(component.addInlineItem).toHaveBeenCalledWith("1");
    });

    it("should not add item on blur for non-last item", () => {
      jest.spyOn(component, "addInlineItem");
      component.mode = FormMode.edit;
      component.setUserItemMap();

      component.addInlineItemOnBlur("1", 0);

      expect(component.addInlineItem).not.toHaveBeenCalled();
    });

    it("should not add item on blur for invalid form", () => {
      jest.spyOn(component, "addInlineItem");
      component.mode = FormMode.edit;
      component.setUserItemMap();

      const userItems = component.userItemMap.get("1");
      const lastIndex = userItems!.length - 1;
      const lastItem = userItems![lastIndex];
      const formGroup = component.receiptItems.at(lastItem.arrayIndex);
      formGroup.get("name")?.setErrors({ required: true });

      component.addInlineItemOnBlur("1", lastIndex);

      expect(component.addInlineItem).not.toHaveBeenCalled();
    });

    it("should handle undefined user items in addInlineItemOnBlur", () => {
      jest.spyOn(component, "addInlineItem");
      component.userItemMap = new Map();

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
      component.form = createFormWithItems(multipleItemsUser);
      component.setUserItemMap();

      const lastItemData = component.userItemMap.get("1")![1];
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
      component.form = createFormWithItems(multipleItemsUser);
      component.setUserItemMap();

      const lastItemData = component.userItemMap.get("1")![1];
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
      component.form = createFormWithItems(multipleItemsUser);
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
      component.form = createFormWithItems(multipleItemsUser);
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
      component.form = createFormWithItems(multipleItemsUser);
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
      component.form = createFormWithItems(singleItemUser);
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
      component.userItemMap = new Map();

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
      component.form = createFormWithItems(resolvedItems);

      const result = component.allUserItemsResolved("1");

      expect(result).toBe(true);
    });

    it("should return false when some items open in allUserItemsResolved", () => {
      const mixedItems = [
        { id: 1, name: "Item 1", amount: "10.50", chargedToUserId: 1, status: ItemStatus.Resolved, receiptId: 1 } as Item,
        { id: 2, name: "Item 2", amount: "8.25", chargedToUserId: 1, status: ItemStatus.Open, receiptId: 1 } as Item,
      ];
      component.form = createFormWithItems(mixedItems);

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
      component.form = createFormWithItems(itemsWithStringUserId);

      const userItems = (component as any).getItemsForUser("123");

      expect(userItems.length).toBe(1);
    });
  });

  describe("Edge Cases", () => {
    it("should handle missing route data", () => {
      // Create a new component instance directly and inject null data via activatedRoute
      const mockActivatedRoute = {
        snapshot: {
          data: null
        }
      };

      const newComponent = new ShareListComponent(mockActivatedRoute as any);
      newComponent.form = createFormWithItems(mockItems);

      // Component tries to access data["receipt"] which throws when data is null
      expect(() => newComponent.ngOnInit()).toThrow();
    });

    it("should handle form without receiptItems", () => {
      component.form = new FormGroup({
        otherField: new FormControl("value")
      });

      // The component doesn't clear the map when no receiptItems, so we expect it to stay unchanged
      const originalMapSize = component.userItemMap.size;
      component.setUserItemMap();

      expect(component.userItemMap.size).toBe(originalMapSize);
    });

    it("should handle invalid user IDs", () => {
      const itemsWithInvalidUserId = [
        { id: 1, name: "Item 1", amount: "10.50", chargedToUserId: 0, status: ItemStatus.Open, receiptId: 1 } as Item,
        { id: 2, name: "Item 2", amount: "15.75", chargedToUserId: -1, status: ItemStatus.Open, receiptId: 1 } as Item,
      ];
      component.form = createFormWithItems(itemsWithInvalidUserId);

      component.setUserItemMap();

      expect(component.userItemMap.has("0")).toBe(true);
      expect(component.userItemMap.has("-1")).toBe(true);
    });

    it("should handle empty user maps", () => {
      component.userItemMap = new Map();
      // Also clear the form controls to match the empty userMap
      component.form = new FormGroup({
        receiptItems: new FormArray([]),
        amount: new FormControl("0")
      });

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
      const originalSize = component.userItemMap.size;

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

      expect(component.userItemMap.size).toBe(originalSize + 1);
      expect(component.userItemMap.has("4")).toBe(true);
    });
  });
});
