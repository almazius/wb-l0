all: downDocker main utils/publisher

main:
	go mod download
	docker compose build
	docker compose up -d
	sleep 3
	go run cmd/main.go

publisher:
	#go mod tidy
	go run utils/publisher/main.go

clearVolumes:
	cd consumer
	docker compose down --values

downDocker:
	docker compose down
