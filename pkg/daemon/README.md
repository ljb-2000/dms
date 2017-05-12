# API Usage
## Get container(s) metrics
REQUEST
```
GET /api/metrics/:id HTTP/1.1
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
## Get container logs
REQUEST
```
GET /api/logs/:id HTTP/1.1
```
RESPONSE
```
HTTP/1.1 200 OK
Content-Type: application/json
{
  "logs": "logs"
}
```
## Get stopped containers
REQUEST
```
GET /api/stopped HTTP/1.1
```
RESPONSE
```
HTTP/1.1 200 OK
Content-Type: application/json
{
  "stopped": ["container"]
}
```
## Get launched containers
REQUEST
```
GET /api/launched HTTP/1.1
```
RESPONSE
```
HTTP/1.1 200 OK
Content-Type: application/json
{
  "launched": ["container"]
}
```