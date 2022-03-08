#!/bin/bash

go build

# exit right away if build fails
buildExitCode=$?
if [[ $buildExitCode -ne 0 ]]; then
  exit $buildExitCode
fi

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
