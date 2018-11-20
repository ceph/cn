#!/bin/bash
set -e

VERSION="$1"
TARBALL="cn-vendor-$VERSION.tar.xz"

fatal() {
  echo "$@"
  exit 1
}

if [ -z "$VERSION" ]; then
  fatal "A version is required as parameter"
fi

CURRENT_DIR=$(pwd)
TMPDIR=$(mktemp -d /tmp/cn-vendor.XXXX)
pushd "$TMPDIR" &>/dev/null
  mkdir src
  pushd src &>/dev/null
    wget -q https://github.com/ceph/cn/archive/v"$VERSION".tar.gz
    tar -xf "v$VERSION.tar.gz"
    pushd "cn-$VERSION" &>/dev/null
      GOPATH=$TMPDIR make prepare
      GOPATH=$TMPDIR dep status
      tar -cJf "$CURRENT_DIR/$TARBALL" vendor/
    popd
  popd
popd

echo "Vendor tarball for version $VERSION is available : $TARBALL"
rm -rf "$TMPDIR"
