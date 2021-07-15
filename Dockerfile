FROM golang AS builder

RUN mkdir -p /workdir
WORKDIR /workdir
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN make go-build

####
FROM registry.access.redhat.com/ubi8/ubi-minimal:latest

COPY --from=builder /workdir/_output/* /usr/local/bin/

ENTRYPOINT ["/usr/local/bin/book2k8sevent"]

