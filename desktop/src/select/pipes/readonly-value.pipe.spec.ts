import { TestBed } from "@angular/core/testing";
import { OptionDisplayPipe } from "./option-display.pipe";
import { ReadonlyValuePipe } from "./readonly-value.pipe";

describe("ReadonlyValuePipe", () => {
  let pipe: ReadonlyValuePipe;
  let optionDisplayPipe: OptionDisplayPipe;

  beforeEach(() => {
    TestBed.configureTestingModule({
      declarations: [ReadonlyValuePipe, OptionDisplayPipe],
      providers: [ReadonlyValuePipe, OptionDisplayPipe]
    });

    pipe = TestBed.inject(ReadonlyValuePipe);
    optionDisplayPipe = TestBed.inject(OptionDisplayPipe);
  });

  it("should create an instance", () => {
    expect(pipe).toBeTruthy();
  });

  it("should find option by value and transform it using optionDisplayPipe", () => {
    const options = [
      { id: 1, name: "Option 1" },
      { id: 2, name: "Option 2" },
      { id: 3, name: "Option 3" }
    ];

    jest.spyOn(optionDisplayPipe, "transform").mockReturnValue("Option 2");

    const result = pipe.transform(2, options, "name", "id");

    expect(optionDisplayPipe.transform).toHaveBeenCalledWith(1, options[1], "name", undefined);
    expect(result).toBe("Option 2");
  });

  it("should handle string values correctly", () => {
    const options = [
      { id: "a", name: "Option A" },
      { id: "b", name: "Option B" },
      { id: "c", name: "Option C" }
    ];

    jest.spyOn(optionDisplayPipe, "transform").mockReturnValue("Option B");

    const result = pipe.transform("b", options, "name", "id");

    expect(optionDisplayPipe.transform).toHaveBeenCalledWith(1, options[1], "name", undefined);
    expect(result).toBe("Option B");
  });

  it("should handle options without optionValueKey", () => {
    const options = ["Option 1", "Option 2", "Option 3"];

    jest.spyOn(optionDisplayPipe, "transform").mockReturnValue("Option 2");

    const result = pipe.transform("Option 2", options, "name");

    expect(optionDisplayPipe.transform).toHaveBeenCalledWith(1, options[1], "name", undefined);
    expect(result).toBe("Option 2");
  });

  it("should work with optionsDisplayArray", () => {
    const options = [
      { id: 1, name: "Option 1" },
      { id: 2, name: "Option 2" },
      { id: 3, name: "Option 3" }
    ];
    const optionsDisplayArray = ["Display 1", "Display 2", "Display 3"];

    jest.spyOn(optionDisplayPipe, "transform").mockReturnValue("Display 2");

    const result = pipe.transform(2, options, "name", "id", optionsDisplayArray);

    expect(optionDisplayPipe.transform).toHaveBeenCalledWith(1, options[1], "name", optionsDisplayArray);
    expect(result).toBe("Display 2");
  });

  it("should return null", () => {
    const options = [
      { id: 1, name: "Option 1" },
      { id: 2, name: "Option 2" },
      { id: 3, name: "Option 3" }
    ];
    const optionsDisplayArray = ["Display 1", "Display 2", "Display 3"];

    jest.spyOn(optionDisplayPipe, "transform").mockReturnValue("Display 2");

    const result = pipe.transform(5, options, "name", "id", optionsDisplayArray);

    expect(optionDisplayPipe.transform).toHaveBeenCalledTimes(0);
    expect(result).toBe(null);
  });
});
