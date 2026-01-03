import { Component, EventEmitter, Input, Output, TemplateRef } from "@angular/core";
import { BehaviorSubject } from "rxjs";

@Component({
    selector: "app-editable-list",
    templateUrl: "./editable-list.component.html",
    styleUrl: "./editable-list.component.scss",
    standalone: false
})
export class EditableListComponent {
  @Input() public listData: any[] = [];

  @Input() public itemTitleTemplate!: TemplateRef<any>;

  @Input() public itemSubtitleTemplate!: TemplateRef<any>;

  @Input() public trackByKey: string = "";

  @Input() public editTemplate?: TemplateRef<any>;

  @Input() public readonly: boolean = false;

  @Output() public editButtonClicked: EventEmitter<number> = new EventEmitter<number>();

  @Output() public deleteButtonClicked: EventEmitter<number> = new EventEmitter<number>();

  private rowOpen: BehaviorSubject<number | undefined> = new BehaviorSubject<
    number | undefined
  >(undefined);

  public rowOpenObservable = this.rowOpen.asObservable();

  public handleEditButtonClicked(index: number): void {
    this.rowOpen.next(index);
    this.editButtonClicked.emit(index);
  }

  public getCurrentRowOpen(): number | undefined {
    return this.rowOpen.value;
  }

  public handleDeleteButtonClicked(index: number): void {
    this.rowOpen.next(undefined);
    this.deleteButtonClicked.emit(index);
  }

  public openLastRow(): void {
    this.rowOpen.next(this.listData.length - 1);
  }

  public closeRow(): void {
    this.rowOpen.next(undefined);
  }
}
