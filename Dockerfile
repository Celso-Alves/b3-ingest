FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN GOOS=linux go build -o b3-ingest .

FROM alpine:3.20
WORKDIR /app
RUN apk add --no-cache tzdata
COPY --from=builder /app/b3-ingest ./b3-ingest
COPY env-local ./env-local
COPY setenv.sh ./setenv.sh
ENV GIN_MODE=release
EXPOSE 8000 8080
# ENTRYPOINT ["/app/b3-ingest"]
