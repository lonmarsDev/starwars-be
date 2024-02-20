# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

# Name of the binary executable
BINARY_NAME=starwarsservice

all: test build

build:
	$(GOBUILD) -o $(BINARY_NAME) -v ./cmd/starwarsservice

test:
	$(GOTEST) -v ./...

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_NAME).exe

run:
	$(GOBUILD) -o $(BINARY_NAME) -v ./cmd/starwarsservice
	./$(BINARY_NAME)