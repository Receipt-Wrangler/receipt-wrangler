import {
  ReceiptHistory,
  ReceiptHistoryActor,
  ReceiptHistoryChange,
  ReceiptHistoryEntry,
  ReceiptHistorySnapshots,
  ReceiptHistoryValue,
  SystemTaskType,
  User,
} from "../../open-api";
import { formatAppDate } from "../../pipes/app-date.pipe";
import { todayInZone } from "../../utils/app-time-zone";
import { formatStatus } from "../../utils/status.utils";

/*
 * Pure helpers behind the receipt history drawer. The server computes the diff
 * and sends every code (action, actor type, field, kind, value type, before
 * state) as a PLAIN STRING, so each helper here has a fallback for a value it
 * does not know: a newer server must degrade to generic wording, never to a
 * blank row or a thrown error.
 */

/** An entry lists this many changes before "Show N more". */
export const HISTORY_COLLAPSE_AFTER = 4;

/** Comment text longer than this is clipped in its chip; the full text becomes the chip's title. */
export const HISTORY_COMMENT_CLIP_LENGTH = 60;

export const UNKNOWN_USER_NAME = "Unknown user";

/** Shown in a chip whose value is present but empty, so a chip is never blank. */
export const EMPTY_VALUE_TEXT = "(empty)";

export type UserLookup = (userId?: number | null) => User | undefined;

export function userDisplayName(user?: User): string {
  return user?.displayName || user?.username || UNKNOWN_USER_NAME;
}

/** `SOME_CODE` → `Some Code`: the generic label for a code this client has no wording for. */
export function titleCaseCode(code?: string | null): string {
  return code ? formatStatus(code) : "";
}

const ACTION_VERBS: Record<string, string> = {
  CREATED: "created this receipt",
  UPDATED: "edited this receipt",
  STATUS_CHANGED: "changed the status",
  RESOLVED: "resolved this receipt",
  IMAGES_COMMENTS_UPDATED: "updated images and comments",
  FAILED: "tried to save changes",
};

export const DEFAULT_ACTION_VERB = "updated this receipt";

export function historyActionVerb(action?: string | null): string {
  return (action && ACTION_VERBS[action]) || DEFAULT_ACTION_VERB;
}

// Named after the receipt form and the receipts table, so a row reads like the
// column it changed ("Receipt Date", not a bare "Date" beside "Resolved Date").
const FIELD_LABELS: Record<string, string> = {
  NAME: "Name",
  AMOUNT: "Amount",
  DATE: "Receipt Date",
  GROUP: "Group",
  PAID_BY: "Paid By",
  STATUS: "Status",
  RESOLVED_DATE: "Resolved Date",
  CATEGORIES: "Categories",
  TAGS: "Tags",
  COMMENT: "Comment",
  IMAGE: "Image",
  EMAIL_INBOX: "Inbox",
  EMAIL_SENDER: "Sender",
};

function withSubject(prefix: string, subject?: string | null): string {
  return subject ? `${prefix}: ${subject}` : prefix;
}

export function historyFieldLabel(change: ReceiptHistoryChange, findUser: UserLookup): string {
  switch (change.field) {
    case "CUSTOM_FIELD":
      return change.label || "Custom field";
    case "ITEM":
      return withSubject("Item", change.label);
    case "SHARE":
      return withSubject("Share", userDisplayName(findUser(change.userId)));
    case "ITEM_CATEGORIES":
      return withSubject("Item categories", change.label);
    case "ITEM_TAGS":
      return withSubject("Item tags", change.label);
    default:
      return FIELD_LABELS[change.field] ?? titleCaseCode(change.field);
  }
}

export type HistoryActorKind = "user" | "quick-scan" | "email" | "system";

export interface HistoryActorDisplay {
  kind: HistoryActorKind;
  name: string;
  /** The user behind the avatar, when known, so it shows their own colour and initial. */
  user?: User;
  /** A material icon shown in place of an initial. */
  icon?: string;
  /** The avatar's background; unset for a known user, whose own colour applies. */
  color?: string;
}

