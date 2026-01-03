import {
  Component,
  EventEmitter,
  Input,
  OnChanges,
  OnInit,
  Output,
  QueryList,
  SimpleChanges,
  ViewChildren,
  ViewEncapsulation,
} from "@angular/core";
import { AbstractControl, FormArray, FormGroup, } from "@angular/forms";
import { MatExpansionPanel } from "@angular/material/expansion";
import { ActivatedRoute } from "@angular/router";
import { Select } from "@ngxs/store";
import { Observable } from "rxjs";
import { RECEIPT_ITEM_STATUS_OPTIONS } from "src/constants/receipt-status-options";
import { FormMode } from "src/enums/form-mode.enum";
import { InputComponent } from "../../input";
import { Category, Group, GroupRole, Item, ItemStatus, Receipt, Tag, User } from "../../open-api";
import { UserState } from "../../store";
import { buildItemForm } from "../utils/form.utils";

export interface ItemData {
  item: Item;
  arrayIndex: number;
  parentItem?: Item;
  isLinkedItem?: boolean;
  linkedItemIndex?: number;
}

@Component({
  selector: "app-share-list",
  templateUrl: "./share-list.component.html",
  styleUrls: ["./share-list.component.scss"],
  encapsulation: ViewEncapsulation.None,
  standalone: false
})
export class ShareListComponent implements OnInit, OnChanges {
  @ViewChildren("userExpansionPanel")
  public userExpansionPanels!: QueryList<MatExpansionPanel>;

  @ViewChildren("nameField")
  public nameFields!: QueryList<InputComponent>;

  @Select(UserState.users) public users!: Observable<User[]>;

  @Input() public form!: FormGroup;

  @Input() public originalReceipt?: Receipt;

  @Input() public categories: Category[] = [];

  @Input() public tags: Tag[] = [];

  @Input() public selectedGroup: Group | undefined;

  @Input() public triggerAddMode: boolean = false;

  @Output() public itemAdded = new EventEmitter<Item>();

  @Output() public itemRemoved = new EventEmitter<{ item: Item; arrayIndex: number; isLinkedItem?: boolean; linkedItemIndex?: number }>();

  @Output() public allItemsResolved = new EventEmitter<string>();

  public newItemFormGroup: FormGroup = new FormGroup({});

  public userItemMap: Map<string, ItemData[]> = new Map<string, ItemData[]>();

  public isAdding: boolean = false;

  public mode: FormMode = FormMode.view;

  public formMode = FormMode;

  public groupRole = GroupRole;

  public itemStatusOptions = RECEIPT_ITEM_STATUS_OPTIONS;

  public get receiptItems(): FormArray {
    return this.form.get("receiptItems") as FormArray;
  }

  constructor(
    private activatedRoute: ActivatedRoute
  ) {}

  public ngOnInit(): void {
    this.originalReceipt = this.activatedRoute.snapshot.data["receipt"];
    this.mode = this.activatedRoute.snapshot.data["mode"];
    this.setUserItemMap();
  }

  public ngOnChanges(changes: SimpleChanges): void {
    if (changes["triggerAddMode"] && changes["triggerAddMode"].currentValue) {
      this.initAddMode();
    }
  }


  public setUserItemMap(): void {
    const receiptItems = this.form.get("receiptItems");
    if (receiptItems) {
      const items = this.form.get("receiptItems")?.value as Item[];
      const map = new Map<string, ItemData[]>();

      if (items?.length > 0) {
        items.forEach((item, index) => {
          // Add regular share items (those with chargedToUserId)
          const chargedToUserId = item?.chargedToUserId?.toString();
          if (chargedToUserId) {
            const itemData: ItemData = {
              item: item,
              arrayIndex: index,
            };

            if (map.has(chargedToUserId)) {
              const newItems = Array.from(map.get(chargedToUserId) as ItemData[]);
              newItems.push(itemData);
              map.set(chargedToUserId, newItems);
            } else {
              map.set(chargedToUserId, [itemData]);
            }
          }

          // Add linkedItems (split items)
          if (item?.linkedItems && item.linkedItems.length > 0) {
            item.linkedItems.forEach((linkedItem, linkedIndex) => {
              const linkedChargedToUserId = linkedItem?.chargedToUserId?.toString();
              if (linkedChargedToUserId) {
                const linkedItemData: ItemData = {
                  item: linkedItem,
                  arrayIndex: index, // Keep reference to parent's index
                  parentItem: item,
                  isLinkedItem: true,
                  linkedItemIndex: linkedIndex, // Track position within linkedItems array
                };

                if (map.has(linkedChargedToUserId)) {
                  const newItems = Array.from(map.get(linkedChargedToUserId) as ItemData[]);
                  newItems.push(linkedItemData);
                  map.set(linkedChargedToUserId, newItems);
                } else {
                  map.set(linkedChargedToUserId, [linkedItemData]);
                }
              }
            });
          }
        });
      }
      this.userItemMap = map;
    }
  }

