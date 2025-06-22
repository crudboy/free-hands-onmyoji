APP_NAME = free-hands-onmyoji
CGO_CPPFLAGS = -I/Library/Developer/CommandLineTools/SDKs/MacOSX.sdk/usr/include/c++/v1
CGO_LDFLAGS = -L/opt/homebrew/opt/opencv/lib

.PHONY: all build run clean deps

all: build

build:
	CGO_CPPFLAGS="$(CGO_CPPFLAGS)" CGO_LDFLAGS="$(CGO_LDFLAGS)" go build -o $(APP_NAME) .

run: build
	./$(APP_NAME)

clean:
	rm -f $(APP_NAME)
test:
	CGO_CPPFLAGS="$(CGO_CPPFLAGS)" CGO_LDFLAGS="$(CGO_LDFLAGS)" go test -v ./...
deps:
	go mod tidy
	go mod download
