#!/usr/bin/env bash

printf "starting MongoDB...\n"
docker run --rm -it -d -p "27017:27017" -v "${PWD}"/mongo/data:/data -v "${PWD}"/mongo/etc:/opt/eduid/etc -v "${PWD}"/mongo/db-scripts:/opt/eduid/db-scripts "docker.sunet.se/eduid/mongodb"


printf "starting Redis...\n"
docker run --rm -it -d -p "6379:6379" -v "${PWD}"/redis-data:/data redis
