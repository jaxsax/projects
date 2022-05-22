FROM golang:1.18.2 AS golang_builder

WORKDIR /workspace
ADD go.mod go.sum ./
RUN go mod download -x

ADD . .
RUN go build -x -o bin/bot ./cmd/simplifiedbot/

FROM ubuntu:focal-20220426

WORKDIR /workspace

# Add Tini
ENV TINI_VERSION v0.19.0
ADD https://github.com/krallin/tini/releases/download/${TINI_VERSION}/tini /tini
RUN chmod +x /tini
ENTRYPOINT ["/tini", "--"]

COPY --from=golang_builder /workspace/bin/bot ./

CMD ["./bot"]  