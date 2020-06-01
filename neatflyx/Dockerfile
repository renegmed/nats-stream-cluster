FROM golang:1.14.0-alpine

WORKDIR src/neatflyx

COPY . .
#RUN pwd
#RUN ls -als
#${GOPATH}/src/comm-nats-streaming/neatflyx
#WORKDIR ${GOPATH}/src/comm-nats-streaming/neatflyx

RUN pwd
RUN ls -als

# RUN make build/docker

# end of first stage, beginning of the second one
# FROM alpine:3.10

#COPY --from=0 /neatflyx /neatflyx

RUN go mod vendor 

RUN go build -o neatflyx main.go server.go

CMD ["./neatflyx"]
