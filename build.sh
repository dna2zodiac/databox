#!/bin/bash


SELF=$(cd `dirname $0`; pwd)

if [ ! -d $SELF/local/src/github.com/dna2zodiac ]; then
   mkdir -p $SELF/local/src/github.com/dna2zodiac
   pushd $SELF/local/src/github.com/dna2zodiac
   ln -s ../../../.. databox
   popd
fi

pushd $SELF/local
GOPATH=`pwd` go install github.com/dna2zodiac/databox/cmd/databox-webserver
popd
