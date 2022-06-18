FROM golang:1.18.2 AS golang_builder

ENV FDB_VERSION 7.1.7
ADD "https://github.com/apple/foundationdb/releases/download/${FDB_VERSION}/foundationdb-clients_${FDB_VERSION}-1_amd64.deb" .
RUN dpkg -i foundationdb-clients_${FDB_VERSION}-1_amd64.deb \
    && rm -f foundationdb-clients_${FDB_VERSION}-1_amd64.deb

WORKDIR /workspace
ADD go.mod go.sum ./
RUN go mod download -x

ADD cmd ./cmd/
ADD internal ./internal/
RUN go build -o bin/bot ./cmd/simplifiedbot/

FROM ubuntu:focal-20220426

WORKDIR /workspace

ENV LAST_UPDATED 2022-05-22

# Add Tini
ENV TINI_VERSION v0.19.0
ADD https://github.com/krallin/tini/releases/download/${TINI_VERSION}/tini /tini
RUN chmod +x /tini
ENTRYPOINT ["/tini", "--"]

ENV FDB_VERSION 7.1.7
ADD "https://github.com/apple/foundationdb/releases/download/${FDB_VERSION}/foundationdb-clients_${FDB_VERSION}-1_amd64.deb" .
RUN dpkg -i foundationdb-clients_${FDB_VERSION}-1_amd64.deb \
    && rm -f foundationdb-clients_${FDB_VERSION}-1_amd64.deb

RUN apt-get update && apt-get install -y \
    ca-certificates \
 && rm -rf /var/lib/apt/lists/*

COPY --from=golang_builder /workspace/bin/bot ./

CMD ["./bot"]  