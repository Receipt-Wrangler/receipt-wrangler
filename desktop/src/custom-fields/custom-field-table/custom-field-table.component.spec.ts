import { AsyncPipe } from "@angular/common";
import { provideHttpClient } from "@angular/common/http";
import { provideHttpClientTesting } from "@angular/common/http/testing";
import { ComponentFixture, TestBed } from "@angular/core/testing";
import { provideNoopAnimations } from "@angular/platform-browser/animations";
import { ActivatedRoute } from "@angular/router";
import { NgxsModule } from "@ngxs/store";
import { SharedUiModule } from "../../shared-ui/shared-ui.module";
import { CustomFieldTableState } from "../../store/custom-field-table.state";
import { TableModule } from "../../table/table.module";

import { CustomFieldTableComponent } from "./custom-field-table.component";

describe("CustomFieldTableComponent", () => {
  let component: CustomFieldTableComponent;
  let fixture: ComponentFixture<CustomFieldTableComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [CustomFieldTableComponent],
      imports: [NgxsModule.forRoot([CustomFieldTableState]), AsyncPipe, SharedUiModule, TableModule],
      providers: [
        provideHttpClient(),
        provideHttpClientTesting(),
        provideNoopAnimations(),
        {
          provide: ActivatedRoute,
          useValue: {}
        }
      ]
    })
      .compileComponents();

    fixture = TestBed.createComponent(CustomFieldTableComponent);
    component = fixture.componentInstance;
  });

  it("should create", () => {
    expect(component).toBeTruthy();
  });
});
