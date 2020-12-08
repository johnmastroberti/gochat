build:
	go build ./msg
	go build ./db
	go build ./cmd/client
	go build ./cmd/server

install:
	go install ./cmd/client
	go install ./cmd/server
