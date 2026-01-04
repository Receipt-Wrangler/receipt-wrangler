import {
  Component,
  EventEmitter,
  Input,
  Output,
  TemplateRef,
} from '@angular/core';

@Component({
    selector: 'app-form-list',
    templateUrl: './form-list.component.html',
    styleUrls: ['./form-list.component.scss'],
    standalone: false
})
export class FormListComponent {
  @Input() public array: any[] = [];

  @Input() public itemDisplayTemplate!: TemplateRef<any>;

  @Input() public itemEditTemplate!: TemplateRef<any>;

  @Input() public nothingToDisplayText: string = '';

  @Input() public addButtonText: string = '';

  @Input() public headerText: string = '';

  @Input() public disabled: boolean = false;

  @Output() public addButtonClicked: EventEmitter<void> = new EventEmitter();

  @Output() public itemDoneButtonClicked: EventEmitter<number> =
    new EventEmitter<number>();

  @Output() public itemCancelButtonClicked: EventEmitter<number> =
    new EventEmitter<number>();

  @Output() public itemDeleteButtonClicked: EventEmitter<number> =
    new EventEmitter<number>();

  public editingIndex: number = -1;

  public onAddButtonClicked(): void {
    this.editingIndex = this.array.length;
    this.addButtonClicked.emit();
  }

  public onDoneButtonClicked(index: number): void {
    this.itemDoneButtonClicked.emit(index);
  }

  public onItemCancelButtonClicked(index: number): void {
    this.resetEditingIndex();
    this.itemCancelButtonClicked.emit(index);
  }

  public resetEditingIndex(): void {
    this.editingIndex = -1;
  }

  public onItemEditButtonClicked(index: number): void {
    this.editingIndex = index;
  }

  public onItemDeleteButtonClicked(index: number): void {
    this.itemDeleteButtonClicked.emit(index);
  }
}
