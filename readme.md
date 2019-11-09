# Pump2 
Pump2 is a gRPC-based tool for building ML Docker images. 
Currently supports custom image builds for both tensorflow and pytorch frameworks.

![Golang](https://img.shields.io/badge/Language%20-go-green)
[![Go Report Card](https://goreportcard.com/badge/github.com/NJUPT-ISL/Pump2)](https://goreportcard.com/report/github.com/NJUPT-ISL/Pump2)

### Architecture
![Arch](https://github.com/NJUPT-ISL/Pump2/blob/master/img/Pump2.jpg)

### Development Environment

- go 1.13
- gRPC 1.24.0
- docker 19.05
- cobra 0.0.5
- gopsutil 2.19.10

### Feature
- Setting ML framework and version
- Packaging third-party python dependent libraries
- Support for TLS authentication encryption
- Builder clustering operation
- Builder health monitoring and operational status monitoring
- Support for building images using GPU

### Get Started
Make sure you have docker running on your host before doing the following.

- Set the Config YAML

You could create a YAML file called `pump2.yaml` in `$HOME`.
You can generate this file by executing the following command:
```shell script
pump2 gen
```
Get more usage with the `-h` option.
The contents of the file may look like this.
```yaml
pump2:
  serverip: 0.0.0.0
  serverport: 5020
  tls:
    tlskey: rpc
    tlscrt: rpc
```
Then modify the parameters you run according to your own needs.
- Run the Pump2 server

You can run the Pump2 server by executing the following command:
```shell script
pump2 run -f $HOME/pump2.yaml
```

