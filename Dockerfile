FROM golang:1.8.0

WORKDIR /www/go/src/msgbird
COPY . .

RUN apt-get update
RUN apt-get install git curl netcat iproute2 net-tools telnet -y
RUN apt-get install golang-go build-essential -y

RUN make init_linux
RUN make install
RUN make test
RUN make build_linux

ENTRYPOINT make run
