# Docker monitoring service
This service allows you to monitor the loading of the docker container in real time with the ability to view load charts.
## View load charts
```
$ cd cmd/daemon/
$ go run daemon.go
$ http://localhost:4222/dashboard
```
## Documentation
[Daemon](https://github.com/lavrs/dms/tree/master/pkg/daemon/README.md)<br>
[Client](https://github.com/lavrs/dms/tree/master/pkg/client/README.md)

## Load testing with k6
[k6 docs](https://docs.k6.io/docs)
##### Run tests
```
$ k6 run k6lt.js
```
