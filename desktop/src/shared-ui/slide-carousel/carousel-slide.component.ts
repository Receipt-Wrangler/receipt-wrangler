import { Component, computed, inject, InjectionToken, Signal } from "@angular/core";

/**
 * What a slide needs from the carousel it sits in. A token rather than an
 * import of `SlideCarouselComponent`, which imports this file for its content
 * query; the two would otherwise import each other.
 */
export interface SlideCarouselHost {
  readonly slides: Signal<readonly CarouselSlideComponent[]>;
  readonly activeSlide: Signal<number>;
}

export const SLIDE_CAROUSEL = new InjectionToken<SlideCarouselHost>("SLIDE_CAROUSEL");

/**
 * One slide of an `app-slide-carousel`. Every slide stays rendered; the
 * inactive ones are faded out and taken out of flow (see the stylesheet), so a
 * query over the slides' content always lines up with the slide index.
 */
@Component({
  selector: "app-carousel-slide",
  template: "<ng-content></ng-content>",
  styleUrl: "./carousel-slide.component.scss",
  host: {
    class: "carousel-item",
    "[class.active]": "active()",
    "[attr.aria-hidden]": "!active()",
  },
})
export class CarouselSlideComponent {
  private readonly carousel = inject(SLIDE_CAROUSEL);

  public readonly active = computed(
    () => this.carousel.slides().indexOf(this) === this.carousel.activeSlide()
  );
}
