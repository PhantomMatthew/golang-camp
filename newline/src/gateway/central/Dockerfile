FROM alpine
WORKDIR /app
COPY src/gateway/central/micro-gateway /app/micro-gateway
RUN chmod +x /app/micro-gateway
ENTRYPOINT ["/app/micro-gateway", "--registry=consul", "--registry_address=:8500", "web"]

