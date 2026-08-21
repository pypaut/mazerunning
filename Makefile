docker-build:
	podman build -t runscrapper -f docker/runscrapper/Dockerfile .

docker-run:
	podman run \
		--env-file .env \
		--rm --name runscrapper runscrapper
