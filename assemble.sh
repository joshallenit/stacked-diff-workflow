#!/bin/bash
 
trap ctrl_c INT

# Avoid say command on ctrl-c
function ctrl_c() {
    exit 1
}

# Note: lintInternalDebug is omitted as it takes over 10 minutes.

set -x # show executing lines

./gradlew assembleInternalDebug \
    compileReleaseUnitTestSources \
    compileInternalDebugUnitTestSources \
    app:assembleInternalDebugAndroidTest \
    detekt \
    "$@"

retVal=$?

set +x

if [ $retVal -eq 0 ]; then
    say assemble complete &
else
    say not assembled &
fi
exit $retVal