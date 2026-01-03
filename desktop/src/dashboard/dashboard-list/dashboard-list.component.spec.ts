import { ScrollingModule } from "@angular/cdk/scrolling";
import { ComponentFixture, TestBed } from "@angular/core/testing";
import { MatListModule } from "@angular/material/list";
import { NoopAnimationsModule } from "@angular/platform-browser/animations";
import { RouterModule } from "@angular/router";
import { DashboardListComponent } from "./dashboard-list.component";

describe("DashboardListComponent", () => {
  let component: DashboardListComponent;
  let fixture: ComponentFixture<DashboardListComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [DashboardListComponent],
      imports: [
        ScrollingModule,
        NoopAnimationsModule,
        MatListModule,
        RouterModule.forRoot([])
      ]
    }).compileComponents();

    fixture = TestBed.createComponent(DashboardListComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it("should create", () => {
    expect(component).toBeTruthy();
  });

  it("should initialize with default values", () => {
    expect(component.items).toEqual([]);
    expect(component.noItemFoundText).toBe("");
    expect(component.itemSize).toBe(67);
    expect(component.buildRouterLinkString({} as any)).toEqual("");
  });

  it("should emit endOfListReached when reaching end of list", () => {
    // Setup spy on output emitter
    const emitSpy = jest.spyOn(component.endOfListReached, "emit");

    // Mock items array
    component.items = Array(10).fill({});
    fixture.detectChanges();

    // Simulate scroll viewport reaching end
    component.cdkVirtualScrollViewport.setRenderedRange({
      start: 5,
      end: 10
    });

    expect(emitSpy).toHaveBeenCalled();
  });

  it("should not emit endOfListReached when not at end of list", () => {
    const emitSpy = jest.spyOn(component.endOfListReached, "emit");

    component.items = Array(10).fill({});
    fixture.detectChanges();

    component.cdkVirtualScrollViewport.setRenderedRange({
      start: 0,
      end: 5
    });

    expect(emitSpy).not.toHaveBeenCalled();
  });

  it("should properly handle buildRouterLinkString function input", () => {
    const mockLinkFn = (item: any) => `/test/${item.id}`;
    component.buildRouterLinkString = mockLinkFn;

    const testItem = { id: 123 };
    expect(component.buildRouterLinkString(testItem)).toBe("/test/123");
  });

  it("should clean up subscriptions on destroy", () => {
    const subscription = jest.spyOn(
      component.cdkVirtualScrollViewport.renderedRangeStream,
      "subscribe"
    );

    component.ngAfterViewInit();
    fixture.destroy();

    expect(subscription).toHaveBeenCalled();
  });
});
