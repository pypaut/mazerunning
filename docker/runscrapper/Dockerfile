FROM golang:1.26 AS builder
WORKDIR /app

COPY ./runscrapper/go.mod ./runscrapper/go.sum* ./
RUN go mod download

COPY ./runscrapper/*.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o runscrapper .

FROM alpine:latest
WORKDIR /app

COPY --from=builder /app/runscrapper .
CMD ["./runscrapper"]
