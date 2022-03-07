#!/bin/bash

go build

# exit right away if build fails
buildExitCode=$?
if [[ $buildExitCode -ne 0 ]]; then
  exit $buildExitCode
fi

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

test_fails_getting_no_repo_plan() {
  assertErrorCode 4 -p examples/bad/no-repo-project ls
  assertContains "Failed executing git command \`clone" "$result"
}

test_fails_loading_invalid_shuttle_configuration() {
  assertErrorCode 2 -p examples/bad/yaml-invalid ls
  assertContains "Failed to parse shuttle configuration" "$result"
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
  ./shuttle version &>/dev/null
  ./shuttle version --commit &>/dev/null
}

test_template_local_path() {
  assertErrorCode 0 -p examples/moon-base template ../custom-template.tmpl -o Dockerfile-custom GO_VERSION=1.16
}

test_template_local_path_alternate_delims() {
  result=$(./shuttle -p examples/moon-base template ../custom-template-alternate-delims.tmpl --delims ">>,<<")
  assertEquals "FROM earth-united/moon-base" "$result"
}

test_run_sub_dir_say() {
  result=$(cd examples/stepping-stone/sub-dir && ./../../../shuttle run say)
  pwd=$(pwd)
  normalizedOutput="${result##"$pwd"}"
  if [[ ! "$normalizedOutput" == "/examples/stepping-stone" ]]; then
    fail "Expected output to be '/examples/stepping-stone', but it was:\n$normalizedOutput"
  fi
}

# Load and run shUnit2.
. ./shunit2
