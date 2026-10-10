# Local entry points; CI runs the same steps (.github/workflows/ci.yml).
# Needs go, node/npm, python 3.11+ (with PyYAML) and a JDK 21 and PHP 8.2+ with Composer on PATH.

.PHONY: all lint generate test test-contract test-go test-typescript test-python test-kotlin test-php test-live drift

all: lint generate test

lint:
	npx --yes @redocly/cli@2.54.2 lint openapi.yaml

# Regenerates every SDK's models from openapi.yaml.
generate:
	cd sdks/go && go generate ./...
	cd sdks/typescript && npm run generate
	cd sdks/python && ./scripts/generate.sh && python scripts/unasync.py
	python3 scripts/gen_models.py kotlin sdks/kotlin/src/main/kotlin/io/github/pierreemmanuelj/alldebrid/Models.kt
	python3 scripts/gen_models.py php sdks/php/src

test: test-contract test-go test-typescript test-python test-kotlin test-php

test-contract:
	cd contract && go test -count=1 ./...

test-go:
	cd sdks/go && go vet ./... && go test -race -count=1 ./...

test-typescript:
	cd sdks/typescript && npm run check

test-python:
	cd sdks/python && ruff check . && ruff format --check . && mypy && pytest

test-kotlin:
	cd sdks/kotlin && ./gradlew build

test-php:
	cd sdks/php && composer install && vendor/bin/phpunit && vendor/bin/phpstan analyse

# Contract tests against a real account (read-only unless ALLDEBRID_LIVE_WRITE=1).
# The output may show response values; do not run it where it is public.
test-live:
	@test -n "$$ALLDEBRID_API_KEY" || (echo "set ALLDEBRID_API_KEY" && exit 1)
	cd contract && go test -count=1 -run TestLive -v ./...

# Has the official documentation changed since the spec was last reviewed?
drift:
	python3 scripts/upstream_changelog.py
