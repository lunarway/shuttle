#!/bin/bash

go build

# exit right away if build fails
buildExitCode=$?
if [[ $buildExitCode -ne 0 ]]; then
  exit $buildExitCode
fi

function assertErrorCode() {
  local expectedErrorCode=$1
  shift
  result=$(./shuttle "$@" 2>&1)
  result_status=$?
  if [[ $result_status -ne $expectedErrorCode ]]; then
    fail "Status code wasn't $expectedErrorCode but $result_status\nOutput: $result"
  fi
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
