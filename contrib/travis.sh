#!/bin/bash
set -ex


#############
# VARIABLES #
#############
PENULTIMATE_TAG=$(git describe --abbrev=0 --tags "$(git rev-list --tags --skip=1 --max-count=1)") # this is n-1 tag


#############
# FUNCTIONS #
#############
function edit_readme {
    #  we replace the n-1 tag with the last one
    sed -i "s/$PENULTIMATE_TAG/$TRAVIS_TAG/g" README.md
}

function setup_git {
    git config --global user.email "buils@travis-ci.com"
    git config --global user.name "Travis CI"
}

function commit_and_push {
    git commit -s -m "$@"
    git pull origin master --rebase
    git push https://"$GITHUB_TOKEN"@github.com/ceph/cn master
}

function commit_spec_file {
    pushd contrib
        ./tune-spec.sh "$PENULTIMATE_TAG" "$TRAVIS_TAG"
        git add cn.spec
    popd 2>/dev/null
    commit_and_push "Packaging: Update specfile version to $TRAVIS_TAG"
}

function commit_changed_readme {
    git add README.md
    commit_and_push "Readme: Bump the new release tag: $TRAVIS_TAG"
}

function compile_cn {
    make prepare
    rm -rf "$GOPATH"/src/github.com/docker/docker/vendor/github.com/docker/go-connections/nat
    make
    sudo ./cn version
}

function test_cn {
    sudo make tests DEBUG=1
}


########
# MAIN #
########
if [[ "$1" == "compile-run-cn" ]]; then
    compile_cn
    test_cn
fi

if [[ "$1" == "tag-release" ]]; then
    if [ -n "$TRAVIS_TAG" ]; then
        echo "I'm running on tag $TRAVIS_TAG, let's build a new release!"
        ./contrib/release.sh -g "$GITHUB_TOKEN" -t "$TRAVIS_TAG" -p "$PENULTIMATE_TAG" -b "master"
        git checkout master
        edit_readme
        setup_git
        commit_spec_file
        commit_changed_readme
    else
        echo "Not running on a tag, nothing to do!"
    fi
fi