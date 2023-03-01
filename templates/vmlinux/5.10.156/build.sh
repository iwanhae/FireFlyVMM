#!/bin/bash

BUILD_PATH=linux-5.10.156
GZ_FILE=${BUILD_PATH}.tar.gz

CONFIG=kernel.config

curl -L https://cdn.kernel.org/pub/linux/kernel/v5.x/linux-5.10.156.tar.gz -o $GZ_FILE
tar -xvf $GZ_FILE
cp $CONFIG $BUILD_PATH/.config
make -C $BUILD_PATH vmlinux -j16
cp $BUILD_PATH/vmlinux ./vmlinux.bin