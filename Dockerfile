FROM alpine:latest
ADD ./library ./config.yaml /app/
WORKDIR /app
ENTRYPOINT ["/app/library"]