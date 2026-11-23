import { provideZonelessChangeDetection } from "@angular/core";
import { ComponentFixture, TestBed } from "@angular/core/testing";
import { By } from "@angular/platform-browser";
import { FormMode } from "src/enums/form-mode.enum";
import { FileDataView } from "../../open-api";
import { ImageViewerComponent } from "../../shared-ui/image-viewer/image-viewer.component";
import { CarouselComponent } from "./carousel.component";

describe("CarouselComponent", () => {
  let component: CarouselComponent;
  let fixture: ComponentFixture<CarouselComponent>;

  // No encoded image, so each viewer renders nothing: these specs are about
  // which slide is shown, not about the canvas.
  const images = (count: number): FileDataView[] =>
    Array.from({ length: count }, (_, i) => ({ id: i + 1, encodedImage: "" }) as FileDataView);

  const slides = (): HTMLElement[] =>
    Array.from(fixture.nativeElement.querySelectorAll('[data-testid="receipt-image-slide"]'));

  const activeIndex = () => slides().findIndex((slide) => slide.classList.contains("active"));

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [CarouselComponent],
      providers: [provideZonelessChangeDetection()],
    }).compileComponents();

    fixture = TestBed.createComponent(CarouselComponent);
    component = fixture.componentInstance;
    fixture.componentRef.setInput("mode", FormMode.view);
    fixture.componentRef.setInput("hideButtonControls", true);
  });

  it("should create", async () => {
    await fixture.whenStable();

    expect(component).toBeTruthy();
  });

  it("starts on the first image and keeps currentlyShownImageIndex in step", async () => {
    fixture.componentRef.setInput("images", images(3));
    await fixture.whenStable();
    expect(activeIndex()).toBe(0);

    fixture.nativeElement.querySelector('[data-testid="carousel-next"]').click();
    await fixture.whenStable();

    expect(activeIndex()).toBe(1);
    expect(component.currentlyShownImageIndex).toBe(1);
  });

  it("opens on initialIndex", async () => {
    fixture.componentRef.setInput("images", images(3));
    fixture.componentRef.setInput("initialIndex", 2);
    await fixture.whenStable();

    expect(activeIndex()).toBe(2);
    expect(component.currentlyShownImageIndex).toBe(2);
  });

  it("zooms the image being shown", async () => {
    fixture.componentRef.setInput("images", images(3));
    fixture.componentRef.setInput("initialIndex", 1);
    await fixture.whenStable();
    const viewers = fixture.debugElement
      .queryAll(By.directive(ImageViewerComponent))
      .map((viewer) => viewer.componentInstance as ImageViewerComponent);
    const zoomIns = viewers.map((viewer) => jest.spyOn(viewer, "zoomIn").mockImplementation());

    component.zoomIn();

    expect(zoomIns.map((spy) => spy.mock.calls.length)).toEqual([0, 1, 0]);
  });
});
