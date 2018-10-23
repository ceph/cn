#!/bin/bash

#############
# VARIABLES #
#############
GITHUB_USER=ceph
repo=cn
CREATE_TAG=

#############
# FUNCTIONS #
#############
fatal() {
  echo "$@"
  if [ -e "$CHANGELOG" ]; then
    rm -f "$CHANGELOG"
  fi

  #If the tag was created by us, let's delete it
  if [ -n "$CREATE_TAG" ]; then
    git checkout "$GIT_BRANCH"
    git tag | grep -qw "$TAG" && git tag -d "$TAG"
  fi
  exit 1
}

isBinaryExists() {
  command -v "$1" &>/dev/null || fatal "Cannot find $1 binary, please check your environement"
}

isVariableExists() {
  variable_name=$1
  value=${!variable_name}

  if [ -z "$value" ]; then
    fatal "Please define $variable_name"
  fi
}

isGitRepositoryClean() {
  git diff --no-ext-diff --quiet --exit-code
}

isGitTagExists() {
  git tag -l | grep -qw "$1"
}

getGoEnv() {
  eval "$(go env | grep ^GOHOST)"
  [ -n "$GOHOSTARCH" ] || fatal "Cannot determine GOHOSTARCH"
  [ -n "$GOHOSTOS" ] || fatal "Cannot determine GOHOSTOS"
  export LOCAL_ARCH="$GOHOSTOS-$GOHOSTARCH"
}

usage() {
  cat << EOF
  $0 - create a github release and upload files

  -h            : show help message
  -g <token>    : github token (can be defined with GITHUB_TOKEN variable)
  -t <tag>      : tag to be released (can be defined with TAG variable)
  -p <tag>      : previous tag to make the CHANGELOG (can be defined with PTAG variable)
EOF
  exit 2
}

########
# MAIN #
########
isBinaryExists go
isBinaryExists git

getGoEnv
echo "Building on $LOCAL_ARCH"

isGitRepositoryClean || fatal "git repository is not clean, cannot make the release !"

GIT_BRANCH=$(git rev-parse --abbrev-ref HEAD)
GIT_LAST_COMMIT=$(git log --format="%H" -n 1)

if command -v github-release &>/dev/null; then
  echo "Installing github-release"
  go get github.com/aktau/github-release
  isBinaryExists github-release
fi

optspec=":hg:p:t:"
while getopts "$optspec" optchar; do
  case "${optchar}" in
    h)
      usage
      ;;
    g)
      export GITHUB_TOKEN=${OPTARG}
      ;;
    t)
      export TAG=${OPTARG}
      ;;
    p)
      export PTAG=${OPTARG}
      ;;
    *)
      if [ "$OPTERR" != 1 ] || [ "${optspec:0:1}" = ":" ]; then
        echo "Non-option argument: '-${OPTARG}'" >&2
      fi
      ;;
  esac
done

isVariableExists GITHUB_TOKEN
isVariableExists TAG

if [[ ${TAG:0:1} != "v" ]]; then
  fatal "The tag ($TAG) should start with a 'v' like in v1.4.0"
fi

if ! isGitTagExists "$PTAG"; then
  IMPLICIT_PTAG=$(git for-each-ref refs/tags --sort=-taggerdate --format='%(refname)' --count=1 | cut -d '/' -f 3)
  if [ -z "$IMPLICIT_PTAG" ]; then
    fatal "Cannot detect any previous release"
  fi
  while true; do
    echo -n "Does $IMPLICIT_PTAG the git tag to consider for builiding the CHANGELOG ? (yes / no) "
    read -r answer
    # Testing lower case version of the answer
    case ${answer,,} in
      yes)
        PTAG=$IMPLICIT_PTAG
        break
        ;;
      no)
        fatal "Please use the -p option to specify the git tag you want"
        ;;
    esac
  done
fi
isVariableExists PTAG

if ! isGitTagExists "$TAG"; then
  echo "git tag $TAG doesn't exist !"
  while true; do
    echo -n "do you want to tag commit $GIT_BRANCH/$GIT_LAST_COMMIT with tag $TAG ? (yes / no) "
    read -r answer
    # Testing lower case version of the answer
    case ${answer,,} in
      yes)
        CREATE_TAG="yes"
        git tag "$TAG" || fatal "Can't tag with tag $TAG"
        git push >/dev/null || fatal "Can't push branch $GIT_BRANCH"
        break
        ;;
      no)
        fatal "Please create git tag $TAG"
        ;;
    esac
  done
else
  # Be sure we build the exact code associate to this TAG
  git checkout -q "$TAG" || fatal "Cannot checkout tag $TAG"
fi

echo "Building binaries for git tag $TAG"
make clean-all
make prepare
rm -rf "$GOPATH"/src/github.com/docker/docker/vendor/github.com/docker/go-connections/nat
make -s release TAG="$TAG" || fatal "Cannot build ceph-nano !"

rm cn || fatal "Cannot remove cn"
ln -sf cn-*-"$LOCAL_ARCH" cn || fatal "Cannot link cn for $LOCAL_ARCH"

sudo make tests || fatal "Tests are not passing ! Cannot release that !"

# If we did checkout the TAG, we need to return to the previous branch
GIT_CURRENT_BRANCH=$(git rev-parse --abbrev-ref HEAD)
if [ "$GIT_CURRENT_BRANCH" != "$GIT_BRANCH" ]; then
  git checkout -q "$GIT_BRANCH" || fatal "Cannot restore $GIT_BRANCH branch"
fi

CHANGELOG=$(mktemp /tmp/changelog.XXXXX)
echo "Building CHANGELOG between $TAG and $PTAG"
echo "CHANGELOG between version $PTAG and $TAG" > "$CHANGELOG"
git log --oneline "$PTAG".."$TAG" --no-decorate >> "$CHANGELOG"

echo "Creating release $TAG"
github-release release --user $GITHUB_USER --repo $repo --tag "${TAG}" -d - < "$CHANGELOG" || fatal "Cannot create release $TAG"

echo "Uploading CHANGELOG"
github-release upload --user $GITHUB_USER --repo $repo --tag "${TAG}" --name CHANGELOG --file "$CHANGELOG" || fatal "Cannot upload CHANGELOG"
rm -f "$CHANGELOG"

echo "Uploading binaries"
for binary in cn*"$TAG"*; do
  echo "- $binary"
  github-release upload --user $GITHUB_USER --repo $repo --tag "${TAG}" --name "$binary" --file "$binary" || fatal "Cannot upload cn"
done

# Everything went well, let's push the tag to the github repo
if [ -n "$CREATE_TAG" ]; then
  git push origin "$TAG" >/dev/null || fatal "Can't push TAG $TAG"
fi

echo "Release can be browsed at https://github.com/$GITHUB_USER/$repo/releases/tag/$TAG"
