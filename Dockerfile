FROM golang:1.25-alpine AS builder
WORKDIR /app
RUN apk add --no-cache git
RUN go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
COPY . .
RUN go get ./cmd/gen_translator_keys
RUN go install ./cmd/gen_translator_keys \
 && gen_translator_keys
RUN go get ./cmd/sqlgen
RUN go install ./cmd/sqlgen \
 && sqlgen
RUN go get ./cmd/bot
RUN go install ./cmd/bot

FROM alpine:latest
WORKDIR /app
RUN apk add --no-cache tor bash postgresql-client
COPY --from=builder /go/bin/bot /usr/local/bin/bot
COPY --from=builder /app/locales /app/locales
CMD ["bot"]