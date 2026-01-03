import { OperationsPipe } from "./operations.pipe";

describe("OperationsPipe", () => {
  it("create an instance", () => {
    const pipe = new OperationsPipe();
    expect(pipe).toBeTruthy();
  });

  it("return empty array when the value is not recognized", () => {
    const pipe = new OperationsPipe();

    const result = pipe.transform("bad value", true);

    expect(result).toEqual([]);
  });

  it("should return options for date", () => {
    const pipe = new OperationsPipe();

    const result = pipe.transform("date", true);

    expect(result).toEqual([
      "Equals",
      "Greater than",
      "Less than",
      "Between",
      "Within current month"
    ]);
  });

  it("should return options for text", () => {
    const pipe = new OperationsPipe();

    const result = pipe.transform("text", true);

    expect(result).toEqual(["Contains", "Equals"]);
  });

  it("should return options for number", () => {
    const pipe = new OperationsPipe();

    const result = pipe.transform("number", true);

    expect(result).toEqual([
      "Equals",
      "Greater than",
      "Less than",
      "Between"
    ]);
  });

  it("should return options for list", () => {
    const pipe = new OperationsPipe();

    const result = pipe.transform("list", true);

    expect(result).toEqual(["Contains"]);
  });

  it("should return options for users", () => {
    const pipe = new OperationsPipe();

    const result = pipe.transform("users", true);

    expect(result).toEqual(["Contains"]);
  });

  it("should return value options for users", () => {
    const pipe = new OperationsPipe();

    const result = pipe.transform("users", false);

    expect(result).toEqual(["CONTAINS"]);
  });

  it("should return value options for list", () => {
    const pipe = new OperationsPipe();

    const result = pipe.transform("list", false);

    expect(result).toEqual(["CONTAINS"]);
  });

  it("should return value options for number", () => {
    const pipe = new OperationsPipe();

    const result = pipe.transform("number", false);

    expect(result).toEqual([
      "EQUALS",
      "GREATER_THAN",
      "LESS_THAN",
      "BETWEEN"
    ]);
  });

  it("should return value options for text", () => {
    const pipe = new OperationsPipe();

    const result = pipe.transform("text", false);

    expect(result).toEqual(["CONTAINS", "EQUALS"]);
  });

  it("should return value options for date", () => {
    const pipe = new OperationsPipe();

    const result = pipe.transform("date", false);

    expect(result).toEqual([
      "EQUALS",
      "GREATER_THAN",
      "LESS_THAN",
      "BETWEEN",
      "WITHIN_CURRENT_MONTH"
    ]);
  });
});
