package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/golang/protobuf/proto"
	stan "github.com/nats-io/stan.go"
	stanpb "github.com/nats-io/stan.go/pb"
	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"

	pb "github.com/renegmed/nats-stream-watcher/proto"
)

func main() {

	natsURL := os.Getenv("NATS_ADDR")
	clusterID := os.Getenv("NATS_CLUSTER_ID")
	startOption := os.Getenv("START_OPT")
	topicPublishEpisode := os.Getenv("NATS_PUB_EPI_TOPIC")

	// Connect to NATS
	natsClient, err := stan.Connect(clusterID, uuid.NewV4().String(), stan.NatsURL(natsURL))
	if err != nil {
		log.Fatalf("Can't connect: %v.\nMake sure a NATS Streaming Server is running at: %s", err, natsURL)
	}
	defer natsClient.Close()

	// Start NATS subscriptions
	startSubscription(natsClient, topicPublishEpisode, watchEpisode, startOpt(startOption))

	log.Infof("Starting new watcher service")

	// Waiting for signal to shutdown.
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
}

func startSubscription(natsClient stan.Conn, topic string, handler stan.MsgHandler, startOpt stan.SubscriptionOption) {
	durableName := uuid.NewV4().String()

	if _, err := natsClient.QueueSubscribe(topic, durableName, handler, startOpt, stan.DurableName(durableName)); err != nil {
		natsClient.Close()
		log.Fatal(err)
	}
	log.Infof("Started new regular subscription")
}

func watchEpisode(natsMsg *stan.Msg) { // implement function requirement
	log.Debug("Received new post generation queue message.")

	var message pb.PublishEpisodeMessage
	if err := proto.Unmarshal(natsMsg.Data, &message); err != nil {
		log.Errorf("Failed to unmarshal queue message: %v", err)
		return
	}

	log.Printf("Watching on S %02d, E %02d of '%s' on '%s'", message.SeasonNo, message.EpisodeNo, message.SeriesName, message.EpisodeUrl)
}

func startOpt(optString string) stan.SubscriptionOption { // implement stan.SubscriptionOption interface???
	switch optString {
	default:
		return stan.StartAt(stanpb.StartPosition_NewOnly)
	case "MOST_RECENT":
		return stan.StartWithLastReceived()
	case "ALL":
		return stan.DeliverAllAvailable()
	}
}
