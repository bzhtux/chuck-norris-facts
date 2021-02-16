#!/usr/bin/env bash

set -euo pipefail

echo "TRAVIS_TEST_RESULT: ${TRAVIS_TEST_RESULT}"

TRIGGER={TRAVIS_TEST_RESULT:-1}

case $TRIGGER in
    0) curl -X POST \
        -H "Content-Type: application/json" \
        -d "{\"source_type\": \"Tag\", \"source_name\": \"${TRAVIS_TAG}\"}" \
        https://hub.docker.com/api/build/v1/source/6375f864-a015-4ec7-808e-0044240f81f7/trigger/f0175241-7846-4651-b66b-ea7280603026/call/
    ;;
    1) echo "Build failed, no trigger set for docker hub build"
    ;;
esac