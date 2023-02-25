#!/bin/sh

baseName=ArtemisLogParser
rm -rf "Build"
mkdir -p "Build"

buildGo() {
  file="Build/$baseName$3"
  GOOS=$1 GOARCH=$2 go build -ldflags="-s -w" -o $file . \
  && echo build $1/$2 as $baseName$3
}

###########

buildGo windows amd64 .exe
buildGo darwin amd64 -osx64
buildGo darwin arm64 -osxArm
buildGo linux amd64 -linux64
buildGo linux arm -linuxArm

###########

upx --best --lzma Build/*

for f in Build/$baseName-*
do
  echo $f
  mv $f $baseName
  7z a "$f.zip" $baseName
  rm $baseName
done