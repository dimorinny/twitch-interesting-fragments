FROM golang:1.7

WORKDIR /go/src/github.com/dimorinny/twitch-interesting-fragments/

ADD . .

RUN mkdir /static
RUN go get && go install

CMD /go/bin/twitch-interesting-fragments