ARG ENVOY_IMAGE

# Stage 1: Extract envoy binary and required shared libraries from envoy-gloo image
FROM --platform=$TARGETPLATFORM $ENVOY_IMAGE AS builder

# Copy the script that extracts platform-specific shared libraries
COPY copy-envoy-libs.sh /tmp/copy-envoy-libs.sh
RUN chmod +x /tmp/copy-envoy-libs.sh && /tmp/copy-envoy-libs.sh

# Stage 2: Build final image with minimal base
FROM cgr.dev/chainguard/static

ARG GOARCH=amd64

# Copy kgateway binary
COPY kgateway-linux-$GOARCH /usr/local/bin/kgateway

# Copy envoy binary and its dependencies from builder stage
COPY --from=builder /out/ /
COPY --from=builder /usr/local/bin/envoy /usr/local/bin/envoy
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

USER 10101

ENTRYPOINT ["/usr/local/bin/kgateway"]
