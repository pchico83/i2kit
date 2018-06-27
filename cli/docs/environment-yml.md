# AWS Project YAML reference

A project defines a group of services that are interconnected and accessible thought private IPs.

Below is an example of a `project.yml`:

```yaml
administrators:
- user1@gmail.com
users:
- user2@gmail.com
- user3@gmail.com
docker:
  username: okteto
  password: password
provider:
  type: aws
  access_key: ***************
  secret_key: ***************
  region: us-west-2
  subnets:
  - subnet-7c75725
  - subnet-7c75726
  keypair: development
  instance_type: t2.small
  hosted_zone: example.com.
  certificate: arn:aws:acm:us-west-2:062762192540:certificate/12de3ac5-abcd-461a-1234-9e81250b33d8
```

Each key is documented below:

## administrators (required)

A list of administrators that have full access to this project. Each administrator will be able to modify and delete the project, as well as deploy, upgrade and destroy services in this project.

```yaml
users:
- user1@gmail.com
- user2@gmail.com
- user3@gmail.com
```

## users (optional)

A list of users that have access to this project. Each user will be able to deploy, upgrade and destroy services in this project.

```yaml
users:
- user2@gmail.com
- user3@gmail.com
```

## docker (optional)

Docker credentials to access your docker images in Docker Hub. It is composed of the following sub keys:

### username (required)

Your Docker Hub username, also known as Docker ID.

```yaml
username: okteto
```

### password (required)

Your Docker Hub password.

```yaml
password: password
```

## provider (required)

Defines the infrastructure where your services gets deployed. It is composed of the following sub keys:

### type (required)

Valid values are `aws` and `k8`. This document refers to the `aws` type.

```yaml
type: aws
```

### access_key and secret_key (required)

AWS credentials used by Okteto to manage resources in your AWS account.

Okteto needs the following policies to create services in your account: `AmazonEC2FullAccess`, `CloudWatchLogsReadOnlyAccess`, `AmazonRoute53FullAccess` and `CloudFormationPolicy`.

```yaml
access_key: AKIAI2K5O4A7QON6YSEA
secret_key: O/5OOpWOuOGW9LZX5wPsNzICjG9VWJlVURevbqmu
```

### region (required)

The AWS region where your services will be deployed.

```yaml
region: us-west-2
```

### subnets (required)

The list of subnets where your instances will be created.

It is recommended to define several subnets for high availability.

```yaml
subnets:
- subnet-7c75725
- subnet-7c75726
```

### keypair (required)

The keypair injected in your service instances. It can be used to SSH-into your instances.

```yaml
keypair: development
```

### instance_type (optional)

The default instance type for instances running your services.

Each service can overwrite this value in its yaml definition.

```yaml
instance_type: m3.medium
```

### hosted_zone and certificate (optional)

In order to define `HTTPS`/`SSL` endpoints in your services, you will need to specify a hosted zone and an AWS certificate.

The hosted zone is used to create your service DNS entries, instead of using `okteto.net.`.

The `certificate` will be attached to each `HTTPS`/`SSL` port defined by your services.

```yaml
hosted_zone: example.com.
certificate: arn:aws:acm:us-west-2:062762192540:certificate/12de3ac5-abcd-461a-1234-9e81250b33d8
```
