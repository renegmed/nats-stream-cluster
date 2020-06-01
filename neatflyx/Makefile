init-project:
	go mod init github.com/renegmed/nats-stream-neatflyx
up:

test:
	go test -race ./...

build/docker: test
	CGO_ENABLED=0 go build -a -installsuffix cgo -o /neatflyx .
