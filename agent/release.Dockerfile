FROM debian:bookworm-slim AS release

RUN apt-get update && apt-get install -y git

# Use binaly built by goreleaser
COPY agent /usr/local/bin/agent

ENTRYPOINT ["/usr/local/bin/agent"]
