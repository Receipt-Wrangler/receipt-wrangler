import { Component, input, output, viewChild } from "@angular/core";
import { FormArray, FormGroup } from "@angular/forms";
import { FormCommand } from "../../form/index";
import { FormConfig } from "../../interfaces/index";
import { UserShortcut } from "../../open-api/index";
import { EditableListComponent } from "../../shared-ui/editable-list/editable-list.component";
import { MatIcon } from "@angular/material/icon";
import { CardComponent } from "../../shared-ui/card/card.component";
import { InputComponent } from "../../input/input/input.component";
import { IconAutocompleteComponent } from "../../shared-ui/icon-autocomplete/icon-autocomplete.component";
import { DialogFooterComponent } from "../../shared-ui/dialog-footer/dialog-footer.component";
import { FormGetPipe } from "../../pipes/form-get.pipe";
import { InputReadonlyPipe } from "../../pipes/input-readonly.pipe";

@Component({
    selector: "app-user-shortcut",
    templateUrl: "./user-shortcut.component.html",
    styleUrl: "./user-shortcut.component.scss",
    imports: [EditableListComponent, MatIcon, CardComponent, InputComponent, IconAutocompleteComponent, DialogFooterComponent, FormGetPipe, InputReadonlyPipe]
})
export class UserShortcutComponent {

  public readonly editableListComponent = viewChild.required(EditableListComponent);

  public readonly parentForm = input.required<FormGroup>();

  public readonly formConfig = input.required<FormConfig>();

  public readonly originalUserShortcuts = input<UserShortcut[]>([]);

  public readonly formCommand = output<FormCommand>();

  public readonly shortcutDoneClicked = output<void>();

  public readonly shortcutCancelClicked = output<void>();

  public isAddingShortcut = false;


  public get userShortcuts(): FormArray {
    return (this.parentForm()?.get("userShortcuts") as FormArray || new FormArray([]));
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
