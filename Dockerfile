FROM golang:1.24-alpine AS build
WORKDIR /app
RUN go mod download
RUN go build -o bot .

FROM alpine:latest
WORKDIR /app
COPY --from=build /app/bot .
CMD ["./bot"]