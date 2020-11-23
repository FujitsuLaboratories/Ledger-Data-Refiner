# Copyright Fujitsu Corp. All Rights Reserved.
#
# SPDX-License-Identifier: Apache-2.0

FROM postgres:10.4-alpine
MAINTAINER fujitsu "ldr@cn.fujitsu.com"

RUN apk update \
 && apk add jq \
 && apk add sudo \
 && rm -rf /var/cache/apk/*

WORKDIR /opt

# Copy files
COPY DOCKER/postgreSQL/db/refinerpg.sql         /opt/refinerpg.sql
COPY DOCKER/postgreSQL/db/createdb.sh           /docker-entrypoint-initdb.d/createdb.sh
COPY DOCKER/postgreSQL/db/config.json           /opt/config.json