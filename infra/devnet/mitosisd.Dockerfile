FROM golang:1.24-alpine AS mitosisd-builder

SHELL ["/bin/sh", "-ecuxo", "pipefail"]
RUN apk add --no-cache ca-certificates build-base git

WORKDIR /code
COPY ./go.mod ./go.sum /code/
RUN go mod download

COPY ./Makefile /code/
COPY ./api /code/api
COPY ./app /code/app
COPY ./bindings /code/bindings
COPY ./cmd /code/cmd
COPY ./types /code/types
COPY ./x /code/x

RUN LEDGER_ENABLED=false BUILD_TAGS=muslc LINK_STATICALLY=true make build \
  && file ./build/mitosisd \
  && echo "Ensuring mitosisd is statically linked ..." \
  && (file ./build/mitosisd | grep "statically linked")

# ============================================================
FROM ubuntu:22.04

RUN apt-get update
RUN apt-get install -y curl make bash jq sed git xxd

# Foundry
COPY --from=ghcr.io/foundry-rs/foundry:stable /usr/local/bin/forge /usr/local/bin/forge
COPY --from=ghcr.io/foundry-rs/foundry:stable /usr/local/bin/cast /usr/local/bin/cast
COPY --from=ghcr.io/foundry-rs/foundry:stable /usr/local/bin/anvil /usr/local/bin/anvil
COPY --from=ghcr.io/foundry-rs/foundry:stable /usr/local/bin/chisel /usr/local/bin/chisel

COPY --from=mitosisd-builder /code/build/mitosisd /usr/bin/mitosisd
COPY --from=mitosisd-builder /code/build/midevtool /usr/bin/midevtool
