.PHONY: run
all: clean sentry run
sentry:
		go build .
run:
		./sentry
clean:
		rm -f ./sentry