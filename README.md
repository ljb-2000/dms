# Docker container monitoring service
This service allows you to monitor the loading of the docker container in real time
## API Documentation
REQUEST
```
GET /stats/:id
```
RESPONSE
```
200 OK
Content-Type: application/json
{
  "data": [
    {
      "Container": "",
      "Name": "ss8",
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
## CLI Documentation
```
NAME:
   dms - docker container monitoring service

USAGE:
   dms [global options] command [command options] [arguments...]

VERSION:
   0.1.0

COMMANDS:
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   -p value, --port value  set daemon port (default: "8080")
   --help, -h              show help
   --version, -v           print the version
```
