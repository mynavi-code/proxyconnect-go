SRCS = main.go
NAME = proxyconnect-go

build: build/$(NAME)

build/$(NAME): $(SRCS)
	go build -o build/$(NAME) $(SRCS)

.PHONY: clean
clean:
	rm build/*
	rm -rf dist/*

.PHONY: dist
dist:
	mkdir -p dist
	mkdir -p dist/darwin-amd64
	mkdir -p dist/darwin-386
	mkdir -p dist/linux-amd64
	mkdir -p dist/linux-386
	mkdir -p dist/windows-amd64
	mkdir -p dist/windows-386
	GOOS=darwin  GOARCH=amd64 go build -o dist/darwin-amd64/$(NAME)      $(SRCS)
	GOOS=darwin  GOARCH=386   go build -o dist/darwin-386/$(NAME)        $(SRCS)
	GOOS=linux   GOARCH=amd64 go build -o dist/linux-amd64/$(NAME)       $(SRCS)
	GOOS=linux   GOARCH=386   go build -o dist/linux-386/$(NAME)         $(SRCS)
	GOOS=windows GOARCH=amd64 go build -o dist/windows-amd64/$(NAME).exe $(SRCS)
	GOOS=windows GOARCH=386   go build -o dist/windows-386/$(NAME).exe   $(SRCS)
