#!/bin/bash
# Functional tests
err_file=""
tmp_dir=/tmp

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

function runCn() {
  caller=${FUNCNAME[1]}
  err_file=$(getTempFile $caller)
  ./cn "$@" &>"$err_file"
  if [[ "$?" -eq 0 ]]; then
    printf '%-20s : SUCCESS\n' ${caller}
  else
    printf '%-20s : ERROR\n' ${caller}
    fatal
    exit 1
  fi
  deleteFile $err_file
}

function start {
  runCn start -d $tmp_dir
}

function help {
  runCn -h
}

function stop {
  runCn stop
}

function status {
  runCn status
}

function restart {
  runCn restart
}

function logs {
  runCn logs
}

function purge {
  runCn purge --yes-i-am-sure
}

function update {
  runCn update
}

function version {
  runCn version
}

function s3_mb {
  runCn s3 mb aaa
}

function s3_rb {
  runCn s3 rb aaa
}

function s3_put {
  dd if=/dev/zero of=dd_file bs=1048576 count=10 &>/dev/null || fatal "Cannot run dd"
  runCn s3 put dd_file aaa
  deleteFile dd_file
}

function s3_get {
  runCn s3 get aaa/dd_file get_file
  deleteFile  get_file
}

function s3_ls {
  runCn s3 ls aaa
}

function s3_la {
  runCn s3 la
}

function s3_info {
  runCn s3 info aaa/dd_file
}

function s3_du {
  runCn s3 du aaa/dd_file
}

function s3_mv {
  runCn s3 mv aaa/dd_file aaa/dd_file2
}

function s3_sync {
  runCn s3 sync $tmp_dir aaa
}

function main(){
  version
  update
  purge
  logs
  restart
  status
  stop
  start
  version
  update
  status
  logs

  s3_mb
  s3_rb
  s3_mb
  s3_put
  s3_get
  s3_ls
  s3_la
  s3_info
  s3_du
  s3_mv
  s3_sync

  ## rb all

  restart

  # s3 again

  status
  version
  update
  purge
}

main
