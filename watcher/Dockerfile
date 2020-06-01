FROM golang:1.14.0-alpine

WORKDIR src/watcher
COPY . . 
# ${GOPATH}/src/comm-nats-streaming/watcher
#WORKDIR ${GOPATH}/src/comm-nats-streaming/watcher

#RUN make build/docker

# end of first stage, beginning of the second one
#FROM alpine:3.10

#COPY --from=0 /watcher /watcher
RUN go mod vendor
RUN go build -o watcher main.go

CMD ["./watcher"]
