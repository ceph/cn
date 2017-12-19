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
file_extension=""

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
    captionForFailure=""
  else
    failed
  fi
}

function runCn() {
  err_file=$(getTempFile $lastTest)
  local command="$@"
  captionForFailure="Failed with $command: $captionForFailure"
  ./cn "$@" &>"$err_file"
  runCnStatus=$?

  if [ -n "$runCnVerbose" ]; then
    cat $err_file
  fi

  deleteFile $err_file
  return $runCnStatus
}

function countS3Objects {
  local bucket=$1
  local captionForFailure="Counting $bucket objects"
  runCnVerbose="True" runCn s3 ls $bucket | grep -a "s3://" | wc -l
}

function isS3ObjectExists {
  local item=$1
  local bucket
  local file
  bucket=$(echo ${item%/*})
  file=$(echo ${item#*/})
  local captionForFailure="Checking if ($item) $bucket/$file exists"
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

function s3_put {
  local file=$1
  local bucket=$2
  runCn s3 put ${file} $bucket
  isS3ObjectExists ${bucket}/${file}
}

function test_s3_put_10MB {
  start_test
  captionForFailure="Cannot run dd" dd if=/dev/zero of=${file} bs=1048576 count=10 &>/dev/null
  s3_put ${file} ${bucket}
  deleteFile ${file}
  reportSuccess
}

function test_s3_put_custom {
  start_test
  local upload_count=$1
  local file_size=$2
  captionForFailure="Cannot run dd" dd if=/dev/zero of=${file} bs=$file_size count=1 &>/dev/null
  new_file=$file
  for i in $(seq 1 $upload_count); do
    new_file=$file.$i
    mv ${file} ${new_file}
    s3_put ${new_file} ${bucket}
    mv $new_file $file
  done
  deleteFile ${file}
  reportSuccess
}

function test_s3_put_50x_4K {
  start_test
  initial_count=$(countS3Objects $bucket)
  test_s3_put_custom 50 4096
  final_count=$(countS3Objects $bucket)
  delta=$(($final_count - $initial_count))
  captionForFailure="delta is $delta"
  [ "$delta" -eq 50 ];
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
  runCn s3 del $bucket/$file${file_extension}
  ! isS3ObjectExists ${bucket}/${file}${file_extension}
  reportSuccess
}

function test_s3_del_custom {
  start_test
  local upload_count=$1
  for i in `seq 1 $upload_count`; do
    test_s3_del $bucket/$file.${i}
  done
  reportSuccess
}

function test_s3_del_50x {
  start_test
  local initial_count=$(countS3Objects $bucket)
  test_s3_del_custom 50
  local final_count=$(countS3Objects $bucket)
  local delta=$(($initial_count - $final_count))
  captionForFailure="delta is $delta"
  [ "$delta" -eq 50 ];
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
  source=${1-$bucket/$file}
  dest=${2-$bucket/${file}.new}
  runCn s3 mv $source $dest
  isS3ObjectExists $dest
  reportSuccess
}

function test_s3_mv_custom {
  start_test
  source=${1-$bucket/$file}
  count=$2
  local bucket=$(echo ${source%/*})
  local file=$(echo ${source#*/})
  local initial_count=$(countS3Objects $bucket)
  for loop in $(seq 1 $count); do
    test_s3_mv ${bucket}/${file}.${loop}${file_extension} ${bucket}/${file}.${loop}
  done
  local final_count=$(countS3Objects $bucket)
  local delta=$(($final_count - $initial_count))
  captionForFailure="delta is $delta"
  # It's weird but mv actually copy the file....
  [ "$delta" -eq $count ];
  reportSuccess
}

function test_s3_mv_50x {
  start_test
  object=${1-$bucket/$file}
  test_s3_mv_custom "${object}" 50
  reportSuccess
}

function test_s3_mv_50x_after_copy {
  start_test
  file_extension=".copy" test_s3_mv_custom ${bucket}/${file} 50
  reportSuccess
}

function test_s3_cp {
  start_test
  source=${1-$file}
  dest=${2-$source}.copy
  runCn s3 cp $bucket/${source} $bucket/$dest
  isS3ObjectExists ${bucket}/${dest}
  reportSuccess
}

function test_s3_cp_custom {
  start_test
  source=$1
  count=$2
  initial_count=$(countS3Objects $bucket)
  for loop in $(seq 1 $count); do
    test_s3_cp ${file} ${file}.$loop
  done
  final_count=$(countS3Objects $bucket)
  delta=$(($final_count - $initial_count))
  captionForFailure="delta is $delta"
  [ "$delta" -eq 50 ];
  reportSuccess
}

function test_s3_cp_50x {
  start_test
  test_s3_cp_custom ${file} 50
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

function test_s3_create_10_buckets {
  local bucket
  start_test
  for i in `seq 1 10`; do
    bucket=bucket_$i
    test_s3_mb
  done
  reportSuccess
}

function test_s3_delete_10_buckets {
  local bucket
  start_test
  for i in `seq 1 10`; do
    bucket=bucket_$i
    test_s3_rb
  done
  reportSuccess
}

function report_configuration {
  OS=$(uname -s)
  case "$OS" in
    "Linux")
      CPU=$(grep "model name" /proc/cpuinfo | cut -d ":" -f 2 | head -1)
      RAM=$(free -h | grep "Mem" | awk '{print $2}')
      ;;
    "Darwin")
      CPU=$(sysctl -n machdep.cpu.brand_string)
      RAM=$(sysctl -a | grep hw.memsize)
      RAM=$(($RAM / 1024 / 1024))
      RAM="${RAM}G"
      ;;
    *)
      fatal "Unsupported platform $OS"
      ;;
  esac
  echo "Running tests on host $(hostname)"
  echo "OS Type   = $OS"
  echo "CPU Model = $CPU"
  echo "Total RAM = $RAM"
  echo "System    = $(uname -a)"
  echo
}


function main() {
  set -e
  trap failed 0

  report_configuration

  # Arguments given on the cli are test names run in sequence
  if [ $# -gt 0 ]; then
    test_purge
    test_start
    test_s3_mb
    for cli_test in "$@"; do
      $cli_test
    done
  else
    for test in version update purge logs restart status stop start version update status logs; do
      test_$test
    done

    for test in create_10_buckets delete_10_buckets mb put_50x_4K del_50x put_10MB get ls la info du cp_50x mv_50x_after_copy ; do
      test_s3_$test
    done

    file_extension=".copy" test_s3_del_50x
    test_s3_sync
    test_s3_rb

    test_restart

    # s3 again

    for test in status version update purge; do
      test_$test
    done
  fi
  trap - 0
}

main "$@"
