# Service YAML reference

A service is a set of containers running the same container image.

Below is an example of a `service.yml`:

```yaml
name: test
replicas: 3
stateful: false
public: true
instance_type: t2.small
containers:
  nginx:
    image: nginx:alpine
    environment:
      - NAME=VALUE
    ports:
      - https:443:http:80
    command: start.sh
```

Each key is documented below:

## name (required)

The name of the service. Every service on the same project will be able to ping this service using its name.

```yaml
name: api
```

Also, a global DNS entry is created for this service. For example, if a service is named *api* and it belongs to the project *staging*, the DNS *api.staging.okteto.net* resolves to the service endpoint.

Service names should be unique per project.

Service names only accept alphanumeric characters.

## replicas (optional)

The number of instances running this service (default: `1`).

```yaml
replicas: 3
```

## stateful (optional)

Make it true for services running a single instance which cannot be accessed through a load balancer, such as databases (default: `false`).

```yaml
stateful: true
```

## public (optional)

Make it true for services accessible from outside of this project (default: `false`).

```yaml
public: true
```

## instance_type (optional)

For AWS projects, it is the instance type that will be used to create the service instances.

If it is not specified, it takes the default value specified at the project level, or `t2.small` if none of these values is set.

```yaml
instance_type: m3.medium
```

## containers (required)

It is the set of containers running this service.

The most common case is to have a single containers, but there are other use cases where more than a single containers is needed (what it is also known as Pods).

Each container might have the following keys:

### image (required)

The image used to deploy this container. This is the only mandatory key.

```yaml
image: drupal
image: nginx:alpine
image: my.registry.com/redis
```

### environment (optional)

A list of environment variables to add in the container at launch. They are represented using an array format.

```yaml
environment:
  - USER=user1
  - PASSWORD=password1
```

### command (optional)

Overrides the default command in the container image.

```yaml
command: echo 'Hello Okteto!'
```

### ports (optional)

Defines the ports that are accessible from other services. They specify the load balancer port number and protocol, and the container port and protocol. For example, the next port defines an external `443` port accepting `https` requests that get redirected into the container `80` port using `http`.

```yaml
ports:
  - https:443:http:80
```

The certificate to be used for `https` and `ssl` protocols is defined at the project level.
