import { Component, contentChildren, effect, input, model } from "@angular/core";
import { CarouselSlideComponent, SLIDE_CAROUSEL } from "./carousel-slide.component";

/**
 * A carousel of `app-carousel-slide`s that crossfades between them.
 *
 * Built on Bootstrap's own carousel classes, so the global Bootstrap stylesheet
 * draws the controls and indicators exactly as it did for the ngx-bootstrap
 * carousel this replaced. Bootstrap's carousel JavaScript is not loaded; the
 * `data-bs-target` attribute on the indicators is only there because Bootstrap's
 * CSS selects them by it.
 *
 * Navigation wraps. There is no autoplay.
 *
 * @example
 * <app-slide-carousel [(activeSlide)]="index">
 *   @for (image of images; track image) {
 *     <app-carousel-slide>...</app-carousel-slide>
 *   }
 * </app-slide-carousel>
 */
@Component({
  selector: "app-slide-carousel",
  templateUrl: "./slide-carousel.component.html",
  styleUrl: "./slide-carousel.component.scss",
  providers: [{ provide: SLIDE_CAROUSEL, useExisting: SlideCarouselComponent }],
})
export class SlideCarouselComponent {
  /** Index of the slide being shown. Two-way bindable. */
  public readonly activeSlide = model(0);

  public readonly showIndicators = input(true);

  /** The previous/next buttons over the slides. */
  public readonly showControls = input(true);

  public readonly slides = contentChildren(CarouselSlideComponent);

  constructor() {
    // When the slide being shown is removed and the index now points past the
    // end (or the owner binds one that does), fall back to the first slide and
    // tell the owner through the two-way binding — an owner that keeps the
    // index, like the Quick Scan dialog, acts on it. A sync to the binding's
    // owner, not a derivation: there is no value to compute without writing it.
    effect(() => {
      const count = this.slides().length;
      const index = this.activeSlide();
      if (count > 0 && (index < 0 || index >= count)) {
        this.activeSlide.set(0);
      }
    });
  }

  public select(index: number): void {
    const count = this.slides().length;
    if (count === 0) {
      return;
    }
    this.activeSlide.set(((index % count) + count) % count);
  }

  public next(): void {
    this.select(this.activeSlide() + 1);
  }

  public previous(): void {
    this.select(this.activeSlide() - 1);
  }

  /**
   * Enter/Space and the arrow keys move between slides, but only while the
   * carousel itself has focus. The ngx-bootstrap carousel acted on keys from
   * anything inside it, so Space typed into a Quick Scan field was swallowed
   * and switched images, and Enter on its Previous control went forward.
   */
  public onKeydown(event: KeyboardEvent): void {
    if (event.target !== event.currentTarget) {
      return;
    }

    switch (event.key) {
      case "Enter":
      case " ":
        event.preventDefault();
        this.next();
        break;
      case "ArrowLeft":
        this.previous();
        break;
      case "ArrowRight":
        this.next();
        break;
    }
  }
}
