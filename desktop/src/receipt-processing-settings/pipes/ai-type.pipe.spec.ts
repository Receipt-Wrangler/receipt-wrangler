import { AiType } from "../../open-api";
import { AiTypePipe } from "./ai-type.pipe";

describe("AiTypePipe", () => {
  it("create an instance", () => {
    const pipe = new AiTypePipe();
    expect(pipe).toBeTruthy();
  });

  it("transform OPEN_AI_CUSTOM ", () => {
    const pipe = new AiTypePipe();
    const result = pipe.transform(AiType.OpenAiCustom);

    expect(result).toEqual("OpenAI Custom");
  });

  it("transform OPEN_AI", () => {
    const pipe = new AiTypePipe();
    const result = pipe.transform(AiType.OpenAi);

    expect(result).toEqual("OpenAI");
  });

  it("transform GEMINI", () => {
    const pipe = new AiTypePipe();
    const result = pipe.transform(AiType.Gemini);

    expect(result).toEqual("Gemini");
  });
});
