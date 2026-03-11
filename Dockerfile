FROM golang:1.25.1-alpine3.22 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ENV CGO_ENABLED=false
ENV GOOSE=linux
RUN go build -o study-platform cmd/study-platform/main.go

FROM alpine:3.22
WORKDIR /app

COPY --from=builder /app/study-platform .
COPY --from=builder /app/migrations ./migrations

EXPOSE 3800

CMD ["./study-platform"]
