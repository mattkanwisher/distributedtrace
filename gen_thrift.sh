#!/bin/sh

mkdir -p gen/

echo "====> DOWNLOADING LATEST ZIPKIN THRIFTs"
cd gen
rm *.thrift
wget -nv \
  "https://raw.githubusercontent.com/twitter/zipkin/master/zipkin-thrift/src/main/thrift/com/twitter/zipkin/scribe.thrift" \
  "https://raw.githubusercontent.com/twitter/zipkin/master/zipkin-thrift/src/main/thrift/com/twitter/zipkin/zipkinCollector.thrift" \
  "https://raw.githubusercontent.com/twitter/zipkin/master/zipkin-thrift/src/main/thrift/com/twitter/zipkin/zipkinDependencies.thrift" \
  "https://raw.githubusercontent.com/twitter/zipkin/master/zipkin-thrift/src/main/thrift/com/twitter/zipkin/zipkinCore.thrift" \
  "https://raw.githubusercontent.com/twitter/zipkin/master/zipkin-thrift/src/main/thrift/com/twitter/zipkin/zipkinQuery.thrift"
cd ..

echo
echo "====> GENERATING GO BINDINGS"
thrift -v -r -out gen -gen go:package_prefix=github.com/mattkanwisher/distributedtrace/gen/,package_name=gen gen/zipkinCollector.thrift
thrift -v -r -out gen -gen go:package_prefix=github.com/mattkanwisher/distributedtrace/gen/,package_name=gen gen/zipkinCore.thrift

echo
echo "====> REMOVING UNUSED PACKAGES"
rm -rv gen/**/*-remote

echo
echo "====> BUILDING"
cd gen
go get -v ./...
go build -v ./...
go fmt ./...
goimports -w=true .
cd ..

