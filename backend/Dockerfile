FROM golang:1.21 as builder

ARG SERVICE=api

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main ./cmd/$SERVICE/main.go

FROM scratch

COPY --from=builder /app/main /app/main

ARG COPY_PATH
COPY $COPY_PATH $COPY_PATH

CMD ["/app/main"]