FROM golang:alpine as builder

ENV GO111MODULE=on

WORKDIR /app

RUN apk add build-base

COPY go.mod .

RUN go mod download

COPY . .

RUN go build -o questions-rest-api .

FROM alpine

ENV PORT=3000

RUN apk update && apk add sqlite

COPY --from=builder /app/questions-rest-api .
COPY --from=builder /app/database ./database
COPY --from=builder /app/swagger.yaml .

CMD ["/questions-rest-api"]

