import { Pipe, PipeTransform } from "@angular/core";
import { FormGroup } from "@angular/forms";
import { ItemData } from "./share-list/share-list.component";

@Pipe({
  name: "userTotalWithPercentage",
  pure: false,
  standalone: false
})
export class UserTotalWithPercentagePipe implements PipeTransform {
  public transform(
    userItemMap: Map<string, ItemData[]>,
    userId: string,
    form: FormGroup
  ): { total: number; percentage: number } {
    const items = userItemMap.get(userId) as ItemData[];
    const userTotal = items?.length > 0 && items
      ? items.map((item) => Number(item.item.amount)).reduce((a, b) => a + b)
      : 0;

    const receiptTotal = Number(form.get("amount")?.value || 0);
    const percentage = receiptTotal > 0 ? (userTotal / receiptTotal) * 100 : 0;

    return {
      total: userTotal,
      percentage: Math.round(percentage * 100) / 100 // Round to 2 decimal places
    };
  }
}
