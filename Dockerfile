FROM golang:alpine

WORKDIR /go/src/message-queue
COPY . /go/src/message-queue
RUN go build -o ./dist/message-queue

EXPOSE 8080
EXPOSE 1234
EXPOSE 27017
ENTRYPOINT ["./dist/message-queue"]
