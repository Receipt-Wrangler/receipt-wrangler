import { ComponentFixture, TestBed } from '@angular/core/testing';
import { MatSortModule } from '@angular/material/sort';
import { MatTableModule } from '@angular/material/table';

import { TableComponent } from './table.component';
import { TableColumn } from '../table-column.interface';

describe('TableComponent', () => {
  let component: TableComponent;
  let fixture: ComponentFixture<TableComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [TableComponent],
      imports: [MatTableModule, MatSortModule],
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
    component.columns = columns;
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
    component.columns = columns;
    component.ngOnChanges({ columns: {} } as any);

    expect(component.defaultSort).toEqual({
      active: '',
      direction: '',
    });
  });
});
