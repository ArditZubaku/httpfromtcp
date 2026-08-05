benchmark:
	docker run --rm --platform=linux/amd64 ghcr.io/six-ddc/plow \
  -c 5000 \
  -n 100000 \
  http://host.docker.internal:42069/
