FROM golang:1.20.2-alpine3.17

# timezone
# RUN cp /usr/share/zoneinfo/Asia/Jakarta /etc/localtime && echo "Asia/Jakarta" > /etc/timezone

# RUN apt update && apt upgrade -y && \
#     apt install -y git \
#     make openssh-client

# WORKDIR /app/cmd/httprest

# RUN curl -fLo install.sh https://raw.githubusercontent.com/cosmtrek/air/master/install.sh \
#     && chmod +x install.sh && sh install.sh && cp ./bin/air /bin/air

# EXPOSE $HTTP_PORT

# WORKDIR /app

# CMD ["air", "air init"]

# WORKDIR /app

# COPY go.mod ./
# RUN go mod download

# COPY *.go ./

# RUN go build -o /svc-authentication

# EXPOSE 8080

# CMD [ "/svc-authentication" ]

RUN apk update && apk add --no-cache git
RUN apk add make

# CMD ["make run_migrate", "make run_seed"]

WORKDIR /app

COPY . .
RUN go install github.com/pressly/goose/v3/cmd/goose@latest
RUN go mod tidy

RUN go build -o /svc-auth

EXPOSE 8080

ENTRYPOINT ["/svc-auth"]