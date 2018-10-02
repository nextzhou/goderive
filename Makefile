install:
	go install -ldflags "-X main.Version=`git describe --always`" .
	go generate ./...
	go install -ldflags "-X main.Version=`git describe --always`" .

test:
	go test -v -race ./...
