all: ui test-harness

.PHONY: test-harness
test-harness:
	go build -o test-harness pkg/cmd/test_harness/main.go

.PHONY: ui
ui:
	cd pkg/test_harness && yarn build

.PHONY: deps
deps:
	cd pkg/test_harness && yarn

test:
	go test ./...

run-test-harness: test-harness
	./test-harness
