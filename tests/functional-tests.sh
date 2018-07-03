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
tests_ran=0
IMAGE_NAME=ceph/daemon
CLUSTER_NAME_BASE=one-cluster
current_cluster_name=""
MAX_CLUSTERS=10

function start_test {
  # If a test starts another test, don't consider a new start
  if [ $nested_tests -eq 0 ]; then
    lastTest=${FUNCNAME[1]}
    start_time=$(date +%s.%N)
    printf '%-35s : ' "${lastTest}"
    tests_ran=$(($tests_ran + 1))
  fi
  nested_tests=$(($nested_tests + 1))
}

function fatal() {
  # Let's remove the installed trap
  # Failing here is not an issue
  trap - 0
  set +e

  if [ -e $err_file ]; then
    cat $err_file
    deleteFile $err_file
  fi

  # Let's try to delete all possible created clusters
  purge

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
    duration=$(echo "$end - $start_time" | bc -l)
    printf 'SUCCESS : %3.2f seconds\n' $duration
  fi
}

function failed() {
  printf 'ERROR  : %s\n' "${captionForFailure}"
  fatal
}

# shellcheck disable=SC2120
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
  local command="$*"
  captionForFailure="Failed with $command: $captionForFailure"
  ./cn "$@" &>"$err_file"
  runCnStatus=$?

  if [ -n "$runCnVerbose" ] || [ -n "$DEBUG" ]; then
    cat $err_file
  fi

  deleteFile $err_file
  return $runCnStatus
}

function countS3Objects {
  local bucket=$1
  local captionForFailure="Counting $bucket objects"
  runCnVerbose="True" runCn s3 ls "$current_cluster_name" $bucket | grep -a "s3://" | wc -l
}

function isS3ObjectExists {
  local item=$1
  local bucket
  local file
  bucket=$(echo ${item%/*})
  file=$(echo ${item#*/})
  local captionForFailure="Checking if ($item) $bucket/$file exists"
  runCnVerbose="True" runCn s3 ls "$current_cluster_name" $bucket | awk '{print $4}' | sed -e "s|s3://$bucket/||g" | grep -qw "$file"
}

function test_start {
  start_test
  for i in $(seq 0 $MAX_CLUSTERS); do
    current_cluster_name="$CLUSTER_NAME_BASE-$i"
    runCn cluster start -d $tmp_dir "$current_cluster_name"
  done
  runCn cluster ls
  reportSuccess
}

function test_help {
  start_test
  runCn -h
  reportSuccess
}

function test_stop {
  start_test
  runCn cluster stop "$current_cluster_name"
  reportSuccess
}

function test_status {
  start_test
  runCn cluster status "$current_cluster_name"
  reportSuccess
}

function test_restart {
  start_test
  for i in $(seq 0 $MAX_CLUSTERS); do
    current_cluster_name="$CLUSTER_NAME_BASE-$i"
    runCn cluster restart "$current_cluster_name"
  done
  reportSuccess
}

function test_logs {
  start_test
  runCn cluster logs "$current_cluster_name"
  reportSuccess
}

function purge() {
  for i in $(seq 0 $MAX_CLUSTERS); do
    current_cluster_name="$CLUSTER_NAME_BASE-$i"
    runCn cluster purge --yes-i-am-sure "$current_cluster_name"
  done
}

function test_purge {
  start_test
  purge
  reportSuccess
}

function test_image_update {
  start_test
  runCn image update $IMAGE_NAME
  reportSuccess
}

function test_image_list {
  start_test
  runCn image ls
  reportSuccess
}

function test_version {
  start_test
  runCn version
  reportSuccess
}

function test_update_check {
  start_test
  runCn update-check
  reportSuccess
}

function test_s3_mb {
  start_test
  runCn s3 mb "$current_cluster_name" $bucket
  reportSuccess
}

function test_s3_rb {
  start_test
  runCn s3 rb "$current_cluster_name" $bucket
  reportSuccess
}

function s3_put {
  local file=$1
  local bucket=$2
  runCn s3 put "$current_cluster_name" ${file} $bucket
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
  captionForFailure="test_s3_put_50x_4K: delta is $delta"
  [ "$delta" -eq 50 ];
  reportSuccess
}


function test_s3_get {
  start_test
  runCn s3 get "$current_cluster_name" $bucket/${file} get_file
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
  runCn s3 del "$current_cluster_name" $bucket/$file${file_extension}
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
  local initial_count
  initial_count=$(countS3Objects $bucket)
  test_s3_del_custom 50
  local final_count
  final_count=$(countS3Objects $bucket)
  local delta=$(($initial_count - $final_count))
  captionForFailure="test_s3_del_50x: delta is $delta"
  [ "$delta" -eq 50 ];
  reportSuccess
}

function test_s3_ls {
  start_test
  runCn s3 ls "$current_cluster_name" $bucket
  reportSuccess
}

function test_s3_la {
  start_test
  runCn s3 la "$current_cluster_name"
  reportSuccess
}

function test_s3_info {
  start_test
  runCn s3 info "$current_cluster_name" $bucket/${file}
  reportSuccess
}

function test_s3_du {
  start_test
  runCn s3 du "$current_cluster_name" $bucket/${file}
  reportSuccess
}

function test_s3_mv {
  start_test
  source=${1-$bucket/$file}
  dest=${2-$bucket/${file}.new}
  runCn s3 mv "$current_cluster_name" $source $dest
  isS3ObjectExists $dest
  reportSuccess
}

function test_s3_mv_custom {
  start_test
  source=${1-$bucket/$file}
  count=$2
  local bucket
  bucket=$(echo ${source%/*})
  local file
  file=$(echo ${source#*/})
  local initial_count
  initial_count=$(countS3Objects $bucket)
  for loop in $(seq 1 $count); do
    test_s3_mv ${bucket}/${file}.${loop}${file_extension} ${bucket}/${file}.${loop}
  done
  local final_count
  final_count=$(countS3Objects $bucket)
  local delta=$(($final_count - $initial_count))
  captionForFailure="test_s3_mv_custom: delta is $delta"
  [ "$delta" -eq 0 ];
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
  runCn s3 cp "$current_cluster_name" $bucket/${source} $bucket/$dest
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
  captionForFailure="test_s3_cp_custom: delta is $delta"
  [ "$delta" -eq "$count" ];
  reportSuccess
}

function test_s3_cp_50x {
  start_test
  test_s3_cp_custom ${file} 50
  reportSuccess
}

function test_s3_sync {
  start_test
  runCn s3 sync "$current_cluster_name" $tmp_dir $bucket
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
    test_start
    test_s3_mb
    for cli_test in "$@"; do
      $cli_test
    done
  else
    for test in version update_check image_update logs restart status stop start version image_update image_list status logs; do
      test_$test
    done

    for test in create_10_buckets delete_10_buckets mb put_50x_4K del_50x put_10MB get ls la info du cp_50x mv_50x_after_copy ; do
      test_s3_$test
    done

    test_s3_sync
    test_s3_rb

    test_restart

    # s3 again

    for test in status version image_update purge; do
      test_$test
    done
  fi
  trap - 0
}

export LC_ALL=
export LANG="en_US.UTF-8"
export LC_NUMERIC="en_US.UTF-8"
global_start_time=$(date +%s.%N)
main "$@"
global_stop_time=$(date +%s.%N)
global_duration=$(echo "$global_stop_time - $global_start_time" | bc -l)
printf "\nRan %d tests in %.2f seconds\n" "$tests_ran" "$global_duration"
