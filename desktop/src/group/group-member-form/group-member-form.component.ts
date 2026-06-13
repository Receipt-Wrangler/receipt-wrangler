import { Component, DestroyRef, OnInit, computed, effect, inject, signal } from "@angular/core";
import { takeUntilDestroyed, toSignal } from "@angular/core/rxjs-interop";
import { FormGroup } from "@angular/forms";
import { MatDialog, MatDialogRef } from "@angular/material/dialog";
import { Store } from "@ngxs/store";
import { catchError, of, startWith } from "rxjs";
import { GroupMember, PermissionScope, Role, RoleService } from "../../open-api";
import { openRolePreviewDialog } from "../../roles/role-preview/role-preview-dialog.component";
import { AuthState } from "../../store";
import { buildGroupMemberForm } from "../utils/group-member.utils";

@Component({
    selector: "app-group-member-form",
    templateUrl: "./group-member-form.component.html",
    styleUrls: ["./group-member-form.component.scss"],
    standalone: false
})
export class GroupMemberFormComponent implements OnInit {
  public userId = this.store.selectSignal(AuthState.userId);

  public headerText: string = "";

  public form: FormGroup = new FormGroup({});

  public usersToOmit: string[] = [];

  public currentGroupMembers: GroupMember[] = [];

  public groupMember: GroupMember | undefined = undefined;

  // The group-role list, the selector options and the currently-selected role
  // (for the preview button) are derived reactively so they update once the
  // roles load asynchronously. Degrade gracefully if the roles can't be read
  // (e.g. a non-admin group owner lacks app.roles.read until permission-based
  // gating ships): the selector is simply empty rather than erroring.
  private readonly roles = toSignal(
    inject(RoleService)
      .getRoles()
      .pipe(catchError(() => of([] as Role[]))),
    { initialValue: [] as Role[] }
  );
  public readonly groupRoleOptions = computed<Role[]>(() =>
    this.roles().filter((role) => role.scope === PermissionScope.Group)
  );
  private readonly selectedGroupRoleId = signal<number | null>(null);
  public readonly selectedRole = computed<Role | null>(
    () => this.roles().find((role) => role.id === this.selectedGroupRoleId()) ?? null
  );

  constructor(
    private matDialogRef: MatDialogRef<GroupMemberFormComponent>,
    private store: Store,
    private matDialog: MatDialog,
    private destroyRef: DestroyRef
  ) {
    // A new member defaults to the configured default group role once the roles
    // resolve; an existing member's role is already pre-filled from the loaded
    // member. Syncing a derived value into the reactive form is an imperative
    // side effect, so an effect() is the sanctioned tool here.
    effect(() => {
      const roles = this.roles();
      if (!this.groupMember && this.form.get("groupRoleId")?.value == null) {
        const defaultRole = roles.find(
          (role) => role.scope === PermissionScope.Group && role.isDefault
        );
        if (defaultRole) {
          this.form.get("groupRoleId")?.setValue(defaultRole.id);
        }
      }
    });
  }

  public ngOnInit(): void {
    this.form = buildGroupMemberForm(this.groupMember);
    this.setUsersToOmit();
    this.listenToGroupRoleChanges();
  }

  private listenToGroupRoleChanges(): void {
    const control = this.form.get("groupRoleId");
    control?.valueChanges
      .pipe(startWith(control.value), takeUntilDestroyed(this.destroyRef))
      .subscribe((id) =>
        this.selectedGroupRoleId.set(id == null ? null : Number(id))
      );
  }

  private setUsersToOmit(): void {
    const userId = this.store.selectSnapshot(AuthState.userId);
    this.usersToOmit = [
      userId.toString(),
      ...this.currentGroupMembers.map((c) => c.userId.toString()),
    ];
  }

  public previewRole(): void {
    const role = this.selectedRole();
    if (role) {
      openRolePreviewDialog(this.matDialog, role);
    }
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
