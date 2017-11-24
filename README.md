# i2kit
i2kit is an immutable infrastructure (i2) deployment tool. It transforms k8 pods in virtual machines using linuxkit, and uses cloud provider technology to support networking, persistence and service discovery.

The selling point is to have the goodness of docker for local development, ci and distribution of content, but keeping the robustness and performance of classic cloud vendor technologies. i2kit does not require a central service, eliminating the complexity and the abstraction layer of a cluster management tool.

# Implementation Details

The first prototype focuses on AWS, using VPC for networking, ELBs for exposing k8 deployments (aká a set of pods), Route53 CNAMES for k8 services and deployment endpoints and EBS for persistency. In other words, a k8 deployment is transformed into a linuxkit AMI, an auto scalability group with the desired number of instances, and a ELB configured for the ports defined in the k8 deployment. Also, a Route53 CNAME for `deployment-name.i2kit.com` is created that resolves to the deployment ELB.
i2kit also supports the deployment of k8 services by creating a CNAME that resolves to the ELBs of the deployment matching the the k8 service selector. In order to find these ELBs, i2kit uses AWS tags, tagging every k8 deployment with its labels.

# Getting Started

Make sure you have the `linuxkit` tool installed:

```
go get -u github.com/linuxkit/linuxkit/src/cmd/linuxkit
```

and the `aws-cli` configured with your credentials. Now, build the `i2kit` binary:

```
go build -o /usr/local/bin/i2kit
```

`AWS_SHARED_CREDENTIALS_FILE` pointing to AWS file credentials (`AWS_CREDENTIALS`).
`AWS_REGION` environment variable as this is used by the AWS Go SDK.
S3 bucket in AWS named `linuxkit` // TODO: will be configurable soon
