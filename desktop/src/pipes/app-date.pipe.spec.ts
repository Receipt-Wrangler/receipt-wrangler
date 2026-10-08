import { Component, provideZonelessChangeDetection, signal } from "@angular/core";
import { TestBed } from "@angular/core/testing";
import { NgxsModule, Store } from "@ngxs/store";
import { SystemSettingsState } from "../store/system-settings.state";
import { SetTimeZone } from "../store/system-settings.state.actions";
import { AppDatePipe, formatAppDate } from "./app-date.pipe";

@Component({
  selector: "app-date-host",
  standalone: true,
  imports: [AppDatePipe],
  template: `<span id="instant">{{ instant() | appDate: "yyyy-MM-dd HH:mm" }}</span><span id="calendar">{{ day() | appDate: "mediumDate" : "calendar" }}</span>`,
})
class HostComponent {
  public readonly instant = signal<string | null>("2026-10-01T01:00:00Z");
  public readonly day = signal<string | null>("2026-09-30T00:00:00Z");
}

describe("AppDatePipe", () => {
  const format = (value: any, fmt: string, mode: "instant" | "calendar", zone: string) =>
    formatAppDate(value, fmt, mode, zone, "en-US");

  describe("instant mode", () => {
    it("formats in the app zone, not the browser's", () => {
      expect(format("2026-10-01T01:00:00Z", "yyyy-MM-dd HH:mm", "instant", "UTC")).toBe("2026-10-01 01:00");
      expect(format("2026-10-01T01:00:00Z", "yyyy-MM-dd HH:mm", "instant", "America/New_York")).toBe(
        "2026-09-30 21:00"
      );
    });

    // Either side of the Mar 8 2026 spring-forward: one zone, two offsets.
    it("applies the offset in force at each instant across a DST boundary", () => {
      expect(format("2026-03-08T06:30:00Z", "HH:mm", "instant", "America/New_York")).toBe("01:30");
      expect(format("2026-03-08T07:30:00Z", "HH:mm", "instant", "America/New_York")).toBe("03:30");
      expect(format("2026-03-08T07:30:00Z", "HH:mm", "instant", "UTC")).toBe("07:30");
    });

    it("accepts a Date", () => {
      expect(format(new Date("2026-10-01T01:00:00Z"), "MMM d", "instant", "America/New_York")).toBe("Sep 30");
    });
  });

  describe("calendar mode", () => {
    it("shows the stored UTC day in every zone", () => {
      for (const zone of ["UTC", "America/New_York", "Pacific/Kiritimati"]) {
        expect(format("2026-09-30T00:00:00Z", "mediumDate", "calendar", zone)).toBe("Sep 30, 2026");
      }
    });

    it("shows a bare day as written", () => {
      expect(format("2026-09-30", "mediumDate", "calendar", "America/New_York")).toBe("Sep 30, 2026");
    });
  });

  it("returns null for an empty or unparseable value", () => {
    expect(format(null, "short", "instant", "UTC")).toBeNull();
    expect(format(undefined, "short", "calendar", "UTC")).toBeNull();
    expect(format("", "short", "instant", "UTC")).toBeNull();
    expect(format("nope", "short", "instant", "UTC")).toBeNull();
  });

  describe("in a template", () => {
    beforeEach(async () => {
      await TestBed.configureTestingModule({
        imports: [HostComponent, NgxsModule.forRoot([SystemSettingsState])],
        providers: [provideZonelessChangeDetection()],
      }).compileComponents();
    });

    const text = (fixture: any, id: string): string =>
      fixture.nativeElement.querySelector(`#${id}`).textContent.trim();

    it("defaults to mediumDate and the UTC zone", () => {
      const pipe = TestBed.runInInjectionContext(() => new AppDatePipe());

      expect(pipe.transform("2026-10-01T01:00:00Z")).toBe("Oct 1, 2026");
    });

    // The zone is store state an admin can change with the page open; every
    // instant on screen follows it, and calendar days do not move.
    it("re-renders when the app time zone changes", async () => {
      const fixture = TestBed.createComponent(HostComponent);
      await fixture.whenStable();
      expect(text(fixture, "instant")).toBe("2026-10-01 01:00");
      expect(text(fixture, "calendar")).toBe("Sep 30, 2026");

      TestBed.inject(Store).dispatch(new SetTimeZone("America/New_York"));
      await fixture.whenStable();

      expect(text(fixture, "instant")).toBe("2026-09-30 21:00");
      expect(text(fixture, "calendar")).toBe("Sep 30, 2026");
    });

    it("re-renders when the value changes", async () => {
      const fixture = TestBed.createComponent(HostComponent);
      await fixture.whenStable();

      fixture.componentInstance.instant.set("2026-03-08T07:30:00Z");
      await fixture.whenStable();

      expect(text(fixture, "instant")).toBe("2026-03-08 07:30");
    });
  });
});
