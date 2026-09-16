#!/usr/bin/env bash
# Regenerates the keyboard-inset demo GIFs.
# See mobile/CLAUDE.md -> "Recording a demo GIF" -> Route 2.
#
# Unlike record_tap_target_demo.sh this needs no ffmpeg, xdotool, Xvfb or a
# Linux desktop build: the frames are captured inside a widget test and encoded
# to GIF in pure Dart. It has to be that way -- there is no software keyboard on
# Linux desktop, so `viewInsets.bottom` is permanently 0 there and the bug this
# demonstrates cannot be reproduced on that target at all.
#
# The explicit path is what keeps this out of CI: `flutter test` with no
# arguments scans only test/, so the demo never runs (and never rewrites the
# committed GIFs) on a normal push.
set -euo pipefail
cd "$(dirname "$0")/.."

flutter test tool/keyboard_demo/keyboard_demo_test.dart

ls -la tool/*-keyboard.gif
