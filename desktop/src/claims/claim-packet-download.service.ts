import { HttpContext, HttpErrorResponse } from "@angular/common/http";
import { inject, Injectable } from "@angular/core";
import { Store } from "@ngxs/store";
import { catchError, from, map, Observable, switchMap, throwError } from "rxjs";
import { SUPPRESS_ERROR_TOAST } from "../interceptors/http-interceptor";
import { Claim, ReimbursementClaimService } from "../open-api";
import { SystemSettingsState } from "../store/system-settings.state";
import { toCalendarDay, todayInZone } from "../utils/app-time-zone";
import { downloadFile, filenameFromContentDisposition } from "../utils/file";
import { CLAIM_NOT_FOUND_MESSAGE } from "./claim.resolver";

/** What the packet holds, next to the plain CSV export. Shown on the download menu. */
export const CLAIM_PACKET_DESCRIPTION =
  "PDF for the finance office, original receipt images, CSVs and a SHA-256 manifest";

export const CLAIM_PACKET_FAILED_MESSAGE = "The claim packet could not be built. Try again in a moment.";

/** nginx gives up on /api/ after 60 s, and nothing is sent until the zip is complete. */
export const CLAIM_PACKET_TIMEOUT_MESSAGE =
  "The claim packet took too long to build. Try again, or split the claim into smaller ones.";

export const CLAIM_PACKET_OFFLINE_MESSAGE = "The server could not be reached, so no packet was downloaded.";

/** A failed packet download, carrying the message to show the user. */
export class ClaimPacketError extends Error {}

/**
 * Downloads a claim's packet (`GET /claim/{claimId}/packet`): a zip of a PDF,
 * the original images, CSVs and a manifest of SHA-256 hashes.
 *
 * The request opts out of the HTTP interceptor's toasts (`SUPPRESS_ERROR_TOAST`)
 * because its error body is a Blob the interceptor cannot read: a failure is
 * reported once, by the caller, from the `ClaimPacketError` this emits.
 */
@Injectable({ providedIn: "root" })
export class ClaimPacketDownloadService {
  private readonly claimService = inject(ReimbursementClaimService);
  private readonly store = inject(Store);

  /** Downloads the packet and emits the file name it was saved under. */
  public download(claim: Pick<Claim, "id" | "name">): Observable<string> {
    const context = new HttpContext().set(SUPPRESS_ERROR_TOAST, true);

    return this.claimService.getClaimPacket(claim.id, "response", false, { context }).pipe(
      map((response) => {
        if (!response.body) {
          throw new ClaimPacketError(CLAIM_PACKET_FAILED_MESSAGE);
        }

        // The server names the file; the fallback is the same rule, for a proxy
        // that strips the header.
        const filename = filenameFromContentDisposition(
          response.headers.get("Content-Disposition"),
          claimPacketFilename(claim.name, this.today())
        );
        downloadFile(response.body, filename);
        return filename;
      }),
      catchError((error) =>
        error instanceof ClaimPacketError
          ? throwError(() => error)
          : from(claimPacketErrorMessage(error)).pipe(
              switchMap((message) => throwError(() => new ClaimPacketError(message)))
            )
      )
    );
  }

  private today(): string {
    return toCalendarDay(todayInZone(this.store.selectSnapshot(SystemSettingsState.timeZone)));
  }
}

/**
 * `<ClaimName>_packet_<yyyy-MM-dd>.zip`, mirroring the server's
 * `claimPacketFilename`: the name reduced to runs of ASCII letters and digits
 * joined by `_`, at most 60 characters, `claim` when nothing is left.
 */
export function claimPacketFilename(claimName: string, day: string): string {
  let safe = claimName.replace(/[^A-Za-z0-9]+/g, "_").replace(/^_+|_+$/g, "");
  if (safe.length > 60) {
    safe = safe.slice(0, 60).replace(/_+$/, "");
  }

  return `${safe || "claim"}_packet_${day}.zip`;
}

/**
 * The message for a failed packet request. The server's `{ errorMsg }` (the size
 * bound's 400 says how many images are over) arrives as a Blob because the call
 * asks for one, so it is read and parsed here. A 403 is the resolver's wording:
 * a missing claim and a forbidden one are the same answer.
 */
export async function claimPacketErrorMessage(error: unknown): Promise<string> {
  if (!(error instanceof HttpErrorResponse)) {
    return CLAIM_PACKET_FAILED_MESSAGE;
  }

  if (error.status === 403) {
    return CLAIM_NOT_FOUND_MESSAGE;
  }

  const serverMessage = await errorMsgFromBody(error.error);
  if (serverMessage) {
    return serverMessage;
  }

  switch (error.status) {
    case 0:
      return CLAIM_PACKET_OFFLINE_MESSAGE;
    case 504:
      return CLAIM_PACKET_TIMEOUT_MESSAGE;
    default:
      return CLAIM_PACKET_FAILED_MESSAGE;
  }
}

async function errorMsgFromBody(body: unknown): Promise<string | undefined> {
  let parsed: unknown = body;
  if (body instanceof Blob) {
    try {
      parsed = JSON.parse(await readBlobText(body));
    } catch {
      // Not JSON: a proxy's HTML error page, say. The status decides instead.
      return undefined;
    }
  }

  const message = (parsed as { errorMsg?: unknown } | null)?.errorMsg;
  return typeof message === "string" && message.trim() ? message : undefined;
}

// FileReader rather than Blob.text(): the same in every browser, and jsdom's Blob
// has no text().
function readBlobText(blob: Blob): Promise<string> {
  return new Promise((resolve, reject) => {
    const reader = new FileReader();
    reader.onload = () => resolve(reader.result as string);
    reader.onerror = () => reject(reader.error);
    reader.readAsText(blob);
  });
}
