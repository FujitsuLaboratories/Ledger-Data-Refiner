#!/bin/bash

# Copyright Fujitsu Corp. All Rights Reserved.
#
# SPDX-License-Identifier: Apache-2.0

echo "Reading variables from JSON file..."

export USER=$(jq .postgreSQL.username ./config.json)
export DATABASE=$(jq .postgreSQL.database ./config.json)
export PASSWORD=$(jq .postgreSQL.passwd ./config.json)

echo "USER=${USER}"
echo "DATABASE=${DATABASE}"
echo "PASSWORD=${PASSWORD}"

echo "Executing SQL scripts, OS="$OSTYPE

case $OSTYPE in
darwin*) psql postgres -v dbname=$DATABASE -v user=$USER -v passwd=$PASSWORD -f ./explorerpg.sql ;
psql postgres -v dbname=$DATABASE -v user=$USER -v passwd=$PASSWORD -f ./updatepg.sql ;;
linux*)
if [ $(id -un) = 'postgres' ]; then
  PSQL="psql"
else
  PSQL="sudo -u postgres psql"
fi;
${PSQL} -v dbname=$DATABASE -v user=$USER -v passwd=$PASSWORD -f ./refinerpg.sql ;;
esac

