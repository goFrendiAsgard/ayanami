test:
	go test -race ./... -coverprofile=profile.out -count=1 -covermode=atomic

testv:
	go test -race ./... -v -coverprofile=profile.out -count=1 -covermode=atomic

coverage:
	go tool cover -html=profile.out

cleantest:
	rm -R .test-*

testgenerate:
	mkdir -p .test-gen && go build && ./ayanami init -p whatever -r github.com/whoever/whatever -d .test-gen && go build -o .test-gen/whatever/generator/whatever .test-gen/whatever/generator/main.go

testscaffold:
	.test-gen/whatever/generator/whatever scaffold

testbuild:
	.test-gen/whatever/generator/whatever build
