#!/bin/bash

echo "ELK : "
grep bulk nohup.out | grep elk | awk '{ total+=$NF;}END{print total;}'

echo "InfluxDB / Kapacitor : "
grep bulk nohup.out | grep influx | awk '{ total+=$NF;}END{print total;}'

