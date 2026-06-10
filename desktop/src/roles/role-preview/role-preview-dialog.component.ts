import { CommonModule } from "@angular/common";
import { Component, Inject, computed, inject } from "@angular/core";
import { toSignal } from "@angular/core/rxjs-interop";
import { MatButtonModule } from "@angular/material/button";
import {
  MAT_DIALOG_DATA,
  MatDialog,
  MatDialogModule,
  MatDialogRef,
} from "@angular/material/dialog";
import { MatIconModule } from "@angular/material/icon";
import { catchError, of } from "rxjs";
import { DEFAULT_DIALOG_CONFIG } from "../../constants/dialog.constant";
import {
  PermissionDescriptor,
  PermissionScope,
  PermissionService,
  Role,
} from "../../open-api";
import { SharedUiModule } from "../../shared-ui/shared-ui.module";
import {
  friendlyActionLabel,
  friendlyResourceName,
  iconForResource,
  resourceKeyOf,
} from "../role-presets";

export interface RolePreviewDialogData {
  role: Role;
}

interface PreviewRow {
  key: string;
  label: string;
  description: string;
}

interface PreviewGroup {
  resourceKey: string;
  name: string;
  icon: string;
  rows: PreviewRow[];
}

/**
 * Read-only preview of a role (app or group). Shows the role's scope, optional
 * description and its granted permissions grouped by resource. Opened from the
 * user and group-member forms so an admin can see what a role grants before
 * assigning it. Standalone so it can be opened from any feature module.
 */
@Component({
  selector: "app-role-preview-dialog",
  standalone: true,
  imports: [
    CommonModule,
    MatButtonModule,
    MatDialogModule,
    MatIconModule,
    SharedUiModule,
  ],
  templateUrl: "./role-preview-dialog.component.html",
  styleUrls: ["./role-preview-dialog.component.scss"],
})
export class RolePreviewDialogComponent {
  public readonly role: Role;

  // The role already carries its permission keys; the registry only supplies
  // human-readable labels/descriptions, so a failed load degrades gracefully to
  // the raw keys (via the fallback labels in `groups`).
  private readonly registry = toSignal(
    inject(PermissionService)
      .getPermissions()
      .pipe(catchError(() => of([] as PermissionDescriptor[]))),
    { initialValue: [] as PermissionDescriptor[] },
  );

  public readonly isGroup = computed<boolean>(
    () => this.role.scope === PermissionScope.Group,
  );

  public readonly scopeLabel = computed<string>(() =>
    this.isGroup() ? "Group role" : "Application role",
  );

  public readonly scopeIcon = computed<string>(() =>
    this.isGroup() ? "workspaces" : "apps",
  );

  public readonly permissionCount = computed<number>(
    () => this.role.permissions?.length ?? 0,
  );

  /** The role's permissions grouped by resource, with labels from the registry. */
  public readonly groups = computed<PreviewGroup[]>(() => {
    const descriptorByKey = new Map(
      this.registry().map((descriptor) => [descriptor.key, descriptor]),
    );
    const groups = new Map<string, PreviewGroup>();

    for (const key of this.role.permissions ?? []) {
      const resourceKey = resourceKeyOf(key);
      let group = groups.get(resourceKey);
      if (!group) {
        group = {
          resourceKey,
          name: friendlyResourceName(resourceKey),
          icon: iconForResource(resourceKey),
          rows: [],
        };
        groups.set(resourceKey, group);
      }

      const descriptor = descriptorByKey.get(key);
      group.rows.push({
        key,
        label: descriptor?.label ?? friendlyActionLabel(key),
        description: descriptor?.description ?? "",
      });
    }

    return [...groups.values()];
  });

  constructor(@Inject(MAT_DIALOG_DATA) data: RolePreviewDialogData) {
    this.role = data.role;
  }
}

/** Opens the shared role-preview dialog for the given role. */
export function openRolePreviewDialog(
  dialog: MatDialog,
  role: Role,
): MatDialogRef<RolePreviewDialogComponent> {
  return dialog.open(RolePreviewDialogComponent, {
    ...DEFAULT_DIALOG_CONFIG,
    // Focus the dialog container, not the trailing "Close" button. The preview's
    // only tabbable control sits at the bottom of a long, scrollable content
    // area, so the default "first-tabbable" focus would scroll it into view and
    // open the dialog at the bottom.
    autoFocus: "dialog",
    data: { role } as RolePreviewDialogData,
  });
}
