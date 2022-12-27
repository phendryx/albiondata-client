#!/usr/bin/env bash

set -eo pipefail

# apt-get update && sudo apt-get install -y libpcap-dev zip

sudo dpkg --add-architecture arm64                      \
 && sudo dpkg --add-architecture armel                      \
 && sudo dpkg --add-architecture armhf                      \
 && sudo dpkg --add-architecture i386                       \
 && sudo dpkg --add-architecture mips                       \
 && sudo dpkg --add-architecture mipsel                     \
 && sudo dpkg --add-architecture powerpc                    \
 && sudo dpkg --add-architecture ppc64el                    \
 && sudo apt-get update                                     \
 && sudo apt-get install -y -q                              \
        autoconf                                       \
        automake                                       \
        autotools-dev                                  \
        bc                                             \
        binfmt-support                                 \
        binutils-multiarch                             \
        binutils-multiarch-dev                         \
        build-essential                                \
        clang                                          \
        crossbuild-essential-arm64                     \
        crossbuild-essential-armel                     \
        crossbuild-essential-armhf                     \
        crossbuild-essential-mipsel                    \
        crossbuild-essential-ppc64el                   \
        curl                                           \
        devscripts                                     \
        gdb                                            \
        git-core                                       \
        libtool                                        \
        llvm                                           \
        mercurial                                      \
        multistrap                                     \
        patch                                          \
        software-properties-common                     \
        subversion                                     \
        wget                                           \
        xz-utils                                       \
        cmake                                          \
        qemu-user-static                               \
        libxml2-dev                                    \
        lzma-dev                                       \
        openssl                                        \
        libssl-dev                                     \

sudo apt-get update && sudo apt-get install -y --no-install-recommends \
		g++ \
		gcc \
		libc6-dev \
		make \
		pkg-config \
        ca-certificates \
        wget \  
        git \      
	build-essential \
	mingw-w64 \
	nsis

export OSXCROSS_NO_INCLUDE_PATH_WARNINGS=1
export MACOSX_DEPLOYMENT_TARGET=10.6
export CC=/usr/osxcross/bin/o64-clang
export CXX=/usr/osxcross/bin/o64-clang++
export GOOS=darwin
export GOARCH=amd64 CGO_ENABLED=1
go build -ldflags "-s -w -X main.version=$GITHUB_REF_NAME" albiondata-client.go


gzip -k9 albiondata-client
mv albiondata-client.gz update-darwin-amd64.gz


# Creates a zipped folder with a run.command file that runs the client under sudo
TEMP="albiondata-client"
ZIPNAME="albiondata-client-amd64-mac.zip"
rm -rfv ./scripts/$TEMP
rm -rfv ./$ZIPNAME
rm -rfv ./scripts/update-darwin-amd64.zip
mkdir -v ./scripts/$TEMP
cp -v albiondata-client ./scripts/$TEMP/albiondata-client-executable
cd scripts
cp -v run.command ./$TEMP/run.command
chown -Rv ${USER}:${USER} ./$TEMP
chmod -v 777 ./$TEMP/*
zip -v ../$ZIPNAME -r ./"$TEMP"
ls -la

# In theory the following works to create an app but there was a permissions issue when opening on the mac
# APP_NAME="Albion Data Client"
# TEMP="$APP_NAME".app
# ZIPNAME="albiondata-client-amd64-mac.zip"

# rm -rfv ./scripts/"$TEMP"
# rm -rfv ./scripts/"$ZIPNAME"
# mkdir -pv ./scripts/"$TEMP"/Contents/MacOS
# cp -v albiondata-client-darwin-10.6-amd64 ./scripts/"$TEMP"/Contents/MacOS/"$APP_NAME"
# chown -Rv ${USER}:${USER} ./scripts/"$TEMP"
# chmod -v 777 ./scripts/"$TEMP"/*

# cd scripts
# zip -v ../$ZIPNAME -r ./"$TEMP"
