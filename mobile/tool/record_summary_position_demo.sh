#!/usr/bin/env bash
# Regenerates the receipt-summary placement GIF.
# See mobile/CLAUDE.md -> "Recording a demo GIF" -> Route 2.
#
# Like record_keyboard_demo.sh this needs no ffmpeg, xdotool, Xvfb or a Linux
# desktop build: frames are captured inside a widget test and encoded to GIF in
# pure Dart.
#
# Unlike BOTH existing demos it needs no production seam either -- no debug flag,
# no source patched and restored under a trap. The position is data on the summary
# response, so the two panels are the identical real tree fed two responses that
# differ in one field. Verify `git diff lib/` is empty after running it; it should
# never have been non-empty.
#
# The explicit path is what keeps this out of CI: `flutter test` with no arguments
# scans only test/, so the demo never runs (and never rewrites the committed GIF)
# on a normal push.
set -euo pipefail
cd "$(dirname "$0")/.."

flutter test tool/demo_capture/summary_position_demo_test.dart

ls -la tool/receipt-summary-position.gif
