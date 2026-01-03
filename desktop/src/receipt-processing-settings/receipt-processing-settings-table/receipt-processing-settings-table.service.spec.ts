import { provideHttpClientTesting } from "@angular/common/http/testing";
import { TestBed } from "@angular/core/testing";
import { NgxsModule } from "@ngxs/store";
import { TableModule } from "../../table/table.module";

import { ReceiptProcessingSettingsTableService } from "./receipt-processing-settings-table.service";
import { provideHttpClient, withInterceptorsFromDi } from "@angular/common/http";

describe("ReceiptProcessingSettingsTableService", () => {
  let service: ReceiptProcessingSettingsTableService;

  beforeEach(() => {
    TestBed.configureTestingModule({
    imports: [NgxsModule.forRoot(), TableModule],
    providers: [provideHttpClient(withInterceptorsFromDi()), provideHttpClientTesting()]
});
    service = TestBed.inject(ReceiptProcessingSettingsTableService);
  });

  it("should be created", () => {
    expect(service).toBeTruthy();
  });
});
