import { provideZonelessChangeDetection } from "@angular/core";
import { TestBed } from "@angular/core/testing";
import { ActivatedRouteSnapshot, Router, RouterStateSnapshot } from "@angular/router";
import { Observable, of, throwError } from "rxjs";
import { ReportTemplate } from "../../open-api";
import { SnackbarService } from "../../services";
import { ReportRunnerService } from "../services/report-runner.service";
import { reportTemplateResolver } from "./report-template.resolver";

describe("reportTemplateResolver", () => {
  let runner: { getTemplate: jest.Mock };
  let router: { navigate: jest.Mock };
  let snackbar: { error: jest.Mock };

  function route(id: string | null): ActivatedRouteSnapshot {
    return { paramMap: { get: () => id } } as unknown as ActivatedRouteSnapshot;
  }

  function run(id: string | null): Observable<ReportTemplate | null> {
    return TestBed.runInInjectionContext(
      () => reportTemplateResolver(route(id), {} as RouterStateSnapshot) as Observable<ReportTemplate | null>,
    );
  }

  beforeEach(() => {
    runner = { getTemplate: jest.fn() };
    router = { navigate: jest.fn() };
    snackbar = { error: jest.fn() };
    TestBed.configureTestingModule({
      providers: [
        provideZonelessChangeDetection(),
        { provide: ReportRunnerService, useValue: runner },
        { provide: Router, useValue: router },
        { provide: SnackbarService, useValue: snackbar },
      ],
    });
  });

  it("resolves the template for the route id", (done) => {
    const template = { id: 7, name: "R" } as ReportTemplate;
    runner.getTemplate.mockReturnValue(of(template));

    run("7").subscribe((result) => {
      expect(runner.getTemplate).toHaveBeenCalledWith(7);
      expect(result).toBe(template);
      done();
    });
  });

  it("redirects to the list and resolves null when the template can't be loaded", (done) => {
    runner.getTemplate.mockReturnValue(throwError(() => new Error("404")));

    run("99").subscribe((result) => {
      expect(result).toBeNull();
      expect(snackbar.error).toHaveBeenCalled();
      expect(router.navigate).toHaveBeenCalledWith(["/reports"]);
      done();
    });
  });

  it("redirects without calling the API for a non-numeric id", (done) => {
    run("abc").subscribe((result) => {
      expect(result).toBeNull();
      expect(runner.getTemplate).not.toHaveBeenCalled();
      expect(router.navigate).toHaveBeenCalledWith(["/reports"]);
      done();
    });
  });

  it("redirects without calling the API for a missing id", (done) => {
    run(null).subscribe((result) => {
      expect(result).toBeNull();
      expect(runner.getTemplate).not.toHaveBeenCalled();
      expect(router.navigate).toHaveBeenCalledWith(["/reports"]);
      done();
    });
  });
});
