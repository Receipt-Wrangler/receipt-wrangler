import { SystemTaskType } from "../open-api";
import { SYSTEM_TASK_TYPE_OPTIONS } from "./system-task-type-options";

describe("SYSTEM_TASK_TYPE_OPTIONS", () => {
  const offered = SYSTEM_TASK_TYPE_OPTIONS.map((option) => option.value);

  it("offers Receipt Uploaded, which a manual create lists at the top level", () => {
    expect(offered).toContain(SystemTaskType.ReceiptUploaded);
  });

  it("omits the types only ever recorded as children", () => {
    expect(offered).not.toContain(SystemTaskType.ChatCompletion);
    expect(offered).not.toContain(SystemTaskType.OcrProcessing);
  });
});
