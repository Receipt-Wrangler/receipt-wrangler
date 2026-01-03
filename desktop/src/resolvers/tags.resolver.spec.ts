import { provideHttpClientTesting } from "@angular/common/http/testing";
import { TestBed } from "@angular/core/testing";
import { ResolveFn } from "@angular/router";
import { ApiModule, Tag, TagService } from "../open-api";
import { tagResolverFn } from "./tags.resolver";
import { provideHttpClient, withInterceptorsFromDi } from "@angular/common/http";

describe("TagsResolverService", () => {
  const executeResolver: ResolveFn<Tag[]> = (
    ...resolverParameters
  ) =>
    TestBed.runInInjectionContext(() =>
      tagResolverFn(...resolverParameters)
    );

  beforeEach(() => {
    TestBed.configureTestingModule({
    imports: [ApiModule],
    providers: [provideHttpClient(withInterceptorsFromDi()), provideHttpClientTesting()]
});
  });

  it("should call tag service", () => {
    const serviceSpy = jest.spyOn(TestBed.inject(TagService), "getAllTags");
    executeResolver({} as any, {} as any);
    expect(serviceSpy).toHaveBeenCalled();
  });
});
