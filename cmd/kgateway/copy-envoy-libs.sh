#!/bin/sh
set -eu

# Note: Since envoy-gloo is only available for amd64, we always copy amd64 libraries
# The envoy binary is only used for xDS validation in the controller, and validation
# only happens on the controller's native arch (typically amd64 in production).
# The controller itself is built for $TARGETPLATFORM, but envoy binary is always amd64.

mkdir -p /out/lib/x86_64-linux-gnu /out/lib64

# Required shared libraries (ldd output from envoy binary on amd64 machine)
# Using explicit checks to provide clear error messages if libraries are missing
for lib in libm.so.6 librt.so.1 libdl.so.2 libpthread.so.0 libc.so.6; do
  if [ ! -f "/lib/x86_64-linux-gnu/$lib" ]; then
    echo "ERROR: Required library /lib/x86_64-linux-gnu/$lib not found in envoy-gloo image" >&2
    exit 1
  fi
  cp "/lib/x86_64-linux-gnu/$lib" "/out/lib/x86_64-linux-gnu/$lib"
done

if [ ! -f "/lib64/ld-linux-x86-64.so.2" ]; then
  echo "ERROR: Required library /lib64/ld-linux-x86-64.so.2 not found in envoy-gloo image" >&2
  exit 1
fi
cp /lib64/ld-linux-x86-64.so.2 /out/lib64/ld-linux-x86-64.so.2
