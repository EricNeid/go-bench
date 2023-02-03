# SPDX-FileCopyrightText: 2021 Eric Neidhardt
# SPDX-License-Identifier: CC0-1.0

DIR := ${CURDIR}
GO_IMAGE := golang:1.20.0-alpine

.PHONY: build-windows
build-windows:
	docker run -it --rm \
		-e GOOS=windows \
		-e GOARCH=amd64 \
		-w /app -v ${DIR}:/app \
		${GO_IMAGE} \
		go build -o ./out/ ./cmd/gobench/


.PHONY: build-linux
build-linux:
	docker run -it --rm \
		-e GOOS=linux \
		-e GOARCH=amd64 \
		-w /app -v ${DIR}:/app \
		${GO_IMAGE} \
		go build -o ./out/ ./cmd/gobench/


.PHONY: test
test:
	docker run -it --rm \
		-w /app -v ${DIR}:/app \
		${GO_IMAGE} \
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
