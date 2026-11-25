# Step 1: build

FROM golang:1.26 AS builder
WORKDIR /src

COPY go.mod go.sum* ./
RUN go mod download

COPY *.go ./
COPY templates/ ./templates/

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-s -w" \
    -o /app/server .

# Step 2: run

FROM gcr.io/distroless/static-debian12:nonroot

COPY --from=builder --chown=nonroot:nonroot /app/server /app/server
USER nonroot:nonroot
EXPOSE 8080
ENTRYPOINT ["/app/server"]
