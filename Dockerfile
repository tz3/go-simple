FROM golang:1.20-alpine

WORKDIR /app

COPY . .

RUN go build -o main ./cmd/main.go

ENV DB_CONNECTION_STRING="postgres://postgres:secret@postgres:5432/test?sslmode=disable"

RUN apk add --no-cache postgresql-client

CMD ["sh", "-c", "until psql $DB_CONNECTION_STRING -c 'select 1'; do sleep 1; done; ./main"]
