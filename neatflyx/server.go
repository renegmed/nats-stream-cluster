package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/golang/protobuf/proto"
	stan "github.com/nats-io/stan.go"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	pb "github.com/renegmed/nats-stream-neatflyx/proto"
)

type publishRequest struct {
	SeriesName string `json:"series_name,omitempty"`
	SeasonNo   int    `json:"season_no,omitempty"`
	EpisodeNo  int    `json:"episode_no,omitempty"`
	EpisodeURL string `json:"episode_url,omitempty"`
}

// var (
// 	//topicPublishEpisode = "episodes:publish"
// 	topicPublishEpisode = os.Getenv("NATS_PUB_EPI_TOPIC")
// )

type server struct {
	natsClient stan.Conn
}

func (s server) HandlePublishEpisode(rw http.ResponseWriter, req *http.Request) {
	topicPublishEpisode := os.Getenv("NATS_PUB_EPI_TOPIC")
	switch req.Method {
	case "POST":
		s.publishEpisode(rw, req, topicPublishEpisode)
	default:
		log.Errorf("Invalid reques method: %s", req.Method)
		http.Error(rw, "Invalid request", http.StatusBadRequest)
	}
}

func (s server) publishEpisode(rw http.ResponseWriter, req *http.Request, topic string) {
	var pubReq publishRequest
	if err := json.NewDecoder(req.Body).Decode(&pubReq); err != nil {
		log.Errorf("Failed to read request: %v", err)
		http.Error(rw, "Invalid request", http.StatusBadRequest)
		return
	}

	message := &pb.PublishEpisodeMessage{
		SeriesName: pubReq.SeriesName,
		SeasonNo:   int64(pubReq.SeasonNo),
		EpisodeNo:  int64(pubReq.EpisodeNo),
		EpisodeUrl: pubReq.EpisodeURL,
	}

	if err := s.publishMessage(topic, message); err != nil {
		log.Errorf("Failed to publish message onto queue '%s': %v", topic, err)
		http.Error(rw, "", http.StatusInternalServerError)
		return
	}

	log.Printf("Publishing on S%02dE%02d of '%s' on '%s'\n", message.SeasonNo, message.EpisodeNo, message.SeriesName, message.EpisodeUrl)
	fmt.Fprint(rw, "Post publication is pending")
}
func (s server) publishMessage(topic string, msg proto.Message) error {
	bs, err := proto.Marshal(msg)
	if err != nil {
		return errors.Wrap(err, "failed to marshal proto message")
	}

	if err := s.natsClient.Publish(topic, bs); err != nil {
		return errors.Wrap(err, "failed to publish message")
	}

	return nil
}
