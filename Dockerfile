FROM golang:alpine AS builder

WORKDIR /build

ADD go.mod .

COPY . .

RUN go build -o api ./src

FROM alpine

WORKDIR /build

COPY --from=builder /build/api      api
# COPY --from=builder /build/.env     .env
COPY --from=builder /build/pages    pages
COPY --from=builder /build/static   static

CMD ["./api"]