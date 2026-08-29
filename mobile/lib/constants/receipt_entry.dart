/// User-facing copy for the receipt-entry affordances (the Scan/Add bottom-nav
/// slot, its long-press menu, the receipts-screen overflow menu, and the
/// "Quick Scan unavailable" banner on the manual form).
///
/// Centralised because the same reason is surfaced in more than one place — the
/// snackbar shown when the Quick Scan sheet is opened directly, and the banner
/// shown when the Scan tap falls through to manual entry — and the two drifting
/// apart would read as two different problems.

const quickScanAiDisabledMessage =
    "A configured Receipt Processing Settings is required to use Quick Scan. "
    "Contact your administrator for more information.";

const quickScanNoPermissionMessage =
    "You don't have Quick Scan permission here, so this opens manual entry "
    "instead.";

/// [quickScanNoPermissionMessage]'s wording when the group being acted on is
/// known (i.e. the user is inside a single group rather than on group-select).
String quickScanNoPermissionMessageForGroup(String groupName) =>
    "You don't have Quick Scan permission in $groupName, so this opens manual "
    "entry instead.";

const quickScanUnavailableTitle = "Quick Scan unavailable";

const noReceiptEntryPermissionMessage =
    "You don't have permission to add receipts here.";

const cameraDeniedFallbackMessage =
    "Camera access is off — pick from your gallery instead.";

const galleryUnavailableMessage =
    "Couldn't open the gallery on this device.";

const addManualReceiptLabel = "Add Manual Receipt";
const quickScanLabel = "Quick Scan";
const uploadFromGalleryLabel = "Upload from Gallery";
const enterDetailsManuallyLabel = "Enter details manually instead";

/// Extra key carrying a [QuickScanBlockedReason] into `/receipts/add`, so the
/// form can explain why the Scan tap landed there. Only the tap path sets it —
/// a deliberate "Add Manual Receipt" never shows the banner.
const quickScanBlockedReasonExtraKey = "quickScanBlockedReason";

/// Extra key carrying the group name the blocked reason refers to.
const quickScanBlockedGroupExtraKey = "quickScanBlockedGroup";
