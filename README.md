# i2kit

The simplicity of containers, the confidence of virtual machines.

*i2kit* combines the simplicity of containers to develop your applications with the security and robustness of virtual machines for production environments at scale.

*i2kit* won't force a development environment on you, keep using what works for you.

Once you're to deploy, define your application deployment behavior using a YAML file (à la docker-compose) and *i2kit* will make the rest for you. *i2kit* builds an extremely lightweight virtual machine for each of your containers/pods using linuxkit, which not only brings the workload isolation and security advantages of virtual machines, but also seamlessly plug into proven cloud vender technology. For example, if AWS is the cloud provider of your choice, *i2kit* will:

- Balance traffic to your Container VMs using ELBs.
- Enforce fault tolerance using autoscalability groups.
- Provide service discovery using Route 53 domains.
- Connect your containers using a VPC network.
- Persist your data volumes using EBS.

*i2kit* does not require a central control plane to manage your running applications (à la Docker Swarm or Kubernetes), which not only reduces configuration, maintenance and infrastructure costs, but also eliminates critical runtime dependencies in your applications.

Check our [academic paper](https://github.com/pchico83/i2kit/tree/master/cli/docs/paper.pdf) for full details about the i2kit behavior.

## Getting Started

Now, build the `i2kit` binary:

```
go build -o /usr/local/bin/i2kit
```

Execute commands by running:

```
i2kit deploy -s service.yml -e environment.yml
i2kit destroy -f service.yml -e environment.yml
```

where `service.yml` is the path to your [Service Manifest]((https://github.com/pchico83/i2kit/tree/master/cli/docs/service-taml.md)) and `environment.yml` is the path to your [Environment Manifest]((https://github.com/pchico83/i2kit/tree/master/cli/docs/environment-yml.md)).

In particular, you will need to have a domain owned by AWS and a hosted zone in this domain.

Once you have configure your `environment.yml`, you can deploy the services in the `cli/samples` folder:

- `nginx.yml`  is a simple service running nginx.
- `voting.yml`, `results.yml` and `redis.yml` is the i2kit equivalent of the well known [docker voting sample application](https://github.com/tutumcloud/voting-demo/blob/master/tutum.yml).
