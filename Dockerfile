FROM golang:1.15-alpine as builder
RUN apk add --update make
RUN mkdir /build
ADD . /build/
WORKDIR /build
RUN make build

FROM alpine:latest
COPY --from=builder /build/main /app/
COPY --from=builder /build/internal/db/migrations /app/migrations
WORKDIR /app
CMD ["/app/main"]