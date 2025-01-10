FROM golang:1.23.3-bookworm AS devlopment

WORKDIR /agent/src/

COPY go.mod go.sum /agent/src/

RUN go mod download

COPY . .

RUN cd /agent/src/cmd/agent && go build -o /agent/bin/agent

ENV PATH="/agent/bin:$PATH"
