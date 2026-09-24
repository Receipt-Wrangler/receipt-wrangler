import { filenameFromContentDisposition } from "./file";

describe("file.utils", () => {
  describe("filenameFromContentDisposition", () => {
    it("reads the quoted form", () => {
      expect(
        filenameFromContentDisposition('attachment; filename="receipt.png"')
      ).toBe("receipt.png");
    });

    it("reads the bare form", () => {
      expect(filenameFromContentDisposition("attachment; filename=receipt.png")).toBe(
        "receipt.png"
      );
    });

    // The server sanitizes an email attachment's name to a basename but does not
    // strip quotes, so mime.FormatMediaType escapes them. Reading only to the
    // first quote truncated the name and lost its extension.
    it("keeps a name containing an escaped quote whole", () => {
      expect(
        filenameFromContentDisposition('attachment; filename="receipt\\"final.pdf"')
      ).toBe('receipt"final.pdf');
    });

    // FormatMediaType emits this form INSTEAD of the plain parameter for any
    // non-ASCII name, so without it every such download landed on the fallback.
    it("reads and percent-decodes the RFC 5987 extended form", () => {
      expect(
        filenameFromContentDisposition(
          "attachment; filename*=utf-8''re%C3%A7u%20final.pdf"
        )
      ).toBe("reçu final.pdf");
    });

    it("prefers the extended form when both are present", () => {
      expect(
        filenameFromContentDisposition(
          `attachment; filename="recu.pdf"; filename*=utf-8''re%C3%A7u.pdf`
        )
      ).toBe("reçu.pdf");
    });

    // A malformed percent sequence must not throw out of a download handler.
    it("falls back past an undecodable extended value", () => {
      expect(
        filenameFromContentDisposition(
          `attachment; filename="recu.pdf"; filename*=utf-8''re%zzu.pdf`
        )
      ).toBe("recu.pdf");
    });

    it("tolerates whitespace around the separator", () => {
      expect(
        filenameFromContentDisposition('attachment; filename = "receipt.png"')
      ).toBe("receipt.png");
    });

    it("falls back when the header is absent, empty or nameless", () => {
      expect(filenameFromContentDisposition(null)).toBe("download");
      expect(filenameFromContentDisposition("")).toBe("download");
      expect(filenameFromContentDisposition("attachment")).toBe("download");
      expect(filenameFromContentDisposition('attachment; filename="   "')).toBe(
        "download"
      );
    });

    it("honours a caller-supplied fallback", () => {
      expect(filenameFromContentDisposition(null, "source-file")).toBe("source-file");
    });
  });
});
