import { provideZonelessChangeDetection } from "@angular/core";
import { ComponentFixture, TestBed } from "@angular/core/testing";
import { ImageCanvasComponent } from "./image-canvas.component";

const STAGE_WIDTH = 600;
const STAGE_HEIGHT = 400;
const NATURAL_WIDTH = 1200;
const NATURAL_HEIGHT = 1600;

/** Fit is limited by height here: 400 / 1600 = 0.25. */
const FIT_SCALE = STAGE_HEIGHT / NATURAL_HEIGHT;

/**
 * jsdom implements neither `PointerEvent` nor pointer capture, so pointer events
 * are faked from `MouseEvent` — which carries every property the component reads
 * — and the capture calls are stubbed per element.
 */
function pointerEvent(
  type: string,
  clientX: number,
  clientY: number,
  pointerId = 1,
  button = 0,
): MouseEvent {
  const event = new MouseEvent(type, { bubbles: true, cancelable: true, clientX, clientY, button });
  Object.defineProperty(event, "pointerId", { value: pointerId });

  return event;
}

describe("ImageCanvasComponent", () => {
  let fixture: ComponentFixture<ImageCanvasComponent>;
  let component: ImageCanvasComponent;
  let host: HTMLElement;
  let image: HTMLImageElement;

  /** The image's transform, parsed back into the viewport it encodes. */
  const viewport = (): { x: number; y: number; scale: number } => {
    const match = /translate\((-?[\d.]+)px, (-?[\d.]+)px\) scale\(([\d.]+)\)/.exec(
      image.style.transform,
    );

    return match
      ? { x: Number(match[1]), y: Number(match[2]), scale: Number(match[3]) }
      : { x: NaN, y: NaN, scale: NaN };
  };

  const handle = (corner: string): HTMLElement =>
    host.querySelector(`[data-corner="${corner}"]`) as HTMLElement;

  const drag = async (
    target: HTMLElement,
    from: [number, number],
    to: [number, number],
    release = true,
  ) => {
    target.dispatchEvent(pointerEvent("pointerdown", from[0], from[1]));
    host.dispatchEvent(pointerEvent("pointermove", to[0], to[1]));
    if (release) {
      host.dispatchEvent(pointerEvent("pointerup", to[0], to[1]));
    }
    await fixture.whenStable();
  };

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [ImageCanvasComponent],
      providers: [provideZonelessChangeDetection()],
    }).compileComponents();

    fixture = TestBed.createComponent(ImageCanvasComponent);
    component = fixture.componentInstance;
    fixture.componentRef.setInput("src", "data:image/png;base64,abc");
    await fixture.whenStable();

    host = fixture.nativeElement as HTMLElement;
    image = host.querySelector("img") as HTMLImageElement;

    // jsdom lays nothing out and decodes no images, so every measurement the
    // component reads has to be stubbed.
    Object.defineProperty(host, "clientWidth", { value: STAGE_WIDTH, configurable: true });
    Object.defineProperty(host, "clientHeight", { value: STAGE_HEIGHT, configurable: true });
    host.getBoundingClientRect = () => ({ left: 0, top: 0, width: STAGE_WIDTH, height: STAGE_HEIGHT }) as DOMRect;
    host.setPointerCapture = jest.fn();
    host.releasePointerCapture = jest.fn();
    host.hasPointerCapture = jest.fn().mockReturnValue(true);
    Object.defineProperty(image, "naturalWidth", { value: NATURAL_WIDTH, configurable: true });
    Object.defineProperty(image, "naturalHeight", { value: NATURAL_HEIGHT, configurable: true });

    image.dispatchEvent(new Event("load"));
    await fixture.whenStable();
  });

  it("should create", () => {
    expect(component).toBeTruthy();
  });

  it("fits the image to the stage on load and centres it", () => {
    const view = viewport();

    expect(view.scale).toBeCloseTo(FIT_SCALE);
    // 1200 * 0.25 = 300 wide in a 600 stage, so 150 either side; the height fills.
    expect(view.x).toBeCloseTo((STAGE_WIDTH - NATURAL_WIDTH * FIT_SCALE) / 2);
    expect(view.y).toBeCloseTo(0);
  });

  it("renders a handle at each corner", () => {
    expect(host.querySelectorAll("[data-corner]").length).toEqual(4);
    ["nw", "ne", "se", "sw"].forEach((corner) => expect(handle(corner)).toBeTruthy());
  });

  it("positions the handle frame over the image", () => {
    const frame = host.querySelector(".rw-image-canvas__frame") as HTMLElement;

    expect(frame.style.width).toEqual(`${NATURAL_WIDTH * FIT_SCALE}px`);
    expect(frame.style.height).toEqual(`${NATURAL_HEIGHT * FIT_SCALE}px`);
  });

  // Pulling a corner scales about the corner opposite it, so that one stays put -
  // in the axis the image overflows. The other axis obeys the centring rule below.
  it("scales about the opposite corner when a handle is pulled", async () => {
    const before = viewport();
    const anchorY = before.y + NATURAL_HEIGHT * before.scale;

    // Drag the top-left handle up and to the left: the image grows.
    await drag(handle("nw"), [before.x, before.y], [before.x - 200, before.y - 200]);

    const after = viewport();
    expect(after.scale).toBeGreaterThan(before.scale);
    expect(after.y + NATURAL_HEIGHT * after.scale).toBeCloseTo(anchorY, 0);
  });

  // An axis narrower than the stage is centred rather than left wherever a
  // gesture put it, so a letterboxed image cannot drift into a corner. It is why
  // the anchor above is asserted on one axis: at the fit scale this image is
  // 300px wide in a 600px stage, so horizontally it stays centred as it grows,
  // and only edge-clamps once it is wider than the stage.
  it("keeps an axis narrower than the stage centred as it grows", async () => {
    const centred = (scale: number) => (STAGE_WIDTH - NATURAL_WIDTH * scale) / 2;

    await drag(handle("nw"), [150, 0], [100, -100]);

    const after = viewport();
    expect(NATURAL_WIDTH * after.scale).toBeLessThan(STAGE_WIDTH);
    expect(after.x).toBeCloseTo(centred(after.scale));
  });

  // A distorted receipt is never wanted, so a corner drag is uniform.
  it("keeps the aspect ratio locked while resizing", async () => {
    const before = viewport();
    await drag(handle("se"), [0, 0], [600, 10]);

    const after = viewport();
    const frame = host.querySelector(".rw-image-canvas__frame") as HTMLElement;
    const ratio = Number.parseFloat(frame.style.width) / Number.parseFloat(frame.style.height);

    expect(after.scale).not.toEqual(before.scale);
    expect(ratio).toBeCloseTo(NATURAL_WIDTH / NATURAL_HEIGHT);
  });

  it("pans a zoomed image and keeps its edges against the stage", async () => {
    component.zoomIn();
    component.zoomIn();
    await fixture.whenStable();

    const before = viewport();
    await drag(host, [300, 200], [300, 120]);

    const after = viewport();
    expect(after.y).toBeLessThan(before.y);

    // Dragging far past the edge stops with the image edge on the stage edge.
    await drag(host, [300, 200], [300, 5000]);
    expect(viewport().y).toEqual(0);
  });

  // The image should never be draggable out of view.
  it("centres an axis that is smaller than the stage instead of letting it drift", async () => {
    await drag(host, [300, 200], [40, 200]);

    expect(viewport().x).toBeCloseTo((STAGE_WIDTH - NATURAL_WIDTH * FIT_SCALE) / 2);
  });

  it("will not zoom out past the fit", async () => {
    component.zoomOut();
    component.zoomOut();
    await fixture.whenStable();

    expect(viewport().scale).toBeCloseTo(FIT_SCALE);
  });

  it("zooms in and back out again", async () => {
    component.zoomIn();
    await fixture.whenStable();
    const zoomed = viewport().scale;
    expect(zoomed).toBeGreaterThan(FIT_SCALE);

    component.zoomOut();
    await fixture.whenStable();
    expect(viewport().scale).toBeCloseTo(FIT_SCALE);
  });

  // The viewer this replaced multiplied deltaY by -0.000001, so a notch moved the
  // scale by 0.0001 and wheel zoom did nothing at all.
  it("zooms a meaningful amount on the wheel, anchored at the cursor", async () => {
    const before = viewport();
    const cursorY = 200;

    host.dispatchEvent(
      new WheelEvent("wheel", { deltaY: -100, clientX: 300, clientY: cursorY, cancelable: true }),
    );
    await fixture.whenStable();

    const after = viewport();
    expect(after.scale).toBeGreaterThan(before.scale * 1.1);
    // The image point under the cursor did not move. Asserted vertically, the
    // axis that overflows the stage - see the centring rule above.
    expect((cursorY - before.y) / before.scale).toBeCloseTo((cursorY - after.y) / after.scale, 1);
  });

  it("returns to the fit on double-click", async () => {
    component.zoomIn();
    component.zoomIn();
    await fixture.whenStable();
    expect(viewport().scale).toBeGreaterThan(FIT_SCALE);

    host.dispatchEvent(new MouseEvent("dblclick", { bubbles: true }));
    await fixture.whenStable();

    expect(viewport().scale).toBeCloseTo(FIT_SCALE);
  });

  it("ignores a non-primary button", async () => {
    const before = viewport();
    host.dispatchEvent(pointerEvent("pointerdown", 300, 200, 1, 2));
    host.dispatchEvent(pointerEvent("pointermove", 300, 100));
    await fixture.whenStable();

    expect(viewport()).toEqual(before);
  });

  it("stops tracking once the drag ends", async () => {
    component.zoomIn();
    component.zoomIn();
    await fixture.whenStable();

    await drag(host, [300, 200], [300, 150]);
    const after = viewport();

    host.dispatchEvent(pointerEvent("pointermove", 300, 50));
    await fixture.whenStable();

    expect(viewport()).toEqual(after);
  });

  // A handle sits over the image, which starts a pan on pointerdown.
  it("does not start a pan when a handle is grabbed", async () => {
    const before = viewport();
    handle("se").dispatchEvent(pointerEvent("pointerdown", 100, 100));
    host.dispatchEvent(pointerEvent("pointermove", 140, 140));
    await fixture.whenStable();

    // Resized about the top-left, so the origin is unchanged - a pan would have
    // moved it.
    expect(viewport().x).toEqual(before.x);
    expect(viewport().y).toEqual(before.y);
  });
});
