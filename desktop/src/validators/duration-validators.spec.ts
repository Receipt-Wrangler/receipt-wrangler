import { FormControl } from "@angular/forms";
import { durationValueValidator } from "./duration-validators";

describe("durationValueValidator", () => {
  const validate = (value: any, max = 720) =>
    durationValueValidator(max)(new FormControl(value));

  it("accepts whole numbers within range", () => {
    expect(validate(1)).toBeNull();
    expect(validate(24)).toBeNull();
    expect(validate(720)).toBeNull();
    expect(validate(30, 30)).toBeNull();
  });

  // The API field is a Go int, so a decimal fails json.Unmarshal rather than
  // returning a field-level error.
  it("rejects fractional values", () => {
    expect(validate(1.5)?.["duration"]).toBe("Must be a whole number.");
    expect(validate("2.5")?.["duration"]).toBe("Must be a whole number.");
  });

  it("rejects values below 1", () => {
    expect(validate(0)?.["duration"]).toBe("Must be at least 1.");
    expect(validate(-4)?.["duration"]).toBe("Must be at least 1.");
  });

  it("rejects values above the max, naming the max in the message", () => {
    expect(validate(721)?.["duration"]).toBe("Must be at most 720.");
    expect(validate(31, 30)?.["duration"]).toBe("Must be at most 30.");
  });

  it("rejects non-numeric input", () => {
    expect(validate("abc")?.["duration"]).toBe("Must be a whole number.");
  });

  // Emptiness belongs to Validators.required, so this validator stays quiet and
  // the input renders the standard "<label> is required." message instead.
  it("leaves emptiness to Validators.required", () => {
    expect(validate(null)).toBeNull();
    expect(validate(undefined)).toBeNull();
    expect(validate("")).toBeNull();
  });
});
