import { Component, Input, OnChanges, OnInit, SimpleChanges, ViewChild, } from "@angular/core";
import { FormControl } from "@angular/forms";
import { Store } from "@ngxs/store";
import { AutocomleteComponent } from "src/autocomplete/autocomlete/autocomlete.component";
import { GroupMemberUserService } from "src/services/group-member-user.service";
import { User } from "../../open-api";
import { UserState } from "../../store";

@Component({
    selector: "app-user-autocomplete",
    templateUrl: "./user-autocomplete.component.html",
    styleUrls: ["./user-autocomplete.component.scss"],
    providers: [GroupMemberUserService],
    standalone: false
})
export class UserAutocompleteComponent implements OnInit, OnChanges {
  constructor(
    private store: Store,
    private groupMemberUserService: GroupMemberUserService
  ) {}

  @ViewChild(AutocomleteComponent)
  public autocompleteComponent!: AutocomleteComponent;

  @Input() public inputFormControl!: FormControl;

  @Input() public label = "";

  @Input() public multiple: boolean = false;

  @Input() public readonly: boolean = false;

  @Input() public usersToOmit: string[] = [];

  @Input() public optionValueKey?: string;

  @Input() public groupId?: string;

  @Input() public selectGroupMembersOnly: boolean = false;

  public users: User[] = [];

  public ngOnChanges(changes: SimpleChanges): void {
    if (changes["groupId"]) {
      this.updateValueOnGroupChange(changes["groupId"].currentValue);
    }

    if (changes["usersToOmit"]) {
      this.filterUsers();
    }
  }

  private updateValueOnGroupChange(groupId?: string): void {
    if (groupId) {
      this.users = this.groupMemberUserService.getUsersInGroup(groupId);
    } else {
      this.users = [];
      this.autocompleteComponent?.clearFilter();
    }
  }

  public ngOnInit(): void {
    if (this.users.length === 0 && this.usersToOmit.length === 0) {
      this.users = this.store.selectSnapshot(UserState.users);
    }
  }

  private filterUsers(): void {
    this.users = this.store
      .selectSnapshot(UserState.users)
      .filter((u) => !this.usersToOmit.includes(u.id.toString()));
  }

  public displayWith(id?: number): string {
    if (id) {
      const user = this.store.selectSnapshot(
        UserState.getUserById(id.toString())
      );
      return user?.displayName ?? "";
    }
    return "";
  }
}
