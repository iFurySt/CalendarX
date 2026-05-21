PROJECT ?=
SLUG ?=
CALENDARX ?= go run ./cmd/calendarx

.PHONY: init new-history new-plan ci test generate build

init:
	@if [ -z "$(PROJECT)" ]; then echo "usage: make init PROJECT=my-project"; exit 1; fi
	./scripts/init-project.sh "$(PROJECT)"

new-history:
	@if [ -z "$(SLUG)" ]; then echo "usage: make new-history SLUG=my-change"; exit 1; fi
	./scripts/new-history.sh "$(SLUG)"

new-plan:
	@if [ -z "$(SLUG)" ]; then echo "usage: make new-plan SLUG=my-plan"; exit 1; fi
	./scripts/new-exec-plan.sh "$(SLUG)"

test:
	go test ./...

generate:
	$(CALENDARX) generate

build:
	$(CALENDARX) build

ci:
	go test ./...
	go run ./cmd/calendarx --help >/dev/null
