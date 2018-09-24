#!/bin/bash

go build

function assertRun() {
  local result
  result=$(./shuttle "$@" 2>&1)
  local result_status=$?
  if [[ $result_status -gt 0 ]]; then
    fail "Status code was $result_status\n$result"
  fi
}

function assertErrorCode() {
  local expectedErrorCode=$1
  shift
  local result
  result=$(./shuttle "$@" 2>&1)
  local result_status=$?
  if [[ $result_status -ne $expectedErrorCode ]]; then
    fail "Status code wasn't $expectedErrorCode but $result_status\nOutput: $result"
  fi
}

test_moon_base_builds() {
  assertRun -p examples/moon-base run build tag=test
}

test_venus_base_fails() {
  assertErrorCode 2 -p examples/venus-base run build tag=test
}

# Load and run shUnit2.
. ./shunit2