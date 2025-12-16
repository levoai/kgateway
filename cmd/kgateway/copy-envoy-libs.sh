#!/bin/sh
set -eu

if [ "${TARGETPLATFORM}" = "linux/amd64" ] ; then
  mkdir -p /out/lib/x86_64-linux-gnu /out/lib64

  # Required shared libraries (ldd output from envoy binary on amd64 machine)
  cp /lib/x86_64-linux-gnu/libm.so.6 /out/lib/x86_64-linux-gnu/libm.so.6
  cp /lib/x86_64-linux-gnu/librt.so.1 /out/lib/x86_64-linux-gnu/librt.so.1
  cp /lib/x86_64-linux-gnu/libdl.so.2 /out/lib/x86_64-linux-gnu/libdl.so.2
  cp /lib/x86_64-linux-gnu/libpthread.so.0 /out/lib/x86_64-linux-gnu/libpthread.so.0
  cp /lib/x86_64-linux-gnu/libc.so.6 /out/lib/x86_64-linux-gnu/libc.so.6
  cp /lib64/ld-linux-x86-64.so.2 /out/lib64/ld-linux-x86-64.so.2
elif [ "${TARGETPLATFORM}" = "linux/arm64" ] ; then
  mkdir -p /out/lib/aarch64-linux-gnu
  
  # Required shared libraries (ldd output from envoy binary on arm64 machine)
  cp /lib/aarch64-linux-gnu/libm.so.6 /out/lib/aarch64-linux-gnu/libm.so.6
  cp /lib/aarch64-linux-gnu/librt.so.1 /out/lib/aarch64-linux-gnu/librt.so.1
  cp /lib/aarch64-linux-gnu/libdl.so.2 /out/lib/aarch64-linux-gnu/libdl.so.2
  cp /lib/aarch64-linux-gnu/libpthread.so.0 /out/lib/aarch64-linux-gnu/libpthread.so.0
  cp /lib/aarch64-linux-gnu/libc.so.6 /out/lib/aarch64-linux-gnu/libc.so.6
  cp /lib/ld-linux-aarch64.so.1 /out/lib/ld-linux-aarch64.so.1
else
  echo "Unsupported TARGETPLATFORM: ${TARGETPLATFORM:-}"
  exit 1
fi
