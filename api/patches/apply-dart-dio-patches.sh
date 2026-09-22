#!/bin/bash
#
# Re-applies the known dart-dio generator patches to a freshly generated mobile client.
#
# The openapi-generator emits four things the app cannot use as generated. Two are outright
# invalid Dart and fail the build; the other two are the `fallback: true` annotations that keep
# a closed built_value EnumClass from throwing on an unrecognized wire value -- which fails the
# WHOLE enclosing payload, and for GroupReceiptSettings that payload is AppData, i.e. LOGIN.
# Those two compile perfectly well without the patch, so a regen that dropped them used to be
# caught only by `flutter test` (mobile/test/models/*_ingest_test.dart), long after the fact and
# only if someone ran it.
#
# Running this from generate-client.sh is what makes the patches part of generation rather than
# a step in a checklist. The important behaviour is the third case in apply(): when neither the
# original nor the patched text is present, the generator's output has changed shape and this
# script FAILS rather than quietly handing back an unpatched client.
#
# Adding a patch: add an apply() call. The match text must appear EXACTLY once in the file.

set -euo pipefail

output_dir=${1:?usage: apply-dart-dio-patches.sh <mobile-client-dir>}

python3 - "$output_dir" <<'PY'
import io
import os
import sys

output_dir = sys.argv[1]

FALLBACK_NOTE = (
    "// Applied by api/patches/apply-dart-dio-patches.sh -- do not hand-edit, and do not\n"
    "  // drop it as generator noise. Without `fallback: true` the generated _$valueOf throws\n"
    "  // on an unrecognized wire value, failing the WHOLE enclosing payload rather than the\n"
    "  // one field. See mobile/CLAUDE.md for the two outages that came of it.\n"
    "  "
)

# (file, text the generator emits, text the app needs)
PATCHES = [
    # quickScanDefaultStatus is typed ReceiptStatus, but the default is emitted as a bare string.
    (
        "lib/src/model/user_preferences.dart",
        "..quickScanDefaultStatus = 'OPEN'",
        "..quickScanDefaultStatus = ReceiptStatus.OPEN",
    ),
    # A bare $ in a non-raw Dart string is an interpolation.
    (
        "lib/src/model/system_settings.dart",
        "..currencyDisplay = '$'",
        "..currencyDisplay = r'$'",
    ),
    # The two enums whose unknown-value tolerance is load-bearing. See mobile/CLAUDE.md ->
    # "Known dart-dio default-value regressions" for why only these two are opened up.
    (
        "lib/src/model/receipt_status.dart",
        "@BuiltValueEnumConst(wireName: r'')",
        FALLBACK_NOTE + "@BuiltValueEnumConst(wireName: r'', fallback: true)",
    ),
    (
        "lib/src/model/receipt_summary_position.dart",
        "@BuiltValueEnumConst(wireName: r'BOTTOM')",
        FALLBACK_NOTE + "@BuiltValueEnumConst(wireName: r'BOTTOM', fallback: true)",
    ),
]

applied = 0
skipped = 0

for relative_path, original, replacement in PATCHES:
    path = os.path.join(output_dir, relative_path)

    if not os.path.isfile(path):
        sys.stderr.write(
            "Error: dart-dio patch target is missing: %s\n"
            "       The generator no longer emits this file. Re-derive the patch.\n" % path)
        sys.exit(1)

    source = io.open(path, encoding="utf-8").read()

    # Test the PATCHED form first. An already-patched file no longer contains the original
    # text, so checking that first would take the "no longer applies" branch below and fail a
    # perfectly good tree -- which is what makes re-running this script safe.
    if replacement in source:
        skipped += 1
        continue

    count = source.count(original)
    if count == 0:
        sys.stderr.write(
            "Error: dart-dio patch no longer applies in %s\n"
            "       Expected to find: %s\n"
            "       The generator's output has changed shape. Re-derive this patch by hand\n"
            "       rather than skipping it -- see mobile/CLAUDE.md.\n" % (relative_path, original))
        sys.exit(1)
    if count > 1:
        sys.stderr.write(
            "Error: dart-dio patch is ambiguous in %s -- matched %d times, expected 1.\n"
            "       Narrow the match text.\n" % (relative_path, count))
        sys.exit(1)

    io.open(path, "w", encoding="utf-8").write(source.replace(original, replacement))
    applied += 1

print("dart-dio patches: %d applied, %d already present" % (applied, skipped))
PY
