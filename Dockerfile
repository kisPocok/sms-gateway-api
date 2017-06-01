FROM golang:1.8.0

WORKDIR /go/src/github.com/kisPocok/sms-gateway-api/
COPY . .

RUN apt-get update
RUN apt-get install git curl netcat iproute2 net-tools telnet -y
RUN apt-get install golang-go build-essential -y
RUN go get github.com/Masterminds/glide

RUN make init_linux
RUN make deps
RUN make test
RUN make build_linux

ENTRYPOINT make run
