# Build Stage
FROM golang:latest AS build

WORKDIR /app

COPY . .


RUN go mod download
RUN go build -o main .

# Final Stage
FROM golang:latest

WORKDIR /app

COPY --from=build /app/main .
COPY .env .

EXPOSE 8080

CMD ["./main"]
