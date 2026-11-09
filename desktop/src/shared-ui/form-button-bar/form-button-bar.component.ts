import { Component, HostBinding, input } from "@angular/core";
import { FormMode } from "src/enums/form-mode.enum";
import { NgClass } from "@angular/common";

/**
 * How far, in px, the bar keeps clear of each viewport edge. The bar is fixed to
 * the viewport, so a host whose own fixed panel covers part of it (the receipt
 * history drawer) passes the band that is left, and the bar shrinks to it.
 */
export interface FormButtonBarInset {
  left: number;
  right: number;
}

@Component({
    selector: "app-form-button-bar",
    templateUrl: "./form-button-bar.component.html",
    styleUrls: ["./form-button-bar.component.scss"],
    imports: [NgClass]
})
export class FormButtonBarComponent {
  public readonly mode = input<FormMode>();

  public readonly justifyContentEnd = input(true);

  /** Null keeps the bar across the whole viewport, as every form but one uses it. */
  public readonly inset = input<FormButtonBarInset | null>(null);

  public formMode = FormMode;

  // Automatically add class to component host for CSS targeting
  @HostBinding('class.form-button-bar-active') 
  get isActive(): boolean {
    const mode = this.mode();
    return mode === FormMode.edit || mode === FormMode.add;
  }
}
