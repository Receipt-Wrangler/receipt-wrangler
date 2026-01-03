import { CUSTOM_ELEMENTS_SCHEMA } from "@angular/core";
import { ComponentFixture, TestBed } from "@angular/core/testing";
import { ActivatedRoute } from "@angular/router";
import { AccordionModule } from "ngx-bootstrap/accordion";

import {
  ReceiptProcessingSettingsChildSystemTaskAccordionComponent
} from "./receipt-processing-settings-child-system-task-accordion.component";

describe("ReceiptProcessingSettingsChildSystemTaskAccordionComponent", () => {
  let component: ReceiptProcessingSettingsChildSystemTaskAccordionComponent;
  let fixture: ComponentFixture<ReceiptProcessingSettingsChildSystemTaskAccordionComponent>;
  
  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ReceiptProcessingSettingsChildSystemTaskAccordionComponent],
      providers: [
        {
          provide: ActivatedRoute,
          useValue: {
            snapshot: {
              data: {
                prompts: [],
              }
            }
          }
        }
      ],
      imports: [AccordionModule],
      schemas: [CUSTOM_ELEMENTS_SCHEMA]
    })
      .compileComponents();

    fixture = TestBed.createComponent(ReceiptProcessingSettingsChildSystemTaskAccordionComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it("should create", () => {
    expect(component).toBeTruthy();
  });
});
