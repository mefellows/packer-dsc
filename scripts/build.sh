#!/bin/bash
#
# This script builds the application from source for multiple platforms.
set -e

# Get the parent directory of where this script is.
SOURCE="${BASH_SOURCE[0]}"
while [ -h "$SOURCE" ] ; do SOURCE="$(readlink "$SOURCE")"; done
DIR="$( cd -P "$( dirname "$SOURCE" )/.." && pwd )"

# Change into that directory
cd $DIR

# Get the git commit
GIT_COMMIT=$(git rev-parse HEAD)
GIT_DIRTY=$(test -n "`git status --porcelain`" && echo "+CHANGES" || true)

# If its dev mode, only build for ourself
if [ "${TF_DEV}x" != "x" ]; then
    XC_OS=${XC_OS:-$(go env GOOS)}
    XC_ARCH="386 amd64"
fi

# Determine the arch/os combos we're building for
XC_ARCH=${XC_ARCH:-"386 amd64 arm"}
XC_OS=${XC_OS:-linux darwin windows freebsd openbsd}

# Delete the old dir
echo "==> Removing old directory..."
rm -f bin/*
rm -rf pkg/*
mkdir -p bin/

# Build!
echo "==> Building..."
set +e
gox \
    -os="${XC_OS}" \
    -arch="${XC_ARCH}" \
    -output "pkg/{{.OS}}_{{.Arch}}/packer-{{.Dir}}" \
    ./plugin/provisioner-dsc
set -e

# Move all the compiled things to the $GOPATH/bin
GOPATH=${GOPATH:-$(go env GOPATH)}
case $(uname) in
    CYGWIN*)
        GOPATH="$(cygpath $GOPATH)"
        ;;
esac
OLDIFS=$IFS
IFS=: MAIN_GOPATH=($GOPATH)
IFS=$OLDIFS

# Copy our OS/Arch to the bin/ directory
echo "==> Copying binaries for this platform..."
DEV_PLATFORM="./pkg/$(go env GOOS)_$(go env GOARCH)"
for F in $(find ${DEV_PLATFORM} -mindepth 1 -maxdepth 1 -type f); do
    cp ${F} bin/
    cp ${F} ${MAIN_GOPATH}/bin/
done

# If its dev mode, update the shasums in parity.rb
if [ "${TF_DEV}x" != "x" ]; then
    # Update parity.rb
    echo "==> Updating ./scripts/parity.rb with latest SHASUMS"
    HASH32=$(shasum -a 1 pkg/darwin_386/packer-provisioner-dsc | cut -d" " -f 1)
    HASH64=$(shasum -a 1 pkg/darwin_amd64/packer-provisioner-dsc | cut -d" " -f 1)
    sed -i "9s/sha1 '\(.*\)'/sha1 '${HASH32}'/g" scripts/packer-provisioner-dsc.rb
    sed -i "12s/sha1 '\(.*\)'/sha1 '${HASH64}'/g" scripts/packer-provisioner-dsc.rb
fi

# Done!
echo
echo "==> Results:"
ls -hl bin/
