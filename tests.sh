#!/bin/bash

go build

function assertRun() {
  result=$(./shuttle "$@" 2>&1)
  result_status=$?
  if [[ $result_status -gt 0 ]]; then
    fail "Status code was $result_status\n$result"
  fi
}

function assertErrorCode() {
  local expectedErrorCode=$1
  shift
  result=$(./shuttle "$@" 2>&1)
  result_status=$?
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

test_can_get_variable_from_local_plan() {
  result=$(./shuttle -p examples/moon-base get docker.image 2>&1)
  assertEquals "earth-united/moon-base" "$result"
}

test_can_get_variable_from_repo_plan() {
  result=$(./shuttle -p examples/repo-project get docker.image 2>&1)
  assertEquals "shuttle/repo-project" "$result"
}

test_fails_getting_no_repo_plan() {
  assertErrorCode 4 -p examples/bad/no-repo-project ls
  if [[ ! "$result" =~ "Failed executing git command \`clone" ]]; then
    fail "Expected output to contain 'Failed executing git command \`clone', but it was:\n$result"
  fi
}

test_get_a_boolean() {
  result=$(./shuttle -p examples/moon-base get run-as-root 2>&1)
  if [[ "$result" != "false" ]]; then
    fail "Expected output to be 'false', but it was:\n$result"
  fi
}

# Load and run shUnit2.
. ./shunit2