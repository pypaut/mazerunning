#!/usr/bin/env bash
set -euo pipefail

COMPOSE_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

cd "$COMPOSE_DIR"

echo "1. Retrieving secret from Scaleway..."
# MR_SECRET_ID must be in the environment
SECRET_JSON=$(scw secret version access secret-id="$MR_SECRET_ID" revision=latest raw=true)

if [ -z "$SECRET_JSON" ]; then
  echo "Error: could not retrieve secret (empty)." >&2
  exit 1
fi

echo "2. Downloading latest Docker image..."
podman compose pull

echo "3. Updating Podman secrets..."
printf '%s' "$SECRET_JSON" | jq -r 'to_entries[] | "\(.key)=\(.value)"' | while IFS='=' read -r key value; do
  [ -z "$key" ] && continue
  printf '%s' "$value" | podman secret create --replace "$key" - >/dev/null
done

echo "4. Restarting containers..."
podman compose up -d --force-recreate --remove-orphans

echo "5. Cleaning orphan images..."
podman image prune -f

echo "Successfully deployed."
