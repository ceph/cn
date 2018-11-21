#!/bin/bash
set -e

PRE_VERSION=${1:-}
VERSION=${2:-}
RELEASE=${RELEASE:-1}
SPEC_FILE=${SPEC_FILE:-cn.spec}
TODAY=$(LANG=C date "+%a %b %d %Y")

function fatal() {
    echo "$@"
    exit -1
}

if [ ! -f "$SPEC_FILE" ]; then
    fatal "$SPEC_FILE doesn't exist !"
fi

if [ -z "$PRE_VERSION" ]; then
    fatal "PRE_VERSION must be defined !"
fi

if [ -z "$VERSION" ]; then
    fatal "VERSION must be defined !"
fi

RPM_VERSION=$(echo ${VERSION//v/})
TEMP=$(mktemp /tmp/tunespec.XXXXXXXXX)
START=$TEMP.start
CHANGELOG=$TEMP.changelog

CHANGELOG_POS=$(grep -n %changelog "$SPEC_FILE" | cut -d ":" -f 1)

# Extracting the header of the spec file
head -"$CHANGELOG_POS" "$SPEC_FILE" > "$TEMP"

# Replace the global variables
sed -i "s|%global source_version.*|%global source_version $RPM_VERSION|g" "$TEMP"
sed -i "s|%global tag.*|%global tag $RELEASE|g" "$TEMP"

# Adding the new changelog
echo "* $TODAY  Erwan Velu <evelu@redhat.com> - $RPM_VERSION-$RELEASE" >> "$TEMP"

#The commit id is replaced by a dash to keep the spec file spirit
git log --oneline "$PRE_VERSION..$VERSION" --no-decorate | sed -e 's/^[a-z0-9]* \(\s*\)/- \1/g' >> "$TEMP"

# Pasting the previous changelog
tail -n +$((CHANGELOG_POS+1)) "$SPEC_FILE" >> "$TEMP"

# Replacing the specfile with the new one
mv "$TEMP" "$SPEC_FILE"

# Cleaning
rm -f "$START" "$CHANGELOG"
