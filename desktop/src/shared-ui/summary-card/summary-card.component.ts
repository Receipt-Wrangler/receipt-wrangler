import { Component, effect, input, signal, untracked } from "@angular/core";
import { take, tap } from "rxjs";
import { UserService } from "../../open-api";
import { CardComponent } from "../card/card.component";
import { NgTemplateOutlet, KeyValuePipe } from "@angular/common";
import { MatList, MatListItem } from "@angular/material/list";
import { CustomCurrencyPipe } from "../../pipes/custom-currency.pipe";
import { UserPipe } from "../../pipes/user.pipe";


@Component({
    selector: "app-summary-card",
    templateUrl: "./summary-card.component.html",
    styleUrls: ["./summary-card.component.scss"],
    imports: [CardComponent, NgTemplateOutlet, MatList, MatListItem, KeyValuePipe, CustomCurrencyPipe, UserPipe]
})
export class SummaryCardComponent {
  constructor(private userService: UserService) {
    effect(() => {
      const groupId = this.groupId();
      const receiptIds = this.receiptIds();
      untracked(() => this.buildOweMap(groupId, receiptIds));
    });
  }

  public readonly headerText = input<string>("");

  public readonly groupId = input<string | number>("");

  public readonly receiptIds = input<number[]>([]);

  public usersOweMap = signal(new Map<string, string>());
  public userOwesMap = signal(new Map<string, string>());

  private buildOweMap(groupId: string | number, receiptIds: number[]): void {
    if (!groupId) {
      return;
    }

    let id: any = Number.parseInt(groupId as any) || (groupId as any);
    if (receiptIds.length > 0) {
      id = undefined;
    }

    this.userService
      .getAmountOwedForUser(
        id,
        receiptIds
      )
      .pipe(
        take(1),
        tap((result) => {
          const newUserOwes = new Map<string, string>();
          const newUsersOwe = new Map<string, string>();

          Object.keys(result).forEach((k) => {
            const key = k.toString();
            if (Number(result[k]) > 0) {
              newUserOwes.set(key, result[k].toString());
            } else {
              const parsed = Number.parseFloat(result[k]);
              newUsersOwe.set(key, Math.abs(parsed).toString());
            }
          });

          this.userOwesMap.set(newUserOwes);
          this.usersOweMap.set(newUsersOwe);
        })
      )
      .subscribe();
  }
}
