import { CdkDragDrop, moveItemInArray } from '@angular/cdk/drag-drop';
import { Component, Inject, OnInit } from '@angular/core';
import { MAT_DIALOG_DATA, MatDialogRef } from '@angular/material/dialog';
import { DEFAULT_RECEIPT_TABLE_COLUMNS, ReceiptTableColumnConfig } from '../../interfaces';

interface ColumnConfigItem extends ReceiptTableColumnConfig {
  displayName: string;
}

@Component({
  selector: 'app-column-configuration-dialog',
  templateUrl: './column-configuration-dialog.component.html',
  styleUrls: ['./column-configuration-dialog.component.scss'],
  standalone: false
})
export class ColumnConfigurationDialogComponent implements OnInit {
  public columns: ColumnConfigItem[] = [];

  private readonly columnDisplayNames: { [key: string]: string } = {
    'created_at': 'Added At',
    'date': 'Receipt Date',
    'name': 'Name',
    'paid_by_user_id': 'Paid By',
    'amount': 'Amount',
    'categories': 'Categories',
    'tags': 'Tags',
    'status': 'Status',
    'resolved_date': 'Resolved Date'
  };

  constructor(
    private dialogRef: MatDialogRef<ColumnConfigurationDialogComponent>,
    @Inject(MAT_DIALOG_DATA) public data: { currentColumns?: ReceiptTableColumnConfig[] }
  ) {}

  ngOnInit(): void {
    this.initializeColumns();
  }

  private initializeColumns(): void {
    const currentColumns = this.data.currentColumns || DEFAULT_RECEIPT_TABLE_COLUMNS;
    
    this.columns = currentColumns
      .sort((a, b) => a.order - b.order)
      .map(col => ({
        ...col,
        displayName: this.columnDisplayNames[col.matColumnDef] || col.matColumnDef
      }));
  }

  public toggleColumnVisibility(column: ColumnConfigItem): void {
    column.visible = !column.visible;
  }

  public drop(event: CdkDragDrop<ColumnConfigItem[]>): void {
    moveItemInArray(this.columns, event.previousIndex, event.currentIndex);
    
    this.columns.forEach((column, index) => {
      column.order = index;
    });
  }

  public resetToDefaults(): void {
    this.columns = DEFAULT_RECEIPT_TABLE_COLUMNS.map(col => ({
      ...col,
      displayName: this.columnDisplayNames[col.matColumnDef] || col.matColumnDef
    }));
  }

  public saveConfiguration(): void {
    const result: ReceiptTableColumnConfig[] = this.columns.map(({ displayName, ...col }) => col);
    this.dialogRef.close(result);
  }

  public cancel(): void {
    this.dialogRef.close(null);
  }

  public get visibleColumnsCount(): number {
    return this.columns.filter(col => col.visible).length;
  }

  public canToggleOff(column: ColumnConfigItem): boolean {
    return this.visibleColumnsCount > 1 || !column.visible;
  }
}