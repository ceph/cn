#!/bin/bash
# Functional tests

function start {
  err_file=/tmp/err_cn_start
  ./cn start -d /tmp >"$err_file" 2>&1
  if [[ "$?" -eq 0 ]]; then
    echo "${FUNCNAME[0]}: SUCCESS"
  else
    echo "${FUNCNAME[0]}: ERROR"
    cat "$err_file"
    exit 1
  fi
  rm -f "$err_file"
}

function help {
  err_file=/tmp/err_cn_help
  ./cn -h >"$err_file" 2>&1
  if [[ "$?" -eq 0 ]]; then
    echo "${FUNCNAME[0]}: SUCCESS"
  else
    echo "${FUNCNAME[0]}: ERROR"
    cat "$err_file"
    exit 1
  fi
  rm -f "$err_file"
}

function stop {
  err_file=/tmp/err_cn_stop
  ./cn stop >"$err_file" 2>&1
  if [[ "$?" -eq 0 ]]; then
    echo "${FUNCNAME[0]}: SUCCESS"
  else
    echo "${FUNCNAME[0]}: ERROR"
    cat "$err_file"
    exit 1
  fi
  rm -f "$err_file"
}

function status {
  err_file=/tmp/err_cn_status
  ./cn status >"$err_file" 2>&1
  if [[ "$?" -eq 0 ]]; then
    echo "${FUNCNAME[0]}: SUCCESS"
  else
    echo "${FUNCNAME[0]}: ERROR"
    cat "$err_file"
    exit 1
  fi
  rm -f "$err_file"
}

function restart {
  err_file=/tmp/err_cn_restart
  ./cn restart >"$err_file" 2>&1
  if [[ "$?" -eq 0 ]]; then
    echo "${FUNCNAME[0]}: SUCCESS"
  else
    echo "${FUNCNAME[0]}: ERROR"
    cat "$err_file"
    exit 1
  fi
  rm -f "$err_file"
}

function logs {
  err_file=/tmp/err_cn_logs
  ./cn logs >"$err_file" 2>&1
  if [[ "$?" -eq 0 ]]; then
    echo "${FUNCNAME[0]}: SUCCESS"
  else
    echo "${FUNCNAME[0]}: ERROR"
    cat "$err_file"
    exit 1
  fi
  rm -f "$err_file"
}

function purge {
  err_file=/tmp/err_cn_purge
  ./cn purge --yes-i-am-sure >"$err_file" 2>&1
  if [[ "$?" -eq 0 ]]; then
    echo "${FUNCNAME[0]}: SUCCESS"
  else
    echo "${FUNCNAME[0]}: ERROR"
    cat "$err_file"
    exit 1
  fi
  rm -f "$err_file"
}

function update {
  err_file=/tmp/err_cn_update
  ./cn update >"$err_file" 2>&1
  if [[ "$?" -eq 0 ]]; then
    echo "${FUNCNAME[0]}: SUCCESS"
  else
    echo "${FUNCNAME[0]}: ERROR"
    cat "$err_file"
    exit 1
  fi
  rm -f "$err_file"
}

function version {
  err_file=/tmp/err_cn_version
  ./cn version >"$err_file" 2>&1
  if [[ "$?" -eq 0 ]]; then
    echo "${FUNCNAME[0]}: SUCCESS"
  else
    echo "${FUNCNAME[0]}: ERROR"
    cat "$err_file"
    exit 1
  fi
  rm -f "$err_file"
}

function s3_mb {
  err_file=/tmp/err_cn_s3_mb
  ./cn s3 mb aaa >"$err_file" 2>&1
  if [[ "$?" -eq 0 ]]; then
    echo "${FUNCNAME[0]}: SUCCESS"
  else
    echo "${FUNCNAME[0]}: ERROR"
    cat "$err_file"
    exit 1
  fi
  rm -f "$err_file"
}

function s3_rb {
  err_file=/tmp/err_cn_s3_rb
  ./cn s3 rb aaa >"$err_file" 2>&1
  if [[ "$?" -eq 0 ]]; then
    echo "${FUNCNAME[0]}: SUCCESS"
  else
    echo "${FUNCNAME[0]}: ERROR"
    cat "$err_file"
    exit 1
  fi
  rm -f "$err_file"
}

function s3_put {
  err_file=/tmp/err_cn_s3_put
  dd if=/dev/zero of=/tmp/ooo bs=1m count=10
  ./cn s3 put ooo aaa >"$err_file" 2>&1
  if [[ "$?" -eq 0 ]]; then
    echo "${FUNCNAME[0]}: SUCCESS"
  else
    echo "${FUNCNAME[0]}: ERROR"
    cat "$err_file"
    exit 1
  fi
  rm -f "$err_file" /tmp/ooo
}

function s3_get {
  err_file=/tmp/err_cn_s3_get
  ./cn s3 get aaa/ooo iii >"$err_file" 2>&1
  if [[ "$?" -eq 0 ]]; then
    echo "${FUNCNAME[0]}: SUCCESS"
  else
    echo "${FUNCNAME[0]}: ERROR"
    cat "$err_file"
    exit 1
  fi
  rm -f "$err_file" /tmp/iii
}

function s3_ls {
  err_file=/tmp/err_cn_s3_ls
  ./cn s3 ls aaa >"$err_file" 2>&1
  if [[ "$?" -eq 0 ]]; then
    echo "${FUNCNAME[0]}: SUCCESS"
  else
    echo "${FUNCNAME[0]}: ERROR"
    cat "$err_file"
    exit 1
  fi
  rm -f "$err_file"
}

function s3_la {
  err_file=/tmp/err_cn_s3_la
  ./cn s3 la >"$err_file" 2>&1
  if [[ "$?" -eq 0 ]]; then
    echo "${FUNCNAME[0]}: SUCCESS"
  else
    echo "${FUNCNAME[0]}: ERROR"
    cat "$err_file"
    exit 1
  fi
  rm -f "$err_file"
}

function s3_info {
  err_file=/tmp/err_cn_s3_info
  ./cn s3 info aaa/ooo >"$err_file" 2>&1
  if [[ "$?" -eq 0 ]]; then
    echo "${FUNCNAME[0]}: SUCCESS"
  else
    echo "${FUNCNAME[0]}: ERROR"
    cat "$err_file"
    exit 1
  fi
  rm -f "$err_file"
}

function s3_du {
  err_file=/tmp/err_cn_s3_du
  ./cn s3 du aaa/ooo >"$err_file" 2>&1
  if [[ "$?" -eq 0 ]]; then
    echo "${FUNCNAME[0]}: SUCCESS"
  else
    echo "${FUNCNAME[0]}: ERROR"
    cat "$err_file"
    exit 1
  fi
  rm -f "$err_file"
}

function s3_mv {
  err_file=/tmp/err_cn_s3_mv
  ./cn s3 mv aaa/ooo aaa/uuu >"$err_file" 2>&1
  if [[ "$?" -eq 0 ]]; then
    echo "${FUNCNAME[0]}: SUCCESS"
  else
    echo "${FUNCNAME[0]}: ERROR"
    cat "$err_file"
    exit 1
  fi
  rm -f "$err_file"
}

function s3_sync {
  err_file=/tmp/err_cn_s3_sync
  ./cn s3 sync /tmp aaa >"$err_file" 2>&1
  if [[ "$?" -eq 0 ]]; then
    echo "${FUNCNAME[0]}: SUCCESS"
  else
    echo "${FUNCNAME[0]}: ERROR"
    cat "$err_file"
    exit 1
  fi
  rm -f "$err_file"
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
