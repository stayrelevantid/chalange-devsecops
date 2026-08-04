#!/bin/bash
# ==============================================================
# drift-check.sh — Drift Detection: Pods tanpa Resource Limits
# Day 54: Chaos Security Engineering
#
# Deteksi "drift" dari kebijakan OPA Gatekeeper: pod yang berhasil
# masuk meskipun admission control mati (fail-open), atau yang
# diterapkan sebelum policy ada.
#
# Kebijakan yang dicek (harus sinkron dengan constraint
# require-resource-limits): tiap container wajib punya
#   resources.requests.cpu, resources.requests.memory,
#   resources.limits.cpu, resources.limits.memory
#
# Usage:
#   ./drift-check.sh [NAMESPACE]   (default: securebank)
#   ./drift-check.sh --all
# ==============================================================
set -euo pipefail

NS="${1:-securebank}"

if [[ "${NS}" == "--all" ]]; then
  NAMESPACES=$(kubectl get ns -o jsonpath='{range .items[*]}{.metadata.name}{"\n"}{end}')
else
  NAMESPACES="${NS}"
fi

echo "=== DRIFT CHECK: pods tanpa resource requests/limits ==="
FOUND=0

for ns in ${NAMESPACES}; do
  echo ""
  echo "--- namespace: ${ns} ---"

  # List container yang kurang salah satu dari: requests.cpu, requests.memory, limits.cpu, limits.memory
  kubectl get pods -n "${ns}" -o json 2>/dev/null | jq -r '
    .items[] | . as $pod |
    .spec.containers[] | . as $c |
    ((.resources.requests.cpu // "MISSING") as $rcpu |
     (.resources.requests.memory // "MISSING") as $rmem |
     (.resources.limits.cpu // "MISSING") as $lcpu |
     (.resources.limits.memory // "MISSING") as $lmem |
     if $rcpu=="MISSING" or $rmem=="MISSING" or $lcpu=="MISSING" or $lmem=="MISSING" then
       "DRIFT: pod=\($pod.metadata.name) container=\($c.name) " +
       "req.cpu=\($rcpu) req.mem=\($rmem) lim.cpu=\($lcpu) lim.mem=\($lmem)"
     else empty end)
  ' || true

  # Hitung total drift di namespace ini
  COUNT=$(kubectl get pods -n "${ns}" -o json 2>/dev/null | jq -r '
    [.items[] | .spec.containers[] |
      select((.resources.requests.cpu // "MISSING")=="MISSING" or
             (.resources.requests.memory // "MISSING")=="MISSING" or
             (.resources.limits.cpu // "MISSING")=="MISSING" or
             (.resources.limits.memory // "MISSING")=="MISSING")] | length
  ' || echo 0)
  if [[ "${COUNT}" -gt 0 ]]; then
    echo "  ^ ${COUNT} container drift ditemukan"
    FOUND=$((FOUND + COUNT))
  else
    echo "  OK — semua container punya resource requests & limits"
  fi
done

echo ""
if [[ "${FOUND}" -gt 0 ]]; then
  echo "RESULT: ${FOUND} container drift terdeteksi. ⚠️  Periksa — admission control mungkin fail-open / policy baru."
  exit 1
else
  echo "RESULT: 0 drift. Semua pod sesuai kebijakan. ✅"
  exit 0
fi
