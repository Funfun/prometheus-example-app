FROM quay.io/prometheus/busybox:latest

ADD currently-app /bin/currently-app

ENTRYPOINT ["/bin/currently-app"]
