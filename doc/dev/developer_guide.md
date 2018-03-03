# Developer Guide

This is the guide to get more idea about internal details, how to built and contribute/develope on it.

## Dependency Management

We use [dep](https://github.com/golang/dep) to manage dependencies of the project.

```bash
dep ensure
```

## Code Generation

If you make any change in your custom resource type, in other words, in you api, you should regenerate/update client codes.
But please make relevant changes in `update-codegen.sh` script ,as well, if you did any version change , etc..

```bash
./hack/update-codegen.sh
```

## Build

Required tools:
- Docker
- Go 1.9+
- git

### Cassandra Operator Image build

The Dockerfile in project root directory is for building your own cassandra-operator image.

I use [multi-stage](https://docs.docker.com/develop/develop-images/multistage-build/) build for building image, it makes image really efficient and small sized.
(Versions, repository and container registry are just examples, use your own for your development)

```
( under $GOPATH/src/github.com/aslanbekirov/cassandra-operator/ )
$ docker build -t aslanbekirov/cassandra-operator:v0.0.1 .
$ docker push aslanbekirov/cassandra-operator:v0.0.1
```

This will push the image to the registry and you can use this pushed image in your operator-deployment in `example/example-cassandra-operator.yaml`