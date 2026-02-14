#!/usr/bin/env bash
set -euo pipefail

ENV="${ENV:-staging}"
PROFILE="${PROFILE:-homelab-${ENV}}"
MINIKUBE_IP="${MINIKUBE_IP:-}"

if [[ "${ENV}" != "dev" && "${ENV}" != "staging" ]]; then
  echo "Only ENV=dev or ENV=staging is supported for sslip.io updates. Got ENV=${ENV}." >&2
  exit 1
fi

if [[ -z "${MINIKUBE_IP}" ]]; then
  if ! command -v minikube >/dev/null 2>&1; then
    echo "minikube command is required (or provide MINIKUBE_IP)." >&2
    exit 1
  fi
  MINIKUBE_IP="$(minikube ip --profile "${PROFILE}")"
fi

if [[ -z "${MINIKUBE_IP}" ]]; then
  echo "Unable to resolve Minikube IP for profile ${PROFILE}." >&2
  exit 1
fi

IP_DASH="${MINIKUBE_IP//./-}"
echo "Updating sslip host IP to ${IP_DASH} for ENV=${ENV} (PROFILE=${PROFILE})"

updated=0
while IFS= read -r file; do
  if rg -q '\.sslip\.io' "${file}"; then
    before="$(cksum "${file}" | awk '{print $1}')"
    perl -i -pe "s/(?:\\d{1,3}-){3}\\d{1,3}(?=\\.sslip\\.io)/${IP_DASH}/g" "${file}"
    after="$(cksum "${file}" | awk '{print $1}')"
    if [[ "${before}" != "${after}" ]]; then
      echo "updated ${file}"
      updated=$((updated + 1))
    fi
  fi
done < <(
  {
    find k8s -type f -path "*/overlays/${ENV}/*"
    echo Makefile
  } | sort -u
)

echo "Done. ${updated} file(s) updated."
