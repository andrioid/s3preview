FROM golang:latest

COPY . /go/src/github.com/andrioid/s3preview

WORKDIR /go/src/github.com/andrioid/s3preview

RUN go get -d && go build 

EXPOSE 80

# Just demonstrating what environmental variables are possible
ENV AWS_ACCESS_KEY_ID="" AWS_SECRET_ACCESS_KEY="" ASSSET_BUCKET="" PORT=80 ASSET_PREFIX="" PREVIEW_PREFIX=""

CMD ./s3preview