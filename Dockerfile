FROM golang:1.25-alpine AS builder

RUN apk add --no-cache git make

RUN go install github.com/pressly/goose/v3/cmd/goose@latest

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /bin/app ./cmd/api

FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

COPY --from=builder /go/bin/goose /bin/goose
COPY --from=builder /bin/app /bin/app
COPY --from=builder /app/migrations /migrations

EXPOSE 8080

CMD ["/bin/app"]
