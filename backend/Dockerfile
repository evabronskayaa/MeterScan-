FROM golang:1.21 as builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main ./cmd/api-server/main.go

FROM scratch

COPY --from=builder /app/main /app/main

COPY ./api-config.json ./api-config.json

CMD ["/app/main"]