# ---------- Build stage ----------
FROM golang:1.25 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY main.go .

RUN CGO_ENABLED=0 GOOS=linux go build -o dog-api main.go


# ---------- Runtime stage ----------
FROM alpine:3.22

WORKDIR /app

COPY --from=builder /app/dog-api .

EXPOSE 8081

CMD ["./dog-api"]
