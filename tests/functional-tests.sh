#!/bin/bash
# Functional tests
err_file=""
tmp_dir=/tmp
bucket=mybucket
file=dd_file
runCnStatus=0
lastTest=""

function fatal() {
  echo "$@"
  if [ -e $err_file ]; then
    cat $err_file
    deleteFile $err_file
  fi
  exit 1
}

function getTempFile() {
  filename=$tmp_dir/$1.XXXXX
  filename=$(mktemp $filename) &>/dev/null || fatal "Cannot create $filename"
  echo $filename
}

function deleteFile() {
  if [ -e "$1" ]; then
    rm -f $1 || fatal "Cannot delete $1"
  fi
}

function success {
  printf '%-20s : SUCCESS\n' ${lastTest}
}

function failed() {
  printf '%-20s : ERROR\n' ${lastTest}
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
  lastTest=${FUNCNAME[1]}
  err_file=$(getTempFile $lastTest)
  ./cn "$@" &>"$err_file"
  runCnStatus=$?
  deleteFile $err_file
  return $runCnStatus
}

}

function test_start {
  runCn start -d $tmp_dir
  reportSuccess
}

function test_help {
  runCn -h
  reportSuccess
}

function test_stop {
  runCn stop
  reportSuccess
}

function test_status {
  runCn status
  reportSuccess
}

function test_restart {
  runCn restart
  reportSuccess
}

function test_logs {
  runCn logs
  reportSuccess
}

function test_purge {
  runCn purge --yes-i-am-sure
  reportSuccess
}

function test_update {
  runCn update
  reportSuccess
}

function test_version {
  runCn version
  reportSuccess
}

function test_s3_mb {
  runCn s3 mb $bucket
  reportSuccess
}

function test_s3_rb {
  runCn s3 rb $bucket
  reportSuccess
}

function test_s3_put {
  dd if=/dev/zero of=${file} bs=1048576 count=10 &>/dev/null || fatal "Cannot run dd"
  runCn s3 put ${file} $bucket
  deleteFile ${file}
  reportSuccess
}

function test_s3_get {
  runCn s3 get $bucket/${file} get_file
  deleteFile get_file
  reportSuccess
}

function test_s3_del {
  if [ -z "$1" ]; then
    runCn s3 del $bucket/${file}
  else
    runCn s3 del $1
  fi
  reportSuccess
}

function test_s3_ls {
  runCn s3 ls $bucket
  reportSuccess
}

function test_s3_la {
  runCn s3 la
  reportSuccess
}

function test_s3_info {
  runCn s3 info $bucket/${file}
  reportSuccess
}

function test_s3_du {
  runCn s3 du $bucket/${file}
  reportSuccess
}

function test_s3_mv {
  runCn s3 mv $bucket/${file} $bucket/${file}.new
  reportSuccess
}
}

function test_s3_sync {
  runCn s3 sync $tmp_dir $bucket
  reportSuccess
}

function main() {
  set -e
  trap failed 0
  for test in version update purge logs restart status stop start version update status logs; do
    test_$test
  done

  for test in mb rb mb put get ls la info du mv sync; do
    test_s3_$test
  done

  test_s3_del $bucket/${file}.new
  test_s3_rb

  test_restart

  # s3 again

  for test in status version update purge; do
    test_$test
  done
  trap - 0
}

main
