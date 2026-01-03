import { HttpClient, provideHttpClient, withInterceptors } from "@angular/common/http";
import { HttpTestingController, provideHttpClientTesting } from "@angular/common/http/testing";
import { TestBed } from "@angular/core/testing";
import { MatSnackBarModule } from "@angular/material/snack-bar";
import { RouterTestingModule } from "@angular/router/testing";
import { NgxsModule } from "@ngxs/store";
import { ApiModule } from "../open-api";
import { httpInterceptor } from "./http-interceptor"; // Updated import path

describe("httpInterceptor", () => {
  let httpTestingController: HttpTestingController;

  beforeEach(() => {
    TestBed.configureTestingModule({
      imports: [
        ApiModule,
        MatSnackBarModule,
        NgxsModule.forRoot([]),
        RouterTestingModule
      ],
      providers: [
        provideHttpClient(withInterceptors([httpInterceptor])),
        provideHttpClientTesting()
      ]
    });

    httpTestingController = TestBed.inject(HttpTestingController);
  });

  it("should be created", () => {
    // Instead of testing if the interceptor exists, we'll make a request
    // and see if it works as expected
    expect(httpInterceptor).toBeTruthy();
  });

  // You can add functionality tests
  it("should allow HTTP requests to pass through", () => {
    const testUrl = "/test";

    // Make an HTTP request (this will be intercepted)
    const httpClient = TestBed.inject(HttpClient);
    httpClient.get(testUrl).subscribe();

    // Expect one request to the specified URL
    const req = httpTestingController.expectOne(testUrl);
    expect(req.request.method).toEqual("GET");

    // Respond with mock data
    req.flush({});

    // Verify no outstanding requests
    httpTestingController.verify();
  });
});