// The design's actor colours, from variables.scss: $accent-palette 500 / 400
// and $primary-palette 600. Hex because they are handed to app-avatar's
// `color` input, not written into a stylesheet.
const SLATE_500 = "#64748b";
const SLATE_400 = "#94a3b8";
const PRIMARY_600 = "#009efa";

export function historyActorDisplay(
  actor: ReceiptHistoryActor | undefined,
  findUser: UserLookup,
): HistoryActorDisplay {
  switch (actor?.type) {
    case "USER": {
      const user = findUser(actor.userId);
      return user
        ? { kind: "user", name: userDisplayName(user), user }
        : { kind: "user", name: UNKNOWN_USER_NAME, icon: "person", color: SLATE_400 };
    }
    case "QUICK_SCAN":
      return { kind: "quick-scan", name: "Quick Scan", icon: "photo_camera", color: PRIMARY_600 };
    case "EMAIL":
      return { kind: "email", name: "Email Upload", icon: "mail", color: SLATE_500 };
    default:
      return { kind: "system", name: "System", icon: "settings", color: SLATE_400 };
  }
}

export interface HistoryDay {
  /** The day in the app zone, `yyyy-MM-dd`; stable across renders for `@for` tracking. */
  key: string;
  label: string;
  entries: ReceiptHistoryEntry[];
}

const UNKNOWN_DAY_KEY = "unknown";

/**
 * Groups entries by the app-zone calendar day they were recorded on, keeping the
 * server's order (newest first) both across and within days. The label names
 * the year only when it is not the current one, as the design's
 * "Tuesday, Oct 13" does.
 */
export function groupHistoryByDay(
  entries: ReceiptHistoryEntry[],
  zone: string,
  locale: string,
  now: Date = new Date(),
): HistoryDay[] {
  const currentYear = todayInZone(zone, now).getFullYear().toString();
  const days = new Map<string, HistoryDay>();

  for (const entry of entries) {
    const key = formatAppDate(entry.recordedAt, "yyyy-MM-dd", "instant", zone, locale) ?? UNKNOWN_DAY_KEY;
    let day = days.get(key);
    if (!day) {
      const format = key.startsWith(currentYear) ? "EEEE, MMM d" : "EEEE, MMM d, y";
      const label = key === UNKNOWN_DAY_KEY
        ? "Unknown date"
        : (formatAppDate(entry.recordedAt, format, "instant", zone, locale) ?? key);
      day = { key, label, entries: [] };
      days.set(key, day);
    }
    day.entries.push(entry);
  }

  return Array.from(days.values());
}

function counted(count: number, singular: string, plural: string): string {
  return `${count} ${count === 1 ? singular : plural}`;
}

export function historySummary(
  history: Pick<ReceiptHistory, "changeCount" | "saveCount" | "contributorCount">,
): string {
  return [
    counted(history.changeCount ?? 0, "change", "changes"),
    counted(history.saveCount ?? 0, "save", "saves"),
    counted(history.contributorCount ?? 0, "contributor", "contributors"),
  ].join(" · ");
}

export interface HistoryNotice {
  state: "REBUILT" | "INCOMPLETE";
  text: string;
}

/**
 * The notice an entry recorded by an older server version needs: how far its
 * "before" side can be trusted. COMPLETE, absent or unknown states need none.
 */
export function historyBeforeStateNotice(
  entry: Pick<ReceiptHistoryEntry, "beforeState" | "beforeSourceType" | "beforeSourceRecordedAt">,
  zone: string,
  locale: string,
): HistoryNotice | null {
  switch (entry.beforeState) {
    case "REBUILT": {
      const source =
        entry.beforeSourceType === SystemTaskType.ReceiptUpdated
          ? "the previous update"
          : entry.beforeSourceType === SystemTaskType.ReceiptUploaded
            ? "the copy saved when it was created"
            : "an earlier copy";
      const date = formatAppDate(entry.beforeSourceRecordedAt, "mediumDate", "instant", zone, locale);
      return {
        state: "REBUILT",
        text:
          `Recorded by an older version — compared against ${source}${date ? ` (${date})` : ""}. ` +
          "Bulk status changes can also show here.",
      };
    }
    case "INCOMPLETE":
      return {
        state: "INCOMPLETE",
        text: "Recorded by an older version with no earlier complete copy. Some changes may be inaccurate.",
      };
    default:
      return null;
  }
}

