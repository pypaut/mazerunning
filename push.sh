podman buildx build --platform linux/amd64 -t ghcr.io/pypaut/runscrapper:latest ./runscrapper && \
podman push ghcr.io/pypaut/runscrapper:latest
