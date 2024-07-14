# Stage 1: Build the ethos-chain and mitosis
FROM golang:1.22-alpine AS go-builder

SHELL ["/bin/sh", "-ecuxo", "pipefail"]

RUN apk add --no-cache ca-certificates build-base git

WORKDIR /code/mitosis
COPY . /code/mitosis/

RUN go mod download

RUN LEDGER_ENABLED=false BUILD_TAGS=muslc LINK_STATICALLY=true make build-mitosisd \
    && file /code/mitosis/build/mitosisd \
    && echo "Ensuring binary is statically linked ..." \
    && (file /code/mitosis/build/mitosisd | grep "statically linked")

ENV GOBIN=/usr/bin

# Stage 3: Construct the final image
FROM alpine:3.16

COPY --from=go-builder /code/mitosis/build/mitosisd /usr/bin/mitosisd
#COPY --from=go-builder /code/ethos-chain/build/ethosd /usr/bin/ethosd
#COPY --from=go-builder /code/ethos-chain/infra/consumer/v50/consumer.sh /scripts/consumer.sh
#COPY --from=go-builder /code/ethos-chain/infra/consumer/v50/validator.sh /scripts/validator.sh
#
#RUN chmod +x /scripts/consumer.sh
#RUN chmod +x /scripts/validator.sh

# Install dependencies used for Starship
RUN apk add --no-cache curl make bash jq sed

#ARG LOCAL_IP
#ENV LOCAL_IP=${LOCAL_IP}
#
#EXPOSE 1317 26656 26657

#CMD ["/usr/bin/icsd", "version"]
