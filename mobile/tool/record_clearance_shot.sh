#!/usr/bin/env bash
# Regenerates tool/receipt-form-clearance.png — the before/after still for the
# floating-submit-button clearance fix.
# See mobile/CLAUDE.md -> "Recording a demo GIF" -> Route 2.
#
# The keyboard demo can flip a `debugDisableKeyboardLift` seam to record its
# "before" panel. There is no such flag for `submitButtonSpacing` and adding one
# to production for a screenshot is not worth it, so this captures the AFTER
# panel from the real tree, then temporarily strips the spacer and captures
# BEFORE.
#
# The trap is load-bearing: without it a crash between the two runs leaves
# receipt_form.dart patched in the working tree.
set -euo pipefail
cd "$(dirname "$0")/.."

FORM="lib/receipts/widgets/receipt_form.dart"
OUT="${1:-tool/receipt-form-clearance.png}"
WORK="$(mktemp -d)"
# Registered before the backup cp, not after: under `set -e` a failed cp exits
# immediately, and a trap installed below it would never run -- leaking $WORK.
# Replaced by the full cleanup once $BACKUP actually exists.
trap 'rm -rf "$WORK"' EXIT
BACKUP="$WORK/receipt_form.dart.orig"

cp "$FORM" "$BACKUP"
cleanup() {
  cp "$BACKUP" "$FORM"
  rm -rf "$WORK"
}
trap cleanup EXIT

echo "==> capturing AFTER (spacer present)"
CLEARANCE_SHOT_OUT="$WORK/after.png" CLEARANCE_SHOT_LABEL=AFTER \
  flutter test tool/demo_capture/clearance_shot_test.dart

echo "==> stripping submitButtonSpacing"
python3 - "$FORM" <<'PY'
import io, sys
p = sys.argv[1]
s = io.open(p, encoding='utf-8').read()
needle = "          submitButtonSpacing,\n"
assert s.count(needle) == 1, 'expected exactly one submitButtonSpacing in the form'
io.open(p, 'w', encoding='utf-8').write(s.replace(needle, "", 1))
PY

echo "==> capturing BEFORE (spacer removed)"
CLEARANCE_SHOT_OUT="$WORK/before.png" CLEARANCE_SHOT_LABEL=BEFORE \
  flutter test tool/demo_capture/clearance_shot_test.dart

# Restore before stitching, so the working tree is clean even if the stitch
# fails. The trap would also do it; this just makes the window as short as
# possible.
cp "$BACKUP" "$FORM"

echo "==> stitching -> $OUT"
dart run tool/demo_capture/stitch_shot.dart \
  "$WORK/before.png" "$WORK/after.png" "$OUT"

ls -la "$OUT"
