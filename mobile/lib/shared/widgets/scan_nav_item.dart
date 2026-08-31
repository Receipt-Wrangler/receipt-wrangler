import 'package:flutter/material.dart';

import '../functions/receipt_entry.dart';
import '../functions/receipt_entry_availability.dart';
import 'bottom_nav.dart';

/// Identity of the scan/add slot within a bottom nav.
const scanNavDestinationId = "scan";

/// Names the hold gesture for assistive tech. The gesture is undiscoverable on
/// its own, and the receipts-screen overflow menu is its accessible equivalent.
const scanNavLongPressHint = "Press and hold for more ways to add a receipt";

/// Builds the receipt-entry slot for a bottom nav, or `null` when the user can
/// neither quick-scan nor create receipts here.
///
/// **This is the single definition of the add/scan button.** Both navs call it,
/// so the icon, the label, the visibility and what a tap or a hold does are
/// decided once — a nav that mounts the slot cannot end up with different
/// behaviour from another.
///
/// - Quick Scan available → a **Scan** destination whose tap opens the document
///   scanner immediately.
/// - Otherwise → an **Add** destination whose tap opens the manual form, which
///   explains why Quick Scan didn't run.
/// - Neither permission → `null`, and the caller omits the slot entirely rather
///   than offering an action that would only be refused.
///
/// A hold always opens the menu, which carries every entry the user is allowed.
/// [anchorKey] must be attached to the destination so the menu can position
/// itself over the slot; the caller owns it because it has to outlive rebuilds.
NavDestinationItem? buildScanNavItem(
  BuildContext context,
  GlobalKey anchorKey,
) {
  final availability = resolveReceiptEntryAvailability(context);
  if (!availability.isVisible) {
    return null;
  }

  final canQuickScan = availability.canQuickScan;

  return NavDestinationItem(
    id: scanNavDestinationId,
    onLongPress: () => showReceiptEntryMenu(context, anchorKey),
    destination: NavigationDestination(
      key: anchorKey,
      // The hint rides on Semantics rather than the destination's `tooltip`,
      // which is disabled just below -- see there for why.
      icon: Semantics(
        hint: scanNavLongPressHint,
        child: Icon(canQuickScan ? Icons.document_scanner : Icons.add),
      ),
      label: canQuickScan ? "Scan" : "Add",
      // Empty string = no tooltip. A NavigationDestination otherwise shows one
      // on **long press** (falling back to the label), which would put a second
      // long-press recognizer in the gesture arena against the one that opens
      // this slot's menu. Rather than rely on which of them wins, take Material's
      // out of the running.
      tooltip: "",
    ),
  );
}

/// The tap action for the slot [buildScanNavItem] produces.
Future<void> onScanNavItemSelected(BuildContext context) =>
    startScanEntry(context);
