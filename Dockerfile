FROM alpine:latest
ADD ./app ./config.yaml /app/
WORKDIR /app
ENTRYPOINT ["/app/app"]