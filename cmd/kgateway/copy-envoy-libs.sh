#!/bin/sh
set -eu

# Note: Since envoy-gloo is only available for amd64, we always copy amd64 libraries
# The envoy binary is only used for xDS validation in the controller, and validation
# only happens on the controller's native arch (typically amd64 in production).
# The controller itself is built for $TARGETPLATFORM, but envoy binary is always amd64.

mkdir -p /out/lib/x86_64-linux-gnu /out/lib64

# Required shared libraries (ldd output from envoy binary on amd64 machine)
cp /lib/x86_64-linux-gnu/libm.so.6 /out/lib/x86_64-linux-gnu/libm.so.6
cp /lib/x86_64-linux-gnu/librt.so.1 /out/lib/x86_64-linux-gnu/librt.so.1
cp /lib/x86_64-linux-gnu/libdl.so.2 /out/lib/x86_64-linux-gnu/libdl.so.2
cp /lib/x86_64-linux-gnu/libpthread.so.0 /out/lib/x86_64-linux-gnu/libpthread.so.0
cp /lib/x86_64-linux-gnu/libc.so.6 /out/lib/x86_64-linux-gnu/libc.so.6
cp /lib64/ld-linux-x86-64.so.2 /out/lib64/ld-linux-x86-64.so.2
