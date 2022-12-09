#   Copyright (c) 2021 Samsung.
#
#   Licensed under the Apache License, Version 2.0 (the "License");
#   you may not use this file except in compliance with the License.
#   You may obtain a copy of the License at
#
#       http://www.apache.org/licenses/LICENSE-2.0
#
#   Unless required by applicable law or agreed to in writing, software
#   distributed under the License is distributed on an "AS IS" BASIS,
#   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
#   See the License for the specific language governing permissions and
#   limitations under the License.

#-----------------------------------------------------------

FROM nexus3.o-ran-sc.org:10002/o-ran-sc/bldr-ubuntu20-c-go:1.1.0 AS a1-build


#TODO check why defualt golang is not working
ARG GOVERSION="1.18.5"
RUN wget -nv https://dl.google.com/go/go${GOVERSION}.linux-amd64.tar.gz \
     && tar -xf go${GOVERSION}.linux-amd64.tar.gz \
     && mv go /opt/go/${GOVERSION} \
     && rm -f go*.gz


ENV DEFAULTPATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin
ENV PATH=$DEFAULTPATH:/usr/local/go/bin:/opt/go/${GOVERSION}/bin:/root/go/bin

RUN apt-get update -y && apt-get install -y jq

# Update CA certificates
RUN apt update && apt install --reinstall -y \
  ca-certificates \
  && \
  update-ca-certificates

#Install RMR

ARG RMR_VER=4.8.3
ARG RMR_PKG_URL=https://packagecloud.io/o-ran-sc/release/packages/debian/stretch/

RUN wget -nv --content-disposition ${RMR_PKG_URL}/rmr_${RMR_VER}_amd64.deb/download.deb
RUN wget -nv --content-disposition ${RMR_PKG_URL}/rmr-dev_${RMR_VER}_amd64.deb/download.deb
RUN    dpkg -i rmr_${RMR_VER}_amd64.deb  \
        && dpkg -i rmr-dev_${RMR_VER}_amd64.deb \
        && ldconfig


ENV PATH="/usr/local/go/bin:${PATH}"

ENV GOPATH="/go"

RUN mkdir -p /go/bin
RUN mkdir -p /go/src/ws
WORKDIR "/go/src/ws"

# Module prepare (if go.mod/go.sum updated)
COPY go.mod /go/src/ws
COPY go.sum /go/src/ws
RUN GO111MODULE=on go mod download

# build and test
COPY . /go/src/ws
COPY  ./config/config_test.yaml /opt/a1-mediator/

ENV CFG_FILE=/opt/a1-mediator/config.yaml
ENV A1_CONFIG_FILE=/opt/a1-mediator/config_test.yaml


# Build the code
RUN GO111MODULE=on GO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /go/src/ws/cache/go/cmd/a1 cmd/a1.go

# Run unit tests
RUN GO111MODULE=on GO_ENABLED=0 GOOS=linux go test -p 1 -cover ./pkg/resthooks/
RUN GO111MODULE=on GO_ENABLED=0 GOOS=linux go test -p 1 -cover ./pkg/a1/
RUN GO111MODULE=on GO_ENABLED=0 GOOS=linux go test -p 1 -cover ./pkg/policy

RUN gofmt -l $(find cmd/ pkg/  -name '*.go' -not -name '*_test.go')

CMD ["/bin/bash"]


#----------------------------------------------------------
FROM ubuntu:18.04 as a1-mediator

RUN apt-get update -y \
    && apt-get install --reinstall -y sudo openssl ca-certificates ca-cacert wget\
    && apt-get clean && update-ca-certificates

#Install RMR

ARG RMR_VER=4.8.3
ARG RMR_PKG_URL=https://packagecloud.io/o-ran-sc/release/packages/debian/stretch/

RUN wget -nv --content-disposition ${RMR_PKG_URL}/rmr_${RMR_VER}_amd64.deb/download.deb
RUN wget -nv --content-disposition ${RMR_PKG_URL}/rmr-dev_${RMR_VER}_amd64.deb/download.deb
RUN    dpkg -i rmr_${RMR_VER}_amd64.deb  \
        && dpkg -i rmr-dev_${RMR_VER}_amd64.deb \
        && ldconfig

#
# a1-mediator
#
RUN mkdir -p /opt/a1-mediator \
    && chmod -R 755 /opt/a1-mediator

COPY --from=a1-build /go/src/ws/cache/go/cmd/a1 /opt/a1-mediator/a1

COPY  ./config/config.yaml /opt/a1-mediator/

WORKDIR /opt/a1-mediator

ARG RMR_CONFIG=/opt/a1-mediator/config.yaml
ENV CFG_FILE=$RMR_CONFIG
ARG A1_CONFIG=/opt/a1-mediator/config.yaml
ENV A1_CONFIG_FILE=$A1_CONFIG


COPY a1-entrypoint.sh /opt/a1-mediator/
ENTRYPOINT ["/opt/a1-mediator/a1-entrypoint.sh"]

