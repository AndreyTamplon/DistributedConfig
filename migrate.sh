#!/bin/bash
apk add tar
apk add curl
DB_NAME=$1
DB_USER=$2
DB_PASS=$3
DB_HOST=$4
if [ $# -eq 0 ]
  then
    echo "usage: <DB_NAME> <DB_USER> <DB_PASS> <DB_HOST>"
fi
curl -L https://github.com/golang-migrate/migrate/releases/download/v4.14.1/migrate.linux-amd64.tar.gz | tar xvz
./migrate.linux-amd64 -path migrations -database "postgres://${DB_HOST}/${DB_NAME}?sslmode=disable&user=${DB_USER}&password=${DB_PASS}" up

