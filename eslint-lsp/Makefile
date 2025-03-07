# Variables
BUNDLE_NAME := eslintServer.js
BUNDLE_PATH := ./out/eslintServer.js
BIN_NAME := eslint-lsp
BIN_PATH := ./bin/eslint-lsp
GO_BUILD_CMD := go run .
INSTALL_BIN_DIR := /usr/local/bin
INSTALL_LIB_DIR := /usr/local/lib/eslint-lsp

.PHONY: all
all: build

.PHONY: build
build:
	@echo "Building $(BUNDLE_NAME) from source..."
	$(GO_BUILD_CMD)
	@echo "Build complete."

.PHONY: clean
clean:
	@echo "Cleaning up..."
	rm -rf ./out
	@echo "Clean complete."

.PHONY: install
install:
	@echo "Install binary to $(INSTALL_BIN_DIR)..."
	mkdir -p $(INSTALL_LIB_DIR)
	install -m 644 $(BUNDLE_PATH) $(INSTALL_LIB_DIR)/
	install -m 755 $(BIN_PATH) $(INSTALL_BIN_DIR)/$(BIN_NAME)
	@echo "Install complete."

.PHONY: uninstall
uninstall:
	@echo "Uninstalling..."
	rm -r $(INSTALL_BIN_DIR)/$(BIN_NAME)
	rm -rf $(INSTALL_LIB_DIR)
	@echo "Uninstall complete."
