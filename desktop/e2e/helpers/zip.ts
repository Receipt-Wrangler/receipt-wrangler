import { inflateRawSync } from 'node:zlib';

/**
 * Reads every file in a zip, for asserting on a downloaded archive without a zip
 * dependency (the project has none). Walks the central directory, which is
 * authoritative: Go's archive/zip writes data descriptors, so the sizes in each
 * local header are zero. Handles the two methods the server uses, Store (0) and
 * Deflate (8); no ZIP64, no encryption, which a test fixture never needs.
 */
export function readZipEntries(zip: Buffer): Map<string, Buffer> {
  const endOfCentralDirectory = findEndOfCentralDirectory(zip);
  const entryCount = zip.readUInt16LE(endOfCentralDirectory + 10);
  let offset = zip.readUInt32LE(endOfCentralDirectory + 16);

  const entries = new Map<string, Buffer>();
  for (let i = 0; i < entryCount; i++) {
    if (zip.readUInt32LE(offset) !== 0x02014b50) {
      throw new Error(`zip: no central directory header at ${offset}`);
    }

    const method = zip.readUInt16LE(offset + 10);
    const compressedSize = zip.readUInt32LE(offset + 20);
    const nameLength = zip.readUInt16LE(offset + 28);
    const extraLength = zip.readUInt16LE(offset + 30);
    const commentLength = zip.readUInt16LE(offset + 32);
    const localHeader = zip.readUInt32LE(offset + 42);
    const name = zip.toString('utf8', offset + 46, offset + 46 + nameLength);

    if (zip.readUInt32LE(localHeader) !== 0x04034b50) {
      throw new Error(`zip: no local header for ${name}`);
    }
    const dataStart =
      localHeader + 30 + zip.readUInt16LE(localHeader + 26) + zip.readUInt16LE(localHeader + 28);
    const data = zip.subarray(dataStart, dataStart + compressedSize);

    switch (method) {
      case 0:
        entries.set(name, Buffer.from(data));
        break;
      case 8:
        entries.set(name, inflateRawSync(data));
        break;
      default:
        throw new Error(`zip: ${name} uses unsupported method ${method}`);
    }

    offset += 46 + nameLength + extraLength + commentLength;
  }

  return entries;
}

function findEndOfCentralDirectory(zip: Buffer): number {
  // The record is 22 bytes plus a comment of at most 65535.
  const earliest = Math.max(0, zip.length - 22 - 0xffff);
  for (let offset = zip.length - 22; offset >= earliest; offset--) {
    if (zip.readUInt32LE(offset) === 0x06054b50) {
      return offset;
    }
  }

  throw new Error('zip: no end of central directory record');
}
