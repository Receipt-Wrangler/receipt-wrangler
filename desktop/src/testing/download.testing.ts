export interface CapturedDownload {
  blob: Blob;
  filename: string;
}

/**
 * Records every download `downloadFile` (`src/utils/file.ts`) starts during the current test, in
 * order, without the browser.
 *
 * The Angular unit-test builder refuses `vi.mock` of a relative import, so `downloadFile` itself
 * can't be replaced. This stubs the DOM side it drives instead: the object URL it makes for the blob
 * and the anchor it clicks. Call it inside the test; the stubs are restored when the test finishes.
 */
export function captureDownloads(): CapturedDownload[] {
  const downloads: CapturedDownload[] = [];
  const blobsByUrl = new Map<string, Blob>();

  const createObjectUrl = vi.spyOn(URL, "createObjectURL").mockImplementation((blob) => {
    const url = `blob:test/${blobsByUrl.size}`;
    blobsByUrl.set(url, blob as Blob);
    return url;
  });
  const click = vi
    .spyOn(HTMLAnchorElement.prototype, "click")
    .mockImplementation(function (this: HTMLAnchorElement) {
      downloads.push({
        blob: blobsByUrl.get(this.getAttribute("href") ?? "") as Blob,
        filename: this.getAttribute("download") ?? "",
      });
    });

  onTestFinished(() => {
    createObjectUrl.mockRestore();
    click.mockRestore();
  });

  return downloads;
}
