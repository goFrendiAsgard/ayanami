test:
	go test ./... -coverprofile=profile.out -covermode=count

coverage:
	go tool cover -html=profile.out

run:
	go build && ./ayanami

help:
	go build && ./ayanami -h
