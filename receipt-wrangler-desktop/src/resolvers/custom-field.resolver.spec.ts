import { TestBed } from "@angular/core/testing";
import { ResolveFn } from "@angular/router";
import { Observable } from "rxjs";
import { CustomField } from "../open-api/index";
import { customFieldResolverFn } from "./custom-field.resolver";


describe("customFieldResolver", () => {
  const executeResolver: ResolveFn<Observable<CustomField[]>> = (...resolverParameters) =>
    TestBed.runInInjectionContext(() => customFieldResolverFn(...resolverParameters));

  beforeEach(() => {
    TestBed.configureTestingModule({});
  });

  it("should be created", () => {
    expect(executeResolver).toBeTruthy();
  });
});
