#!/usr/bin/env bash
# Regenerates the quick-date-filter demo GIF.
# See mobile/CLAUDE.md -> "Recording a demo GIF" -> Route 2.
#
# Like record_keyboard_demo.sh this needs no ffmpeg, xdotool, Xvfb or a Linux
# desktop build: the frames are captured inside a widget test and encoded to
# GIF in pure Dart. Route 1 is not an option here for a different reason than
# the keyboard demos -- driving the real control needs the Go API behind it,
# which the Claude Code sandbox cannot bring up (ImageMagick 7 from source).
#
# The explicit path is what keeps this out of CI: `flutter test` with no
# arguments scans only test/, so the demo never runs (and never rewrites the
# committed GIF) on a normal push.
set -euo pipefail
cd "$(dirname "$0")/.."

flutter test tool/demo_capture/quick_date_demo_test.dart

ls -la tool/quick-date-filter.gif
