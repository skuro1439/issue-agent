FROM golang:1.23.3-bookworm AS devlopment

WORKDIR /agent

COPY . /agent/src

RUN cd /agent/src/cmd/agent && go build -o /agent/bin/agent

ENV PATH="/agent/bin:$PATH"
