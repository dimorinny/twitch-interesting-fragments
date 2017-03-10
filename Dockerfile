FROM golang:1.8

WORKDIR /go/src/github.com/dimorinny/twitch-interesting-fragments/

ADD . .

RUN go get && go install

CMD /go/bin/twitch-interesting-fragments