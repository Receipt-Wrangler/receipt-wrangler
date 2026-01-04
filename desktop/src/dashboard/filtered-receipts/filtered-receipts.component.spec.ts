import { ScrollingModule } from "@angular/cdk/scrolling";
import { provideHttpClientTesting } from "@angular/common/http/testing";
import { CUSTOM_ELEMENTS_SCHEMA } from "@angular/core";
import { ComponentFixture, TestBed } from "@angular/core/testing";
import { NgxsModule, Store } from "@ngxs/store";
import { of } from "rxjs";
import { ReceiptFilterService } from "src/services/receipt-filter.service";
import { CustomCurrencyPipe } from "../../pipes/custom-currency.pipe";
import { GroupState } from "../../store";
import { FilteredReceiptsComponent } from "./filtered-receipts.component";
import { provideHttpClient, withInterceptorsFromDi } from "@angular/common/http";

describe("FilteredReceiptsComponent", () => {
  let component: FilteredReceiptsComponent;
  let store: Store;
  let fixture: ComponentFixture<FilteredReceiptsComponent>;

  beforeEach(() => {
    TestBed.configureTestingModule({
    declarations: [FilteredReceiptsComponent],
    schemas: [CUSTOM_ELEMENTS_SCHEMA],
    imports: [NgxsModule.forRoot([GroupState]),
        ScrollingModule],
    providers: [CustomCurrencyPipe, provideHttpClient(withInterceptorsFromDi()), provideHttpClientTesting()]
});
    store = TestBed.inject(Store);
    fixture = TestBed.createComponent(FilteredReceiptsComponent);
    component = fixture.componentInstance;
    component.widget = {} as any;
    fixture.detectChanges();
  });

  it("should create", () => {
    expect(component).toBeTruthy();
  });

  it("should get data", () => {
    const serviceSpy = jest.spyOn(
      TestBed.inject(ReceiptFilterService),
      "getPagedReceiptsForGroups"
    ).mockReturnValue(of([] as any));
    store.reset({
      groups: {
        selectedGroupId: "1",
      },
    });
    component.groupId = 1;

    component.ngOnInit();

    expect(serviceSpy).toHaveBeenCalledWith(
      "1",
      undefined,
      undefined,
      undefined,
      undefined,
      {
        page: 1,
        pageSize: 25,
        orderBy: "date",
        sortDirection: "desc",
        filter: undefined,
      }
    );
  });

  it("should get next page of data", () => {
    const serviceSpy = jest.spyOn(
      TestBed.inject(ReceiptFilterService),
      "getPagedReceiptsForGroups"
    ).mockReturnValue(of({ data: [{} as any] } as any));

    store.reset({
      groups: {
        selectedGroupId: "1",
      },
    });
    component.groupId = 1;

    component.receipts = [{} as any, {} as any, {} as any, {} as any];
    component.endOfListReached();

    expect(serviceSpy).toHaveBeenCalledWith(
      "1",
      undefined,
      undefined,
      undefined,
      undefined,
      {
        page: 2,
        pageSize: 25,
        orderBy: "date",
        sortDirection: "desc",
        filter: undefined,
      }
    );
    expect(component.receipts).toEqual([
      {} as any,
      {} as any,
      {} as any,
      {} as any,
      {} as any,
    ]);
  });

  it("should not get next page of data", () => {
    const serviceSpy = jest.spyOn(
      TestBed.inject(ReceiptFilterService),
      "getPagedReceiptsForGroups"
    ).mockReturnValue(of({ data: [{} as any] } as any));

    store.reset({
      groups: {
        selectedGroupId: "1",
      },
    });

    component.receipts = [{} as any, {} as any, {} as any, {} as any];
    component.endOfListReached();

    expect(serviceSpy).toHaveBeenCalledTimes(0);
    expect(component.receipts).toEqual([
      {} as any,
      {} as any,
      {} as any,
      {} as any,
    ]);
  });
});
