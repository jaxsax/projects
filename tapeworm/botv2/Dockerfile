FROM golang:1.19.5 AS golang_builder

WORKDIR /workspace
ADD go.mod go.sum Makefile ./
RUN go mod download -x

ADD cmd ./cmd/
ADD internal ./internal/
COPY db/migrations ./migrations
RUN make build

FROM ubuntu:focal-20220426

WORKDIR /workspace

ENV LAST_UPDATED 2022-05-22

# Add Tini
ENV TINI_VERSION v0.19.0
ADD https://github.com/krallin/tini/releases/download/${TINI_VERSION}/tini /tini
RUN chmod +x /tini
ENTRYPOINT ["/tini", "--"]

RUN apt-get update && apt-get install -y \
    ca-certificates \
 && rm -rf /var/lib/apt/lists/*

COPY --from=golang_builder /workspace/bin/ ./
COPY --from=golang_builder /workspace/migrations /workspace/db/migrations

CMD ["./bot"]  
