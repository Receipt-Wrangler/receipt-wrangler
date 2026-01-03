import { FormOption } from "../../interfaces/form-option.interface";
import { AiType } from "../../open-api";

export const aiTypeOptions: FormOption[] = [
  {
    value: AiType.OpenAi,
    displayValue: "OpenAI"
  },
  {
    value: AiType.OpenAiCustom,
    displayValue: "OpenAI Custom",
  },
  {
    value: AiType.Gemini,
    displayValue: "Gemini",
  },
  {
    value: AiType.Ollama,
    displayValue: "Ollama",
  }
];


