import { Component, provideZonelessChangeDetection, signal } from "@angular/core";
import { ComponentFixture, TestBed } from "@angular/core/testing";
import { CarouselSlideComponent } from "./carousel-slide.component";
import { SlideCarouselComponent } from "./slide-carousel.component";

@Component({
  imports: [SlideCarouselComponent, CarouselSlideComponent],
  template: `
    <app-slide-carousel
      [(activeSlide)]="index"
      [showIndicators]="showIndicators()"
      [showControls]="showControls()"
    >
      @for (item of items(); track item) {
        <app-carousel-slide data-testid="slide">
          <input [attr.data-testid]="'field-' + item" />
        </app-carousel-slide>
      }
    </app-slide-carousel>
  `,
})
class HostComponent {
  public readonly index = signal(0);
  public readonly items = signal(["a", "b", "c"]);
  public readonly showIndicators = signal(true);
  public readonly showControls = signal(true);
}

describe("SlideCarouselComponent", () => {
  let fixture: ComponentFixture<HostComponent>;
  let host: HostComponent;

  const query = <T extends HTMLElement = HTMLElement>(selector: string): T | null =>
    fixture.nativeElement.querySelector(selector);

  const queryAll = (selector: string): HTMLElement[] =>
    Array.from(fixture.nativeElement.querySelectorAll(selector));

  const slides = () => queryAll('[data-testid="slide"]');
  const indicators = () => queryAll('[data-testid="carousel-indicator"]');
  const activeIndex = () => slides().findIndex((slide) => slide.classList.contains("active"));
  const container = () => query(".carousel")!;

  const click = async (selector: string) => {
    query(selector)!.click();
    await fixture.whenStable();
  };

  const press = async (target: HTMLElement, key: string): Promise<KeyboardEvent> => {
    const event = new KeyboardEvent("keydown", { key, bubbles: true, cancelable: true });
    target.dispatchEvent(event);
    await fixture.whenStable();
    return event;
  };

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [HostComponent],
      providers: [provideZonelessChangeDetection()],
    }).compileComponents();

    fixture = TestBed.createComponent(HostComponent);
    host = fixture.componentInstance;
    await fixture.whenStable();
  });

  it("shows the first slide and hides the rest from assistive technology", () => {
    expect(activeIndex()).toBe(0);
    expect(slides().map((slide) => slide.getAttribute("aria-hidden"))).toEqual([
      "false",
      "true",
      "true",
    ]);
  });

  it("wraps forwards and backwards, and reports the index", async () => {
    await click('[data-testid="carousel-next"]');
    await click('[data-testid="carousel-next"]');
    expect(activeIndex()).toBe(2);

    await click('[data-testid="carousel-next"]');
    expect(activeIndex()).toBe(0);
    expect(host.index()).toBe(0);

    await click('[data-testid="carousel-previous"]');
    expect(activeIndex()).toBe(2);
    expect(host.index()).toBe(2);
  });

  it("selects a slide from its indicator, and marks only that one current", async () => {
    indicators()[1].click();
    await fixture.whenStable();

    expect(activeIndex()).toBe(1);
    expect(host.index()).toBe(1);
    expect(indicators().map((indicator) => indicator.getAttribute("aria-current"))).toEqual([
      null,
      "true",
      null,
    ]);
  });

  it("follows the bound index", async () => {
    host.index.set(2);
    await fixture.whenStable();

    expect(activeIndex()).toBe(2);
  });

  it("renders no indicators or controls for a single slide", async () => {
    host.items.set(["a"]);
    await fixture.whenStable();

    expect(indicators()).toHaveLength(0);
    expect(query('[data-testid="carousel-next"]')).toBeNull();
    expect(query('[data-testid="carousel-previous"]')).toBeNull();
  });

  it("leaves out indicators and controls that are switched off", async () => {
    host.showIndicators.set(false);
    host.showControls.set(false);
    await fixture.whenStable();

    expect(indicators()).toHaveLength(0);
    expect(query('[data-testid="carousel-next"]')).toBeNull();
    expect(slides()).toHaveLength(3);
  });

  it("navigates with the keyboard while the carousel itself has focus", async () => {
    await press(container(), "ArrowRight");
    expect(activeIndex()).toBe(1);

    await press(container(), "ArrowLeft");
    expect(activeIndex()).toBe(0);

    const enter = await press(container(), "Enter");
    expect(activeIndex()).toBe(1);
    expect(enter.defaultPrevented).toBe(true);

    const space = await press(container(), " ");
    expect(activeIndex()).toBe(2);
    expect(space.defaultPrevented).toBe(true);
  });

  it("leaves keys typed into a field inside a slide alone", async () => {
    const field = query('[data-testid="field-a"]')!;

    for (const key of [" ", "Enter", "ArrowRight", "ArrowLeft"]) {
      const event = await press(field, key);
      expect(event.defaultPrevented).toBe(false);
    }
    expect(activeIndex()).toBe(0);
    expect(host.index()).toBe(0);
  });

  it("shows the next slide when the shown one is removed", async () => {
    host.index.set(1);
    await fixture.whenStable();

    host.items.set(["a", "c"]);
    await fixture.whenStable();

    expect(host.index()).toBe(1);
    expect(activeIndex()).toBe(1);
    expect(slides()[1].querySelector('[data-testid="field-c"]')).not.toBeNull();
  });

  it("falls back to the first slide, and tells the owner, when the last one shown is removed", async () => {
    host.index.set(2);
    await fixture.whenStable();

    host.items.set(["a", "b"]);
    await fixture.whenStable();

    expect(host.index()).toBe(0);
    expect(activeIndex()).toBe(0);
  });
});
