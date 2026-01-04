import { Pipe, PipeTransform } from "@angular/core";
import { Store } from "@ngxs/store";
import { User } from "../open-api";
import { UserState } from "../store";

@Pipe({
    name: "user",
    standalone: false
})
export class UserPipe implements PipeTransform {
  constructor(private store: Store) {}

  // TODO: fix delete receipt busted
  // TODO: implement e2e
  public transform(userId?: string): User | undefined {
    return this.store.selectSnapshot(UserState.getUserById(userId ?? ""));
  }
}
