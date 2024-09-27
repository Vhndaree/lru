bench: 
	go test -bench . 

bench-store: 
	mkdir -p ./benchmark
	cat ./benchmark/new.txt > ./benchmark/old.txt
	go test -bench . > ./benchmark/new.txt

test: ## Run test
	go test -cover -race -short -v ./...

vet: ## Run go vet against code
	go vet ./...

help:
	@echo 'Usage: make <OPTIONS> ... <TARGETS>'
	@echo ''
	@echo 'Available targets are:'
	@echo ''
	@grep -E '^[ a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
        awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
	@echo ''

.PHONY: bench bench-store test vet
