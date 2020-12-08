build:
	go build ./msg
	go build ./db
	go build ./cmd/gc-client
	go build ./cmd/gc-server

install:
	go install ./cmd/gc-client
	go install ./cmd/gc-server