FROM golang:1.26 AS builder
WORKDIR /app

COPY go.mod go.sum* ./
RUN go mod download

COPY *.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o runscrapper .

FROM alpine:latest
WORKDIR /app

COPY --from=builder /app/runscrapper .
CMD ["./runscrapper"]
