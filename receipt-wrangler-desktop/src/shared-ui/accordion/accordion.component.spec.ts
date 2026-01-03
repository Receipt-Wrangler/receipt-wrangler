import {ComponentFixture, TestBed} from '@angular/core/testing';

import {AccordionComponent} from './accordion.component';
import {AccordionModule} from "ngx-bootstrap/accordion";
import {CUSTOM_ELEMENTS_SCHEMA} from "@angular/core";

describe('AccordionComponent', () => {
  let component: AccordionComponent;
  let fixture: ComponentFixture<AccordionComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [AccordionComponent],
      imports: [
        AccordionModule
      ],
      schemas: [CUSTOM_ELEMENTS_SCHEMA]
    })
      .compileComponents();

    fixture = TestBed.createComponent(AccordionComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
