FROM gcr.io/distroless/static:nonroot

WORKDIR /

COPY /build/self-signed-crt-operator .

ENTRYPOINT ["/self-signed-crt-operator"]
