FROM golang:1.23-alpine AS mitosisd-builder

SHELL ["/bin/sh", "-ecuxo", "pipefail"]
RUN apk add --no-cache ca-certificates build-base git

WORKDIR /code
COPY ./go.mod ./go.sum /code/
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

COPY --from=mitosisd-builder /code/build/mitosisd /usr/bin/mitosisd
