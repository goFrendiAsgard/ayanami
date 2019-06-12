test:
	go test -race ./... -coverprofile=profile.out -count=1 -covermode=atomic

testv:
	go test -race ./... -v -coverprofile=profile.out -count=1 -covermode=atomic

coverage:
	go tool cover -html=profile.out

build:
	./build.sh

cleartestgen:
	rm ./.test-gen/whatever/deployable -R && rm ./.test-gen/whatever/generator -R

testgen:
	mkdir -p .test-gen && go build && ./ayanami init -p whatever -r github.com/whoever/whatever -d .test-gen -e full && cd .test-gen/whatever/generator && make scaffold && make build && cd ../deployable/megazord && go build && ./megazord
