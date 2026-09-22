/**
 * Reads the file name out of a `Content-Disposition` header, for a download whose
 * name only the server knows.
 *
 * Handles both the quoted and bare forms. The fallback is deliberately generic
 * rather than an error: the bytes are the point, and a download that lands under
 * an unhelpful name still beats one that does not happen.
 */
export function filenameFromContentDisposition(
  header: string | null,
  fallback: string = "download"
): string {
  if (!header) {
    return fallback;
  }

  const match = /filename\s*=\s*"([^"]+)"|filename\s*=\s*([^;]+)/i.exec(header);
  const filename = (match?.[1] ?? match?.[2] ?? "").trim();

  return filename.length > 0 ? filename : fallback;
}

export function downloadFile(blob: Blob, filename: string): void {
  const url = window.URL.createObjectURL(blob);
  const link = document.createElement("a");
  link.href = url;
  link.setAttribute("download", filename);
  document.body.appendChild(link);
  link.click();
  document.body.removeChild(link);
}
