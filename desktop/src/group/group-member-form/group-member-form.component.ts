import { Component, DestroyRef, OnInit, computed, signal } from "@angular/core";
import { takeUntilDestroyed } from "@angular/core/rxjs-interop";
import { FormGroup } from "@angular/forms";
import { MatDialog, MatDialogRef } from "@angular/material/dialog";
import { Store } from "@ngxs/store";
import { catchError, of, startWith, take } from "rxjs";
import { GroupMember, PermissionScope, Role, RoleService } from "../../open-api";
import { openRolePreviewDialog } from "../../roles/role-preview/role-preview-dialog.component";
import { AuthState } from "../../store";
import { buildGroupMemberForm, legacyGroupRoleFromRole } from "../utils/group-member.utils";

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
  // roles load asynchronously.
  private readonly roles = signal<Role[]>([]);
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
    private roleService: RoleService,
    private matDialog: MatDialog,
    private destroyRef: DestroyRef
  ) {}

  public ngOnInit(): void {
    this.form = buildGroupMemberForm(this.groupMember);
    this.setUsersToOmit();
    this.loadRoles();
    this.listenToGroupRoleChanges();
  }

  private loadRoles(): void {
    this.roleService
      .getRoles()
      // Degrade gracefully if the roles can't be read (e.g. a non-admin group
      // owner lacks app.roles.read until permission-based gating ships): the
      // selector is simply empty rather than erroring.
      .pipe(
        take(1),
        catchError(() => of([] as Role[])),
        takeUntilDestroyed(this.destroyRef)
      )
      .subscribe((roles) => {
        this.roles.set(roles);
        // A new member defaults to the configured default group role; an existing
        // member's role is already pre-filled from the loaded member.
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
      // Keep the legacy enum in sync with the chosen modern role so the parent
      // group form's member table and "keep an owner" check (which still read
      // the legacy enum) work for newly assigned members. The backend re-derives
      // it authoritatively on save.
      const role = this.selectedRole();
      if (role) {
        this.form.get("groupRole")?.setValue(legacyGroupRoleFromRole(role));
      }
      this.matDialogRef.close(this.form);
    }
  }
}
