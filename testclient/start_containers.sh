#!/bin/sh

docker pull tutum/influxdb
docker pull grafana/grafana

docker rm -fv influx
docker rm -fv grafana

docker create --name influx -p 8083:8083 -p 8086:8086 --restart=always tutum/influxdb
docker create --name grafana -p 3000:3000 --restart=always grafana/grafana
docker start influx
docker start grafana
boot2docker ip

