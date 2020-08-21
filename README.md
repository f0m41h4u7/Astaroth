# Astaroth

[![Build Status](https://travis-ci.com/f0m41h4u7/Astaroth.svg?token=qkqdG1nMjn7NW6KwV5QR&branch=master)](https://travis-ci.com/f0m41h4u7/Astaroth)
[![Go Report Card](https://goreportcard.com/badge/github.com/f0m41h4u7/Astaroth)](https://goreportcard.com/report/github.com/f0m41h4u7/Astaroth)

##### System monitoring daemon

Collects system metrics and streams averaged data via GRPC.

### Install

Download and install one of the packages available:

```shell
$ rpm -i astaroth-*.src.rpm
or
$ apt install ./astaroth-*.deb
```
Or build from source:

```shell
$ git clone https://github.com/f0m41h4u7/Astaroth
$ cd Astaroth
$ make
```
Start server:

```shell
$ astaroth -port <port to serve> -config /path/to/config.json
```

Config file is used to define which metrics are enabled ([example config](configs/config.json)).

### Client

Build simple client (prints received metrics to stdout):
```shell
$ make client
$ ./client
```
A client subscribes to server, defining an interval to receive metrics and an interval to average server.

### Currently supported metrics

| Metric name   | Linux             | Windows            |
| ------------- | ------            | -------            |
| CPU usage     | :heavy_plus_sign: | :heavy_plus_sign:  |
| Load average  | :heavy_plus_sign: | :heavy_minus_sign: |
| Disk data     | :heavy_plus_sign: | :heavy_minus_sign: |
| Network stats | :heavy_plus_sign: | :heavy_plus_sign:  |
| Top talkers   | :heavy_plus_sign: | :heavy_minus_sign: |
