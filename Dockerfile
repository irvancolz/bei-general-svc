FROM golang:1.20.2-alpine3.17 as builder

RUN apk update && apk add --no-cache git
RUN apk add make

WORKDIR /src

COPY . .
# RUN go install github.com/pressly/goose/v3/cmd/goose@latest
RUN go mod tidy

RUN go build -o /svc-general -mod=mod

FROM golang:1.20.2-alpine3.17

WORKDIR /app

COPY --from=builder /svc-general /svc-general
COPY --from=builder /src/static ./static 

EXPOSE 8080

ENTRYPOINT ["/svc-general"]