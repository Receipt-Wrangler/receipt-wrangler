import { provideHttpClientTesting } from "@angular/common/http/testing";
import { TestBed } from "@angular/core/testing";
import { ResolveFn } from "@angular/router";
import { NgxsModule } from "@ngxs/store";
import { Observable } from "rxjs";
import { ApiModule, Receipt, ReceiptService } from "../open-api";
import { receiptResolverFn } from "./receipt.resolver";
import { provideHttpClient, withInterceptorsFromDi } from "@angular/common/http";

describe("ReceiptResolverService", () => {
  const executeResolver: ResolveFn<Observable<Receipt>> = (
    ...resolverParameters
  ) =>
    TestBed.runInInjectionContext(() =>
      receiptResolverFn(...resolverParameters)
    );


  beforeEach(() => {
    TestBed.configureTestingModule({
    imports: [ApiModule, NgxsModule.forRoot([])],
    providers: [provideHttpClient(withInterceptorsFromDi()), provideHttpClientTesting()]
});
  });

  it("should call receipt service", () => {
    const serviceSpy = jest.spyOn(TestBed.inject(ReceiptService), "getReceiptById");
    executeResolver({ params: { id: 1 } } as any, {} as any);
    expect(serviceSpy).toHaveBeenCalledWith(1);

  });

});
