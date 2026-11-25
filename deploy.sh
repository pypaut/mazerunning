#!/usr/bin/env bash
set -euo pipefail

SECRET_ID="MY_SECRET_ID"
COMPOSE_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

cd "$COMPOSE_DIR"

echo "1. Retrieving secret from Scaleway..."
SECRET_JSON=$(scw secret version access secret-id="$SECRET_ID" revision=latest raw=true)

if [ -z "$SECRET_JSON" ]; then
  echo "Error: could not retrieve secret (empty)." >&2
  exit 1
fi

echo "2. Downloading latest Docker image..."
podman compose pull

echo "3. Restarting containers with secrets..."
(
  eval $(echo "$SECRET_JSON" | jq -r 'to_entries[] | "export \(.key)=\(.value | @sh)"')
  podman compose up -d --remove-orphans
)

echo "4. Cleaning orphan images..."
podman image prune -f

echo "Successfully deployed."
