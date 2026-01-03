import { provideHttpClientTesting } from "@angular/common/http/testing";
import { TestBed } from "@angular/core/testing";
import { CanActivateFn } from "@angular/router";
import { NgxsModule, Store } from "@ngxs/store";
import { Observable, of, take, tap } from "rxjs";
import { ApiModule, GroupRole, ReceiptService } from "../open-api";
import { receiptGuardGuard } from "./receipt-guard.guard";
import { provideHttpClient, withInterceptorsFromDi } from "@angular/common/http";

describe("receiptGuardGuard", () => {
  const executeGuard: CanActivateFn = (...guardParameters) =>
    TestBed.runInInjectionContext(() => receiptGuardGuard(...guardParameters));
  let store: Store;
  let receiptService: ReceiptService;

  beforeEach(() => {
    TestBed.configureTestingModule({
    imports: [ApiModule, NgxsModule.forRoot([])],
    providers: [provideHttpClient(withInterceptorsFromDi()), provideHttpClientTesting()]
});

    store = TestBed.inject(Store);
    receiptService = TestBed.inject(ReceiptService);
  });

  it("should be created", () => {
    expect(executeGuard).toBeTruthy();
  });

  it("should call service with the correct arguments", () => {
    const spy = jest.spyOn(receiptService, "hasAccessToReceipt").mockReturnValue(
      of(false) as any
    );
    const route: any = {
      data: {
        groupRole: GroupRole.Viewer,
      },
      params: {
        id: 1,
      },
    };

    executeGuard(route, {} as any);

    expect(spy).toHaveBeenCalledWith(1, GroupRole.Viewer);
  });

  it("should allow the user through", (done) => {
    jest.spyOn(receiptService, "hasAccessToReceipt").mockReturnValue(
      of(true) as any
    );
    const route: any = {
      data: {
        groupRole: GroupRole.Viewer,
      },
      params: {
        id: 1,
      },
    };

    (executeGuard(route, {} as any) as Observable<boolean>)
      .pipe(
        take(1),
        tap((result) => {
          expect(result).toEqual(true);
          done();
        })
      )
      .subscribe();
  });

  // TODO: Fix this test
  // it('should not allow the user through', (done) => {
  //   jest.spyOn(receiptService, 'hasAccessToReceipt').mockReturnValue(
  //     of(throwError('')) as any
  //   );
  //   const route: any = {
  //     data: {
  //       groupRole: GroupRole.Viewer,
  //     },
  //     params: {
  //       id: 1,
  //     },
  //   };

  //   (executeGuard(route, {} as any) as Observable<boolean>)
  //     .pipe(take(1))
  //     .subscribe((result) => {
  //       expect(result).toEqual(false);
  //       done();
  //     });
  // });
});
