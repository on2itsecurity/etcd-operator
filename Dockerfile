ARG alpinever=3.17
FROM golang:1.19-alpine$alpinever AS build-base
# Install SSL ca certificates.
# Ca-certificates is required to call HTTPS endpoints.
RUN apk update && apk add --no-cache ca-certificates git gcc musl-dev
WORKDIR /go/src/github.com/on2itsecurity/etcd-operator
COPY cmd cmd
COPY pkg pkg
COPY version version
COPY go.* .
RUN go mod vendor

FROM build-base AS release-builder
ARG REVISION=dev

ENV CGO_ENABLED=0
ENV GOOS=linux
RUN mkdir -p /rootfs/usr/local/bin
RUN mkdir -m 1777 /rootfs/tmp
RUN go build --ldflags "-w -s -X 'github.com/on2itsecurity/etcd-operator/version.GitSHA=$REVISION'" -o /rootfs/usr/local/bin/etcd-operator github.com/on2itsecurity/etcd-operator/cmd/operator
RUN go build --ldflags "-w -s -X 'github.com/on2itsecurity/etcd-operator/version.GitSHA=$REVISION'" -o /rootfs/usr/local/bin/etcd-backup-operator github.com/on2itsecurity/etcd-operator/cmd/backup-operator
RUN go build --ldflags "-w -s -X 'github.com/on2itsecurity/etcd-operator/version.GitSHA=$REVISION'" -o /rootfs/usr/local/bin/etcd-restore-operator github.com/on2itsecurity/etcd-operator/cmd/restore-operator
# ldd will sort out all need libraries, we output only the library path, create directories in /rootfs, and copy the libraries to /rootfs
# enable when CGO_ENABLED=1
#RUN ldd /rootfs/usr/local/bin/*-operator | grep "=> /" | awk '{print $3}' | xargs -i sh -c 'mkdir -p $(dirname "/rootfs{}"); cp -a "{}" "/rootfs{}"'

FROM alpine:$alpinever AS env-builder
ENV USER=etcd-operator
ENV UID=1000
RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/srv/etcd-operator" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \
    "${USER}"


FROM build-base AS env-test
ARG KUBERNETES=v1.26.0

ADD https://storage.googleapis.com/kubernetes-release/release/$KUBERNETES/bin/linux/amd64/kubectl /bin/
RUN chmod +x /bin/kubectl

COPY test test
RUN go get -v -t -d ./test/...
RUN go mod vendor
RUN go test ./test/e2e/ -c -o /bin/etcd-operator-e2e --race
RUN go test ./test/e2e/e2eslow -c -o /bin/etcd-operator-e2eslow --race
RUN go test ./test/e2e/upgradetest/  -c -o /bin/etcd-operator-upgradetest --race

FROM alpine:$alpinever as test-e2e
RUN apk add --no-cache bash
COPY hack hack
COPY --from=env-test /bin/etcd-operator-* /bin
COPY --from=env-test /bin/kubectl /bin

FROM build-base as go-test
RUN go test github.com/on2itsecurity/etcd-operator/pkg/...

FROM scratch
# Setup environment with certificates and user
COPY --from=build-base /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=env-builder /etc/passwd /etc/group /etc/
# Copy libraries and compiled binaries
COPY --from=release-builder /rootfs /

USER etcd-operator:etcd-operator

ENTRYPOINT ["operator"]
