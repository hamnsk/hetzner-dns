FROM golang:1.21.1-alpine3.18 AS certs

FROM scratch
WORKDIR /app
COPY hetzner-dns /app/hetzner-dns
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

CMD ["/app/hetzner-dns"]
