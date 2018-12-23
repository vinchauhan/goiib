FROM golang:1.11

RUN go get 
#RUN mkdir $GOPATH/src/github.com/vinchauhan/goiib

#COPY . $GOPATH/src/github.com/vinchauhan/goiib

WORKDIR $GOPATH/src/github.com/vinchauhan/goiib

RUN go get .../.

RUN go install

VOLUME [ "/build:$GOPATH/go/src/goiib/pkg" ]

COPY pkg/goiib /build

