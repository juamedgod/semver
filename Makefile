.PHONY: test cover all lint

BUILD_DIR := ./artifacts

fmtcheck = @if goimports -l $(1) | read var; then echo "goimports check failed for $(1):\n `goimports -d $(1)`"; exit 1; fi

all:
	@$(MAKE) vet
	@$(MAKE) lint
	@$(MAKE) cover

get-build-deps:
	@echo "+ Downloading build dependencies"
	@go get golang.org/x/tools/cmd/goimports
	@go get github.com/golang/lint/golint

vet:
	@go vet .

lint:
	@golint .
	$(call fmtcheck, .)

test:
	@go test .


cover: test
	@mkdir -p $(BUILD_DIR)
	@go test -coverprofile=$(BUILD_DIR)/cover.out
	@go tool cover  -html=$(BUILD_DIR)/cover.out -o=$(BUILD_DIR)/coverage.html

