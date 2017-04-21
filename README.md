[![Go Report Card](https://goreportcard.com/badge/github.com/lavrs/docker-monitoring-service)](https://goreportcard.com/report/github.com/lavrs/docker-monitoring-service) [![Build Status](https://travis-ci.org/lavrs/docker-monitoring-service.svg?branch=master)](https://travis-ci.org/lavrs/docker-monitoring-service)
# Docker monitoring service
This service allows you to monitor the loading of the docker container in real time with the ability to view load charts
### View load charts
1. Run dms daemon
2. Open http://localhost:4222/charts
### API Usage
REQUEST
```
GET /metrics/:id HTTP/1.1
```
RESPONSE
```
HTTP/1.1 200 OK
Content-Type: application/json
{
  "metrics": [
    {
      "Container": "",
      "Name": "container1",
      "ID": "0ddf7dfdedb61c22a47aa032b069cb51f11c7e95a61f210aab2d419829dab46f",
      "CPUPercentage": 0.0023855158363192976,
      "Memory": 581632,
      "MemoryLimit": 8388608,
      "MemoryPercentage": 6.93359375,
      "NetworkRx": 3410,
      "NetworkTx": 998,
      "BlockRead": 0,
      "BlockWrite": 0,
      "PidsCurrent": 2,
      "IsInvalid": false
    }
  ],
  "launched": [
    "container1"
  ],
  "stopped": [
    "container2"
  ],
  "message": "message"
}
```
### CLI Usage (run daemon)
```
NAME:
   dms - docker monitoring service

USAGE:
   dms [global options] command [command options] [arguments...]

VERSION:
   0.1.0

COMMANDS:
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   -p value, --port value  set daemon port (default: "4222")
   -uclt value, --upd-container-list-time value  set update container list interval (default: 3)
   -uct value, --upd-container-time value  set update container interval (default: 1)
   --help, -h              show help
   --version, -v           print the version
```
