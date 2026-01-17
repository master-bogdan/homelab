#!/usr/bin/env sh
set -e

: "${APPS_DIR:?APPS_DIR is required}"
: "${APP:?APP is required}"
: "${SEAWEED_INTERNAL_BASE:?SEAWEED_INTERNAL_BASE is required}"
: "${UI_PUBLIC_BUCKET:?UI_PUBLIC_BUCKET is required}"

KUBECTL="${KUBECTL:-kubectl}"
NS="${SEAWEEDFS_NAMESPACE:-platform}"
POD_PREFIX="${SEAWEED_UPLOAD_POD_PREFIX:-seaweedfs-uploader}"
IMAGE="${SEAWEED_UPLOAD_IMAGE:-curlimages/curl:8.5.0}"
SETUP_CMD="${SEAWEED_UPLOAD_SETUP_CMD:-true}"
WAIT_TIMEOUT="${SEAWEED_UPLOAD_WAIT:-90s}"
CONNECT_TIMEOUT="${SEAWEED_UPLOAD_CONNECT_TIMEOUT:-5}"
MAX_TIME="${SEAWEED_UPLOAD_MAX_TIME:-60}"
OUT_DIRS="${UI_BUILD_OUTPUT_DIRS:-out build dist}"

APP_DIR="$APPS_DIR/$APP"
if [ ! -d "$APP_DIR" ]; then
	echo "ERROR: $APP_DIR not found" >&2
	exit 1
fi

OUT_DIR=""
for dir in $OUT_DIRS; do
	if [ -d "$APP_DIR/$dir" ]; then
		OUT_DIR="$APP_DIR/$dir"
		break
	fi
done
if [ -z "$OUT_DIR" ]; then
	echo "ERROR: No build output dir found for $APP" >&2
	exit 1
fi

POD="${POD_PREFIX}-${APP}"
cleanup() { "$KUBECTL" -n "$NS" delete pod "$POD" --ignore-not-found >/dev/null; }
trap cleanup EXIT
cleanup
"$KUBECTL" -n "$NS" run "$POD" --image="$IMAGE" --restart=Never --command -- \
	sh -c "$SETUP_CMD && sleep 3600"
"$KUBECTL" -n "$NS" wait --for=condition=Ready pod/"$POD" --timeout="$WAIT_TIMEOUT"

echo "Uploading $OUT_DIR -> $SEAWEED_INTERNAL_BASE/$UI_PUBLIC_BUCKET/$APP (in-cluster)"
OUT_BASE=$(basename "$OUT_DIR")
"$KUBECTL" -n "$NS" exec "$POD" -- sh -c "mkdir -p /tmp/upload"
"$KUBECTL" -n "$NS" cp "$OUT_DIR" "$POD:/tmp/upload"
"$KUBECTL" -n "$NS" exec "$POD" -- env \
	SEAWEED_BASE="$SEAWEED_INTERNAL_BASE" \
	BUCKET="$UI_PUBLIC_BUCKET" \
	APP="$APP" \
	UPLOAD_DIR="/tmp/upload/$OUT_BASE" \
	CONNECT_TIMEOUT="$CONNECT_TIMEOUT" \
	MAX_TIME="$MAX_TIME" \
	sh -c 'set -e; cd "$UPLOAD_DIR" && find . -type f | while IFS= read -r f; do
		rel="${f#./}"
		url="$SEAWEED_BASE/$BUCKET/$APP/$rel"
		content_type="application/octet-stream"
		case "$rel" in
			*.html) content_type="text/html; charset=utf-8" ;;
			*.css) content_type="text/css; charset=utf-8" ;;
			*.js|*.mjs) content_type="application/javascript; charset=utf-8" ;;
			*.json|*.map) content_type="application/json; charset=utf-8" ;;
			*.txt) content_type="text/plain; charset=utf-8" ;;
			*.svg) content_type="image/svg+xml" ;;
			*.png) content_type="image/png" ;;
			*.jpg|*.jpeg) content_type="image/jpeg" ;;
			*.gif) content_type="image/gif" ;;
			*.webp) content_type="image/webp" ;;
			*.ico) content_type="image/x-icon" ;;
			*.woff2) content_type="font/woff2" ;;
			*.woff) content_type="font/woff" ;;
			*.ttf) content_type="font/ttf" ;;
			*.otf) content_type="font/otf" ;;
			*.eot) content_type="application/vnd.ms-fontobject" ;;
			*.pdf) content_type="application/pdf" ;;
		esac
		echo "  -> $rel"
		curl -sS --fail -X PUT -H "Expect:" -H "Content-Type: $content_type" \
			--connect-timeout "$CONNECT_TIMEOUT" \
			--max-time "$MAX_TIME" \
			--upload-file "$f" "$url" -o /dev/null
	done'
cleanup
