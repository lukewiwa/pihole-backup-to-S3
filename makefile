.PHONY: build-x86
build-x86:
	go build -o pcb2s3-x86

.PHONY: build-arm
build-arm:
	env GOOS=linux GOARCH=arm64 go build -o pcb2s3-arm
