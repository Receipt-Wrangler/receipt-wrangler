import { Component, computed, Inject, signal } from "@angular/core";
import { MAT_DIALOG_DATA, MatDialogRef } from "@angular/material/dialog";
import { buildSplitDiff, collapseUnchanged, DiffLine, InlineSegments, inlineChange, SplitDiffRow } from "../../utils/line-diff";
import { ReceiptUpdateSnapshots, toJsonLines } from "../../utils/receipt-update-description";
import { FilterTab } from "../filter-bar/filter-tab.interface";

export interface ReceiptUpdateDiffDialogData {
  snapshots: ReceiptUpdateSnapshots;
}

export type ReceiptDiffView = "all" | "changes";

/** One side of a diff row, split so the changed part can be highlighted. */
export interface ReceiptDiffCell extends InlineSegments {
  lineNumber: number;
}

export interface ReceiptDiffRow {
  kind: SplitDiffRow["kind"];
  left?: ReceiptDiffCell;
  right?: ReceiptDiffCell;
}

/**
 * Before/after view of a RECEIPT_UPDATED system task: the two receipt
 * snapshots as pretty-printed JSON, side by side, with changed lines
 * highlighted. "Changes only" collapses the unchanged lines around each change.
 */
@Component({
  selector: "app-receipt-update-diff-dialog",
  templateUrl: "./receipt-update-diff-dialog.component.html",
  styleUrl: "./receipt-update-diff-dialog.component.scss",
  standalone: false,
})
export class ReceiptUpdateDiffDialogComponent {
  public readonly headerText: string;

  public readonly beforeUpdatedAt?: string;

  public readonly afterUpdatedAt?: string;

  /** The pair as readable JSON — the stored description is double-encoded. */
  public readonly copyText: string;

  public readonly addedLineCount: number;

  public readonly removedLineCount: number;

  public readonly tabs: FilterTab[];

  public readonly view = signal<ReceiptDiffView>("all");

  private readonly rows: ReceiptDiffRow[];

  public readonly displayRows = computed(() =>
    this.view() === "changes" ? collapseUnchanged(this.rows) : this.rows
  );

  constructor(
    public dialogRef: MatDialogRef<ReceiptUpdateDiffDialogComponent>,
    @Inject(MAT_DIALOG_DATA) data: ReceiptUpdateDiffDialogData,
  ) {
    const { before, after } = data.snapshots;
    const splitRows = buildSplitDiff(toJsonLines(before), toJsonLines(after));

    this.rows = splitRows.map(toReceiptDiffRow);
    this.addedLineCount = splitRows.filter((row) => row.kind !== "equal" && row.right).length;
    this.removedLineCount = splitRows.filter((row) => row.kind !== "equal" && row.left).length;

    const changedRowCount = splitRows.filter((row) => row.kind !== "equal").length;
    this.tabs = [
      { value: "all", label: "All lines" },
      { value: "changes", label: "Changes only", count: changedRowCount },
    ];

    const name = typeof after["name"] === "string" ? after["name"] : "";
    this.headerText = name ? `Receipt update: ${name}` : "Receipt update";
    this.beforeUpdatedAt = asString(before["updatedAt"]);
    this.afterUpdatedAt = asString(after["updatedAt"]);
    this.copyText = JSON.stringify({ before, after }, null, 2);
  }

  public setView(view: string): void {
    this.view.set(view as ReceiptDiffView);
  }

  public close(): void {
    this.dialogRef.close();
  }
}

function toReceiptDiffRow(row: SplitDiffRow): ReceiptDiffRow {
  if (row.kind === "changed" && row.left && row.right) {
    const segments = inlineChange(row.left.text, row.right.text);
    return {
      kind: row.kind,
      left: { lineNumber: row.left.lineNumber, ...segments.left },
      right: { lineNumber: row.right.lineNumber, ...segments.right },
    };
  }

  return { kind: row.kind, left: wholeLine(row.left), right: wholeLine(row.right) };
}

/** A line with no inline highlight: the row's own tint already marks it. */
function wholeLine(line?: DiffLine): ReceiptDiffCell | undefined {
  return line ? { lineNumber: line.lineNumber, prefix: line.text, changed: "", suffix: "" } : undefined;
}

function asString(value: unknown): string | undefined {
  return typeof value === "string" ? value : undefined;
}
