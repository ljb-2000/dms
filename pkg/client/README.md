# CLI Usage
```
NAME:
   dms - Docker monitoring service

USAGE:
   dms [global options] command [command options] [arguments...]

VERSION:
   0.1.0

COMMANDS:
     stopped   view stopped containers
     launched  view launched containers
     logs      view container logs
     metrics   view container(s) metrics
     help, h   Shows a list of commands or help for one command

GLOBAL OPTIONS:
   -d, --debug             set debug mode
   -a value, --addr value  set daemon address (default: "http://localhost:4222")
   --help, -h              show help
   --version, -v           print the version
```
## Container logs
```
$ dms logs <container_id>
```
## Stopped contaiers
```
$ dms stopped
```
## Launched containers
```
$ dms launched
```
## Container metrics
```
$ dms metrics <container_id> <container_id>
$ dms metrics all
```
