import { Component, Input, OnInit, TemplateRef } from "@angular/core";

@Component({
    selector: "app-form-section",
    templateUrl: "./form-section.component.html",
    styleUrls: ["./form-section.component.scss"],
    standalone: false
})
export class FormSectionComponent implements OnInit {
  @Input() public headerText: string = "";

  @Input() public headerButtonsTemplate?: TemplateRef<any>;

  @Input() public indent: boolean = true;

  @Input() public subtitle: string = "";

  @Input() public collapsed: boolean = false;

  public isCollapsed: boolean = false;

  public ngOnInit(): void {
    this.isCollapsed = this.collapsed;
  }

  public toggleCollapsed(): void {
    this.isCollapsed = !this.isCollapsed;
  }
}
