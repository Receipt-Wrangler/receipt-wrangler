# Fork Roadmap

Features this fork means to explore. Nothing here is built yet.

**Why this fork exists:** a personal record of payments that can be handed over as proof,
exact transactions with their receipt images, independent of the system of whoever is disputing
them. Features are judged by how well they produce that evidence.

## 1. Include images in the receipt export (next up, ~1 day)

Today **Export to CSV** (filtered set or selected rows) zips `receipts.csv` + `items.csv` only;
images can only be downloaded one at a time from each receipt.

- Add an `images/` folder to the export zip, with descriptive names such as
  `2026-10-03_Costco_412_1.jpg`.
- Add an **Images** column to `receipts.csv` listing each row's filenames.
- Possibly an "include images" checkbox so a data-only export stays light.
- Add a **SHA-256** column per image to show files are unchanged since they were recorded.

Where: `api/internal/services/receipt_csv.go` (`GetZippedCsvFiles`), both handlers in
`api/internal/handlers/export.go`. Filtering and permissions are already handled there.

## 2. One-PDF evidence bundle (~2 days)

One page per transaction: date, amount, name and the image together, as a single PDF that is
easy to hand over. Reuse the HTML-to-PDF path the Report Builder already uses
(`NewHtmlToPdfService`, Chromium). Consider a cover page with totals and the date range.

## 3. Discrepancy tracking ("differentials"), later

Recording what was billed against what is actually owed, and surfacing the difference. Likely
built on custom fields plus the existing Report Builder / Receipt Summary totals rather than a new
subsystem. The receipt change history already shows when each receipt was recorded.

## 4. Multicurrency support (large, cross-cutting)

Today currency is **display-only and install-wide**: `SystemSettings.CurrencyDisplay` (symbol,
separators, position) formats every amount the same way. `Receipt.Amount`, `Item.Amount` and
CURRENCY custom field values are bare decimals with no currency attached, so mixed currencies would
silently sum as if they were one.

Likely phases:

- **Record it.** A currency code (ISO 4217) per receipt, with items inheriting the receipt's. A
  default per install or group. The AI extraction prompt returns the currency it reads, and
  displays and exports show the code next to every amount. On the wire it must be a **plain string,
  never an enum** (the closed-Dart-enum outage pattern in `CLAUDE.md`).
- **Never sum across currencies.** Receipt Summary, Report Builder, the who-owes-whom settlement
  and dashboard widgets either group totals by currency or convert. This is where most of the work
  is.
- **Convert (optional).** Store the original amount, its currency, the rate used and its date on
  the receipt rather than re-deriving them later, so a total can be explained in a dispute. That
  needs a rate source.
- **Precision.** `Receipt.Amount` is `decimal(10,2)`; JPY uses 0 decimal places and KWD/BHD use 3.

Touches all three components plus a `swagger.yml` change with both clients regenerated in the same
change. Rough size: "record it + group totals by currency" is a week or so. Conversion is extra.

## 5. Review queue for uncertain scans

When a quick scan or email upload can't confidently read something (an unclear currency, a smudged
total, an ambiguous date), flag the receipt with the specific concern instead of silently guessing.
The user then works through a review queue to resolve each concern.

What already exists:

- A `NEEDS_ATTENTION` receipt status, and per-group/per-user default statuses for quick scan and
  email (`QuickScanDefaultStatus`, `EmailDefaultReceiptStatus`). Setting the default to
  `NEEDS_ATTENTION` and filtering on it is a crude queue available today, with no code.
- Editable prompts (`api/internal/services/prompts.go`). The default prompt says "Omit any value if
  not found with confidence", so uncertainty is currently **dropped silently**, not reported.
- Receipt comments, which can carry a note to the reviewer.

Likely phases:

- **MVP (~1-2 days, backend only).** The prompt asks for a `concerns` array (field + question). If
  it is non-empty, the receipt is saved as `NEEDS_ATTENTION` with the concerns posted as a system
  comment. Both clients get this through the existing status filter and comments. Reuse the
  existing status rather than adding one: a new value in the status enum would break
  already-released mobile builds (closed-enum outages, see `CLAUDE.md`).
- **Deterministic checks alongside the AI's.** An LLM is a poor judge of its own uncertainty, so add
  rule-based flags: missing currency, line items not summing to the total, a future date, a likely
  duplicate (same store, date and amount).
- **Proper review UI (~1 week).** Store flags in their own table (field, question, resolved by and
  when) and show each on the receipt form with a control to answer it. Resolving the last flag
  moves the receipt out of the queue. Optionally rerun extraction with the answers added to the
  prompt (the true "ask back").

Pairs naturally with multicurrency (item 4): "currency unclear" is the first flag to build.

## Considered and deferred: Google Drive mirror

Nightly copy into `Receipts/<yyyy-mm>/` with a monthly spreadsheet linking to each image.
Moderate effort (OAuth with the `drive.file` scope and a published app, a sync job, a table of
what's been uploaded). **Deleting local copies after upload is out of scope**: images are read
from local disk throughout the API, so it would need a storage-backend rewrite. Revisit only if a
backup offsite becomes the priority; `rclone` on the VM is the no-code stopgap.
