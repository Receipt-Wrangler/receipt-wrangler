import { provideNoopAnimations } from '@angular/platform-browser/animations';
import { CommonModule } from '@angular/common';
import { Component, signal, TemplateRef, viewChild } from '@angular/core';
import { ComponentFixture, TestBed } from '@angular/core/testing';
import { MatSortModule } from '@angular/material/sort';
import { MatTableModule } from '@angular/material/table';

import { TableComponent } from './table.component';
import { MatTableDataSource } from '@angular/material/table';
import { TableColumn } from '../table-column.interface';

/** Renders one templated column so the template's context can be inspected. */
@Component({
  standalone: false,
  template: `
    <ng-template #cell let-column="column">{{ column.columnHeader }}</ng-template>
    <app-table
      [columns]="columns()"
      [displayedColumns]="displayedColumns()"
      [dataSource]="dataSource"
    ></app-table>
  `,
})
class ColumnContextHostComponent {
  readonly cell = viewChild.required<TemplateRef<any>>('cell');
  dataSource = new MatTableDataSource<any>([{ id: 1 }]);
  columns = signal<TableColumn[]>([]);
  // Empty until the column exists: mat-table throws on a displayed id it has no
  // definition for, and the template ref only resolves after the first pass.
  displayedColumns = signal<string[]>([]);
}

describe('TableComponent', () => {
  let component: TableComponent;
  let fixture: ComponentFixture<TableComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [TableComponent, ColumnContextHostComponent],
      imports: [CommonModule, MatTableModule, MatSortModule],
      // The expandedDetail row carries an animation trigger, so rendering the
      // table (rather than only its class) needs an animations provider.
      providers: [provideNoopAnimations()],
    }).compileComponents();

    fixture = TestBed.createComponent(TableComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('should find the first default sort direction and set it', () => {
    const columns: TableColumn[] = [
      {
        columnHeader: 'header',
        matColumnDef: 'column',
        sortable: false,
        defaultSortDirection: 'asc',
      },
      {
        columnHeader: 'header',
        matColumnDef: 'column1',
        sortable: false,
      },
      {
        columnHeader: 'header',
        matColumnDef: 'column2',
        sortable: false,
      },
    ];
    fixture.componentRef.setInput('columns', columns);
    component.ngOnChanges({ columns: {} } as any);

    expect(component.defaultSort).toEqual({
      active: 'column',
      direction: 'asc',
    });
  });

  it('should not find a default sort column', () => {
    const columns: TableColumn[] = [
      {
        columnHeader: 'header',
        matColumnDef: 'column',
        sortable: false,
      },
      {
        columnHeader: 'header',
        matColumnDef: 'column1',
        sortable: false,
      },
      {
        columnHeader: 'header',
        matColumnDef: 'column2',
        sortable: false,
      },
    ];
    fixture.componentRef.setInput('columns', columns);
    component.ngOnChanges({ columns: {} } as any);

    expect(component.defaultSort).toEqual({
      active: '',
      direction: '',
    });
  });

  // A cell template gets the whole column, not just the row, so one template can
  // serve many columns that differ only in the definition they render (the
  // receipts table draws every custom field this way).
  it('hands the column to a cell template', () => {
    const hostFixture = TestBed.createComponent(ColumnContextHostComponent);
    hostFixture.detectChanges();

    hostFixture.componentInstance.columns.set([
      {
        columnHeader: 'From the column',
        matColumnDef: 'column',
        sortable: false,
        template: hostFixture.componentInstance.cell(),
      },
    ]);
    hostFixture.componentInstance.displayedColumns.set(['column']);
    hostFixture.detectChanges();

    expect(
      (hostFixture.nativeElement as HTMLElement).querySelector('td')?.textContent
    ).toContain('From the column');
  });

  describe('indexFor', () => {
    // A cell template's `index` must point at the row's position in the ORIGINAL
    // data, so an actions column still edits/deletes the right record.

    it('maps an id-bearing row to its position in the source data', () => {
      const data = [{ id: 7 }, { id: 8 }, { id: 9 }];
      fixture.componentRef.setInput('dataSource', new MatTableDataSource(data));
      component.ngOnChanges({ dataSource: {} as any });

      expect(component.indexFor(data[0])).toBe(0);
      expect(component.indexFor(data[2])).toBe(2);
    });

    it('resolves a row WITHOUT an id to its own position', () => {
      // GroupMember has a composite key and no id. These rows used to all collapse
      // onto a single `undefined` key in rowIndexes, so every one reported the LAST
      // row's index — edit and DELETE in an actions column hit the wrong record.
      const data = [{ userId: 1 }, { userId: 2 }, { userId: 3 }];
      fixture.componentRef.setInput('dataSource', new MatTableDataSource(data));
      component.ngOnChanges({ dataSource: {} as any });

      expect(component.indexFor(data[0])).toBe(0);
      expect(component.indexFor(data[1])).toBe(1);
      expect(component.indexFor(data[2])).toBe(2);
    });

    it('handles a mix of id-bearing and id-less rows per row', () => {
      const data = [{ userId: 1 }, { id: 42 }, { userId: 3 }];
      fixture.componentRef.setInput('dataSource', new MatTableDataSource(data));
      component.ngOnChanges({ dataSource: {} as any });

      expect(component.indexFor(data[0])).toBe(0);
      expect(component.indexFor(data[1])).toBe(1);
      expect(component.indexFor(data[2])).toBe(2);
    });

    it('treats a null id as id-less rather than as a key', () => {
      const data = [{ id: null }, { id: undefined }];
      fixture.componentRef.setInput('dataSource', new MatTableDataSource(data));
      component.ngOnChanges({ dataSource: {} as any });

      expect(component.indexFor(data[0])).toBe(0);
      expect(component.indexFor(data[1])).toBe(1);
    });
  });
});