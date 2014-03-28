
SOURCES=framer.go framer_test.go

all: $(SOURCES)
	go get -d && go test

