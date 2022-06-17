FROM golang:1.18.2 AS golang_builder

WORKDIR /workspace
ADD go.mod go.sum ./
RUN go mod download -x

ADD . .
RUN go build -o bin/bot ./cmd/simplifiedbot/

FROM node:17.9.0-slim AS fe_builder

WORKDIR /workspace/app

ADD app/package.json app/yarn.lock .
RUN yarn

ADD app/ .
RUN yarn run build:prod

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

COPY --from=golang_builder /workspace/bin/bot ./
COPY --from=fe_builder /workspace/app/out/ ./frontend

CMD ["./bot"]  