#!/usr/bin/env bash

set -eu

cd "$(dirname "$0")/.."

repo_name=${1:-}
if [ -z "$repo_name" ]; then
  echo "the repository name is required" >&2
  exit 1
fi

mkdir -p bin
curl -L -o bin/cc-test-reporter https://codeclimate.com/downloads/test-reporter/test-reporter-0.6.3-linux-amd64
chmod a+x bin/cc-test-reporter
export PATH="$PWD/bin:$PATH"
bash scripts/test-code-climate.sh "$repo_name"
