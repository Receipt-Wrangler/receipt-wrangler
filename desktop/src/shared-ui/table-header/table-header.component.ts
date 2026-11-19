import { Component, input } from '@angular/core';

@Component({
    selector: 'app-table-header',
    templateUrl: './table-header.component.html',
    styleUrls: ['./table-header.component.scss']
})
export class TableHeaderComponent {
  public readonly headerText = input<string>('');

  public readonly subtitle = input<string>('');
}
