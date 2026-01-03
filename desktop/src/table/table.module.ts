import { CommonModule } from "@angular/common";
import { NgModule } from "@angular/core";
import { MatCheckboxModule } from "@angular/material/checkbox";
import { MatIcon } from "@angular/material/icon";
import { MatPaginatorModule } from "@angular/material/paginator";
import { MatSortModule } from "@angular/material/sort";
import { MatTableModule } from "@angular/material/table";
import { ButtonModule } from "../button";
import { TableComponent } from "./table/table.component";
import { RowExpandablePipe } from './table/row-expandable.pipe';

@NgModule({
  declarations: [TableComponent, RowExpandablePipe],
  imports: [
    CommonModule,
    MatTableModule,
    MatSortModule,
    MatPaginatorModule,
    MatCheckboxModule,
    MatIcon,
    ButtonModule,
  ],
  exports: [TableComponent],
})
export class TableModule {}