  public initAddMode(): void {
    this.isAdding = true;
    this.newItemFormGroup = buildItemForm(
      undefined,
      this.originalReceipt?.id?.toString(),
      true,
      false
    );
  }

  public exitAddMode(): void {
    this.isAdding = false;
    this.newItemFormGroup = new FormGroup({});
  }

  public submitNewItemFormGroup(): void {
    if (this.newItemFormGroup.valid) {
      const newItem = this.newItemFormGroup.value as Item;
      this.itemAdded.emit(newItem);
      this.exitAddMode();
    }
  }

  public removeItem(itemData: ItemData): void {
    this.itemRemoved.emit({ 
      item: itemData.item, 
      arrayIndex: itemData.arrayIndex,
      isLinkedItem: itemData.isLinkedItem,
      linkedItemIndex: itemData.linkedItemIndex
    });
  }

  public addInlineItem(userId: string, event?: MouseEvent): void {
    if (event) {
      event?.stopImmediatePropagation();
    }

    if (this.mode !== FormMode.view) {
      const newItem = {
        name: "",
        chargedToUserId: Number(userId),
      } as Item;
      this.itemAdded.emit(newItem);
    }
  }

  public addInlineItemOnBlur(userId: string, index: number): void {
    const userItems = this.userItemMap.get(userId);
    if (userItems && userItems.length - 1 === index) {
      const item = userItems.at(index) as ItemData;
      const itemInput = this.receiptItems.at(item?.arrayIndex);
      if (itemInput.valid) {
        const activeElement = document.activeElement as HTMLElement;
        this.addInlineItem(userId);
      }
    }
  }

  public checkLastInlineItem(userId: string): void {
    if (this.mode !== FormMode.view) {
      const items = this.userItemMap.get(userId);
      if (items && items.length > 1) {
        const lastItem = items[items.length - 1];
        const formGroup = this.receiptItems.at(lastItem.arrayIndex);
        const nameValue = formGroup.get("name")?.value;
        const amountValue = formGroup.get("amount")?.value;

        if (formGroup.pristine && (!nameValue || nameValue.trim() === "") && (!amountValue || amountValue === 0)) {
          this.itemRemoved.emit({ 
            item: lastItem.item, 
            arrayIndex: lastItem.arrayIndex,
            isLinkedItem: lastItem.isLinkedItem,
            linkedItemIndex: lastItem.linkedItemIndex
          });
        }
      }
    }
  }

  public resolveAllItemsClicked(event: MouseEvent, userId: string): void {
    event.stopImmediatePropagation();
    const filtered = this.getItemsForUser(userId);

    filtered.forEach((i) =>
      i.patchValue({
        status: ItemStatus.Resolved,
      })
    );
    this.allItemsResolved.emit(userId);
  }

  public allUserItemsResolved(userId: string): boolean {
    const userItems = this.getItemsForUser(userId);
    return userItems.every(
      (i) => i.get("status")?.value === ItemStatus.Resolved
    );
  }

  private getItemsForUser(userId: string): AbstractControl[] {
    return this.receiptItems.controls.filter(
      (i) => i.get("chargedToUserId")?.value?.toString() === userId
    );
  }

  public getFormControlPath(itemData: ItemData, fieldName: string): string {
    if (!itemData || !fieldName) {
      return '';
    }
    
    if (itemData.isLinkedItem && itemData.linkedItemIndex !== undefined) {
      // For linkedItems: receiptItems.parentIndex.linkedItems.linkedIndex.fieldName
      return `receiptItems.${itemData.arrayIndex}.linkedItems.${itemData.linkedItemIndex}.${fieldName}`;
    } else {
      // For regular items: receiptItems.arrayIndex.fieldName
      return `receiptItems.${itemData.arrayIndex}.${fieldName}`;
    }
  }
}
