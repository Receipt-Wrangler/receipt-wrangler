import { provideHttpClientTesting } from "@angular/common/http/testing";
import { TestBed } from "@angular/core/testing";
import { NgxsModule } from "@ngxs/store";
import { ApiModule } from "../open-api";
import { ReceiptFilterService } from "./receipt-filter.service";
import { provideHttpClient, withInterceptorsFromDi } from "@angular/common/http";

describe("ReceiptFilterService", () => {
  let service: ReceiptFilterService;

  beforeEach(() => {
    TestBed.configureTestingModule({
    imports: [ApiModule, NgxsModule.forRoot([])],
    providers: [provideHttpClient(withInterceptorsFromDi()), provideHttpClientTesting()]
});
    service = TestBed.inject(ReceiptFilterService);
  });

  it("should be created", () => {
    expect(service).toBeTruthy();
  });
});
