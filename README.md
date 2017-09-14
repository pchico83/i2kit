# i2kit
i2kit is an immutable infrastructure (i2) deployment tool. It transforms k8 pods in virtual machines using linuxkit, and uses cloud provider technology to support networking, persistence and service discovery. This prototype is focus on AWS, using VPC for networking, ELBs for service discovery and EBS for persistency. The selling point is to have the goodness of docker for local dev, ci and distribution, and be compatible with k8 specs but eliminating the complexity and the abstraction layer of a cluster management tool.

# Getting Started

Make sure you have `moby` and `linuxkit` tools installed:

```
go get -u github.com/moby/tool/cmd/moby
go get -u github.com/linuxkit/linuxkit/src/cmd/linuxkit
```

Build the `i2kit` binary:

```
go build -o /usr/local/bin/i2kit
```
