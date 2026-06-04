BINARY = aac

.PHONY: build test clean fmt run

build:
	go build -o $(BINARY) .

test:
	go test ./...

clean:
	go clean
	rm -f $(BINARY) $(BINARY).test *.tar.gz *.zip

fmt:
	gofmt -l .

run:
	@if [ -z "$(FILE)" ]; then echo "Usage: make run FILE=file.aa"; exit 1; fi
	./$(BINARY) $(FILE)
