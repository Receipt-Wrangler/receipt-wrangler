import { Component, input } from "@angular/core";
import { AccordionPanel } from "./accordion-panel.interface";
import { MatAccordion, MatExpansionPanel, MatExpansionPanelHeader, MatExpansionPanelTitle, MatExpansionPanelDescription } from "@angular/material/expansion";
import { NgTemplateOutlet } from "@angular/common";

@Component({
    selector: "app-accordion",
    templateUrl: "./accordion.component.html",
    styleUrl: "./accordion.component.scss",
    imports: [MatAccordion, MatExpansionPanel, MatExpansionPanelHeader, MatExpansionPanelTitle, MatExpansionPanelDescription, NgTemplateOutlet]
})
export class AccordionComponent {
  public readonly panels = input<AccordionPanel[]>([]);
}
