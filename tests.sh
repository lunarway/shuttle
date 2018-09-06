#!/bin/bash

go build

local failed_tests=0

function test {
  local bold=$'\e[01m'
  local green=$'\e[32m'
  local red=$'\e[31m'
  local reset=$'\e[00m'

  local testname=$1
  shift
  echo -n "├─ ${bold}$testname$reset"
  local result
  result=$(./shuttle "$@" 2>&1)
  local result_status=$?
  if [[ $result_status -eq 0 ]]; then
    echo -e "\e[0K\r├─ ${bold}$testname ── ${green}success$reset"
  else
    failed_tests=$((failed_tests + 1))
    echo -e "\e[0K\r${bold}${red}├─ ${bold}${red}$testname ── failed $reset"
    echo "$result" | sed "s/^/${bold}${red}│$reset /"
    echo "${bold}${red}├───────────────────────────────────────────────────────$reset"
  fi
}

echo "Running tests:"
test 'can run build' -p examples/moon-base run build test
test 'cant run undefined' -p examples/moon-bitch run build test

exit $failed_tests