APP_NAME = orihon
OUTPUT   = index.html
MARKDOWN_FILES := $(shell find content -type f -name '*.md')
TEMPLATE_FILES := $(shell find template -type f)

$(OUTPUT): $(MARKDOWN_FILES) $(TEMPLATE_FILES)
	go run ./cmd/main.go -input=./content -output=./$(OUTPUT)

all: $(OUTPUT)

build:
	go build -o $(APP_NAME) cmd/main.go

run: build
	./$(APP_NAME) -input=./content -output=./$(OUTPUT)

clean:
	rm -f $(OUTPUT) $(APP_NAME) $(APP_NAME).txtar

test:
	go test -v ./...

txtar:
	find . -name .git -prune -o -type f \
		-not -name "$(APP_NAME).txtar" \
		-not -name "$(APP_NAME)" \
		-not -name "*.png" \
		-not -name "*.jpg" \
		-not -name "*.jpeg" \
		-not -name "*.gif" \
		-not -name "*.svg" \
		-print \
	| xargs go run golang.org/x/exp/cmd/txtar@latest > $(APP_NAME).txtar

.PHONY: all build run clean test txtar
