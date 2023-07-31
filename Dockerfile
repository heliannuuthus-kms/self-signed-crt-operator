# Build the manager binary
FROM golang:1.20 as builder
ARG TARGETOS
ARG TARGETARCH

WORKDIR /workspace

COPY . .

RUN go env -w GOPROXY=https://goproxy.cn,direct &&  \
    go mod download && \
    CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH} go build -a -o manager cmd/main.go

FROM gcr.io/distroless/static:nonroot
WORKDIR /
COPY --from=builder /workspace/manager .
USER 65532:65532

ENTRYPOINT ["/manager"]
