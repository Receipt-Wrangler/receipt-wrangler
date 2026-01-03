import { OcrEngine } from "../../open-api";
import { OcrEnginePipe } from "./ocr-engine.pipe";

describe("OcrEnginePipe", () => {
  it("create an instance", () => {
    const pipe = new OcrEnginePipe();
    expect(pipe).toBeTruthy();
  });

  it("transform TESSERACT", () => {
    const pipe = new OcrEnginePipe();
    const result = pipe.transform(OcrEngine.Tesseract);

    expect(result).toEqual("Tesseract");
  });

  it("transform EASY_OCR", () => {
    const pipe = new OcrEnginePipe();
    const result = pipe.transform(OcrEngine.EasyOcr);

    expect(result).toEqual("EasyOCR");
  });
});
