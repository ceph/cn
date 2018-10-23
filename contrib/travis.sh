#!/bin/bash
set -ex


#############
# VARIABLES #
#############
PENULTIMATE_TAG=$(git describe --abbrev=0 --tags "$(git rev-list --tags --skip=1 --max-count=1)") # this is n-1 tag
LAST_COMMIT_SHORT_SHA1=$(git log --pretty=format:'%h' -n 1)


#############
# FUNCTIONS #
#############
function edit_readme {
    #  we replace the n-1 tag with the last one
    sed -i "s/$PENULTIMATE_TAG/$TRAVIS_TAG/g" README.md

    # we replace the curl line with the new sha1
    sed -i "s|\\(^curl.*\\)-[0-9a-f]\\{5,40\\}-\\(.*-a[mr][dm]64.*\\)|\\1-$LAST_COMMIT_SHORT_SHA1-\\2|g" README.md
}

function commit_changed_readme {
    git config --global user.email "buils@travis-ci.com"
    git config --global user.name "Travis CI"
    git add README.md
    git commit -s -m "Bump README with the new release tag: $TRAVIS_TAG"
    git push https://"$GITHUB_TOKEN"@github.com/ceph/cn master
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
        ./contrib/release.sh -g "$GITHUB_TOKEN" -t "$TRAVIS_TAG" -p "$PENULTIMATE_TAG"
        git checkout master
        edit_readme
        commit_changed_readme
    else
        echo "Not running on a tag, nothing to do!"
    fi
fi