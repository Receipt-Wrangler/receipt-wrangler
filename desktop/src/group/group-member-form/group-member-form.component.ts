import { Component, OnInit } from "@angular/core";
import { FormGroup } from "@angular/forms";
import { MatDialogRef } from "@angular/material/dialog";
import { Select, Store } from "@ngxs/store";
import { Observable } from "rxjs";
import { FormOption } from "src/interfaces/form-option.interface";
import { GroupMember } from "../../open-api";
import { AuthState } from "../../store";
import { GROUP_ROLE_OPTIONS } from "../role-options";
import { buildGroupMemberForm } from "../utils/group-member.utils";

@Component({
    selector: "app-group-member-form",
    templateUrl: "./group-member-form.component.html",
    styleUrls: ["./group-member-form.component.scss"],
    standalone: false
})
export class GroupMemberFormComponent implements OnInit {
  @Select(AuthState.userId) public userId!: Observable<string>;

  public headerText: string = "";

  public form: FormGroup = new FormGroup({});

  public roleOptions: FormOption[] = GROUP_ROLE_OPTIONS;

  public usersToOmit: string[] = [];

  public currentGroupMembers: GroupMember[] = [];

  public groupMember: GroupMember | undefined = undefined;

  constructor(
    private matDialogRef: MatDialogRef<GroupMemberFormComponent>,
    private store: Store
  ) {}

  public ngOnInit(): void {
    this.form = buildGroupMemberForm(this.groupMember);
    this.setUsersToOmit();
  }

  private setUsersToOmit(): void {
    const userId = this.store.selectSnapshot(AuthState.userId);
    this.usersToOmit = [
      userId.toString(),
      ...this.currentGroupMembers.map((c) => c.userId.toString()),
    ];
  }

  public closeModal(): void {
    this.matDialogRef.close();
  }

  public submit(): void {
    if (this.form.valid) {
      this.matDialogRef.close(this.form);
    }
  }
}
