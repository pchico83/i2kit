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

Check our [paper](https://github.com/pchico83/i2kit/tree/master/cli/docs/paper.pdf) for full details about the i2kit behavior.

## Getting Started

Now, build the `i2kit` binary:

```
go build -o /usr/local/bin/i2kit
```

Create an `environment.yml` configuration as documented [here](https://github.com/pchico83/i2kit/tree/master/cli/docs/environment-yml.md).

In particular, you will need to have a [hosted zone](https://docs.aws.amazon.com/Route53/latest/DeveloperGuide/CreatingHostedZone.html) configured in Route53.

Once you have your environment ready, execute the following command:

```
i2kit deploy -s samples/nginx/service.yml -e environment.yml
```

After the command finishes, you can browse to *nginx.your_hosted_zone_here* to verify that nginx is running as expected.

Finally, destroy your service by executing the following command:

```
i2kit destroy -s samples/nginx/service.yml -e environment.yml
```

[Service Manifest](https://github.com/pchico83/i2kit/tree/master/cli/docs/service-yml.md) gives more information about how to create your own Service Manifests.

[Environment Manifest](https://github.com/pchico83/i2kit/tree/master/cli/docs/environment-yml.md) gives more information about how to create your own Environment Manifests.
