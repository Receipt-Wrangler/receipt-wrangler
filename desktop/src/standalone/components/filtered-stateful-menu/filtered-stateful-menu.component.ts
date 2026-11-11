import { CdkMenu, CdkMenuTrigger } from "@angular/cdk/menu";
import { Component, Input, computed, input, output } from "@angular/core";
import { FormControl } from "@angular/forms";
import { toSignal } from "@angular/core/rxjs-interop";
import { MatCheckbox } from "@angular/material/checkbox";
import { MatMenuItem } from "@angular/material/menu";
import { BaseButtonComponent } from "../../../button/base-button/base-button.component";
import { StatefulMenuItem } from "./stateful-menu-item";
import { ButtonComponent } from "../../../button/button/button.component";
import { InputComponent } from "../../../input/input/input.component";

@Component({
  selector: "app-filtered-stateful-menu",
  imports: [
    CdkMenuTrigger,
    CdkMenu,
    ButtonComponent, InputComponent, MatMenuItem, MatCheckbox,
  ],
  templateUrl: "./filtered-stateful-menu.component.html",
  styleUrl: "./filtered-stateful-menu.component.scss",
})
export class FilteredStatefulMenuComponent extends BaseButtonComponent {
  public readonly items = input<StatefulMenuItem[]>([]);

  public readonly filterFunc = input((item: StatefulMenuItem, filter: string) => item.displayValue.toLowerCase().includes(filter?.toLowerCase() ?? ""));

  public readonly filterLabel = input("Filter options");

  @Input() public headerText = "";

  public readonly readonly = input(false);

  public readonly itemSelected = output<StatefulMenuItem>();

  public filterFormControl = new FormControl("");

  private filterValue = toSignal(this.filterFormControl.valueChanges, { initialValue: "" });

  public filteredItems = computed(() => {
    const filter = this.filterValue() ?? "";
    return this.filterItems(this.items(), filter);
  });

  public onItemSelected(item: StatefulMenuItem, event: MouseEvent) {
    event.stopPropagation();
    event.stopImmediatePropagation();
    event.preventDefault();

    if (!this.readonly()) {
      this.itemSelected.emit(item);
    }
  }

  public resetFilter(): void {
    this.filterFormControl.setValue("");
  }

  public filterItems(items: StatefulMenuItem[], filterString: string): StatefulMenuItem[] {
    if (!filterString) {
      return Array.from(items);
    }
    return items.filter(item => this.filterFunc()(item, filterString));
  }
}
