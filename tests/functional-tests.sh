#!/bin/bash
# Functional tests
err_file=""
tmp_dir=/tmp
bucket=mybucket
file=dd_file
runCnVerbose=""
runCnStatus=0
lastTest=""
captionForFailure=""
start_time=0
nested_tests=0 # How many test_* are nested

function start_test {
  # If a test starts another test, don't consider a new start
  if [ $nested_tests -eq 0 ]; then
    lastTest=${FUNCNAME[1]}
    start_time=$(date +%s.%N)
    printf '%-35s : ' "${lastTest}"
  fi
  nested_tests=$(($nested_tests + 1))
}

function fatal() {
  if [ -e $err_file ]; then
    cat $err_file
    deleteFile $err_file
  fi
  exit 1
}

function getTempFile() {
  filename=$tmp_dir/$1.XXXXX
  local captionForFailure="Cannot create $filename"
  filename=$(mktemp $filename) &>/dev/null
  echo $filename
}

function deleteFile() {
  local captionForFailure="Cannot delete file $1"
  if [ -e "$1" ]; then
    rm -f $1
  fi
}

function success {
  nested_tests=$(($nested_tests - 1))
  # Until we reach the initial test, don't print anything
  if [ $nested_tests -eq 0 ]; then
    end=$(date +%s.%N)
    duration=$(echo "$end - $start_time" | bc -l | sed -e "s/\./,/g")
    printf 'SUCCESS : %3.2f seconds\n' $duration
  fi
}

function failed() {
  printf 'ERROR  : %s\n' "${captionForFailure}"
  fatal
}

function reportSuccess {
  if [ -n "$1" ]; then
    returnCode=$1
  else
    returnCode=$runCnStatus
  fi

  if [[ "$returnCode" -eq 0 ]]; then
    success
  else
    failed
  fi
}

function runCn() {
  err_file=$(getTempFile $lastTest)
  ./cn "$@" &>"$err_file"
  runCnStatus=$?

  if [ -n "$runCnVerbose" ]; then
    cat $err_file
  fi

  deleteFile $err_file
  return $runCnStatus
}

function isS3ObjectExists {
  local bucket=$1
  local file=$2
  local captionForFailure="Checking if $bucket/$file exists"
  runCnVerbose="True" runCn s3 ls $bucket | awk '{print $4}' | sed -e "s|s3://$bucket/||g" | grep -qw "$file"
}

function test_start {
  start_test
  runCn start -d $tmp_dir
  reportSuccess
}

function test_help {
  start_test
  runCn -h
  reportSuccess
}

function test_stop {
  start_test
  runCn stop
  reportSuccess
}

function test_status {
  start_test
  runCn status
  reportSuccess
}

function test_restart {
  start_test
  runCn restart
  reportSuccess
}

function test_logs {
  start_test
  runCn logs
  reportSuccess
}

function test_purge {
  start_test
  runCn purge --yes-i-am-sure
  reportSuccess
}

function test_update {
  start_test
  runCn update
  reportSuccess
}

function test_version {
  start_test
  runCn version
  reportSuccess
}

function test_s3_mb {
  start_test
  runCn s3 mb $bucket
  reportSuccess
}

function test_s3_rb {
  start_test
  runCn s3 rb $bucket
  reportSuccess
}

function test_s3_put {
  start_test
  captionForFailure="Cannot run dd" dd if=/dev/zero of=${file} bs=1048576 count=10 &>/dev/null
  runCn s3 put ${file} $bucket
  isS3ObjectExists ${bucket} ${file}
  deleteFile ${file}
  reportSuccess
}

function test_s3_get {
  start_test
  runCn s3 get $bucket/${file} get_file
  deleteFile get_file
  reportSuccess
}

function test_s3_del {
  start_test
  local bucket=$bucket
  local file=$file
  if [ $# -eq 1 ]; then
    bucket=$(echo ${1%/*})
    file=$(echo ${1#*/})
  fi
  runCn s3 del $bucket/$file
  ! isS3ObjectExists ${bucket} ${file}
  reportSuccess
}

function test_s3_ls {
  start_test
  runCn s3 ls $bucket
  reportSuccess
}

function test_s3_la {
  start_test
  runCn s3 la
  reportSuccess
}

function test_s3_info {
  start_test
  runCn s3 info $bucket/${file}
  reportSuccess
}

function test_s3_du {
  start_test
  runCn s3 du $bucket/${file}
  reportSuccess
}

function test_s3_mv {
  start_test
  runCn s3 mv $bucket/${file} $bucket/${file}.new
  isS3ObjectExists ${bucket} ${file}.new
  reportSuccess
}

function test_s3_cp {
  start_test
  runCn s3 cp $bucket/${file} $bucket/${file}.copy
  isS3ObjectExists ${bucket} ${file}.copy
  reportSuccess
}

function test_s3_sync {
  start_test
  runCn s3 sync $tmp_dir $bucket
  reportSuccess
}

#function test_template {
#start_test
#runCn
#reportSuccess
#}
function main() {
  set -e
  trap failed 0
  for test in version update purge logs restart status stop start version update status logs; do
    test_$test
  done

  for test in mb rb mb put get ls la info du cp mv sync; do
    test_s3_$test
  done

  test_s3_del $bucket/${file}.new
  test_s3_del $bucket/${file}.copy
  test_s3_rb

  test_restart

  # s3 again

  for test in status version update purge; do
    test_$test
  done
  trap - 0
}

main
