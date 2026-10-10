#!/bin/sh
# Regenerates src/open_alldebrid/models.py from the spec.
# Needs datamodel-code-generator==0.82.0 and ruff on PATH.
set -eu
cd "$(dirname "$0")/.."
datamodel-codegen \
  --input ../../openapi.yaml \
  --input-file-type openapi \
  --output-model-type typing.TypedDict \
  --target-python-version 3.11 \
  --use-standard-collections \
  --use-union-operator \
  --use-schema-description \
  --use-field-description \
  --disable-timestamp \
  --formatters ruff-format \
  --custom-file-header '# Generated from openapi.yaml by datamodel-code-generator. DO NOT EDIT.' \
  --output src/open_alldebrid/models.py
