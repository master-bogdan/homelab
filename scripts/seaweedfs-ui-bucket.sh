#!/usr/bin/env sh
set -e

: "${SEAWEEDFS_MASTER_ADDR:?SEAWEEDFS_MASTER_ADDR is required}"
: "${SEAWEEDFS_FILER_ADDR:?SEAWEEDFS_FILER_ADDR is required}"

KUBECTL="${KUBECTL:-kubectl}"
NS="${SEAWEEDFS_NAMESPACE:-platform}"
BUCKET="${UI_PUBLIC_BUCKET:-public}"
ANON_ACTIONS="${SEAWEED_ANON_ACTIONS:-Read,Write}"

"$KUBECTL" -n "$NS" exec -i statefulset/seaweedfs-master -- sh -c \
	"echo \"s3.bucket.create --name $BUCKET\" | weed shell -master=\"$SEAWEEDFS_MASTER_ADDR\" -filer=\"$SEAWEEDFS_FILER_ADDR\"" || true
"$KUBECTL" -n "$NS" exec -i statefulset/seaweedfs-master -- sh -c \
	"echo \"s3.configure --user anonymous --buckets $BUCKET --actions $ANON_ACTIONS --apply true\" | weed shell -master=\"$SEAWEEDFS_MASTER_ADDR\" -filer=\"$SEAWEEDFS_FILER_ADDR\"" || true
