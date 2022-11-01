FROM golang:1.19-bullseye AS builder

COPY go.mod go.sum main.go /build/
RUN cd /build && go build .

FROM busybox

COPY --from=builder /build/kustomize-sops /sops-decrypt
