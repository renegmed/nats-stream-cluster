package main

import (
	"net/http"
	"os"

	"github.com/gorilla/mux"
	stan "github.com/nats-io/stan.go"
	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
)

func main() {
	serverPort := os.Getenv("SERVER_ADDR")
	natsURL := os.Getenv("NATS_ADDR")
	clusterID := os.Getenv("NATS_CLUSTER_ID")

	if serverPort == "" {
		log.Fatalf("Error. Please provide server port number, env SERVER_ADDR.")
	}
	// Connect to NATS-Streaming
	natsClient, err := stan.Connect(clusterID, uuid.NewV4().String(), stan.NatsURL(natsURL))
	if err != nil {
		log.Fatalf("Can't connect: %v.\nMake sure a NATS Streaming Server is running at: %s", err, serverPort)
	}
	defer natsClient.Close()

	srv := server{
		natsClient: natsClient,
	}

	// Serve HTTP
	r := mux.NewRouter()
	r.HandleFunc("/publish", srv.HandlePublishEpisode)

	log.Infof("Starting HTTP server on '%s'", serverPort)

	if err := http.ListenAndServe(serverPort, r); err != nil {
		log.Fatal(err)
	}
}
