#!/usr/bin/env bash
# Runs the Flutter integration_test suite against a local iOS Simulator.
#
# Usage:
#   cd mobile && ./run-e2e-ios.sh                                     # all tests
#   cd mobile && ./run-e2e-ios.sh integration_test/smoke_login_test.dart
#
# Resolves $E2E_IOS_DEVICE (default: iPhone 15) to a simulator UDID, boots it
# if necessary, then runs one `flutter drive` per spec. The simulator is left
# booted on exit -- shut it down (xcrun simctl shutdown <udid>) for a cold
# boot next time.
#
# Requires:
#   - Xcode + iOS Simulator (xcrun simctl on PATH)
#   - flutter on PATH (or installed under ~/Documents/flutter/bin, auto-discovered)
#   - coreutils on PATH for `gtimeout` (brew install coreutils)
#   - Go API running locally on :8081 (cd api && go run main.go)
#   - The two e2e users seeded (see mobile/CLAUDE.md "Prerequisites")
#
# Env vars (sourced from api/dev/switch-to-sqlite.sh if present):
#   E2E_ADMIN_USERNAME, E2E_ADMIN_PASSWORD, E2E_USER_USERNAME, E2E_USER_PASSWORD
#
# Overrides:
#   E2E_IOS_DEVICE        Simulator device name (default: "iPhone 15"). Multiple
#                         sims with the same name across iOS runtimes get
#                         resolved by the first match in `simctl list` order;
#                         set E2E_IOS_UDID to skip name lookup entirely.
#   E2E_IOS_UDID          Exact simulator UDID. Overrides E2E_IOS_DEVICE when set.
#   E2E_MOBILE_BASE_URL   API base URL the app hits (default: http://localhost:8081/api;
#                         iOS Simulator shares the host network, so localhost works)

set -euo pipefail

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$script_dir"

# --- locate flutter ----------------------------------------------------------
if ! command -v flutter >/dev/null 2>&1; then
  for candidate in "$HOME/Documents/flutter/bin" "$HOME/flutter/bin" "/opt/flutter/bin" "/usr/local/flutter/bin"; do
    if [[ -x "$candidate/flutter" ]]; then
      export PATH="$candidate:$PATH"
      break
    fi
  done
fi
command -v flutter >/dev/null 2>&1 || { echo "flutter not on PATH" >&2; exit 1; }
command -v xcrun >/dev/null 2>&1 || { echo "xcrun not on PATH (install Xcode command line tools)" >&2; exit 1; }
command -v gtimeout >/dev/null 2>&1 || { echo "gtimeout not found; brew install coreutils" >&2; exit 1; }
command -v python3 >/dev/null 2>&1 || { echo "python3 not on PATH (needed to safely build the dart-define JSON)" >&2; exit 1; }

# --- credentials -------------------------------------------------------------
# Source switch-to-sqlite.sh only when none of the four E2E_* creds are already
# in the env. The script `export`s all four unconditionally, so sourcing it in
# any context where a caller already populated even one of them (CI, or a
# partial override locally) would silently clobber that value with the dev
# default. The trailing `: "${VAR:?...}"` checks below catch the
# partial-population case with a clear error rather than letting defaults paper
# over a missing var.
env_script="../api/dev/switch-to-sqlite.sh"
if [[ -z "${E2E_ADMIN_USERNAME:-}" \
   && -z "${E2E_ADMIN_PASSWORD:-}" \
   && -z "${E2E_USER_USERNAME:-}" \
   && -z "${E2E_USER_PASSWORD:-}" \
   && -f "$env_script" ]]; then
  # shellcheck disable=SC1090
  source "$env_script"
fi
: "${E2E_ADMIN_USERNAME:?set E2E_ADMIN_USERNAME or source api/dev/switch-to-sqlite.sh}"
: "${E2E_ADMIN_PASSWORD:?set E2E_ADMIN_PASSWORD or source api/dev/switch-to-sqlite.sh}"
: "${E2E_USER_USERNAME:?set E2E_USER_USERNAME or source api/dev/switch-to-sqlite.sh}"
: "${E2E_USER_PASSWORD:?set E2E_USER_PASSWORD or source api/dev/switch-to-sqlite.sh}"

mobile_base_url="${E2E_MOBILE_BASE_URL:-http://localhost:8081/api}"
device_name="${E2E_IOS_DEVICE:-iPhone 15}"

# --- resolve simulator UDID --------------------------------------------------
# Match "<name> (" so "iPhone 15" doesn't catch "iPhone 15 Pro" -- after the
# device name the next literal char in the simctl output is the opening paren
# of the UDID, while sibling models have a model qualifier in between.
# `simctl list devices available` only lists devices on installed runtimes.
if [[ -n "${E2E_IOS_UDID:-}" ]]; then
  udid="$E2E_IOS_UDID"
  echo "==> Using simulator UDID from E2E_IOS_UDID: $udid"
