gen_proto:
	protoc --go_out=${GOPATH}/src proto/episode-publish.proto --grpc-gateway_out=logtostderr=true:.

.PHONY: proto
proto:
	protoc -I proto/ proto/*.proto --go_out=plugins=grpc:proto

up:
	docker-compose up --build -d 
.PHONY: up

cluster-info:
	curl http://127.0.0.1:8222/varz
	curl http://127.0.0.1:8222/routez
.PHONY: cluster-info

publish:
	curl -X POST http://127.0.0.1:9000/publish -d '{"series_name":"Dexter", "season_no":1,"episode_no":1,"episode_url":"https://hbo.com/dexter/1/1"}'
.PHONY: publish 

tail-api:
	docker logs neatflyx -f
tail-curious:
	docker logs watcher_curious -f
tail-patient:
	docker logs watcher_patient -f
tail-binge:
	docker logs watcher_binge -f
.PHONY: tail-neatflyx tail-curious tail-patient tail-binge