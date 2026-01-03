import { TestBed } from "@angular/core/testing";
import { ResolveFn } from "@angular/router";
import { ReceiptProcessingSettings } from "../open-api";

import { receiptProcessingSettingsResolver } from "./receipt-processing-settings.resolver";

describe("receiptProcessingSettingsResolver", () => {
  const executeResolver: ResolveFn<ReceiptProcessingSettings> = (...resolverParameters) =>
    TestBed.runInInjectionContext(() => receiptProcessingSettingsResolver(...resolverParameters));

  beforeEach(() => {
    TestBed.configureTestingModule({});
  });

  it("should be created", () => {
    expect(executeResolver).toBeTruthy();
  });
});
