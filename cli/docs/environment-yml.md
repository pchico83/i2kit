# Environment YAML reference
An environment is an infrastructure configuration shared between a set of services.

Below is an example of an `environment.yml`:

```
# your docker hub credentials
name: test
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
  security_group: sg-dffe41a3
  keypair: development
  hosted_zone: example.com.
  instance_type: t2.small
  certificate: arn:aws:acm:us-west-2:062762192540:certificate/12de3ac5-abcd-461a-1234-9e81250b33d8
```
