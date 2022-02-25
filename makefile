.PHONY: build-x86
build-x86:
	go build -o pb2s3-x86

.PHONY: build-arm64
build-arm64:
	env GOOS=linux GOARCH=arm64 go build -o pb2s3-arm64

.PHONY: build-arm32
build-arm32:
	env GOOS=linux GOARCH=arm GOARM=5 go build -o pb2s3-arm32
