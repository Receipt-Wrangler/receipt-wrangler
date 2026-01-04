import { Component, EventEmitter, Input, Output, ViewChild } from "@angular/core";
import { FormArray, FormGroup } from "@angular/forms";
import { FormCommand } from "../../form/index";
import { FormConfig } from "../../interfaces/index";
import { UserShortcut } from "../../open-api/index";
import { EditableListComponent } from "../../shared-ui/editable-list/editable-list.component";

@Component({
    selector: "app-user-shortcut",
    templateUrl: "./user-shortcut.component.html",
    styleUrl: "./user-shortcut.component.scss",
    standalone: false
})
export class UserShortcutComponent {

  @ViewChild(EditableListComponent) public editableListComponent!: EditableListComponent;

  @Input() public parentForm!: FormGroup;

  @Input() public formConfig!: FormConfig;

  @Input() public originalUserShortcuts: UserShortcut[] = [];

  @Output() public formCommand = new EventEmitter<FormCommand>();

  @Output() public shortcutDoneClicked = new EventEmitter<void>();

  @Output() public shortcutCancelClicked = new EventEmitter<void>();

  public isAddingShortcut = false;


  public get userShortcuts(): FormArray {
    return (this.parentForm?.get("userShortcuts") as FormArray || new FormArray([]));
  }

  public removeShortcut(index: number): void {
    this.formCommand.emit({
      path: `userShortcuts`,
      command: "removeAt",
      payload: index
    });
  }

  public emitShortcutDoneClicked(): void {
    this.shortcutDoneClicked.emit();
  }

  public emitShortcutCancelClicked(): void {
    this.shortcutCancelClicked.emit();
  }
}
