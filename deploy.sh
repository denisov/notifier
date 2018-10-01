set -x

go test ./... && \
    now --public && \
    now alias