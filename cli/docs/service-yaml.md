# Service YAML reference
A service is a set of virtual machines running the same container/pod.

Below is an example of a `service.yml`:

```
name: test
replicas: 3
containers:
  nginx:
    image: nginx:alpine
    environment:
      - NAME=VALUE
    ports:
      - https:443:http:80
    command: start.sh
```
