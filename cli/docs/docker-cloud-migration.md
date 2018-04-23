# Migrate from docker cloud

Cluster and application management services in Docker Cloud are [shutting down on May 21](https://docs.docker.com/docker-cloud/migration/). If you're running your applications there, it is time to start migrating your deployments.  

# Why i2kit?
i2kit is really simple to use. Write a manifest, configure an AWS account, run `i2kit` in the command line and your applications will be running in AWS in minutes. No annoying control plane to setup and manage.

# Migrate to i2kit

At a high level, you would need to:
1. Prepare your AWS account
1. Convert your docker cloud yaml files to i2kit yaml files
1. Deploy your services using `i2kit`
1. Point your application CNAMES to new service endpoints.

# Voting-app example
We are going to use [docker's voting app](https://github.com/dockersamples/example-voting-app) for this example.
In the [`dockercloud.yml`](https://github.com/dockersamples/example-voting-app/blob/master/dockercloud.yml), the voting app is defined as a stack of six microservices:

- vote: Web front-end that displays voting options
- redis: In-memory k/v store that collects votes
- worker: Stores votes in database
- db: Persistent store for votes
- result: Web server that pulls and displays results from database
- lb: Container-based load balancer

## Prepare your AWS account
i2kit requires the following resources:
- an access key and a secret key
- a vpc with at least 1 subnet (we recommend 3)
- a keypair
- a hosted zone registered in route53
- An SSL certificate that matches your hosted zone's name (optional)
- Permissions to create, update and destroy route53 entries, auto scaling groups, elastic load balancers, security groups, ec2 instances, elastic ips, cloudwatch logs, IAM profiles and cloud formation stacks.

## Create your environment manifest
Follow the instructions [available here](environment-yml.md) to create your `environment.yml` file. It will be similar to the one displayed below:

```
# your docker hub credentials, if using private images
docker:
  username: i2kit-tester
  password: *******

# your AWS configuration
provider:
  name: test
  access_key: ***************
  secret_key: ***************
  region: us-west-2
  subnets:
  - subnet-7c75725
  - subnet-7c75726
  keypair: development
  hosted_zone: example.com.
  size: t2.small
```

## Convert your docker-cloud.yml to i2kit service manifests
To launch your application with i2kit, you need to convert your `docker-cloud.yml` manifest to i2kit service manifests. Once you have every service converted, you can start to deploy and test them one by one.

Unlike docker cloud, i2kit works with `services`. A `service` is a set of virtual machines running the same list of containers. All containers defined on the same service manifest will run in the same virtual machine. For most cases, every first level key in your docker cloud manifest will map to an i2kit's service manifest.

The full example is included [as a sample](../samples/voting-demo).

### db service
The `docker-cloud.yml` file defines the db service as follows:
```
db:
  image: 'postgres:9.4'
  restart: always
```

For i2kit, you'd need to define a service manifest file (e.g. `db.yml`) that contains the name, the image, the number of replicas, and the ports to expose:

```
name: db
replicas: 1
stateful: true
containers:
  db:
    image: postgres:9.4
    ports:
      - tcp:5432:tcp:5432
```

i2kit will automaticaly configure the container to restart, so that policy is not necessary.

When deployed, i2kit will create a load balancer, and configure the containers to use your hosted zone as the primary DNS zone. This would allow the rest of the services on the same VPCto resolve both `db` and `db.HOSTEDZONE` to the load balancer.

### redis service
`docker-cloud.yml` file:

```
redis:
  image: 'redis:latest'
  restart: always
```

i2kit's `redis.yml`:

```
name: redis
replicas: 1
stateful: true
containers:
  redis:
    image: redis:alpine
    ports:
      - tcp:6379:tcp:6379
```

### vote service
The Docker Cloud stackfile for the vote service defines an image, a restart policy, and a specific number of Pods (replicas: 5). This tells docker cloud to always have 5 healthy instances of the pod:

```
vote:
  autoredeploy: true
  image: 'docker/example-voting-app-vote:latest'
  restart: always
  target_num_containers: 5
```

i2kit supports the same feature via the `replicas` attribute:

```
name: vote
replicas: 3
containers:
  vote:
    image: docker/example-voting-app-vote:latest
    ports:
      - http:80:http:80
```

In this example, `vote.yml` is telling i2kit to deploy 3 instances, with a load balancer on port 80, and to configure an autoscaling group to ensure that there are always 3 instances online and healthy. If an instance is destroyed, or is deemed unhealthy, AWS will automatically relaunch it for your.

Currently i2kit doesn't support autoredeploys, so we ignore this attribute for now. This can be implemented by your CI/CD pipeline invoking the `i2kit deploy` command with an updated manifest.

### worker service

`docker-cloud.yml` file:

```
worker:
  autoredeploy: true
  image: 'docker/example-voting-app-worker:latest'
  restart: always
  target_num_containers: 3
```

i2kit's `worker.yml`
```
name: worker
replicas: 1
containers:
  worker:
    image: docker/example-voting-app-worker:latest
```

### result service

`docker-cloud.yml` file:

```
result:
  autoredeploy: true
  image: 'docker/example-voting-app-result:latest'
  ports:
    - '80:80'
  restart: always
```

i2kit's `result.yml`
```
name: result
replicas: 1
containers:
  result:
    image: docker/example-voting-app-result:latest
    ports:
      - http:80:http:80
```

### lb service
The lb service is not required in i2kit. i2kit will automatically create an `elastic load balancer` for your service when required.

## Test the app with i2kit
To test the application, you need to deploy each service one by one. The command will wait for the service to be deployed, and for its healthcheck to be succesfull.

```
i2kit deploy -s db.yml -e environment.yml
i2kit deploy -s redis.yml -e environment.yml
i2kit deploy -s vote.yml -e environment.yml
i2kit deploy -s worker.yml -e environment.yml
i2kit deploy -s result.yml -e environment.yml
```

Once all the commands have finished, validate that it works by browsing to `http://vote.$YOURHOSTEDZONE` and `http://results.$YOURHOSTEDZONE`

To destroy the created resources, run:
```
i2kit destroy -s db.yml -e environment.yml
i2kit destroy -s redis.yml -e environment.yml
i2kit destroy -s vote.yml -e environment.yml
i2kit destroy -s worker.yml -e environment.yml
i2kit destroy -s result.yml -e environment.yml
```

Once your services has been validated, remember point your application CNAMES to new service endpoints. They will all follow the `protocol:$SERVICENAME.$YOURHOSTEDZONE:$PORT` format.
