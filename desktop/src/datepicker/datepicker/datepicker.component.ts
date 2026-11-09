import { Component } from "@angular/core";
import { BaseInputComponent } from "../../base-input";
import { MatFormField, MatLabel, MatInput, MatSuffix, MatError } from "@angular/material/input";
import { ReactiveFormsModule } from "@angular/forms";
import { MatDatepickerInput, MatDatepickerToggle, MatDatepicker } from "@angular/material/datepicker";
import { AsyncPipe } from "@angular/common";

@Component({
    selector: "app-datepicker",
    templateUrl: "./datepicker.component.html",
    styleUrls: ["./datepicker.component.scss"],
    imports: [MatFormField, MatLabel, MatInput, ReactiveFormsModule, MatDatepickerInput, MatDatepickerToggle, MatSuffix, MatDatepicker, MatError, AsyncPipe]
})
export class DatepickerComponent extends BaseInputComponent {}
