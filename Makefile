test:
	go test -race ./... -coverprofile=profile.out -count=1 -covermode=atomic

testv:
	go test -race ./... -v -coverprofile=profile.out -count=1 -covermode=atomic

coverage:
	go tool cover -html=profile.out

build:
	./build.sh

deletemegazord:
	rm ./.test-gen/whatever/deployable -R && rm ./.test-gen/whatever/generator -R

runmegazord:
	mkdir -p .test-gen && go build && ./ayanami init -p whatever -r github.com/whoever/whatever -d .test-gen && cd .test-gen/whatever/generator && go build && ./whatever scaffold && ./whatever build && cd ../deployable/megazord && go build && ./megazord
