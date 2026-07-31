import { Component, DestroyRef, OnInit, computed, effect, inject, signal } from "@angular/core";
import { takeUntilDestroyed, toSignal } from "@angular/core/rxjs-interop";
import { FormGroup } from "@angular/forms";
import { MatDialog, MatDialogRef } from "@angular/material/dialog";
import { Store } from "@ngxs/store";
import { catchError, of, startWith, take, tap } from "rxjs";
import {
  Category,
  CategoryService,
  GroupMember,
  GroupsService,
  PermissionScope,
  Role,
  RoleService,
  Tag,
  TagService,
} from "../../open-api";
import { openRolePreviewDialog } from "../../roles/role-preview/role-preview-dialog.component";
import { GrantSelection } from "../../shared-ui/grant-picker/grant-picker.component";
import {
  buildMemberGrantRow,
  MemberGrantRow,
  saveChangedMemberGrants,
} from "../../shared-ui/grant-picker/member-grant-assignment";
import { AuthState, GroupState } from "../../store";
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
    () =>
      this.groupRoleOptions().find(
        (role) => role.id === this.selectedGroupRoleId()
      ) ?? null
  );

  constructor(
    private matDialogRef: MatDialogRef<GroupMemberFormComponent>,
    private store: Store,
    private matDialog: MatDialog,
    private destroyRef: DestroyRef,
    private categoryService: CategoryService,
    private tagService: TagService,
    private groupsService: GroupsService
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

  // ----- Per-member category/tag assignment -----
  // Only offered for a membership that already exists on the server: grants are
  // written through their own endpoint, keyed by (group, user), so there is
  // nothing to write until the member has been saved. A member being added here
  // is only persisted when the parent group form is submitted.
  public readonly categoryPool = signal<Category[]>([]);
  public readonly tagPool = signal<Tag[]>([]);
  private readonly groups = this.store.selectSignal(GroupState.groups);
  private readonly persistedMember = signal<GroupMember | null>(null);

  public readonly grantRow = computed<MemberGrantRow | null>(() => {
    const member = this.persistedMember();
    if (!member) {
      return null;
    }

    const group = this.groups().find((candidate) => candidate.id === member.groupId);
    if (!group) {
      return null;
    }

    return buildMemberGrantRow(group, member, this.groupRoleOptions());
  });

  private readonly editedGrants = new Map<number, GrantSelection>();

  public ngOnInit(): void {
    this.form = buildGroupMemberForm(this.groupMember);
    this.setUsersToOmit();
    this.listenToGroupRoleChanges();
    this.resolvePersistedMember();
  }

  /**
   * Finds the saved membership behind this dialog, if any. The dialog is handed
   * the parent form array's value, which carries no grants — the authoritative
   * copy (with grants) lives on the group in the store.
   */
  private resolvePersistedMember(): void {
    const groupId = this.groupMember?.groupId;
    const userId = this.groupMember?.userId;
    if (!groupId || !userId) {
      return;
    }

    const member = this.groups()
      .find((group) => group.id === groupId)
      ?.groupMembers?.find((candidate) => candidate.userId === userId);
    if (!member) {
      return;
    }

    this.persistedMember.set(member);
    this.loadGrantPools();
  }

  private loadGrantPools(): void {
    this.categoryService
      .getAllCategories()
      .pipe(take(1), catchError(() => of([] as Category[])), takeUntilDestroyed(this.destroyRef))
      .subscribe((categories) => this.categoryPool.set(categories));
    this.tagService
      .getAllTags()
      .pipe(take(1), catchError(() => of([] as Tag[])), takeUntilDestroyed(this.destroyRef))
      .subscribe((tags) => this.tagPool.set(tags));
  }

  public onGrantsChange(groupId: number, selection: GrantSelection): void {
    this.editedGrants.set(groupId, selection);
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
    if (!this.form.valid) {
      return;
    }

    const row = this.grantRow();
    const member = this.persistedMember();
    if (!row || !member) {
      this.matDialogRef.close(this.form);
      return;
    }

    // Grants have their own endpoint and their own permission, so they are saved
    // here and now rather than riding the parent group form's submit. Only a
    // changed assignment issues a request.
    saveChangedMemberGrants(this.groupsService, member.userId, [row], this.editedGrants)
      .pipe(
        take(1),
        tap(() => this.matDialogRef.close(this.form)),
        catchError(() => {
          // The interceptor surfaces the failure; keep the dialog open so the
          // admin can correct the selection rather than losing it.
          return of(undefined);
        })
      )
      .subscribe();
  }
}
