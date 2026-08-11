import { Pipe, PipeTransform } from "@angular/core";
import { AiType } from "../../open-api/index";

@Pipe({
  name: "urlHint"
})
export class UrlHintPipe implements PipeTransform {
  public transform(aiType: AiType): string {
    if (aiType === AiType.OpenAiCustom) {
      return "Full OpenAI-compatible base url, including any path. /chat/completions is appended automatically. Azure endpoints must include their path, e.g. https://your-resource.services.ai.azure.com/openai/v1 — a bare origin will not work.";
    }

    if (aiType === AiType.Ollama) {
      return "Full url of the Ollama chat endpoint, used exactly as entered, e.g. http://localhost:11434/api/chat.";
    }

    return "";
  }
}
