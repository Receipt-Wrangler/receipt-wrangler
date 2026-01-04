import { TestBed } from "@angular/core/testing";
import { PrettyJsonPipe } from "./pretty-json.pipe";

describe("PrettyJsonPipe", () => {
  let pipe: PrettyJsonPipe;

  beforeEach(() => {
    TestBed.configureTestingModule({
      declarations: [PrettyJsonPipe],
      imports: [],
      providers: [PrettyJsonPipe]
    });

    pipe = TestBed.inject(PrettyJsonPipe);
  });

  it("create an instance", () => {
    expect(pipe).toBeTruthy();
  });

  it("should return empty string when value is not json", () => {
    const result = pipe.transform("not json");
    expect(result).toEqual("");
  });

  it("should return sanitized html when value is json", () => {
    const json = {
      key: "value"
    };
    const jsonString = JSON.stringify(json);
    const result = pipe.transform(jsonString);
    expect(result === "").toBe(false);
  });
});
