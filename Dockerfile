FROM golang:1.12 AS builder

# Variables used to build the binary
# If VERSION is not passed as a build arg, go with dev
# If USER is not passed as a build arg, go with dohernandez
ARG VERSION=dev
ARG USER=dohernandez


WORKDIR /go/src/github.com/dohernandez/swiss-army-knife

COPY . .

RUN echo "Building binary VERSION ${VERSION} with USER ${USER}"
RUN make deps-vendor build


FROM ubuntu:bionic

LABEL quay.expires-after=8w

RUN groupadd -r dohernandez && useradd --no-log-init -r -g dohernandez dohernandez
USER dohernandez

COPY --from=builder --chown=dohernandez:dohernandez \
    /go/src/github.com/dohernandez/swiss-army-knife/bin/swiss-army-knife /bin/swiss-army-knife
