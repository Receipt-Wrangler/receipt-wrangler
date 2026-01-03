import { Component, OnInit } from "@angular/core";
import { ActivatedRoute } from "@angular/router";
import { FormConfig } from "src/interfaces";
import { FormMode } from "src/enums/form-mode.enum";

@Component({
    selector: "app-api-keys",
    templateUrl: "./api-keys.component.html",
    styleUrls: ["./api-keys.component.scss"],
    standalone: false
})
export class ApiKeysComponent implements OnInit {
  public formConfig!: FormConfig;
  public formMode = FormMode;

  constructor(
    private route: ActivatedRoute
  ) {}

  public ngOnInit(): void {
    this.formConfig = this.route?.snapshot?.data?.["formConfig"];
  }
}