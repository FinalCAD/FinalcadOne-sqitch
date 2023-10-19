FROM golang:1.21-alpine AS builder
ARG APPUSER=appuser

ENV USER=${APPUSER}
ENV UID=1001

RUN adduser -D -g "" -H -s "/sbin/nologin" -u "${UID}" "${USER}"
WORKDIR /extension

COPY go.mod go.sum ./

RUN go mod download

COPY internal ./internal
COPY cmd ./cmd
COPY pkg ./pkg

RUN go mod tidy && \
    go mod verify && \
    CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -v -o /bin/sqitch-config cmd/sqitch-config/main.go

FROM scratch
ARG APPUSER=appuser

WORKDIR /

COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

USER ${APPUSER}:${APPUSER}

COPY --from=builder --chown=${APPUSER}:${APPUSER} /bin/sqitch-config /sqitch-config

ENTRYPOINT ["/sqitch-config"]
