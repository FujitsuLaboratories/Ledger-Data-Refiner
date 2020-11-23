# Copyright Fujitsu Corp. All Rights Reserved.
#
# SPDX-License-Identifier: Apache-2.0

FROM golang:1.13.5-alpine3.10 as builder
MAINTAINER fujitsu "ldr@cn.fujitsu.com"
WORKDIR /opt

ENV GOPROXY https://goproxy.io
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -a -o ledgerdata-refiner main.go

FROM alpine:3.10 AS final

WORKDIR /opt
COPY --from=builder /opt/. /opt

EXPOSE 30052

CMD /opt/ledgerdata-refiner