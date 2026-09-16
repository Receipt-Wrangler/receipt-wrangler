import { provideHttpClientTesting } from "@angular/common/http/testing";
import { TestBed } from "@angular/core/testing";
import { NgxsModule, Store } from "@ngxs/store";
import { of } from "rxjs";
import { ApiModule, FilterOperation, ReceiptService, ReceiptSummary } from "../open-api";
import { buildDefaultReceiptFilter, ReceiptTableState } from "../store/receipt-table.state";
import { SetReceiptFilterData } from "../store/receipt-table.actions";
import { ReceiptFilterService } from "./receipt-filter.service";
import { provideHttpClient, withInterceptorsFromDi } from "@angular/common/http";

describe("ReceiptFilterService", () => {
  let service: ReceiptFilterService;
  let store: Store;
  let receiptService: ReceiptService;

  beforeEach(() => {
    TestBed.configureTestingModule({
    imports: [ApiModule, NgxsModule.forRoot([ReceiptTableState])],
    providers: [provideHttpClient(withInterceptorsFromDi()), provideHttpClientTesting()]
});
    service = TestBed.inject(ReceiptFilterService);
    store = TestBed.inject(Store);
    receiptService = TestBed.inject(ReceiptService);
  });

  it("should be created", () => {
    expect(service).toBeTruthy();
  });

  describe("getReceiptSummaryForGroup", () => {
    /**
     * The summary's whole promise is that it describes the CURRENT filter. This adapter is the only
     * place that promise is kept - it reads the filter from the store rather than taking one - so a
     * regression here would silently report totals for the whole group while the table showed a
     * filtered page. The component specs mock this method, so nothing else covers it.
     */
    function seedStatusFilter(): void {
      const filter = buildDefaultReceiptFilter();
      filter.status = { operation: FilterOperation.Equals, value: ["OPEN"] };

      store.dispatch(
        new SetReceiptFilterData({
          page: 3,
          pageSize: 25,
          orderBy: "name",
          sortDirection: "asc" as any,
          filter,
        } as any)
      );
    }

    it("forwards the active filter from the store", () => {
      const spy = jest
        .spyOn(receiptService, "getReceiptSummaryForGroup")
        .mockReturnValue(of({} as ReceiptSummary) as any);

      seedStatusFilter();
      service.getReceiptSummaryForGroup("1", 1).subscribe();

      const command = spy.mock.calls[0][1] as any;
      expect(command.filter.status.value).toEqual(["OPEN"]);
      expect(command.filter.status.operation).toBe(FilterOperation.Equals);
    });

    it("converts the group id to a number", () => {
      const spy = jest
        .spyOn(receiptService, "getReceiptSummaryForGroup")
        .mockReturnValue(of({} as ReceiptSummary) as any);

      service.getReceiptSummaryForGroup("7", 7).subscribe();

      // The caller holds the group id as a route string; the generated client types it numeric, and
      // a string would serialize into the url as-is rather than failing loudly.
      expect(spy.mock.calls[0][0]).toBe(7);
      expect(typeof spy.mock.calls[0][0]).toBe("number");
    });

    it("preserves a configurationGroupId that differs from the group id", () => {
      const spy = jest
        .spyOn(receiptService, "getReceiptSummaryForGroup")
        .mockReturnValue(of({} as ReceiptSummary) as any);

      // The All-group case: the data spans every group, the configuration comes from one of them.
      service.getReceiptSummaryForGroup("1", 42).subscribe();

      expect((spy.mock.calls[0][1] as any).configurationGroupId).toBe(42);
    });

    it("returns the summary from the client untouched", (done) => {
      const summary = { enabled: true, configurationGroupId: 1 } as ReceiptSummary;
      jest
        .spyOn(receiptService, "getReceiptSummaryForGroup")
        .mockReturnValue(of(summary) as any);

      service.getReceiptSummaryForGroup("1", 1).subscribe((result) => {
        expect(result).toBe(summary);
        done();
      });
    });
  });
});
