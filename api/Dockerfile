FROM golang:1.12

ENV GO111MODULE=on

WORKDIR /go/src/app
COPY . .

RUN go get github.com/pilu/fresh
RUN go build

CMD [ "fresh" ]
