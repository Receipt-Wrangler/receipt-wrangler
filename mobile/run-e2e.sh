#!/usr/bin/env bash
# Runs the Flutter integration_test suite against the Linux desktop target.
#
# Usage:
#   cd mobile && ./run-e2e.sh               # all tests in integration_test/
#   cd mobile && ./run-e2e.sh integration_test/smoke_login_test.dart
#
# Requires:
#   - flutter on PATH with `linux` desktop enabled
#       (flutter config --enable-linux-desktop)
#   - Go API running locally on :8081 (cd api && go run main.go)
#   - The two e2e users seeded against the local DB (see mobile/CLAUDE.md)
#
# Env vars (sourced from api/dev/switch-to-sqlite.sh if available):
#   E2E_ADMIN_USERNAME, E2E_ADMIN_PASSWORD, E2E_USER_USERNAME, E2E_USER_PASSWORD
#   E2E_MOBILE_BASE_URL (optional) — the API base URL the app hits. Defaults to
#     http://localhost:8081/api. The desktop suite's E2E_BASE_URL points at
#     the Angular dev server (:4200 via proxy), which the mobile app cannot
#     use — hence the separate variable.

set -euo pipefail

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$script_dir"

env_script="../api/dev/switch-to-sqlite.sh"
if [[ -f "$env_script" ]]; then
  # shellcheck disable=SC1090
  source "$env_script"
fi

: "${E2E_ADMIN_USERNAME:?set E2E_ADMIN_USERNAME or source api/dev/switch-to-sqlite.sh}"
: "${E2E_ADMIN_PASSWORD:?set E2E_ADMIN_PASSWORD or source api/dev/switch-to-sqlite.sh}"
: "${E2E_USER_USERNAME:?set E2E_USER_USERNAME or source api/dev/switch-to-sqlite.sh}"
: "${E2E_USER_PASSWORD:?set E2E_USER_PASSWORD or source api/dev/switch-to-sqlite.sh}"

mobile_base_url="${E2E_MOBILE_BASE_URL:-http://localhost:8081/api}"

command -v python3 >/dev/null 2>&1 || { echo "python3 not on PATH (needed to safely build the dart-define JSON)" >&2; exit 1; }

tmp_env="$(mktemp --suffix=.json)"
trap 'rm -f "$tmp_env"' EXIT

# Build the JSON via python3 instead of a heredoc so any value containing
# quotes/backslashes/newlines gets properly escaped. The heredoc form would
# emit invalid JSON that `flutter --dart-define-from-file` would reject.
# Matches the same approach in run-e2e-android.sh and run-e2e-ios.sh.
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

# Flutter desktop apps need a display. If DISPLAY isn't set and xvfb-run is
# available, run the test under Xvfb so it works headlessly (containers, CI).
runner=()
if [[ -z "${DISPLAY:-}" ]] && command -v xvfb-run >/dev/null 2>&1; then
  runner=(xvfb-run -a)
fi

# Build the list of test files to run.
#   - User passed specific files/dirs?  Use those verbatim.
#   - Default (no args)?                Discover *_test.dart and run each
#                                       in its own `flutter test` invocation.
# The per-file loop is needed because back-to-back integration_test runs
# in a single `flutter test integration_test/` invocation fail on Linux
# desktop with "Error waiting for a debug connection: The log reader
# stopped unexpectedly" -- the second app launch can't acquire the xvfb
# display the first one held. Splitting into separate processes lets each
# fully tear down before the next starts.
if [[ $# -gt 0 ]]; then
  targets=("$@")
else
  mapfile -t targets < <(find integration_test -maxdepth 2 -name '*_test.dart' | sort)
fi

if [[ ${#targets[@]} -eq 0 ]]; then
  echo "No *_test.dart files found under integration_test/" >&2
  exit 1
fi

exit_code=0
for target in "${targets[@]}"; do
  echo "==> ${runner[*]:-flutter test} $target"
  if ! "${runner[@]}" flutter test "$target" -d linux \
       --dart-define-from-file="$tmp_env"; then
    exit_code=1
  fi
done
exit "$exit_code"
