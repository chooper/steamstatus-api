FROM golang:1.6
ADD . /go/src/github.com/chooper/steamstatus-api
RUN go get github.com/chooper/steamstatus-api
RUN go install github.com/chooper/steamstatus-api
EXPOSE 10000
ENV PORT=10000
ENTRYPOINT /go/bin/steamstatus-api web
