BUILD_DIR := .build

.PHONY: \
	dep \
	test

$(BUILD_DIR)/vendor: Gopkg.toml Gopkg.lock
	dep ensure

deps: $(BUILD_DIR)/vendor
	@echo All dependencies successfully installed.

test:
	go test -v -cover github.com/Icenium/go-personal-data-filter/filter
