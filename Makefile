test:
	go test -race ./... -coverprofile=profile.out -count=1 -covermode=atomic

testv:
	go test -race ./... -v -coverprofile=profile.out -count=1 -covermode=atomic

coverage:
	go tool cover -html=profile.out

cleantest:
	rm -R .test-*

build:
	./build.sh

testgenerate:
	mkdir -p .test-gen && go build && ./ayanami init -p whatever -r github.com/whoever/whatever -d .test-gen && cd .test-gen/whatever/generator && go build && ./whatever scaffold && ./whatever build
