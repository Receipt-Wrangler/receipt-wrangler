import 'package:flutter/material.dart';

import '../../utils/receipt_history.dart';

/// The 36px avatar beside a history entry: the user's own colour and initial,
/// or an icon for Quick Scan, Email Upload, System and an unknown user.
///
/// Not the shared `UserAvatar`, which uses the theme colour rather than the
/// user's `defaultAvatarColor` and asserts the user exists -- history regularly
/// names users who were deleted or are hidden from the viewer.
class ReceiptHistoryActorAvatar extends StatelessWidget {
  const ReceiptHistoryActorAvatar({super.key, required this.actor});

  final HistoryActorDisplay actor;

  static const double size = 36;

  @override
  Widget build(BuildContext context) {
    final background = actor.color ??
        parseAvatarColor(actor.user?.defaultAvatarColor) ??
        Theme.of(context).colorScheme.primary;
    final initial = actor.name.isNotEmpty ? actor.name[0].toUpperCase() : '?';

    // The name is printed beside it, so the avatar is decoration.
    return ExcludeSemantics(
      child: Container(
        width: size,
        height: size,
        alignment: Alignment.center,
        decoration: BoxDecoration(color: background, shape: BoxShape.circle),
        child: actor.icon != null
            ? Icon(actor.icon, size: 20, color: Colors.white)
            : Text(
                initial,
                style: const TextStyle(
                  color: Colors.white,
                  fontSize: 15,
                  fontWeight: FontWeight.w600,
                ),
              ),
      ),
    );
  }
}
