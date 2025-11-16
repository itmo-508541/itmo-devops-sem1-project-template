FROM golang:1.25-alpine AS builder
WORKDIR /srv
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . ./
RUN go build -tags dist -o ./bin/app ./cmd/main.go

FROM alpine:3.22 AS prod
WORKDIR /srv
COPY --from=builder /srv/bin/* ./

CMD ["/srv/docker-entrypoint.sh"]
