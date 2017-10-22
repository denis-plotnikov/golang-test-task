FROM golang

ADD . /go/src/golang-test-task

RUN go get golang.org/x/net/html 
RUN go install golang-test-task

ENTRYPOINT /go/bin/golang-test-task

EXPOSE 50000
