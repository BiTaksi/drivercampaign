ready: generate-mock-all lint test-unit test-integration

wire:
	wire ./...

test-unit:
	go test ./internal/... ./pkg/... -race -coverprofile=coverage.out -covermode=atomic -v

test-integration:
	go test -tags integration ./internal/handler/... -race -coverprofile=coverage_integration.out -coverpkg=./internal/handler/... -covermode=atomic -v

i18n-merge:
	goi18n merge -outdir=./locale ./locale/active.*.toml ./locale/translate.*.toml

lint:
	golangci-lint run

lint-fix:
	golangci-lint run --fix

generate-mock-all:
	mockgen -source=./pkg/requestclient/requestclient.go -destination=./pkg/requestclient/mocks/requestclient_mock.go -package=mocks
	mockgen -source=./pkg/tokenizer/tokenizer.go -destination=./pkg/tokenizer/mocks/tokenizer_mock.go -package=mocks
	mockgen -source=./pkg/kafka/producer.go -destination=./pkg/kafka/mocks/producer_mock.go -package=mocks
	mockgen -source=./pkg/eventmanager/manager.go -destination=./pkg/eventmanager/mocks/manager_mock.go -package=mocks
	mockgen -source=./pkg/httpclient/httpclient.go -destination=./pkg/httpclient/mocks/httpclient_mock.go -package=mocks
	mockgen -source=./pkg/nrclient/nrclient.go -destination=./pkg/nrclient/mocks/nrclient_mock.go -package=mocks

