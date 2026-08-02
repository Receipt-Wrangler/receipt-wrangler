import { animate, state, style, transition, trigger } from "@angular/animations";
import { LiveAnnouncer } from "@angular/cdk/a11y";
import { SelectionModel } from "@angular/cdk/collections";
import { Component, Input, OnChanges, SimpleChanges, input, output, viewChild } from "@angular/core";
import { MatPaginator, PageEvent } from "@angular/material/paginator";
import { MatSort, Sort } from "@angular/material/sort";
import { MatTableDataSource } from "@angular/material/table";
import { TableColumn } from "../table-column.interface";

@Component({
  selector: "app-table",
  templateUrl: "./table.component.html",
  styleUrls: ["./table.component.scss"],
  animations: [
    trigger("detailExpand", [
      state("collapsed,void", style({ height: "0px", minHeight: "0" })),
      state("expanded", style({ height: "*" })),
      transition("expanded <=> collapsed", animate("225ms cubic-bezier(0.4, 0.0, 0.2, 1)")),
    ]),
  ],
  standalone: false
})
export class TableComponent implements OnChanges {
  public readonly sort = viewChild.required(MatSort);
  public readonly paginator = viewChild.required(MatPaginator);

  public readonly columns = input<TableColumn[]>([]);
  public readonly displayedColumns = input<string[]>([]);
  public readonly dataSource = input(new MatTableDataSource<any>([]));
  public readonly pagination = input<boolean>(false);
  public readonly selectionCheckboxes = input<boolean>(false);
  public readonly page = input<number>(0);
  public readonly pageSize = input<number>(50);
  @Input() public length: number = 0;
  public readonly expandedRowTemplate = input<any>();
  public readonly rowExpandable = input<(row: any) => boolean>(() => true);

  public readonly sorted = output<Sort>();
  public readonly pageChange = output<PageEvent>();

  public defaultSort?: Sort;

  public selection = new SelectionModel<any>(true, []);

  public expandedElement: any;

  public rowIndexes: { [key: number]: number } = {};
  private rowIndexesByReference = new Map<any, number>();

  constructor(private _liveAnnouncer: LiveAnnouncer) {}

  public ngOnChanges(changes: SimpleChanges): void {
    if (changes["columns"]) {
      const column = this.columns().find((c) => c.defaultSortDirection);
      if (column) {
        this.defaultSort = {
          active: column.matColumnDef,
          direction: column.defaultSortDirection ?? "desc",
        };
        this.sort().sort({
          id: column.matColumnDef,
          start: column.defaultSortDirection as any,
          disableClear: true,
        });
      } else {
        this.defaultSort = {
          active: "",
          direction: "",
        };
      }
    }

    if (changes["dataSource"]) {
      this.setRowIndexes();
    }
  }

  /**
   * The index a cell template should act on: the row's position in the ORIGINAL
   * data, so an action still targets the right record after the table is sorted.
   *
   * Rows WITHOUT an id (e.g. GroupMember, whose key is composite userId+groupId)
   * fall back to the rendered index. They used to all collapse onto a single
   * `undefined` key in rowIndexes, so every such row reported the LAST row's
   * index — meaning edit/delete in an actions column hit the wrong record.
   */
  public indexFor(element: any): number {
    if (element?.id !== undefined && element?.id !== null) {
      const mapped = this.rowIndexes[element.id];
      if (mapped !== undefined) {
        return mapped;
      }
    }
    // Reference lookup — the table renders these exact objects, so this resolves
    // an id-less row to its real position. Precomputed rather than an indexOf
    // scan, because the template calls this for every rendered cell.
    return this.rowIndexesByReference.get(element) ?? -1;
  }

  // Two maps because the rows have two possible identities: by id where there is
  // one, by object reference for composite-key rows (e.g. GroupMember).
  //
  // Both rely on callers REPLACING the dataSource rather than mutating its `data`
  // in place — which every consumer does (`dataSource.set(new MatTableDataSource(...))`),
  // and which is what makes ngOnChanges fire. An in-place mutation would leave both
  // maps stale and indexFor pointing at the wrong record.
  private setRowIndexes(): void {
    const indexes: { [key: number]: number } = {};
    const byReference = new Map<any, number>();
    this.dataSource().data.forEach((row, index) => {
      if (row.id !== undefined && row.id !== null) {
        indexes[row.id] = index;
      }
      byReference.set(row, index);
    });

    this.rowIndexes = indexes;
    this.rowIndexesByReference = byReference;
  }

  public isAllSelected() {
    const numSelected = this.selection.selected.length;
    const numRows = this.dataSource().data.length;
    return numSelected === numRows;
  }

  public toggleAllRows() {
    if (this.isAllSelected()) {
      this.selection.clear();
      return;
    }

    this.selection.select(...this.dataSource().data);
  }

  /** Announce the change in sort state for assistive technology. */
  announceSortChange(sortState: Sort) {
    // This example uses English messages. If your application supports
    // multiple language, you would internationalize these strings.
    // Furthermore, you can customize the message to add additional
    // details about the values being sorted.
    if (sortState.direction) {
      this._liveAnnouncer.announce(`Sorted ${sortState.direction}ending`);
    } else {
      this._liveAnnouncer.announce("Sorting cleared");
    }

    this.sorted.emit(sortState);
  }

  public pageChanged(pageEvent: PageEvent): void {
    this.selection.clear();
    this.pageChange.emit(pageEvent);
  }

  public expanderClicked(event: MouseEvent, row: any): void {
    this.expandedElement = this.expandedElement === row ? null : row;
    event.stopPropagation();
  }
}
