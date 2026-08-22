# Stage 1: Build app
FROM golang:1.26-bookworm AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/cmd/expense-tracker/exe /app/cmd/expense-tracker

# Stage 2: Small image
FROM alpine:3.23
WORKDIR /app
COPY --from=builder app/cmd/expense-tracker/exe /app
CMD ["/app/exe"]