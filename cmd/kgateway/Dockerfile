ARG ENVOY_IMAGE

# Extract envoy binary and required shared libraries from envoy-gloo image.
# Always use linux/amd64 since envoy-gloo is not multi-arch compatible.
FROM --platform=linux/amd64 $ENVOY_IMAGE AS builder

COPY copy-envoy-libs.sh /tmp/copy-envoy-libs.sh
RUN chmod +x /tmp/copy-envoy-libs.sh && /tmp/copy-envoy-libs.sh

FROM gcr.io/distroless/static:nonroot

ARG GOARCH=amd64

COPY kgateway-linux-$GOARCH /usr/local/bin/kgateway
COPY --from=builder /out/lib/ /lib/
COPY --from=builder /out/lib64/ /lib64/
COPY --from=builder /usr/local/bin/envoy /usr/local/bin/envoy
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

USER 10101

ENTRYPOINT ["/usr/local/bin/kgateway"]
