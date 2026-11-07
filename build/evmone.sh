#!/usr/bin/env bash

set -e

if [[ "$OSTYPE" != "linux"* ]]; then
    echo "This script is only currently configured to work on Linux. Please see \"https://github.com/ipsilon/evmone\" documentation for instructions to build in other environments."
    exit 1
fi

mkdir -p build/_workspace/evmone
[[ -f build/_workspace/evmone/evmone-0.5.0-linux-x86_64.tar.gz ]] && exit 0
wget -O build/_workspace/evmone/evmone-0.5.0-linux-x86_64.tar.gz.tmp https://github.com/ipsilon/evmone/releases/download/v0.5.0/evmone-0.5.0-linux-x86_64.tar.gz
# The SHA-256 the release publishes beside this file. The download keeps a temporary
# name until it passes, so a failed check cannot satisfy the existence test above.
echo "784037cb212fd434eb7a13f3c5f87c9fb3d28818d0c54924e775cc09cbabf9cf  build/_workspace/evmone/evmone-0.5.0-linux-x86_64.tar.gz.tmp" | sha256sum -c -
mv build/_workspace/evmone/evmone-0.5.0-linux-x86_64.tar.gz.tmp build/_workspace/evmone/evmone-0.5.0-linux-x86_64.tar.gz
tar xzvf build/_workspace/evmone/evmone-0.5.0-linux-x86_64.tar.gz -C build/_workspace/evmone/
