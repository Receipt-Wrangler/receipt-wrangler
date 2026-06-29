import { provideHttpClient, withInterceptorsFromDi } from "@angular/common/http";
import { provideHttpClientTesting } from "@angular/common/http/testing";
import { TestBed } from "@angular/core/testing";
import { ResolveFn } from "@angular/router";
import { NgxsModule, Store } from "@ngxs/store";
import { Observable, of } from "rxjs";
import { ApiModule, CustomField, CustomFieldService, Permission } from "../open-api";
import { AuthState } from "../store";
import { SetPermissions } from "../store/auth.state.actions";
import { customFieldResolverFn } from "./custom-field.resolver";

describe("customFieldResolver", () => {
  let store: Store;
  let customFieldService: CustomFieldService;

  const executeResolver: ResolveFn<Observable<CustomField[]>> = (...resolverParameters) =>
    TestBed.runInInjectionContext(() => customFieldResolverFn(...resolverParameters));

  beforeEach(() => {
    TestBed.configureTestingModule({
      imports: [ApiModule, NgxsModule.forRoot([AuthState])],
      providers: [provideHttpClient(withInterceptorsFromDi()), provideHttpClientTesting()],
    });
    store = TestBed.inject(Store);
    customFieldService = TestBed.inject(CustomFieldService);
  });

  it("should be created", () => {
    expect(executeResolver).toBeTruthy();
  });

  it("does not call the catalog without app.custom-fields.read", (done) => {
    const serviceSpy = jest.spyOn(customFieldService, "getPagedCustomFields");
    store.dispatch(new SetPermissions([], {}));

    (executeResolver({} as any, {} as any) as Observable<CustomField[]>).subscribe((result) => {
      expect(result).toEqual([]);
      expect(serviceSpy).not.toHaveBeenCalled();
      done();
    });
  });

  it("calls the catalog with app.custom-fields.read", () => {
    const serviceSpy = jest
      .spyOn(customFieldService, "getPagedCustomFields")
      .mockReturnValue(of({ data: [] }) as any);
    store.dispatch(new SetPermissions([Permission.AppCustomFieldsRead], {}));

    (executeResolver({} as any, {} as any) as Observable<CustomField[]>).subscribe();

    expect(serviceSpy).toHaveBeenCalled();
  });
});
