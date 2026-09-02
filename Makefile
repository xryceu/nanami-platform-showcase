.DEFAULT_GOAL := verify

.PHONY: test test-race vet fmt-check frontend-install frontend-verify electron-install electron-verify verify

test:
	go test ./...

test-race:
	go test -race ./...

vet:
	go vet ./...

fmt-check:
	@test -z "$$(gofmt -l pkg)" || (gofmt -l pkg && exit 1)

frontend-install:
	npm --prefix frontend ci

frontend-verify:
	npm --prefix frontend run verify

electron-install:
	npm --prefix electron ci

electron-verify:
	npm --prefix electron run verify

verify: fmt-check vet test-race frontend-verify electron-verify
