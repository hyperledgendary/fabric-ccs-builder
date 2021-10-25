# Copyright IBM Corp. All Rights Reserved.
#
# SPDX-License-Identifier: Apache-2.0
ARG GO_VER=1.14.4
ARG ALPINE_VER=3.12

FROM golang:${GO_VER}-alpine${ALPINE_VER} as golang

WORKDIR /go/src/github.com/hyperledgendary/fabric-ccs-builder

#ENV GOPROXY=https://goproxy.cn,direct
COPY . .

RUN go install -v ./cmd/...

FROM alpine as fabric-ccs-builder

RUN mkdir -p /go/bin
WORKDIR /go/bin
COPY --from=golang /go/bin/build /go/bin/build
COPY --from=golang /go/bin/detect /go/bin/detect
COPY --from=golang /go/bin/release /go/bin/release

