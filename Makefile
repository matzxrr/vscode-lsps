# Makefile for managing extracted VS Code language servers

.PHONY: all extract install-deps build clean

# Language servers to manage
LSP_SERVERS := css-lsp html-lsp json-lsp eslint-lsp

# Default target
all: extract

# Extract language servers from VS Code repositories
extract:
	@echo "Extracting language servers from VS Code repositories..."
	@go run tools/extractor/main.go

# Install node dependencies for all language servers
install-deps:
	@echo "Installing Node.js dependencies for all language servers..."
	@for server in $(LSP_SERVERS); do \
		echo "Installing dependencies for $$server..."; \
		cd $$server && npm install && cd ..; \
	done

# Install a specific language server's dependencies
install-deps-%:
	@echo "Installing dependencies for $*..."
	@cd $* && npm install

# Build all language servers
build: install-deps
	@echo "Building all language servers..."
	@for server in $(LSP_SERVERS); do \
		echo "Building $$server..."; \
		cd $$server && npm run compile && cd ..; \
	done

# Build a specific language server
build-%: install-deps-%
	@echo "Building $*..."
	@cd $* && npm run compile

# Run a specific language server
run-%: build-%
	@echo "Running $*..."
	@cd $* && npm start

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@for server in $(LSP_SERVERS); do \
		echo "Cleaning $$server..."; \
		rm -rf $$server/node_modules $$server/out $$server/dist; \
	done

# Helper to create a package.json for a server that doesn't have one
create-package-json-%:
	@echo "Creating package.json for $*..."
	@echo '{\n\
  "name": "$*",\n\
  "version": "1.0.0",\n\
  "description": "Extracted VS Code $* language server",\n\
  "main": "out/node/main.js",\n\
  "scripts": {\n\
    "compile": "tsc -p .",\n\
    "watch": "tsc -w -p .",\n\
    "start": "node out/node/main.js"\n\
  },\n\
  "dependencies": {\n\
    "vscode-languageserver": "^8.0.2",\n\
    "vscode-languageserver-textdocument": "^1.0.8"\n\
  },\n\
  "devDependencies": {\n\
    "typescript": "^4.9.5"\n\
  }\n\
}' > $*/package.json

# Helper to create a tsconfig.json for a server that doesn't have one
create-tsconfig-%:
	@echo "Creating tsconfig.json for $*..."
	@echo '{\n\
  "compilerOptions": {\n\
    "target": "es2020",\n\
    "module": "commonjs",\n\
    "moduleResolution": "node",\n\
    "sourceMap": true,\n\
    "outDir": "out",\n\
    "rootDir": "src",\n\
    "strict": true,\n\
    "esModuleInterop": true\n\
  },\n\
  "include": ["src"],\n\
  "exclude": ["node_modules"]\n\
}' > $*/tsconfig.json

# Create workspace file if needed
create-workspace:
	@echo "Creating VS Code workspace file..."
	@echo '{\n\
  "folders": [\n\
    { "path": "." },\n\
    { "path": "css-lsp" },\n\
    { "path": "html-lsp" },\n\
    { "path": "json-lsp" },\n\
    { "path": "eslint-lsp" }\n\
  ],\n\
  "settings": {}\n\
}' > vscode-lsps.code-workspace

# Help command
help:
	@echo "Available commands:"
	@echo "  make extract         - Extract language servers from VS Code repositories"
	@echo "  make install-deps    - Install Node.js dependencies for all language servers"
	@echo "  make install-deps-X  - Install dependencies for a specific language server (e.g., make install-deps-css-lsp)"
	@echo "  make build           - Build all language servers"
	@echo "  make build-X         - Build a specific language server (e.g., make build-css-lsp)"
	@echo "  make run-X           - Run a specific language server (e.g., make run-css-lsp)"
	@echo "  make clean           - Clean build artifacts"
	@echo "  make create-package-json-X - Create package.json for a specific server"
	@echo "  make create-tsconfig-X     - Create tsconfig.json for a specific server"
	@echo "  make create-workspace      - Create VS Code workspace file"
	@echo "  make help            - Show this help message"
