#!/usr/bin/env bash
# Post-deploy check: the expected version is live and a user can create,
# read and progress a lead. Used by `make smoke` and the CD workflow.
set -euo pipefail

APP_URL="${APP_URL:-http://localhost:8088}"
EXPECTED_VERSION="${EXPECTED_VERSION:-}"

fail() { echo "FAIL: $*" >&2; exit 1; }

echo "== readiness"
for i in $(seq 1 30); do
  curl -sf "$APP_URL/readyz" >/dev/null && break
  [ "$i" = 30 ] && fail "service not ready at $APP_URL"
  sleep 2
done
curl -s "$APP_URL/readyz"; echo

echo "== version"
version_json=$(curl -sf "$APP_URL/version")
echo "$version_json"
if [ -n "$EXPECTED_VERSION" ]; then
  echo "$version_json" | grep -q "\"version\":\"$EXPECTED_VERSION\"" \
    || fail "expected version $EXPECTED_VERSION to be running"
fi

echo "== create lead"
email="smoke-$(date +%s)@example.com"
created=$(curl -sf -X POST "$APP_URL/leads" -H 'Content-Type: application/json' \
  -d "{\"name\":\"Smoke Test\",\"email\":\"$email\",\"phone\":\"9876543210\",\"source\":\"smoke-test\"}")
echo "$created"
id=$(echo "$created" | sed -n 's/.*"id":"\([^"]*\)".*/\1/p')
[ -n "$id" ] || fail "no lead id in create response"

echo "== read lead"
curl -sf "$APP_URL/leads/$id" | grep -q "$email" || fail "could not read lead $id"

echo "== move lead NEW -> CONTACTED"
curl -sf -X PATCH "$APP_URL/leads/$id/status" -H 'Content-Type: application/json' \
  -d '{"status":"CONTACTED"}' | grep -q '"status":"CONTACTED"' || fail "status update failed"

echo "== invalid transition is rejected"
code=$(curl -s -o /dev/null -w '%{http_code}' -X PATCH "$APP_URL/leads/$id/status" \
  -H 'Content-Type: application/json' -d '{"status":"CONVERTED"}')
[ "$code" = 400 ] || fail "expected 400 for CONTACTED -> CONVERTED, got $code"

echo "== cleanup"
curl -sf -X DELETE "$APP_URL/leads/$id" >/dev/null || fail "delete failed"

echo "SMOKE TEST PASSED (lead $id)"
