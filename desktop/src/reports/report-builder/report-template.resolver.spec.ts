import type { Mock } from "vitest";
import { provideZonelessChangeDetection } from "@angular/core";
import { TestBed } from "@angular/core/testing";
import { ActivatedRouteSnapshot, Router, RouterStateSnapshot } from "@angular/router";
import { firstValueFrom, Observable, of, throwError } from "rxjs";
import { ReportTemplate } from "../../open-api";
import { SnackbarService } from "../../services";
import { ReportRunnerService } from "../services/report-runner.service";
import { reportTemplateResolver } from "./report-template.resolver";

describe("reportTemplateResolver", () => {
  let runner: { getTemplate: Mock };
  let router: { navigate: Mock };
  let snackbar: { error: Mock };

  function route(id: string | null): ActivatedRouteSnapshot {
    return { paramMap: { get: () => id } } as unknown as ActivatedRouteSnapshot;
  }

  function run(id: string | null): Observable<ReportTemplate | null> {
    return TestBed.runInInjectionContext(
      () => reportTemplateResolver(route(id), {} as RouterStateSnapshot) as Observable<ReportTemplate | null>,
    );
  }

  beforeEach(() => {
    runner = { getTemplate: vi.fn() };
    router = { navigate: vi.fn() };
    snackbar = { error: vi.fn() };
    TestBed.configureTestingModule({
      providers: [
        provideZonelessChangeDetection(),
        { provide: ReportRunnerService, useValue: runner },
        { provide: Router, useValue: router },
        { provide: SnackbarService, useValue: snackbar },
      ],
    });
  });

  it("resolves the template for the route id", async () => {
    const template = { id: 7, name: "R" } as ReportTemplate;
    runner.getTemplate.mockReturnValue(of(template));

    const result = await firstValueFrom(run("7"));

    expect(runner.getTemplate).toHaveBeenCalledWith(7);
    expect(result).toBe(template);
  });

  it("redirects to the list and resolves null when the template can't be loaded", async () => {
    runner.getTemplate.mockReturnValue(throwError(() => new Error("404")));

    const result = await firstValueFrom(run("99"));

    expect(result).toBeNull();
    expect(snackbar.error).toHaveBeenCalled();
    expect(router.navigate).toHaveBeenCalledWith(["/reports"]);
  });

  it("redirects without calling the API for a non-numeric id", async () => {
    const result = await firstValueFrom(run("abc"));

    expect(result).toBeNull();
    expect(runner.getTemplate).not.toHaveBeenCalled();
    expect(router.navigate).toHaveBeenCalledWith(["/reports"]);
  });

  it("redirects without calling the API for a missing id", async () => {
    const result = await firstValueFrom(run(null));

    expect(result).toBeNull();
    expect(runner.getTemplate).not.toHaveBeenCalled();
    expect(router.navigate).toHaveBeenCalledWith(["/reports"]);
  });

  it("redirects without calling the API for an id past the safe-integer range", async () => {
    // Number("9007199254740993") rounds to 9007199254740992, so requesting it would
    // hit a different template; isSafeInteger rejects it before the call.
    const result = await firstValueFrom(run("9007199254740993"));

    expect(result).toBeNull();
    expect(runner.getTemplate).not.toHaveBeenCalled();
    expect(router.navigate).toHaveBeenCalledWith(["/reports"]);
  });
});
