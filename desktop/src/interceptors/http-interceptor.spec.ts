import { HttpClient, provideHttpClient, withInterceptors } from "@angular/common/http";
import { HttpTestingController, provideHttpClientTesting } from "@angular/common/http/testing";
import { TestBed } from "@angular/core/testing";
import { MatSnackBarModule } from "@angular/material/snack-bar";
import { ActivatedRoute } from "@angular/router";
import { Store } from "@ngxs/store";
import { ApiModule } from "../open-api";
import { SnackbarService } from "../services";
import { httpInterceptor } from "./http-interceptor";

const FORBIDDEN_MESSAGE = "You do not have permission to perform this action.";

describe("httpInterceptor", () => {
  let httpTestingController: HttpTestingController;
  let httpClient: HttpClient;
  let snackbarService: SnackbarService;
  let errorSpy: jest.SpyInstance;
  let isLoggedIn: boolean;
  let activatedRouteStub: { snapshot: { queryParams: Record<string, unknown> } };

  beforeEach(() => {
    isLoggedIn = true;
    activatedRouteStub = { snapshot: { queryParams: {} } };

    TestBed.configureTestingModule({
      imports: [ApiModule, MatSnackBarModule],
      providers: [
        provideHttpClient(withInterceptors([httpInterceptor])),
        provideHttpClientTesting(),
        { provide: Store, useValue: { selectSnapshot: () => isLoggedIn } },
        { provide: ActivatedRoute, useValue: activatedRouteStub },
      ],
    });

    httpTestingController = TestBed.inject(HttpTestingController);
    httpClient = TestBed.inject(HttpClient);
    snackbarService = TestBed.inject(SnackbarService);
    errorSpy = jest.spyOn(snackbarService, "error").mockImplementation(() => {});
  });

  afterEach(() => {
    httpTestingController.verify();
  });

  it("should be created", () => {
    expect(httpInterceptor).toBeTruthy();
  });

  it("should allow HTTP requests to pass through", () => {
    httpClient.get("/test").subscribe();

    const req = httpTestingController.expectOne("/test");
    expect(req.request.method).toEqual("GET");
    req.flush({});
  });

  it("shows a forbidden toast and propagates the error on a mutation 403 when logged in", () => {
    let erroredStatus: number | undefined;
    httpClient
      .post("/test", {})
      .subscribe({ error: (e) => (erroredStatus = e.status) });

    httpTestingController
      .expectOne("/test")
      .flush({}, { status: 403, statusText: "Forbidden" });

    expect(errorSpy).toHaveBeenCalledTimes(1);
    expect(errorSpy).toHaveBeenCalledWith(FORBIDDEN_MESSAGE);
    expect(erroredStatus).toEqual(403);
  });

  it("does not toast a 403 on a background GET read", () => {
    httpClient.get("/test").subscribe({ error: () => {} });

    httpTestingController
      .expectOne("/test")
      .flush({}, { status: 403, statusText: "Forbidden" });

    expect(errorSpy).not.toHaveBeenCalled();
  });

  it("does not toast a mutation 403 while in queue mode", () => {
    activatedRouteStub.snapshot.queryParams = { queueMode: "true" };

    httpClient.post("/test", {}).subscribe({ error: () => {} });

    httpTestingController
      .expectOne("/test")
      .flush({}, { status: 403, statusText: "Forbidden" });

    expect(errorSpy).not.toHaveBeenCalled();
  });

  it("does not toast a 403 when the token is no longer valid", () => {
    isLoggedIn = false;

    httpClient.post("/test", {}).subscribe({ error: () => {} });

    httpTestingController
      .expectOne("/test")
      .flush({}, { status: 403, statusText: "Forbidden" });

    expect(errorSpy).not.toHaveBeenCalled();
  });

  it("surfaces server errorMsg payloads via a toast", () => {
    httpClient.get("/test").subscribe({ error: () => {} });

    httpTestingController
      .expectOne("/test")
      .flush(
        { errorMsg: "Something broke" },
        { status: 400, statusText: "Bad Request" }
      );

    expect(errorSpy).toHaveBeenCalledTimes(1);
    expect(errorSpy).toHaveBeenCalledWith("Something broke");
  });
});
