import { FormControl, FormGroup } from "@angular/forms";
import { ItemData } from "./share-list/share-list.component";
import { UserTotalWithPercentagePipe } from "./user-total-with-percentage.pipe";

describe("UserTotalWithPercentagePipe", () => {
  let pipe: UserTotalWithPercentagePipe;

  beforeEach(() => {
    pipe = new UserTotalWithPercentagePipe();
  });

  it("should create an instance", () => {
    expect(pipe).toBeTruthy();
  });

  it("should calculate total and percentage correctly", () => {
    const mockItems: ItemData[] = [
      { item: { amount: 25 } as any, arrayIndex: 0 },
      { item: { amount: 15 } as any, arrayIndex: 1 }
    ];

    const userItemMap = new Map<string, ItemData[]>();
    userItemMap.set("user1", mockItems);

    const form = new FormGroup({
      amount: new FormControl(100)
    });

    const result = pipe.transform(userItemMap, "user1", form);

    expect(result.total).toBe(40);
    expect(result.percentage).toBe(40);
  });

  it("should return 0 for empty items", () => {
    const userItemMap = new Map<string, ItemData[]>();
    userItemMap.set("user1", []);

    const form = new FormGroup({
      amount: new FormControl(100)
    });

    const result = pipe.transform(userItemMap, "user1", form);

    expect(result.total).toBe(0);
    expect(result.percentage).toBe(0);
  });

  it("should handle division by zero", () => {
    const mockItems: ItemData[] = [
      { item: { amount: 25 } as any, arrayIndex: 0 }
    ];

    const userItemMap = new Map<string, ItemData[]>();
    userItemMap.set("user1", mockItems);

    const form = new FormGroup({
      amount: new FormControl(0)
    });

    const result = pipe.transform(userItemMap, "user1", form);

    expect(result.total).toBe(25);
    expect(result.percentage).toBe(0);
  });

  it("should round percentage to 2 decimal places", () => {
    const mockItems: ItemData[] = [
      { item: { amount: 33.333 } as any, arrayIndex: 0 }
    ];

    const userItemMap = new Map<string, ItemData[]>();
    userItemMap.set("user1", mockItems);

    const form = new FormGroup({
      amount: new FormControl(100)
    });

    const result = pipe.transform(userItemMap, "user1", form);

    expect(result.total).toBe(33.333);
    expect(result.percentage).toBe(33.33);
  });
});