/** How each value type is turned into text: the app's own pipes, injected by the caller. */
export interface HistoryValueFormatters {
  currency: (value: string) => string;
  /** A calendar day, shown as its stored day whatever the zone. */
  calendarDate: (value: string) => string | null;
  /** An instant, shown in the app zone. */
  instant: (value: string) => string | null;
  findUser: UserLookup;
}

export function formatHistoryValue(value: ReceiptHistoryValue, formatters: HistoryValueFormatters): string {
  const raw = value.value ?? "";
  switch (value.type) {
    case "TEXT":
      return raw;
    case "CURRENCY":
      return raw ? formatters.currency(raw) : "";
    case "DATE":
      return raw ? (formatters.calendarDate(raw) ?? raw) : "";
    case "INSTANT":
      return raw ? (formatters.instant(raw) ?? raw) : "";
    case "USER":
      return userDisplayName(formatters.findUser(value.userId));
    case "STATUS":
      return titleCaseCode(raw);
    case "BOOLEAN":
      return raw === "true" ? "Yes" : raw === "false" ? "No" : raw;
    default:
      return raw;
  }
}

export interface HistoryChip {
  text: string;
  /** The full text when `text` was clipped. */
  title?: string;
}

export interface HistoryChangeRow {
  /** The raw codes, for `data-field` / `data-kind`. */
  field: string;
  kind: string;
  label: string;
  from?: HistoryChip;
  to?: HistoryChip;
  arrow: boolean;
}

function clip(text: string): string {
  return text.length > HISTORY_COMMENT_CLIP_LENGTH
    ? `${text.slice(0, HISTORY_COMMENT_CLIP_LENGTH - 3)}…`
    : text;
}

function buildChip(
  change: ReceiptHistoryChange,
  value: ReceiptHistoryValue | undefined,
  formatters: HistoryValueFormatters,
): HistoryChip | undefined {
  if (!value) {
    return undefined;
  }

  const text = formatHistoryValue(value, formatters);
  if (change.field !== "COMMENT") {
    return { text: text || EMPTY_VALUE_TEXT };
  }

  // A comment reads as its author and its (clipped) words: Maya Chen: “Took the bag fee off…”
  const author = change.userId ? `${userDisplayName(formatters.findUser(change.userId))}: ` : "";
  const clipped = clip(text);
  return {
    text: `${author}“${clipped}”`,
    title: clipped === text ? undefined : `${author}“${text}”`,
  };
}

/**
 * One change as the drawer renders it. ADDED shows only the new value and
 * REMOVED only the old; CHANGED (and a kind this client does not know) shows
 * whichever sides were sent, with the arrow between two.
 */
export function buildHistoryChangeRow(
  change: ReceiptHistoryChange,
  formatters: HistoryValueFormatters,
): HistoryChangeRow {
  const from = change.kind === "ADDED" ? undefined : buildChip(change, change.from, formatters);
  const to = change.kind === "REMOVED" ? undefined : buildChip(change, change.to, formatters);

  return {
    field: change.field,
    kind: change.kind,
    label: historyFieldLabel(change, formatters.findUser),
    from,
    to,
    arrow: !!from && !!to,
  };
}

function parseSnapshot(text?: string | null): unknown {
  if (text === undefined || text === null || text === "") {
    return null;
  }

  try {
    return JSON.parse(text);
  } catch {
    // Not JSON after all: show what the server sent rather than nothing.
    return text;
  }
}

/**
 * The "View JSON" text: the receipt as created for a creation entry (which has
 * no "before"), or `{ before, after }` for a change, pretty-printed.
 */
export function historySnapshotsJson(snapshots: ReceiptHistorySnapshots): string {
  const after = parseSnapshot(snapshots.after);
  if (snapshots.before === undefined || snapshots.before === null) {
    return JSON.stringify(after, null, 2);
  }

  return JSON.stringify({ before: parseSnapshot(snapshots.before), after }, null, 2);
}
