#!/bin/bash
set -ex


#############
# FUNCTIONS #
#############
function edit_readme {
    sed -i "s/v[0-9].[0-9].[0-9]/v$TRAVIS_TAG/g" ../README.md
}

function commit_change_readme_release {
    git config --global user.email "seb@redhat.com"
    git config --global user.name "Sébastien Han"
    git add README.md
    git commit -s -m "Bump README with the new release tag: $TRAVIS_TAG"
    git push --quiet https://$GITHUB_TOKEN@github.com/ceph/cn master
}

function compile_cn {
    make prepare
    mv "$GOPATH"/src/github.com/docker/docker/vendor/github.com/docker/go-connections/nat "$GOPATH"/src/github.com/docker/docker/vendor/github.com/docker/go-connections/nonat
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
        # ./release.sh -g "$GITHUB_TOKEN" -t "$TRAVIS_TAG"
        # edit_readme
        # commit_change_readme_release
    else
        echo "Not running on a tag, nothing to do!"
    fi
fi