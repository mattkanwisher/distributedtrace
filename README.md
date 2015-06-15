# DistributedTrace 

A GO Zipkin backend, this takes in zipkin protocol and stores it in InfluxDb/Redis. Mysql support coming soon. Extracted from my book http://microservicesingo.com

[![Build Status](https://travis-ci.org/mattkanwisher/distributedtrace.svg)](https://travis-ci.org/mattkanwisher/distributedtrace)
[![GoDoc](https://godoc.org/github.com/mattkanwisher/distributedtrace?status.svg)](https://godoc.org/github.com/mattkanwisher/distributedtrace)

[ZipKin](https://github.com/twitter/zipkin) proxy implementation in Go. Use as a library
or use the CLI subpackage `./zipkin`.

## INSTALL

```sh
$ go get github.com/mattkanwisher/distributedtrace/zipkin
```

## USAGE

Right now the CLI takes configuration parameters from the environment variable.

* `ZIPKIN_ADDR` - The address to listen on. Defaults to `0.0.0.0:9410`
* `ZIPKIN_OUTPUT` - The output to proxy to. `influx|redis|console|null`
* `ZIPKIN_OUTPUT_ADDR` - The output address for the specified output kind.

## EXAMPLE

```sh
$ ZIPKIN_OUTPUT=influx ZIPKIN_OUTPUT_ADDR=http://admin:admin@192.168.59.103:8086/spans
```

Outputs to [InfluxDB](http://influxdb.com) at address `192.168.59.103:8086` with username
`admin` and password `admin` and using the database `spans`.

## CLIENT CONFIGURATION

Uses the TCP scribe-collector interface for clients.

