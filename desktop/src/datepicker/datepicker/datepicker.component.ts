import { Component } from "@angular/core";
import { BaseInputComponent } from "../../base-input";
import { MatFormField, MatLabel, MatInput, MatSuffix, MatError } from "@angular/material/input";
import { ReactiveFormsModule } from "@angular/forms";
import { MatDatepickerInput, MatDatepickerToggle, MatDatepicker } from "@angular/material/datepicker";
import { AsyncPipe } from "@angular/common";
import { provideNativeDateAdapter } from "@angular/material/core";

@Component({
    selector: "app-datepicker",
    templateUrl: "./datepicker.component.html",
    styleUrls: ["./datepicker.component.scss"],
    imports: [MatFormField, MatLabel, MatInput, ReactiveFormsModule, MatDatepickerInput, MatDatepickerToggle, MatSuffix, MatDatepicker, MatError, AsyncPipe],
    // The app's only Material datepicker, so the date adapter is provided here.
    providers: [provideNativeDateAdapter()]
})
export class DatepickerComponent extends BaseInputComponent {}
