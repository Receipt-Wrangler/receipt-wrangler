import { ComponentFixture, TestBed } from "@angular/core/testing";
import { NgxsModule } from "@ngxs/store";
import {
  ReceiptProcessingSettingsTableService
} from "../../receipt-processing-settings/receipt-processing-settings-table/receipt-processing-settings-table.service";
import { BaseTableService } from "../../services/base-table.service";

import { BaseTableComponent } from "./base-table.component";

describe("BaseTableComponent", () => {
  let component: BaseTableComponent<any>;
  let fixture: ComponentFixture<BaseTableComponent<any>>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [BaseTableComponent],
      imports: [NgxsModule.forRoot([])],
      providers: [
        {
          provide: BaseTableService,
          useValue: ReceiptProcessingSettingsTableService
        }
      ]
    })
      .compileComponents();

    fixture = TestBed.createComponent(BaseTableComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it("should create", () => {
    expect(component).toBeTruthy();
  });
});
