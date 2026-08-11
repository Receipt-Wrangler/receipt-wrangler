import { AiType } from "../../open-api/index";
import { UrlHintPipe } from "./url-hint.pipe";

describe("UrlHintPipe", () => {
  let pipe: UrlHintPipe;

  beforeEach(() => {
    pipe = new UrlHintPipe();
  });

  it("create an instance", () => {
    expect(pipe).toBeTruthy();
  });

  it("tells OpenAI-custom users to include the endpoint's path", () => {
    const hint = pipe.transform(AiType.OpenAiCustom);

    expect(hint).toContain("/openai/v1");
    expect(hint).toContain("bare origin will not work");
  });

  it("tells Ollama users the url is used as entered", () => {
    const hint = pipe.transform(AiType.Ollama);

    expect(hint).toContain("Ollama");
    expect(hint).toContain("exactly as entered");
  });

  it("returns no hint for ai types that do not take a url", () => {
    expect(pipe.transform(AiType.OpenAi)).toEqual("");
    expect(pipe.transform(AiType.Gemini)).toEqual("");
  });
});
