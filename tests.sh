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

test_moon_base_builds_with_absolute_path() {
  assertRun -p $(pwd)/examples/moon-base run build tag=test
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

test_can_check_variable_exists() {
  ./shuttle -p examples/moon-base has run-as-root 2>&1
}

test_can_check_variable_exists_with_stdout() {
  result=$(./shuttle -p examples/moon-base has --stdout run-as-root 2>&1)
  if [[ "$result" != "true" ]]; then
    fail "Expected output to be 'true', but it was:\n$result"
  fi
}

test_can_check_variable_doesnt_exist() {
  assertErrorCode 1 -p examples/moon-base has not.a.thing 2>&1
}

test_can_check_variable_doesnt_exist_with_stdout() {
  result=$(./shuttle -p examples/moon-base has --stdout oh.no 2>&1)
  if [[ "$result" != "false" ]]; then
    fail "Expected output to be 'false', but it was:\n$result"
  fi
}


test_can_check_if_script_exists() {
  ./shuttle -p examples/moon-base has --script build 2>&1
}

test_can_check_if_script_does_not_exist() {
  assertErrorCode 1 -p examples/moon-base has --script do_not 2>&1
}

test_can_execute_shuttle_version_without_error() {
  ./shuttle version &> /dev/null
  ./shuttle version --commit &> /dev/null
}


test_run_shell_error_outputs_exit_code() {
  assertErrorCode 4 -p examples/moon-base run crash
  if [[ ! "$result" =~ "Exit code: 1" ]]; then
    fail "Expected output to contain 'Exit code: 1', but it was:\n$result"
  fi
}

# Load and run shUnit2.
. ./shunit2