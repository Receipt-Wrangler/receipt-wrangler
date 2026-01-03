import { provideHttpClient, withInterceptorsFromDi } from "@angular/common/http";
import { provideHttpClientTesting } from "@angular/common/http/testing";
import { TestBed } from "@angular/core/testing";

import { ReceiptExportService } from "./receipt-export.service";

describe("ReceiptExportService", () => {
  let service: ReceiptExportService;

  beforeEach(() => {
    TestBed.configureTestingModule({
      providers: [
        provideHttpClient(withInterceptorsFromDi()),
        provideHttpClientTesting(),
      ]
    });
    service = TestBed.inject(ReceiptExportService);
  });

  it("should be created", () => {
    expect(service).toBeTruthy();
  });
});
