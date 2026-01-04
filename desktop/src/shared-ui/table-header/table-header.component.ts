import { Component, Input } from '@angular/core';

@Component({
    selector: 'app-table-header',
    templateUrl: './table-header.component.html',
    styleUrls: ['./table-header.component.scss'],
    standalone: false
})
export class TableHeaderComponent {
  @Input() public headerText: string = '';
}
