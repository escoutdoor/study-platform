.PHONY:build
build:
	CGO_ENABLED=0 GOOS=linux go build -o ./bin/study-platform ./cmd/study-platform/main.go

.PHONY:run
run: build
	./bin/study-platform
