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

function assertContains() {
  local expectedResult=$1
  local actualResult=$2
  if [[ ! "$actualResult" =~ "$expectedResult" ]]; then
    fail "Expected output to contain '$expectedResult', but it was:\n$actualResult"
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
  assertContains "Failed to load shuttle configuration" "$result"
}

test_can_get_variable_from_local_plan() {
  result=$(./shuttle -p examples/moon-base get docker.image 2>&1)
  assertEquals "earth-united/moon-base" "$result"
}

test_plan_from_relative_local_plan() {
  result=$(./shuttle -p examples/moon-base plan 2>&1)
  assertEquals "../station-plan" "$result"
}

test_plan_from_git_plan() {
  result=$(./shuttle -p examples/repo-project plan 2>&1)
  assertEquals "https://github.com/lunarway/shuttle-example-go-plan.git" "$result"
}

test_plan_from_git_plan_with_branch() {
  result=$(./shuttle -p examples/repo-project-branched plan 2>&1)
  assertEquals "https://github.com/lunarway/shuttle-example-go-plan.git#change-build" "$result"
}

test_plan_from_no_plan() {
  result=$(./shuttle -p examples/no-plan-project plan 2>&1)
  assertEquals "" "$result"
}

test_plan_with_template_from_no_plan() {
  result=$(./shuttle -p examples/no-plan-project plan --template '{{.PlanRaw}}' 2>&1)
  assertEquals "false" "$result"
}

test_can_get_variable_from_repo_plan() {
  result=$(./shuttle -p examples/repo-project get docker.destImage 2>&1)
  assertEquals "repo-project" "$result"
}

test_fails_getting_no_repo_plan() {
  assertErrorCode 4 -p examples/bad/no-repo-project ls
  assertContains "Failed executing git command \`clone" "$result"
}

test_fails_loading_invalid_shuttle_configuration() {
  assertErrorCode 2 -p examples/bad/yaml-invalid ls
  assertContains "Failed to parse shuttle configuration" "$result"
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
  assertContains "Exit code: 1" "$result"
}

test_run_shell_error_outputs_script_name() {
  assertErrorCode 4 -p examples/moon-base run crash
  if [[ ! "$result" =~ "crash" ]]; then
    fail "Expected output to contain the script name 'crash', but it was:\n$result"
  fi
}

test_run_shell_error_outputs_missing_arg() {
  assertErrorCode 1 -p examples/moon-base run required-arg
  if [[ ! "$result" =~ "required-arg" ]]; then
    fail "Expected output to contain the script name 'required-arg', but it was:\n$result"
  fi
}

test_template_local_path() {
  assertErrorCode 0 -p examples/moon-base template ../custom-template.tmpl -o Dockerfile-custom
}

test_run_repo_say_branch() {
  result=$(./shuttle -p examples/repo-project-branched run say)
  if [[ ! "$result" =~ "something clever" ]]; then
    fail "Expected output to contain 'something clever', but it was:\n$result"
  fi
}

test_run_repo_say() {
  result=$(./shuttle -p examples/repo-project run say)
  if [[ ! "$result" =~ "something masterly" ]]; then
    fail "Expected output to contain 'something masterly', but it was:\n$result"
  fi
}

test_run_repo_say_tagged() {
  result=$(./shuttle -p examples/repo-project --plan "#tagged" run say)
  if [[ ! "$result" =~ "something tagged" ]]; then
    fail "Expected output to contain 'something tagged', but it was:\n$result"
  fi
}

test_run_repo_say_sha() {
  result=$(./shuttle -p examples/repo-project --plan "#2b52c21" run say)
  if [[ ! "$result" =~ "something minor" ]]; then
    fail "Expected output to contain 'something minor', but it was:\n$result"
  fi
}


# Load and run shUnit2.
. ./shunit2
