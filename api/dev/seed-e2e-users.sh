#!/usr/bin/env bash
# Idempotent seed script for the two e2e test users that mobile/run-e2e-*.sh
# and desktop's Playwright suite expect.
#
# Why this exists:
#   - The API auto-creates a default `admin/admin` user on first startup
#     (repositories/db.go -> CreateUserIfNoneExist) outside `deployEnv=test`.
#   - The mobile e2e tests need `e2e-admin` (role ADMIN) and `e2e-user`
#     (role USER) to log in. These are NOT auto-seeded.
#   - The UI signup path is gated on `enableLocalSignUp`, which is `false`
#     locally, AND would assign USER role anyway because the auto-admin
#     already occupies the "first user = ADMIN" slot in CreateUser.
#   - The admin-protected POST /user/ endpoint accepts an explicit `userRole`
#     and creates whatever you ask for, so we use it to seed both.
#
# Usage:
#   cd api && ./dev/seed-e2e-users.sh
#
# Env overrides (all optional, defaults match switch-to-sqlite.sh):
#   API_BASE_URL          default http://localhost:8081/api
#   ADMIN_USERNAME        default admin   (the auto-created default admin)
#   ADMIN_PASSWORD        default admin
#   E2E_ADMIN_USERNAME    default e2e-admin
#   E2E_ADMIN_PASSWORD    default e2e-admin-password
#   E2E_USER_USERNAME     default e2e-user
#   E2E_USER_PASSWORD     default e2e-user-password

set -euo pipefail

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# Pick up E2E_* defaults from switch-to-sqlite.sh if not already exported.
if [[ -f "$script_dir/switch-to-sqlite.sh" ]]; then
  # shellcheck disable=SC1090
  source "$script_dir/switch-to-sqlite.sh"
fi

api_base_url="${API_BASE_URL:-http://localhost:8081/api}"
admin_username="${ADMIN_USERNAME:-admin}"
admin_password="${ADMIN_PASSWORD:-admin}"
e2e_admin_username="${E2E_ADMIN_USERNAME:-e2e-admin}"
e2e_admin_password="${E2E_ADMIN_PASSWORD:-e2e-admin-password}"
e2e_user_username="${E2E_USER_USERNAME:-e2e-user}"
e2e_user_password="${E2E_USER_PASSWORD:-e2e-user-password}"

command -v curl >/dev/null    || { echo "curl not on PATH" >&2; exit 1; }
command -v python3 >/dev/null || { echo "python3 not on PATH" >&2; exit 1; }

echo "==> logging in as $admin_username at $api_base_url"
# Build the JSON body via python3 instead of bash string interpolation so a
# password containing quotes/backslashes/newlines can't corrupt the payload.
login_body="$(python3 -c 'import json,sys; print(json.dumps({"username":sys.argv[1],"password":sys.argv[2]}))' \
  "$admin_username" "$admin_password")"
login_response="$(curl -fsS --max-time 10 -X POST \
  "$api_base_url/login/?tokensInBody=true" \
  -H 'Content-Type: application/json' \
  -d "$login_body" || true)"

if [[ -z "$login_response" ]]; then
  echo "login failed: empty response (API down? wrong base URL?)" >&2
  exit 1
fi

jwt="$(printf '%s' "$login_response" \
  | python3 -c 'import json,sys; d=json.load(sys.stdin); print(d.get("jwt",""))' \
  || true)"

if [[ -z "$jwt" ]]; then
  echo "login failed: no jwt in response body" >&2
  echo "response: $login_response" >&2
  exit 1
fi

echo "==> fetching existing users"
existing_users_json="$(curl -fsS --max-time 10 \
  "$api_base_url/user/" \
  -H "Authorization: Bearer $jwt")"

existing_usernames="$(printf '%s' "$existing_users_json" \
  | python3 -c 'import json,sys; print("\n".join(u["username"] for u in json.load(sys.stdin)))')"

echo "    existing: $(echo "$existing_usernames" | tr '\n' ' ')"

# create_user_if_missing <username> <password> <displayName> <role>
create_user_if_missing() {
  local username="$1"
  local password="$2"
  local display_name="$3"
  local role="$4"

  if printf '%s\n' "$existing_usernames" | grep -Fxq "$username"; then
    printf '==> %-12s (%-5s) ... [exists]\n' "$username" "$role"
    return 0
  fi

  local body
  body="$(python3 -c 'import json,sys; print(json.dumps({"username":sys.argv[1],"password":sys.argv[2],"displayName":sys.argv[3],"userRole":sys.argv[4],"isDummyUser":False}))' \
    "$username" "$password" "$display_name" "$role")"

  # Use mktemp instead of /tmp/...$$ -- the PID-based path is predictable and
  # vulnerable to a symlink race on shared machines.
  local response http_code response_file
  response_file="$(mktemp -t seed-e2e-users.XXXXXX)"
  response="$(curl -sS --max-time 10 -o "$response_file" -w '%{http_code}' \
    -X POST "$api_base_url/user/" \
    -H "Authorization: Bearer $jwt" \
    -H 'Content-Type: application/json' \
    -d "$body")"
  http_code="$response"
  body="$(cat "$response_file" 2>/dev/null || true)"
  rm -f "$response_file"

  if [[ "$http_code" == "200" ]]; then
    printf '==> %-12s (%-5s) ... [created]\n' "$username" "$role"
  else
    printf '==> %-12s (%-5s) ... [FAILED http %s] %s\n' \
      "$username" "$role" "$http_code" "$body" >&2
    return 1
  fi
}

rc=0
create_user_if_missing "$e2e_admin_username" "$e2e_admin_password" "E2E Admin" "ADMIN" || rc=1
create_user_if_missing "$e2e_user_username"  "$e2e_user_password"  "E2E User"  "USER"  || rc=1

exit $rc
