import { TestBed } from "@angular/core/testing";
import { ResolveFn } from "@angular/router";
import { Prompt } from "../open-api";

import { promptsResolver } from "./prompts.resolver";

describe("promptsResolver", () => {
  const executeResolver: ResolveFn<Prompt[]> = (...resolverParameters) =>
    TestBed.runInInjectionContext(() => promptsResolver(...resolverParameters));

  beforeEach(() => {
    TestBed.configureTestingModule({});
  });

  it("should be created", () => {
    expect(executeResolver).toBeTruthy();
  });
});
