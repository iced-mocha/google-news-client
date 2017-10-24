FROM golang:1.9

RUN go get -u github.com/golang/dep/cmd/dep && go install github.com/golang/dep/cmd/dep

WORKDIR /go/src/github.com/icedmocha/google-news-client
COPY . /go/src/github.com/icedmocha/google-news-client

RUN dep ensure && go install && /bin/bash -c "source workspace.env" 

ENTRYPOINT ["google-news-client"]
