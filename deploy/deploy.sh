#!/bin/bash
# windows /bin/bash, mac /bin/sh
PRO_NAME=$1

rm -rf "$PRO_NAME"
rm -rf "$PRO_NAME".tar.gz

echo "${PRO_NAME}"

mkdir "$PRO_NAME"
mkdir "$PRO_NAME/config"
cp ../config/prod.ini "$PRO_NAME"/config/
cp "run.sh" "$PRO_NAME"

cd .. && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o "$PRO_NAME"
cp "$PRO_NAME" deploy/"$PRO_NAME"/
rm "$PRO_NAME"
cd deploy && tar zcvf "$PRO_NAME".tar.gz "$PRO_NAME" && rm -rf "$PRO_NAME"
