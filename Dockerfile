FROM golang:latest

#RUN mkdir -p /go/src/github.com/andrioid/s3preview

COPY . /go/src/github.com/andrioid/s3preview

WORKDIR /go/src/github.com/andrioid/s3preview

RUN go get -d && go build 

EXPOSE 80

CMD ./s3preview