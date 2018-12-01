.PHONY: run dependencies
all: clean sentry run
sentry: dependencies
		go build .
run: sentry
		./sentry
clean:
		rm -f ./sentry
dependencies:
		go get -d -v ./...
install: dependencies
		go install -v ./...
docker:
		docker build -t sentry .