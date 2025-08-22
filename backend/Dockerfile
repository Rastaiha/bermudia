FROM golang:1.24 AS build

ENV CGO_ENABLED=1

WORKDIR /app
COPY ./ ./

RUN go build -o main main.go

FROM debian:12.11

# Install CA certificates
RUN apt-get update && \
    apt-get install -y --no-install-recommends ca-certificates && \
    update-ca-certificates && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY --from=build /app/main .

EXPOSE 8080

CMD ["./main"]
