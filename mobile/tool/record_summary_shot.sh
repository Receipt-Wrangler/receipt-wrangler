#!/usr/bin/env bash
# Regenerates tool/receipt-summary-bar.png — the reference still for the receipt
# summary bar, in the two configurations that look least alike.
# See mobile/CLAUDE.md -> "Recording a demo GIF" -> Route 2.
#
# A still rather than a frame of the placement GIF: `writeSideBySideGif` quantises
# to an octree palette, which bands badly on flat UI, so a GIF frame is a poor
# thing to judge spacing and weight from. A PNG is what the widget drew.
#
# Left panel is the shape most groups have (a couple of statuses, no currency
# fields); right is the heavy one (every status broken out, two currency fields),
# which is the case that used to clip its own numbers and is the one worth a look
# after any change to the row layout.
set -euo pipefail
cd "$(dirname "$0")/.."

OUT="${1:-tool/receipt-summary-bar.png}"
WORK="$(mktemp -d)"
trap 'rm -rf "$WORK"' EXIT

echo "==> capturing the typical configuration"
SUMMARY_SHOT_VARIANT=typical SUMMARY_SHOT_OUT="$WORK/typical.png" \
  flutter test tool/demo_capture/summary_shot_test.dart

echo "==> capturing the heavy configuration"
SUMMARY_SHOT_VARIANT=heavy SUMMARY_SHOT_OUT="$WORK/heavy.png" \
  flutter test tool/demo_capture/summary_shot_test.dart

echo "==> stitching -> $OUT"
dart run tool/demo_capture/stitch_shot.dart \
  "$WORK/typical.png" "$WORK/heavy.png" "$OUT"

ls -la "$OUT"
