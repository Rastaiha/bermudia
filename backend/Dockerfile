FROM golang:1.24 AS build

ENV CGO_ENABLED=1

WORKDIR /app
COPY ./ ./

RUN go build -o main main.go

FROM debian:12.11

WORKDIR /app

COPY --from=build /app/main .

EXPOSE 8080

CMD ["./main"]
