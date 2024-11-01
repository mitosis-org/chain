FROM golang:1.22-alpine AS ethosd-builder

SHELL ["/bin/sh", "-ecuxo", "pipefail"]
RUN apk add --no-cache ca-certificates build-base git

COPY ./ethos /code

WORKDIR /code/ethos-avs
RUN go mod download
WORKDIR /code/ethos-chain
RUN go mod download

WORKDIR /code/ethos-chain
RUN LEDGER_ENABLED=false BUILD_TAGS=muslc LINK_STATICALLY=true make build \
  && file ./build/ethosd \
  && echo "Ensuring ethosd is statically linked ..." \
  && (file ./build/ethosd | grep "statically linked")

FROM golang:1.22-alpine AS mitosisd-builder

SHELL ["/bin/sh", "-ecuxo", "pipefail"]
RUN apk add --no-cache ca-certificates build-base git

WORKDIR /code
COPY ./go.mod ./go.sum /code/
COPY ./ethos /code/ethos
RUN go mod download

COPY ./Makefile /code/
COPY ./app /code/app
COPY ./cmd /code/cmd
COPY ./tests /code/tests

RUN LEDGER_ENABLED=false BUILD_TAGS=muslc LINK_STATICALLY=true make build-mitosisd \
  && file ./build/mitosisd \
  && echo "Ensuring mitosisd is statically linked ..." \
  && (file ./build/mitosisd | grep "statically linked")

# ============================================================
FROM alpine:3.16

RUN apk add --no-cache curl make bash jq sed

COPY --from=ethereum/client-go:v1.14.11 /usr/local/bin/geth /usr/local/bin/geth
COPY --from=ethosd-builder /code/ethos-chain/build/ethosd /usr/bin/ethosd
COPY --from=mitosisd-builder /code/build/mitosisd /usr/bin/mitosisd
