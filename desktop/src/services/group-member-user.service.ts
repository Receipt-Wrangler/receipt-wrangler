import { Injectable } from "@angular/core";
import { Store } from "@ngxs/store";
import { GroupMember, User } from "../open-api";
import { GroupState, UserState } from "../store";

@Injectable({
  providedIn: "root",
})
export class GroupMemberUserService {
  constructor(private store: Store) {}

  public getUsersInGroup(groupId: string): User[] {
    const users: User[] = [];
    const groupMembers: GroupMember[] = this.store
      .selectSnapshot(GroupState.allGroupMembers)
      .filter((gm) => gm.groupId.toString() === groupId.toString());

    groupMembers.forEach((gm) => {
      const user = this.store.selectSnapshot(
        UserState.findUserById(gm.userId.toString())
      );
      if (user) {
        users.push(user);
      }
    });
    return users;
  }
}
