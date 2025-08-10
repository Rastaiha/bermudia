FROM golang:1.24-alpine AS build

WORKDIR /app
COPY ./ ./
RUN go build -o main main.go

FROM debian:12.11

WORKDIR /app

COPY --from=build /app/main .

EXPOSE 8080

CMD ["./main"]
