# SPDX-FileCopyrightText: 2021 Eric Neidhardt
# SPDX-License-Identifier: CC0-1.0

DIR := ${CURDIR}

.PHONY: build-windows
build-windows:
	docker run -it --rm \
		-e CGO_ENABLED=0 \
		-e GOOS=windows \
		-e GOARCH=amd64 \
		-w /app -v ${DIR}:/app \
		golang:1.19.3-alpine \
		go build -o ./out/ ./cmd/gobench/


.PHONY: build-linux
build-linux:
	docker run -it --rm \
		-e CGO_ENABLED=0 \
		-e GOOS=linux \
		-e GOARCH=amd64 \
		-w /app -v ${DIR}:/app \
		golang:1.19.3-alpine \
		go build -o ./out/ ./cmd/gobench/


.PHONY: test
test:
	docker run -it --rm \
		-e CGO_ENABLED=0 \
		-w /app -v ${DIR}:/app \
		golang:1.19.3-alpine \
		go test ./...


.PHONY: lint
lint:
	docker run -it --rm \
		-e CGO_ENABLED=0 \
		-w /app -v ${DIR}:/app \
		golangci/golangci-lint:v1.50.1 \
		golangci-lint run ./...


.PHONY: clean
clean:
	rm -rf testlocal/
	rm -rf out/
