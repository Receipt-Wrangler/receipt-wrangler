import { Component, Input, OnChanges, SimpleChanges } from "@angular/core";
import { UntilDestroy } from "@ngneat/until-destroy";
import { take, tap } from "rxjs";
import { UserService } from "../../open-api";


@UntilDestroy()
@Component({
    selector: "app-summary-card",
    templateUrl: "./summary-card.component.html",
    styleUrls: ["./summary-card.component.scss"],
    standalone: false
})
export class SummaryCardComponent implements OnChanges {
  constructor(private userService: UserService) {}

  @Input() public headerText: string = "";

  @Input() public groupId: string | number = "";

  @Input() public receiptIds: number[] = [];

  public usersOweMap: Map<string, string> = new Map();
  public userOwesMap: Map<string, string> = new Map();

  public ngOnChanges(changes: SimpleChanges): void {
    if (changes["groupId"] || changes["receiptIds"]) {
      this.buildOweMap();
    }
  }

  private buildOweMap(): void {
    if (!this.groupId) {
      return;
    }

    let id: any = Number.parseInt(this.groupId as any) || (this.groupId as any);
    if (this.receiptIds.length > 0) {
      id = undefined;
    }
    
    this.userService
      .getAmountOwedForUser(
        id,
        this.receiptIds
      )
      .pipe(
        take(1),
        tap((result) => {
          this.userOwesMap = new Map();
          this.usersOweMap = new Map();

          Object.keys(result).forEach((k) => {
            const key = k.toString();
            if (Number(result[k]) > 0) {
              this.userOwesMap.set(key, result[k].toString());
            } else {
              const parsed = Number.parseFloat(result[k]);
              this.usersOweMap.set(key, Math.abs(parsed).toString());
            }
          });
        })
      )
      .subscribe();
  }
}
