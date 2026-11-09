import { ComponentFixture, TestBed } from "@angular/core/testing";
import { provideZonelessChangeDetection } from "@angular/core";
import { ButtonModule } from "../../button";
import { CopyButtonComponent } from "./copy-button.component";

describe("CopyButtonComponent", () => {
  let component: CopyButtonComponent;
  let fixture: ComponentFixture<CopyButtonComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [ButtonModule, CopyButtonComponent],
      providers: [provideZonelessChangeDetection()],
    }).compileComponents();

    fixture = TestBed.createComponent(CopyButtonComponent);
    component = fixture.componentInstance;
    fixture.componentRef.setInput("text", "hello world");
  });

  it("creates", () => {
    expect(component).toBeTruthy();
  });

  it("writes text to clipboard and emits copied on success", async () => {
    const writeText = jest.fn<Promise<void>, [string]>().mockResolvedValue(undefined);
    const originalClipboard = (navigator as any).clipboard;
    (navigator as any).clipboard = { writeText };

    const copiedSpy = jest.fn();
    component.copied.subscribe(copiedSpy);

    try {
      await component.copy();

      expect(writeText).toHaveBeenCalledWith("hello world");
      expect(copiedSpy).toHaveBeenCalledWith("hello world");
      expect(component["isCopied"]()).toBe(true);
    } finally {
      (navigator as any).clipboard = originalClipboard;
    }
  });

  it("does not emit copied when writeText rejects", async () => {
    const writeText = jest.fn<Promise<void>, [string]>().mockRejectedValue(new Error("blocked"));
    const originalClipboard = (navigator as any).clipboard;
    (navigator as any).clipboard = { writeText };

    const copiedSpy = jest.fn();
    component.copied.subscribe(copiedSpy);

    try {
      await component.copy();

      expect(writeText).toHaveBeenCalled();
      expect(copiedSpy).not.toHaveBeenCalled();
      expect(component["isCopied"]()).toBe(false);
    } finally {
      (navigator as any).clipboard = originalClipboard;
    }
  });

  it("is a no-op in non-secure contexts where clipboard is unavailable", async () => {
    const originalClipboard = (navigator as any).clipboard;
    (navigator as any).clipboard = undefined;

    const copiedSpy = jest.fn();
    component.copied.subscribe(copiedSpy);

    try {
      await component.copy();
      expect(copiedSpy).not.toHaveBeenCalled();
      expect(component["isCopied"]()).toBe(false);
    } finally {
      (navigator as any).clipboard = originalClipboard;
    }
  });
});
