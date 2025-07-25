FROM golang:1.24-alpine AS builder

WORKDIR /app/iksctl

COPY ./ ./

RUN go mod download

RUN go build -o ./bin/iksctl ./cmd/iksctl

FROM golang:1.24-alpine

RUN apk --no-cache add ca-certificates
WORKDIR /app

RUN wget --no-check-certificate -O /usr/local/share/ca-certificates/SystematicRootCA.crt https://nexus3.systematicgroup.local/repository/ITM_raw/certificates/SystematicRootCA.crt && \
    update-ca-certificates

COPY --from=builder /app/iksctl/bin/iksctl .

ENTRYPOINT ["./iksctl"]