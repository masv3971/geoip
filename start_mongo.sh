#!/usr/bin/env bash

docker run --rm -it -d --platform linux/amd64 -p "27017:27017" -v ~/Documents/repon/git/geoip/mongo/data:/data -v ~/Documents/repon/git/geoip/mongo/etc:/opt/eduid/etc -v ~/Documents/repon/git/geoip/mongo/db-scripts:/opt/eduid/db-scripts "docker.sunet.se/eduid/mongodb"

