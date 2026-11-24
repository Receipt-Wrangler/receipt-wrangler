import { provideHttpClientTesting } from "@angular/common/http/testing";
import { TestBed } from "@angular/core/testing";
import { Router } from "@angular/router";
import { RouterTestingModule } from "@angular/router/testing";
import { NgxsModule, Store } from "@ngxs/store";
import { firstValueFrom, of, Subject, throwError } from "rxjs";
import { ApiModule, AuthService, Claims } from "../open-api";
import { AuthState, Logout, SetAuthState } from "../store";
import { TokenRefreshService } from "./token-refresh.service";
import { provideHttpClient, withInterceptorsFromDi } from "@angular/common/http";

describe("TokenRefreshService", () => {
  let service: TokenRefreshService;
  let store: Store;
  let authService: AuthService;
  let router: Router;

  const mockClaims: Claims = {
    userId: 1,
    displayName: "Test User",
    defaultAvatarColor: "#CD5C5C",
    username: "testuser",
    iss: "https://receiptWrangler.io",
    exp: Math.floor(Date.now() / 1000) + 1200,
  };

  beforeEach(() => {
    TestBed.configureTestingModule({
      imports: [
        ApiModule,
        NgxsModule.forRoot([AuthState]),
        RouterTestingModule,
      ],
      providers: [
        provideHttpClient(withInterceptorsFromDi()),
        provideHttpClientTesting(),
      ],
    });

    service = TestBed.inject(TokenRefreshService);
    store = TestBed.inject(Store);
    authService = TestBed.inject(AuthService);
    router = TestBed.inject(Router);
  });

  it("should be created", () => {
    expect(service).toBeTruthy();
  });

  describe("refreshToken", () => {
    it("should call authService.getNewRefreshToken and return claims", async () => {
      const spy = vi.spyOn(authService, "getNewRefreshToken").mockReturnValue(
        of(mockClaims as any)
      );

      const claims = await firstValueFrom(service.refreshToken());

      expect(spy).toHaveBeenCalledTimes(1);
      expect(claims).toEqual(mockClaims);
    });

    it("should dispatch SetAuthState on successful refresh", async () => {
      vi.spyOn(authService, "getNewRefreshToken").mockReturnValue(
        of(mockClaims as any)
      );
      const dispatchSpy = vi.spyOn(store, "dispatch");

      await firstValueFrom(service.refreshToken());

      expect(dispatchSpy).toHaveBeenCalledWith(
        new SetAuthState(mockClaims)
      );
    });

    it("should dispatch Logout and clear localStorage on refresh failure", async () => {
      const error = new Error("Token expired");
      vi.spyOn(authService, "getNewRefreshToken").mockReturnValue(
        throwError(() => error)
      );
      const dispatchSpy = vi.spyOn(store, "dispatch").mockReturnValue(of(undefined));
      const clearSpy = vi.spyOn(Storage.prototype, "clear");
      const navigateSpy = vi.spyOn(router, "navigate").mockResolvedValue(true);

      await expect(firstValueFrom(service.refreshToken())).rejects.toBe(error);

      expect(dispatchSpy).toHaveBeenCalledWith(new Logout());
      expect(clearSpy).toHaveBeenCalled();
      expect(navigateSpy).toHaveBeenCalledWith(["/auth/login"]);
      clearSpy.mockRestore();
    });

    it("should navigate to /auth/login on refresh failure", async () => {
      vi.spyOn(authService, "getNewRefreshToken").mockReturnValue(
        throwError(() => new Error("fail"))
      );
      vi.spyOn(store, "dispatch").mockReturnValue(of(undefined));
      const clearSpy = vi.spyOn(Storage.prototype, "clear");
      const navigateSpy = vi.spyOn(router, "navigate").mockResolvedValue(true);

      await expect(firstValueFrom(service.refreshToken())).rejects.toThrow("fail");

      expect(navigateSpy).toHaveBeenCalledWith(["/auth/login"]);
      clearSpy.mockRestore();
    });

    it("should serialize concurrent calls — only one HTTP request fires", () => {
      const subject = new Subject<any>();
      const spy = vi.spyOn(authService, "getNewRefreshToken").mockReturnValue(
        subject.asObservable()
      );

      // Subscribe twice before the first completes
      const results: Claims[] = [];
      service.refreshToken().subscribe((c) => results.push(c));
      service.refreshToken().subscribe((c) => results.push(c));

      // Only one HTTP call should have been made
      expect(spy).toHaveBeenCalledTimes(1);

      // Emit the response
      subject.next(mockClaims);
      subject.complete();

      // Both subscribers should have received the same result
      expect(results.length).toBe(2);
      expect(results[0]).toEqual(mockClaims);
      expect(results[1]).toEqual(mockClaims);
    });

    it("should serialize concurrent calls — SetAuthState dispatched once", () => {
      const subject = new Subject<any>();
      vi.spyOn(authService, "getNewRefreshToken").mockReturnValue(
        subject.asObservable()
      );
      const dispatchSpy = vi.spyOn(store, "dispatch");

      service.refreshToken().subscribe();
      service.refreshToken().subscribe();

      subject.next(mockClaims);
      subject.complete();

      const setAuthCalls = dispatchSpy.mock.calls.filter(
        ([action]) => action instanceof SetAuthState
      );
      expect(setAuthCalls.length).toBe(1);
    });

    it("should reset in-flight state after successful completion allowing new requests", () => {
      const spy = vi.spyOn(authService, "getNewRefreshToken").mockReturnValue(
        of(mockClaims as any)
      );

      // First call
      service.refreshToken().subscribe();
      expect(spy).toHaveBeenCalledTimes(1);

      // Second call after first completes — should make a new HTTP request
      service.refreshToken().subscribe();
      expect(spy).toHaveBeenCalledTimes(2);
    });

    it("should reset in-flight state after error allowing new requests", () => {
      const spy = vi.spyOn(authService, "getNewRefreshToken");
      vi.spyOn(store, "dispatch").mockReturnValue(of(undefined));
      vi.spyOn(router, "navigate").mockResolvedValue(true);
      const clearSpy = vi.spyOn(Storage.prototype, "clear");

      // First call — fails
      spy.mockReturnValue(throwError(() => new Error("fail")));
      service.refreshToken().subscribe({ error: () => {} });
      expect(spy).toHaveBeenCalledTimes(1);

      // Second call after error — should make a new HTTP request
      spy.mockReturnValue(of(mockClaims as any));
      service.refreshToken().subscribe();
      expect(spy).toHaveBeenCalledTimes(2);

      clearSpy.mockRestore();
    });

    it("should propagate error to all concurrent subscribers on failure", () => {
      const subject = new Subject<any>();
      vi.spyOn(authService, "getNewRefreshToken").mockReturnValue(
        subject.asObservable()
      );
      vi.spyOn(store, "dispatch").mockReturnValue(of(undefined));
      vi.spyOn(router, "navigate").mockResolvedValue(true);
      const clearSpy = vi.spyOn(Storage.prototype, "clear");

      const errors: any[] = [];
      service.refreshToken().subscribe({ error: (e) => errors.push(e) });
      service.refreshToken().subscribe({ error: (e) => errors.push(e) });

      const error = new Error("refresh failed");
      subject.error(error);

      expect(errors.length).toBe(2);
      expect(errors[0]).toBe(error);
      expect(errors[1]).toBe(error);

      clearSpy.mockRestore();
    });

    it("should dispatch Logout only once for concurrent failures", () => {
      const subject = new Subject<any>();
      vi.spyOn(authService, "getNewRefreshToken").mockReturnValue(
        subject.asObservable()
      );
      const dispatchSpy = vi.spyOn(store, "dispatch").mockReturnValue(of(undefined));
      vi.spyOn(router, "navigate").mockResolvedValue(true);
      const clearSpy = vi.spyOn(Storage.prototype, "clear");

      service.refreshToken().subscribe({ error: () => {} });
      service.refreshToken().subscribe({ error: () => {} });

      subject.error(new Error("fail"));

      const logoutCalls = dispatchSpy.mock.calls.filter(
        ([action]) => action instanceof Logout
      );
      expect(logoutCalls.length).toBe(1);

      clearSpy.mockRestore();
    });

    it("should not dispatch SetAuthState on failure", () => {
      vi.spyOn(authService, "getNewRefreshToken").mockReturnValue(
        throwError(() => new Error("fail"))
      );
      const dispatchSpy = vi.spyOn(store, "dispatch").mockReturnValue(of(undefined));
      vi.spyOn(router, "navigate").mockResolvedValue(true);
      const clearSpy = vi.spyOn(Storage.prototype, "clear");

      service.refreshToken().subscribe({ error: () => {} });

      const setAuthCalls = dispatchSpy.mock.calls.filter(
        ([action]) => action instanceof SetAuthState
      );
      expect(setAuthCalls.length).toBe(0);

      clearSpy.mockRestore();
    });

    it("should update auth state expiration date on success", async () => {
      const futureExp = Math.floor(Date.now() / 1000) + 3600;
      const claims = { ...mockClaims, exp: futureExp };

      vi.spyOn(authService, "getNewRefreshToken").mockReturnValue(
        of(claims as any)
      );

      await firstValueFrom(service.refreshToken());

      const expirationDate = store.selectSnapshot(
        (state) => state.auth.expirationDate
      );
      expect(expirationDate).toBe(futureExp.toString());
    });

    it("should handle rapid sequential calls correctly", () => {
      const spy = vi.spyOn(authService, "getNewRefreshToken").mockReturnValue(
        of(mockClaims as any)
      );

      const results: Claims[] = [];

      // 5 rapid sequential calls — each completes before the next starts
      for (let i = 0; i < 5; i++) {
        service.refreshToken().subscribe((c) => results.push(c));
      }

      // Each call completes synchronously (of()), so each should start a new request
      expect(spy).toHaveBeenCalledTimes(5);
      expect(results.length).toBe(5);
    });

    it("should handle late subscriber to in-flight refresh", () => {
      const subject = new Subject<any>();
      vi.spyOn(authService, "getNewRefreshToken").mockReturnValue(
        subject.asObservable()
      );

      const results: Claims[] = [];

      // First subscriber
      service.refreshToken().subscribe((c) => results.push(c));

      // Complete the request
      subject.next(mockClaims);
      subject.complete();

      expect(results.length).toBe(1);

      // Late subscriber after completion — shareReplay(1) replays the last value
      // but finalize has reset refreshInFlight$, so this starts a new request
      const subject2 = new Subject<any>();
      vi.spyOn(authService, "getNewRefreshToken").mockReturnValue(
        subject2.asObservable()
      );

      service.refreshToken().subscribe((c) => results.push(c));
      subject2.next(mockClaims);
      subject2.complete();

      expect(results.length).toBe(2);
    });
  });
});
