import { downloadFile } from "../utils/file";
import { captureDownloads } from "./download.testing";

describe("captureDownloads", () => {
  it("records what downloadFile hands the browser", () => {
    const downloads = captureDownloads();
    const blob = new Blob(["bytes"]);

    downloadFile(blob, "receipt.pdf");

    expect(downloads).toHaveLength(1);
    expect(downloads[0].blob).toBe(blob);
    expect(downloads[0].filename).toBe("receipt.pdf");
  });

  it("lets clicks on other links through", () => {
    const downloads = captureDownloads();
    const link = document.createElement("a");
    link.setAttribute("href", "/receipts/7/view");
    const listener = vi.fn((event: Event) => event.preventDefault());
    link.addEventListener("click", listener);

    link.click();

    expect(listener).toHaveBeenCalledTimes(1);
    expect(downloads).toEqual([]);
  });
});
