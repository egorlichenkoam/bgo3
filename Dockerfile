FROM alpine:3

COPY bank /app/

ENTRYPOINT ["/app/bank"]

EXPOSE 9999