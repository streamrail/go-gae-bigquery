#!/bin/sh
echo "running goapp get to fetch dependencies..."
./go_appengine/goapp get -d
echo "dependencies fetched."
exit 0