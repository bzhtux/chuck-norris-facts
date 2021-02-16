#!/usr/bin/env bash

set -euo pipefail

TRIGGER=${TRAVIS_TEST_RESULT:-1}
DOCKER_TAG=${TRAVIS_TAG:-master}

case $TRIGGER in
    0) case $DOCKER_TAG in
            master) SRC_NAME="Branch"
                ;;
            *) SRC_NAME="Tag"
                ;;
        esac
        ;;
    1) echo "Tests failed, no trigger set for docker hub build"
        ;;
esac

curl -X POST \
        -H "Content-Type: application/json" \
        -d "{\"source_type\": \"${SRC_NAME}\", \"source_name\": \"${DOCKER_TAG}\"}" \
        https://hub.docker.com/api/build/v1/source/6375f864-a015-4ec7-808e-0044240f81f7/trigger/f0175241-7846-4651-b66b-ea7280603026/call/