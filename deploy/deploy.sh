#!/bin/bash
PRO_NAME=$1
goos=$2
goarch=$3
cgo=$4

rm -rf "$PRO_NAME"
rm -rf "$PRO_NAME".tar.gz

echo "${PRO_NAME}"

mkdir "$PRO_NAME"
mkdir "$PRO_NAME/config"
cp -r ../config/ "$PRO_NAME"
cp "run.sh" "$PRO_NAME"

if [ "$goos" == "linux" ]; then
  cd .. && CGO_ENABLED=$cgo GOOS=$goos GOARCH=$goarch go build -o "$PRO_NAME"
  cp "$PRO_NAME" deploy/"$PRO_NAME"/
  rm "$PRO_NAME"
fi
if [ "$goos" == "windows" ]; then
#  windows 打包后缀拼接 exe
  cd .. && CGO_ENABLED=$cgo GOOS=$goos GOARCH=$goarch go build -o "$PRO_NAME".exe
  cp "$PRO_NAME".exe deploy/"$PRO_NAME"/
  rm "$PRO_NAME".exe
fi
cd deploy && tar zcvf "$PRO_NAME".tar.gz "$PRO_NAME" && rm -rf "$PRO_NAME"
echo "BUILD SUCCESS"
