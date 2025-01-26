.PHONY: openapi_http lint fmt test efficient_structs mockall coverage

include .env
export

openapi_http:
	@echo "Generating OpenAPI documentation..."
	@./scripts/openapi-http.sh api main

lint:
	@./scripts/lint.sh

fmt:
	goimports -l -w -d -v ./

test:
	@./scripts/test.sh .env
	@./scripts/test.sh .e2e.env

efficient_structs:
	@echo "Fixing structs..."
	@./scripts/structs_efficient.sh

mockall:
	@echo "Generating mocks..."
	@./scripts/mockall.sh

coverage:
	@echo "Generating coverage..."
	@./scripts/coverage.sh
