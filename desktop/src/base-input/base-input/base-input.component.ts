import { Component, OnInit, input, output } from "@angular/core";
import { FormControl } from "@angular/forms";
import { map, Observable, startWith } from "rxjs";
import { InputErrorMessage } from "./input-error-message";

@Component({
    selector: "app-base-input",
    templateUrl: "./base-input.component.html",
    styleUrls: ["./base-input.component.scss"]
})
export class BaseInputComponent implements OnInit {
  public readonly inputFormControl = input<FormControl>(new FormControl());

  public readonly label = input<string>("");

  public readonly additionalErrorMessages = input<{ [key: string]: string }>();

  public readonly readonly = input<boolean>(false);

  public readonly placeholder = input<string>();

  public readonly hint = input<string>();

  public readonly appearance = input<"fill" | "outline">("fill");

  public readonly inputBlur = output<any>();

  public formControlErrors!: Observable<InputErrorMessage[]>;

  public errorMessages: { [key: string]: string } = {};

  public ngOnInit(): void {
    this.errorMessages = {
      required: `${this.label()} is required.`,
      email: `${this.label()} must be a valid email address.`,
      duplicate: `${this.label()} must be unique.`,
      min: `Value must be larger than 0`,
    };

    this.formControlErrors = this.inputFormControl().statusChanges.pipe(
      startWith(this.inputFormControl().status),
      map(() => {
        const errors = this.inputFormControl().errors;
        if (errors) {
          const keys = Object.keys(this.inputFormControl().errors as any);
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

    if (this.additionalErrorMessages()) {
      this.errorMessages = {
        ...this.errorMessages,
        ...this.additionalErrorMessages(),
      };
    }
  }
}
