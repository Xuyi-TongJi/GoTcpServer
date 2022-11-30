#!/bin/bash
# shellcheck disable=SC2035
protoc --go_out=. *.proto
