import { ChangeDetectionStrategy, Component, Inject, computed, inject, signal } from "@angular/core";
import { MAT_DIALOG_DATA, MatDialogRef } from "@angular/material/dialog";
import { Store } from "@ngxs/store";
import { AuthState, GroupState } from "src/store";
import { Permission } from "../../../open-api";
import { CHIP_COLORS, groupInitials } from "../../models/report-chip.util";

export interface AddGroupDialogData {
  selectedGroupIds: string[];
}

interface GroupChoice {
  id: string;
  name: string;
  initials: string;
  color: string;
  selected: boolean;
}

/**
 * Picks the groups a report covers, limited to groups where the user holds
 * group.reports.read (the permission that authorizes generation). Returns the
 * selected group-id array on submit, or undefined on cancel.
 */
@Component({
  selector: "app-add-group-dialog",
  templateUrl: "./add-group-dialog.component.html",
  styleUrls: ["./add-group-dialog.component.scss"],
  standalone: false,
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class AddGroupDialogComponent {
  private readonly store = inject(Store);
  private readonly dialogRef = inject<MatDialogRef<AddGroupDialogComponent, string[]>>(MatDialogRef);
  private readonly selectedIds = signal<Set<string>>(new Set());

  public readonly groups = computed<GroupChoice[]>(() => {
    const selected = this.selectedIds();
    return this.store
      .selectSnapshot(GroupState.groupsWithoutAll)
      .filter(
        (group) =>
          group.id != null &&
          this.store.selectSnapshot(
            AuthState.hasGroupPermission(group.id, Permission.GroupReportsRead)
          )
      )
      .map((group, index) => {
        const id = group.id!.toString();
        return {
          id,
          name: group.name ?? id,
          initials: groupInitials(group.name),
          color: CHIP_COLORS[index % CHIP_COLORS.length],
          selected: selected.has(id),
        };
      });
  });

  constructor(@Inject(MAT_DIALOG_DATA) data: AddGroupDialogData) {
    this.selectedIds.set(new Set(data.selectedGroupIds ?? []));
  }

  public toggle(id: string): void {
    this.selectedIds.update((set) => {
      const next = new Set(set);
      if (next.has(id)) {
        next.delete(id);
      } else {
        next.add(id);
      }
      return next;
    });
  }

  public submit(): void {
    this.dialogRef.close([...this.selectedIds()]);
  }

  public cancel(): void {
    this.dialogRef.close();
  }
}
