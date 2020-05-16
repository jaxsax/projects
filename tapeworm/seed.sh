#!/usr/bin/env bash

docker-compose exec db psql -U postgres -c 'create database tapeworm_bot'

docker cp ./botv2/schema.sql tapeworm_db_1:/tmp/schema.sql
docker-compose exec db bash -c 'psql -U postgres tapeworm_bot < /tmp/schema.sql'
docker-compose exec db psql -U postgres tapeworm_bot -c '\dt'
