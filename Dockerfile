FROM golang:1.25.10-alpine AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -buildvcs=false -o /out/gotaskflow ./cmd/app

FROM alpine:3.22

WORKDIR /app

RUN apk add --no-cache ca-certificates && adduser -D -H appuser

COPY --from=builder /out/gotaskflow ./gotaskflow
COPY web ./web

EXPOSE 8080

USER appuser

CMD ["./gotaskflow"]
