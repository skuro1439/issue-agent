FROM debian:bookworm-slim AS release

# Use binaly built by goreleaser
COPY agent-bin /usr/local/bin/agent

ENTRYPOINT ["/usr/local/bin/agent"]
