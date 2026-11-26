SECRET_ID := $(MR_SECRET_ID)
IMAGE := $(MR_IMAGE)
VM := $(MR_VM)
REMOTE_DIR := $(MR_REMOTE_DIR)

.PHONY: deploy

push:
	podman buildx build --platform linux/amd64 -t ${IMAGE} --push ./runscrapper
	
deploy:
	rsync -avz --delete docker-compose.yml migrations $(VM):$(REMOTE_DIR)/
	ssh ${VM} "${REMOTE_DIR}/deploy.sh"

secrets:
	scw secret version access secret-id="$(SECRET_ID)" revision=latest raw=true \
		| jq -j 'to_entries[] | "\(.key)\u0000\(.value)\u0000"' \
		| xargs -0 -n2 sh -c 'printf "%s" "$$2" | podman secret create --replace "$$1" - >/dev/null' _
