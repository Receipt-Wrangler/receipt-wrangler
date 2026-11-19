import { Component, Input, signal, TemplateRef, input, output } from "@angular/core";
import { MatList, MatListItem, MatListItemTitle, MatListItemLine, MatListItemMeta } from "@angular/material/list";
import { NgTemplateOutlet } from "@angular/common";
import { EditButtonComponent } from "../edit-button/edit-button.component";
import { DeleteButtonComponent } from "../delete-button/delete-button.component";

@Component({
    selector: "app-editable-list",
    templateUrl: "./editable-list.component.html",
    styleUrl: "./editable-list.component.scss",
    imports: [MatList, NgTemplateOutlet, MatListItem, MatListItemTitle, MatListItemLine, MatListItemMeta, EditButtonComponent, DeleteButtonComponent]
})
export class EditableListComponent {
  public readonly listData = input<any[]>([]);

  public readonly itemTitleTemplate = input.required<TemplateRef<any>>();

  public readonly itemSubtitleTemplate = input.required<TemplateRef<any>>();

  public readonly trackByKey = input<string>("");

  @Input() public editTemplate?: TemplateRef<any>;

  public readonly readonly = input<boolean>(false);

  public readonly editButtonClicked = output<number>();

  public readonly deleteButtonClicked = output<number>();

  public readonly rowOpen = signal<number | undefined>(undefined);

  public handleEditButtonClicked(index: number): void {
    this.rowOpen.set(index);
    this.editButtonClicked.emit(index);
  }

  public getCurrentRowOpen(): number | undefined {
    return this.rowOpen();
  }

  public handleDeleteButtonClicked(index: number): void {
    this.rowOpen.set(undefined);
    this.deleteButtonClicked.emit(index);
  }

  public openLastRow(index?: number): void {
    this.rowOpen.set(index ?? this.listData().length - 1);
  }

  public closeRow(): void {
    this.rowOpen.set(undefined);
  }
}
