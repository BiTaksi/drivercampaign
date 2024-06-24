FROM golang:1.22.2-alpine

RUN apk --no-cache add build-base git gcc openssh-client curl

WORKDIR /app

RUN curl -fLo install.sh https://raw.githubusercontent.com/cosmtrek/air/master/install.sh \
    && chmod +x install.sh && sh install.sh && cp ./bin/air /bin/air

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . .

CMD air -c .air.toml
