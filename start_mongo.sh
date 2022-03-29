#!/usr/bin/env bash

docker run --rm -it -d --platform linux/amd64 -p "27017:27017" -v "${PWD}"/mongo/data:/data -v "${PWD}"/mongo/etc:/opt/eduid/etc -v "${PWD}"/mongo/db-scripts:/opt/eduid/db-scripts "docker.sunet.se/eduid/mongodb"