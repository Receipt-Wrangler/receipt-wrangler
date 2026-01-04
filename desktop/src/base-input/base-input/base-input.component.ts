import { Component, EventEmitter, Input, OnInit, Output } from "@angular/core";
import { FormControl } from "@angular/forms";
import { map, Observable, startWith } from "rxjs";
import { BaseInputInterface } from "../base-input.interface";
import { InputErrorMessage } from "./input-error-message";

@Component({
    selector: "app-base-input",
    templateUrl: "./base-input.component.html",
    styleUrls: ["./base-input.component.scss"],
    standalone: false
})
export class BaseInputComponent implements OnInit, BaseInputInterface {
  @Input() public inputFormControl: FormControl = new FormControl();

  @Input() public label: string = "";

  @Input() public additionalErrorMessages?: { [key: string]: string };

  @Input() public readonly: boolean = false;

  @Input() public placeholder?: string;

  @Input() public hint?: string;

  @Input() public appearance: "fill" | "outline" = "fill";

  @Output() public inputBlur: EventEmitter<any> = new EventEmitter<any>(
    undefined
  );

  public formControlErrors!: Observable<InputErrorMessage[]>;

  public errorMessages: { [key: string]: string } = {};

  public ngOnInit(): void {
    this.errorMessages = {
      required: `${this.label} is required.`,
      email: `${this.label} must be a valid email address.`,
      duplicate: `${this.label} must be unique.`,
      min: `Value must be larger than 0`,
    };

    this.formControlErrors = this.inputFormControl.statusChanges.pipe(
      startWith(this.inputFormControl.status),
      map(() => {
        const errors = this.inputFormControl.errors;
        if (errors) {
          const keys = Object.keys(this.inputFormControl.errors as any);
          return keys.map((k: string) => {
            const value = errors[k];
            let message = "";

            if (typeof value === "string") {
              message = value;
            } else if (this.errorMessages[k]) {
              message = this.errorMessages[k];
            }

            return {
              error: k as string,
              message: message,
            };
          });
        } else {
          return [];
        }
      })
    );

    if (this.additionalErrorMessages) {
      this.errorMessages = {
        ...this.errorMessages,
        ...this.additionalErrorMessages,
      };
    }
  }
}
