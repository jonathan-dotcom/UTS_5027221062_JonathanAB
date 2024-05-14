FROM envoyproxy/envoy:dev-4b0495bb6bf09f97292a8b30b2e97b71cda59256
COPY envoy.yaml /etc/envoy/envoy.yaml
RUN chmod go+r /etc/envoy/envoy.yaml