else
  udid="$(xcrun simctl list devices available \
    | grep -E "^ *${device_name} \(" \
    | head -n 1 \
    | grep -oE '[0-9A-F]{8}-[0-9A-F]{4}-[0-9A-F]{4}-[0-9A-F]{4}-[0-9A-F]{12}')"

  if [[ -z "$udid" ]]; then
    echo "No simulator named '$device_name' found. Available devices:" >&2
    xcrun simctl list devices available >&2
    echo "Set E2E_IOS_DEVICE to one of the names above, or E2E_IOS_UDID to a specific UDID." >&2
    exit 1
  fi
fi

# Boot if not already booted. simctl bootstatus blocks until the device is
# Booted+SpringBoard-ready, which is what we want before launching the app.
state="$(xcrun simctl list devices | grep -F "$udid" | grep -oE '\((Booted|Shutdown|Shutting Down)\)' | tr -d '()' | head -n 1)"
if [[ "$state" != "Booted" ]]; then
  echo "==> Booting simulator $udid"
  xcrun simctl boot "$udid"
fi
# `simctl bootstatus -b` blocks until the device is Booted+SpringBoard-ready;
# on an unhealthy sim it can block indefinitely. Cap it at 5 min (same ceiling
# as the Android AVD boot wait) and fail fast if it expires.
if ! gtimeout 300 xcrun simctl bootstatus "$udid" -b >/dev/null; then
  echo "Simulator $udid failed to reach Booted+SpringBoard-ready within 5 minutes." >&2
  exit 1
fi

# `simctl boot` only starts the runtime; the Simulator.app GUI window doesn't
# appear unless we launch it explicitly. Useful both for visibility while a
# test is running and so the user knows the sim is alive.
open -a Simulator --args -CurrentDeviceUDID "$udid" >/dev/null 2>&1 || true

echo "==> Simulator $udid is booted"

# --- ensure dependencies resolved -------------------------------------------
# The per-spec `flutter drive` invocations pass --no-pub (faster reruns, matches
# CI). CI runs `flutter pub get` as a prior step; locally we need to do it
# ourselves once so .dart_tool/package_config.json has integration_test mapped
# -- otherwise the first build fails with
# "Couldn't resolve the package 'integration_test'".
echo "==> flutter pub get"
flutter pub get

# --- dart-define payload -----------------------------------------------------
# Build the JSON via python3 instead of a heredoc so any value containing
# quotes/backslashes/newlines gets properly escaped. The heredoc form would
# emit invalid JSON that `flutter --dart-define-from-file` would reject.
tmp_env="$(mktemp -t run-e2e-ios.XXXXXX)"
trap 'rm -f "$tmp_env"' EXIT
python3 - "$tmp_env" "$mobile_base_url" \
  "$E2E_ADMIN_USERNAME" "$E2E_ADMIN_PASSWORD" \
  "$E2E_USER_USERNAME"  "$E2E_USER_PASSWORD" <<'PY'
import json, sys
out, base, au, ap, uu, up = sys.argv[1:]
with open(out, "w", encoding="utf-8") as f:
    json.dump({
        "E2E_BASE_URL": base,
        "E2E_ADMIN_USERNAME": au,
        "E2E_ADMIN_PASSWORD": ap,
        "E2E_USER_USERNAME": uu,
        "E2E_USER_PASSWORD": up,
    }, f)
PY

# --- target specs ------------------------------------------------------------
if [[ $# -gt 0 ]]; then
  targets=("$@")
else
  targets=()
  while IFS= read -r spec; do targets+=("$spec"); done \
    < <(find integration_test -maxdepth 2 -name '*_test.dart' | sort)
fi
if [[ ${#targets[@]} -eq 0 ]]; then
  echo "No *_test.dart files found under integration_test/" >&2
  exit 1
fi

# --- per-spec drive loop -----------------------------------------------------
# One `flutter drive` per spec -- same reason as Android: the top-level
# GoRouter in main.dart persists location across testWidgets within a single
# flutter process.
#
# Between specs (and between retry attempts): terminate + uninstall by bundle
# id (io.receiptwrangler) and pkill -f dartvm + io.receiptwrangler. iOS flake
# class (flutter#129246, #136222, #153433): without these, the previous spec's
# dart isolate keeps the vmservice port and the next launch silently hangs.
# gtimeout 600 (10min, well above the ~30s legit spec runtime) walks past a
# hung spec instead of consuming the rest of the run.
#
# One retry on failure: the very first spec on a fresh CI runner pays a
# cold pod install + ~150s Xcode build, after which flutter drive's log
# reader can race the debug-server attach and abort with "The log reader
# failed unexpectedly". Subsequent rebuilds are ~55s and don't hit it.
# Retrying once catches that flake without slowing the happy path. The
# first-attempt log stays visible so genuine breakage is still observable.
prepare_sim() {
  xcrun simctl terminate "$udid" io.receiptwrangler >/dev/null 2>&1 || true
  xcrun simctl uninstall "$udid" io.receiptwrangler >/dev/null 2>&1 || true
  pkill -f dartvm >/dev/null 2>&1 || true
  pkill -f io.receiptwrangler >/dev/null 2>&1 || true
}
exit_code=0
for target in "${targets[@]}"; do
  echo ""
  echo "==> $target"
  prepare_sim
  if ! gtimeout 600 flutter drive \
       --no-pub \
       --driver=test_driver/integration_test.dart \
       --target="$target" \
       -d "$udid" \
       --dart-define-from-file="$tmp_env"; then
    echo "==> $target failed on attempt 1; retrying once after 5s..." >&2
    sleep 5
    prepare_sim
    if ! gtimeout 600 flutter drive \
         --no-pub \
         --driver=test_driver/integration_test.dart \
         --target="$target" \
         -d "$udid" \
         --dart-define-from-file="$tmp_env"; then
      exit_code=1
    fi
  fi
done

exit "$exit_code"
