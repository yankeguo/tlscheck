FROM golang:1.22 AS builder
ENV CGO_ENABLED 0
ARG VERSION
WORKDIR /go/src/app
ADD . .
RUN go build -o /tlscheck

FROM alpine:3.20
RUN apk add --no-cache ca-certificates
COPY --from=builder /tlscheck /tlscheck
WORKDIR /data
CMD ["/tlscheck", "config.yaml"]
