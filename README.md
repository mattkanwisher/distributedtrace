# DistributedTrace 

A GO Zipkin backend, this takes in zipkin protocol and stores it in InfluxDb/Redis. Mysql support coming soon. Extracted from my book http://microservicesingo.com

[![GoDoc](https://godoc.org/github.com/mattkanwisher/distributedtrace?status.svg)](https://godoc.org/github.com/mattkanwisher/distributedtrace)

[ZipKin](https://github.com/twitter/zipkin) proxy implementation in Go. Use as a library
or use the CLI subpackage `./zipkin`.

## INSTALL

```sh
$ go get github.com/mattkanwisher/distributedtrace/zipkin
```

## USAGE

Right now the CLI takes configuration parameters on command line only, later on we will support environment variables.

* Influx/Redis/Mysql as the ouput type, currently only Influx supported
* Url to output to
* -l "0.0.0.0:9410" # the bind address

## EXAMPLE

```sh
$ ./bin/zipkin influx  http://admin:admin@192.168.59.103:8086/spans  -l "0.0.0.0:9410"
```

Outputs to [InfluxDB](http://influxdb.com) at address `192.168.59.103:8086` with username
`admin` and password `admin` and using the database `spans`.

## CLIENT CONFIGURATION

Uses the TCP scribe-collector interface for clients.

