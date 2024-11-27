FROM golang:1.23-alpine

RUN apk add --no-cache git && \
    go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@v4.12.2

RUN addgroup -S node && adduser -S node -G node && \
    mkdir -p /app && chown node:node /app
WORKDIR /app

COPY --chown=node:node go.mod ./
COPY --chown=node:node go.sum ./
RUN go mod download

USER node

COPY --chown=node:node . .

RUN go build -o /app /app/cmd/main.go

CMD ["/app/main"]