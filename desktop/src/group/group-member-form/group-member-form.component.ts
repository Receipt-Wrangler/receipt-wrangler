import { Component, DestroyRef, OnInit, computed, effect, inject, signal } from "@angular/core";
import { takeUntilDestroyed, toSignal } from "@angular/core/rxjs-interop";
import { FormGroup } from "@angular/forms";
import { MatDialog, MatDialogRef } from "@angular/material/dialog";
import { Store } from "@ngxs/store";
import { catchError, forkJoin, of, startWith, take, tap } from "rxjs";
import {
  Category,
  CategoryService,
  Group,
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
  private readonly rolesLoaded = signal(false);
  private readonly roles = toSignal(
    inject(RoleService)
      .getRoles()
      .pipe(
        catchError(() => of([] as Role[])),
        tap(() => this.rolesLoaded.set(true))
      ),
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
  private readonly persistedMember = signal<GroupMember | null>(null);
  private readonly poolsLoaded = signal(false);

  /**
   * Whether every async input the assignment section depends on has landed: the
   * membership itself, the category/tag pools, and the role list (which supplies
   * the ceiling).
   *
   * The section is withheld until all of them are ready. Rendering earlier makes
   * it re-render as each one arrives — which flashes an unconstrained picker at
   * the admin before the ceiling narrows it, and briefly shows it empty before
   * the saved assignment resolves.
   */
  public readonly grantsReady = computed<boolean>(
    () => !!this.persistedMember() && this.poolsLoaded() && this.rolesLoaded()
  );

  public readonly grantRow = computed<MemberGrantRow | null>(() => {
    const member = this.persistedMember();
    if (!member) {
      return null;
    }

    // Only the group's id matters here — the dialog already shows which group it
    // belongs to, so a placeholder name keeps this off GroupState entirely.
    return buildMemberGrantRow(
      { id: member.groupId, name: '' } as Group,
      member,
      this.groupRoleOptions(),
    );
  });

  private readonly editedGrants = new Map<number, GrantSelection>();

  public ngOnInit(): void {
    this.form = buildGroupMemberForm(this.groupMember);
    this.setUsersToOmit();
    this.listenToGroupRoleChanges();
    this.resolvePersistedMember();
  }

  /**
   * Loads the saved membership behind this dialog, if any.
   *
   * Fetched from the API rather than read out of GroupState: grants are written
   * through their own endpoint, which does not refresh the store, so a store copy
   * can be stale by exactly the edit the admin just made in the user form. The
   * dialog must show the current assignment, not a cached one.
   */
  private resolvePersistedMember(): void {
    const groupId = this.groupMember?.groupId;
    const userId = this.groupMember?.userId;
    if (!groupId || !userId) {
      return;
    }

    this.groupsService
      .getGroupById(groupId)
      .pipe(take(1), catchError(() => of(null)), takeUntilDestroyed(this.destroyRef))
      .subscribe((group) => {
        const member = group?.groupMembers?.find(
          (candidate: GroupMember) => candidate.userId === userId,
        );
        if (!member) {
          return;
        }

        this.persistedMember.set(member);
        this.loadGrantPools();
      });
  }

  // Both pools are loaded together so the section can be withheld until they are
  // BOTH in, rather than rendering twice as they trickle in separately.
  private loadGrantPools(): void {
    forkJoin({
      categories: this.categoryService
        .getAllCategories()
        .pipe(take(1), catchError(() => of([] as Category[]))),
      tags: this.tagService.getAllTags().pipe(take(1), catchError(() => of([] as Tag[]))),
    })
      .pipe(takeUntilDestroyed(this.destroyRef))
      .subscribe(({ categories, tags }) => {
        this.categoryPool.set(categories);
        this.tagPool.set(tags);
        this.poolsLoaded.set(true);
      });
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
