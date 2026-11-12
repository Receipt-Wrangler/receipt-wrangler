import { HttpErrorResponse, provideHttpClient } from "@angular/common/http";
import { HttpTestingController, provideHttpClientTesting } from "@angular/common/http/testing";
import { TestBed } from "@angular/core/testing";
import { Store } from "@ngxs/store";
import { firstValueFrom } from "rxjs";
import { SUPPRESS_ERROR_TOAST } from "../interceptors/http-interceptor";
import { ApiModule } from "../open-api";
import * as fileUtils from "../utils/file";
import { CLAIM_NOT_FOUND_MESSAGE } from "./claim.resolver";
import {
  CLAIM_PACKET_FAILED_MESSAGE,
  CLAIM_PACKET_OFFLINE_MESSAGE,
  CLAIM_PACKET_TIMEOUT_MESSAGE,
  ClaimPacketDownloadService,
  ClaimPacketError,
  claimPacketErrorMessage,
  claimPacketFilename,
} from "./claim-packet-download.service";

describe("ClaimPacketDownloadService", () => {
  let service: ClaimPacketDownloadService;
  let http: HttpTestingController;
  let downloadSpy: jest.SpyInstance;
  let timeZone: string;

  const claim = { id: 12, name: "Voucher 14: Rome / Lima" };
  const zip = new Blob(["PK"], { type: "application/zip" });

  beforeEach(() => {
    timeZone = "UTC";
    jest.useFakeTimers({ now: new Date("2026-11-12T15:00:00Z"), doNotFake: ["queueMicrotask", "nextTick"] });

    TestBed.configureTestingModule({
      imports: [ApiModule],
      providers: [
        provideHttpClient(),
        provideHttpClientTesting(),
        { provide: Store, useValue: { selectSnapshot: () => timeZone } },
      ],
    });

    service = TestBed.inject(ClaimPacketDownloadService);
    http = TestBed.inject(HttpTestingController);
    downloadSpy = jest.spyOn(fileUtils, "downloadFile").mockImplementation(() => {});
  });

  afterEach(() => {
    http.verify();
    downloadSpy.mockRestore();
    jest.useRealTimers();
  });

  const expectPacketRequest = () => http.expectOne("/api/claim/12/packet");

  it("asks for the claim's packet as a blob, and keeps the interceptor quiet", () => {
    service.download(claim).subscribe();

    const request = expectPacketRequest();
    expect(request.request.method).toBe("GET");
    expect(request.request.responseType).toBe("blob");
    expect(request.request.context.get(SUPPRESS_ERROR_TOAST)).toBe(true);
    request.flush(zip);
  });

  it("saves the packet under the server's file name", async () => {
    const saved = firstValueFrom(service.download(claim));

    expectPacketRequest().flush(zip, {
      headers: { "Content-Disposition": "attachment; filename=Voucher_14_packet_2026-11-11.zip" },
    });

    await expect(saved).resolves.toBe("Voucher_14_packet_2026-11-11.zip");
    expect(downloadSpy).toHaveBeenCalledWith(expect.any(Blob), "Voucher_14_packet_2026-11-11.zip");
  });

  it("builds the server's name itself when the header is missing, dated today in the app zone", async () => {
    // 15:00 UTC on Nov 12 is still Nov 12 in UTC, but already Nov 13 in Tokyo.
    timeZone = "Asia/Tokyo";
    const saved = firstValueFrom(service.download({ id: 12, name: "Voucher 14: Rome / Lima" }));

    expectPacketRequest().flush(zip);

    await expect(saved).resolves.toBe("Voucher_14_Rome_Lima_packet_2026-11-13.zip");
    expect(downloadSpy).toHaveBeenCalledWith(expect.any(Blob), "Voucher_14_Rome_Lima_packet_2026-11-13.zip");
  });

  it("reads the size bound's message out of the blob error body", async () => {
    const saved = firstValueFrom(service.download(claim));

    expectPacketRequest().flush(
      new Blob([JSON.stringify({ errorMsg: "This claim has 61 images; a packet can include at most 60" })], {
        type: "application/json",
      }),
      { status: 400, statusText: "Bad Request" }
    );

    const error = await saved.catch((e) => e);
    expect(error).toBeInstanceOf(ClaimPacketError);
    expect(error.message).toBe("This claim has 61 images; a packet can include at most 60");
    expect(downloadSpy).not.toHaveBeenCalled();
  });

  it("reports a missing or forbidden claim in the resolver's words", async () => {
    const saved = firstValueFrom(service.download(claim));

    expectPacketRequest().flush(new Blob([JSON.stringify({ errorMsg: "Unauthorized" })]), {
      status: 403,
      statusText: "Forbidden",
    });

    await expect(saved).rejects.toThrow(CLAIM_NOT_FOUND_MESSAGE);
  });

  it("explains a proxy timeout whose body is not JSON", async () => {
    const saved = firstValueFrom(service.download(claim));

    expectPacketRequest().flush(new Blob(["<html>504 Gateway Time-out</html>"], { type: "text/html" }), {
      status: 504,
      statusText: "Gateway Time-out",
    });

    await expect(saved).rejects.toThrow(CLAIM_PACKET_TIMEOUT_MESSAGE);
  });

  it("fails cleanly on a 200 with no body", async () => {
    const saved = firstValueFrom(service.download(claim));

    expectPacketRequest().flush(null);

    await expect(saved).rejects.toThrow(CLAIM_PACKET_FAILED_MESSAGE);
    expect(downloadSpy).not.toHaveBeenCalled();
  });
});

describe("claimPacketFilename", () => {
  it.each([
    ["Voucher 14", "Voucher_14_packet_2026-11-12.zip"],
    ["  TA-2026/14 (Rome → Lima)  ", "TA_2026_14_Rome_Lima_packet_2026-11-12.zip"],
    ["Éte à Paris", "te_Paris_packet_2026-11-12.zip"],
    ["!!!", "claim_packet_2026-11-12.zip"],
    ["", "claim_packet_2026-11-12.zip"],
  ])("reduces %j as the server does", (name, expected) => {
    expect(claimPacketFilename(name, "2026-11-12")).toBe(expected);
  });

  it("keeps at most 60 characters, never ending on an underscore", () => {
    const name = "a".repeat(59) + " b" + "c".repeat(20);

    expect(claimPacketFilename(name, "2026-11-12")).toBe("a".repeat(59) + "_packet_2026-11-12.zip");
  });
});

describe("claimPacketErrorMessage", () => {
  it("uses a plain JSON errorMsg as well as a blob one", async () => {
    await expect(
      claimPacketErrorMessage(new HttpErrorResponse({ status: 500, error: { errorMsg: "Error building claim packet" } }))
    ).resolves.toBe("Error building claim packet");
  });

  it.each([
    [new HttpErrorResponse({ status: 0 }), CLAIM_PACKET_OFFLINE_MESSAGE],
    [new HttpErrorResponse({ status: 500, error: new Blob(["{}"]) }), CLAIM_PACKET_FAILED_MESSAGE],
    [new HttpErrorResponse({ status: 500, error: new Blob([JSON.stringify({ errorMsg: "  " })]) }), CLAIM_PACKET_FAILED_MESSAGE],
    [new Error("boom"), CLAIM_PACKET_FAILED_MESSAGE],
  ])("falls back by status when the body says nothing (%#)", async (error, expected) => {
    await expect(claimPacketErrorMessage(error)).resolves.toBe(expected);
  });
});
