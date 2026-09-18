# ---------- Build stage ----------
FROM golang:1.25 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY main.go .

RUN CGO_ENABLED=0 GOOS=linux go build -o cat-api main.go


# ---------- Runtime stage ----------
FROM alpine:3.22

WORKDIR /app

COPY --from=builder /app/cat-api .

EXPOSE 8080

CMD ["./cat-api"]
