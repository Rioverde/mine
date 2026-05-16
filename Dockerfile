FROM golang:1.25-alpine AS builder
WORKDIR /src
COPY . .
RUN go build -o /mine ./cmd/mine

FROM alpine:3.20
COPY --from=builder /mine /mine
COPY config.yaml /config.yaml
ENTRYPOINT ["/mine"]
