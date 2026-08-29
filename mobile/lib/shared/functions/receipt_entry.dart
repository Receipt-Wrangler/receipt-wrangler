import 'package:cunning_document_scanner/cunning_document_scanner.dart'
    show CunningDocumentScannerException;
import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:permission_handler/permission_handler.dart' show openAppSettings;

import '../../constants/receipt_entry.dart';
import '../../interfaces/upload_multipart_file_data.dart';
import '../../utils/group.dart';
import '../../utils/permissions.dart';
import '../../utils/scan.dart';
import '../../utils/snackbar.dart';
import 'quick_scan.dart';
import 'receipt_entry_availability.dart';

/// The actions behind every receipt-entry affordance.
///
/// The Scan/Add bottom-nav slot, its long-press menu and the receipts-screen
/// overflow menu all route through here, so what a tap *does* is decided in one
/// place from one [ReceiptEntryAvailability] reading. Each action re-checks the
/// gate it needs rather than trusting its caller — the affordances are hidden in
/// the states where they wouldn't work, and these checks are the backstop for a
/// stale one.

/// The Scan slot's tap.
///
/// Quick Scan available → straight into the document scanner (no menu in
/// between); blocked → straight into the manual form, carrying the reason so it
/// can be explained there.
Future<void> startScanEntry(BuildContext context) async {
  final availability = resolveReceiptEntryAvailability(context);

  if (!availability.canQuickScan) {
    if (!availability.canCreateManual) {
      // Unreachable through the UI (the slot is hidden in this state), kept as
      // the backstop for a stale permission set.
      showErrorSnackbar(context, noReceiptEntryPermissionMessage);
      return;
    }
    openManualReceipt(context,
        reason: availability.blockedReason, groupName: availability.groupName);
    return;
  }

  final access = await ensureCameraAccess();
  if (!context.mounted) {
    return;
  }

  switch (access) {
    case CameraAccess.granted:
      await openQuickScanFromCamera(context);
    case CameraAccess.denied:
      // The user declined this time but can be asked again, so no settings
      // detour — just carry on with the source that still works.
      await fallBackToGallery(context, offerSettings: false);
    case CameraAccess.permanentlyDenied:
      await fallBackToGallery(context, offerSettings: true);
  }
}

/// Captures with the document scanner, then opens the sheet seeded with the
/// pages. A cancelled scan yields nothing and is treated as "never mind" — no
/// empty sheet, no error.
Future<void> openQuickScanFromCamera(BuildContext context) async {
  final List<UploadMultipartFileData> images;
  try {
    images = await scanImagesMultiPart(100);
  } on CunningDocumentScannerException catch (_) {
    // The scanner re-checks camera permission itself and throws when it is
    // missing. [ensureCameraAccess] normally settles that first, but the grant
    // can be revoked between the two calls (or refused inside the scanner's own
    // prompt) — so land on the same fallback rather than an unhandled error.
    if (context.mounted) {
      await fallBackToGallery(context, offerSettings: true);
    }
    return;
  }

  if (images.isEmpty || !context.mounted) {
    return;
  }
  showQuickScanBottomSheet(context,
      initialImages: buildQuickScanImages(context, images));
}

/// Explains that the camera is unavailable and continues with the gallery, the
/// one image source that is still open to the user.
///
/// [offerSettings] adds the OS settings shortcut, for the states a re-request
/// cannot recover from — there the prompt resolves instantly with no dialog, so
/// without this the tap would look like it did nothing.
Future<void> fallBackToGallery(
  BuildContext context, {
  required bool offerSettings,
}) async {
  showInfoSnackbar(
    context,
    cameraDeniedFallbackMessage,
    action: offerSettings
        ? SnackBarAction(label: "Settings", onPressed: openAppSettings)
        : null,
  );
  await openQuickScanFromGallery(context);
}

/// Picks from the gallery, then opens the sheet seeded with the selection.
Future<void> openQuickScanFromGallery(BuildContext context) async {
  final images = await pickGalleryImages(context);
  if (images.isEmpty || !context.mounted) {
    return;
  }
  showQuickScanBottomSheet(context,
      initialImages: buildQuickScanImages(context, images));
}

/// [getGalleryImages] with the unsupported-platform throw turned into a message.
///
/// It hard-throws off android/ios (`lib/utils/scan.dart`'s
/// `Platform.operatingSystem` switch), which the nav's camera-denied fallback
/// now makes reachable — an unhandled async error there would surface as a red
/// screen rather than an explanation.
Future<List<UploadMultipartFileData>> pickGalleryImages(
    BuildContext context) async {
  try {
    return await getGalleryImages();
  } catch (_) {
    if (context.mounted) {
      showErrorSnackbar(context, galleryUnsupportedPlatformMessage);
    }
    return [];
  }
}

/// Opens the manual receipt form.
///
/// [reason] is passed only by the Scan tap that fell through to manual entry —
/// a deliberate "Add Manual Receipt" leaves it null so the form shows no banner.
void openManualReceipt(
  BuildContext context, {
  QuickScanBlockedReason? reason,
  String? groupName,
}) {
  context.go("/receipts/add", extra: {
    // Preserved because `getGroupId` reads it off `extra` when the route has no
    // `groupId` path parameter.
    "groupId": getGroupId(context),
    if (reason != null) quickScanBlockedReasonExtraKey: reason,
    if (reason != null && groupName != null)
      quickScanBlockedGroupExtraKey: groupName,
  });
}

/// The Scan slot's long-press menu, anchored on the slot itself.
///
/// Every item is gated on the permission the action needs: manual entry on
/// `group.receipts.create`, Quick Scan and gallery upload on the AI flag plus
/// `group.receipts.quick-scan` (the gallery flow feeds Quick Scan, so it needs
/// the quick-scan permission rather than create).
void showReceiptEntryMenu(BuildContext context, GlobalKey anchorKey) {
  final renderObject = anchorKey.currentContext?.findRenderObject();
  if (renderObject is! RenderBox) {
    return;
  }
  final Offset offset = renderObject.localToGlobal(Offset.zero);
  final Size size = renderObject.size;

  final RelativeRect position = RelativeRect.fromLTRB(
    offset.dx,
    offset.dy,
    offset.dx + size.width,
    offset.dy + size.height,
  );

  showMenu(
    context: context,
    position: position,
    items: buildReceiptEntryMenuItems(context),
  );
}

/// The receipt-entry items, shared by the long-press menu and the receipts
/// screen's overflow menu so the two can never offer different actions.
///
/// Returns an empty list when the user can do neither — callers decide whether
/// that means "hide the button" (the overflow) or "explain" (the nav backstop).
List<PopupMenuEntry> buildReceiptEntryMenuItems(BuildContext context) {
  final availability = resolveReceiptEntryAvailability(context);

  return <PopupMenuEntry>[
    if (availability.canQuickScan)
      PopupMenuItem(
        value: 0,
        onTap: () => startScanEntry(context),
        child: const Text(quickScanLabel),
      ),
    if (availability.canCreateManual)
      PopupMenuItem(
        value: 1,
        onTap: () => openManualReceipt(context),
        child: const Text(addManualReceiptLabel),
      ),
    if (availability.canQuickScan)
      PopupMenuItem(
        value: 2,
        onTap: () => openQuickScanFromGallery(context),
        child: const Text(uploadFromGalleryLabel),
      ),
  ];
}
