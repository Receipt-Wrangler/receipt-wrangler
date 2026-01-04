import { CUSTOM_ELEMENTS_SCHEMA } from "@angular/core";
import { ComponentFixture, TestBed } from "@angular/core/testing";
import { PrettyJsonPipe } from "../task-table/pretty-json.pipe";

import { PrettyJsonComponent } from "./pretty-json.component";

describe("PrettyJsonComponent", () => {
  let component: PrettyJsonComponent;
  let fixture: ComponentFixture<PrettyJsonComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [PrettyJsonComponent, PrettyJsonPipe],
      providers: [PrettyJsonPipe],
      schemas: [CUSTOM_ELEMENTS_SCHEMA]
    })
      .compileComponents();

    fixture = TestBed.createComponent(PrettyJsonComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it("should create", () => {
    expect(component).toBeTruthy();
  });
});
