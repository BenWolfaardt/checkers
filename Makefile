mock-expected-keepers:
	mockgen -source=x/checkers/types/expected_keepers.go \
		-package testutil \
		-destination=x/checkers/testutil/expected_keepers_mocks.go

build-linux:
	# For Debian
	# GOOS=linux GOARCH=amd64 go build -o ./build/checkersd-linux-amd64 ./cmd/checkersd/main.go
	# GOOS=linux GOARCH=arm64 go build -o ./build/checkersd-linux-arm64 ./cmd/checkersd/main.go
	# For Alpine
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./build/checkersd-linux-amd64 ./cmd/checkersd/main.go
    CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o ./build/checkersd-linux-arm64 ./cmd/checkersd/main.go
do-checksum-linux:
	cd build && sha256sum \
		checkersd-linux-amd64 checkersd-linux-arm64 \
		> checkers-checksum-linux

build-linux-with-checksum: build-linux do-checksum-linux

build-darwin:
	GOOS=darwin GOARCH=amd64 go build -o ./build/checkersd-darwin-amd64 ./cmd/checkersd/main.go
	GOOS=darwin GOARCH=arm64 go build -o ./build/checkersd-darwin-arm64 ./cmd/checkersd/main.go

build-all: build-linux build-darwin

do-checksum-darwin:
	cd build && sha256sum \
		checkersd-darwin-amd64 checkersd-darwin-arm64 \
		> checkers-checksum-darwin

build-darwin-with-checksum: build-darwin do-checksum-darwin

build-with-checksum: build-linux-with-checksum build-darwin-with-checksum
