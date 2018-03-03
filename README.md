# Cassandra Operator 


**Disclaimer:** This is my side project to learn/experiment and apply all features of Custom Resource Definitions, stateful applications and related technologies in Kubernetes. It is under active development, can be beneficial to be used as template for cassandra operator design

## Project Overview

The Cassandra operator manages Cassandra clusters deployed to Kubernetes and automates tasks related to operating an Cassandra cluster.

- [Create and destroy](#create-and-destroy-an-Cassandra-cluster)
- [Backup and restore a cluster] (Not implemented) 
- [Rolling upgrade](Not implemented)

Read [RBAC docs](./doc/user/rbac.md) for how to setup RBAC rules for Cassandra operator if RBAC is in place.

Read [Developer Guide](./doc/dev/developer_guide.md) for setting up development environment if you want to contribute.

## Requirements

- Kubernetes 1.8+
- Cassandra 3.11+

## Deploy Cassandra operator

```bash
kubectl create -f example/example-cassandra-operator.yaml
```
This is a deployment for operator itself. It will create a custrom resource called `cassandraclusters.cassandra.database.com` which will enable us to create `cassandracluster` objects/resources

## Create and destroy an Cassandra cluster

```bash
$ kubectl create -f example/example-cassandra-cluster.yaml
```
I included some cassandra config in `example/example-cassandra-cluster.yaml` file, but these are just examples to show how you can set them. Obviously, you can use your cassandra image for your projects.

A 3 member Cassandra cluster will be created.

```bash
$ kubectl get pods
NAME                            READY     STATUS    RESTARTS   AGE
cassandra-0       1/1       Running   0          1m
cassandra-1       1/1       Running   0          1m
cassandra-2       1/1       Running   0          1m
```
Destroy Cassandra cluster:

```bash
$ kubectl delete -f example/example-cassandra-cluster.yaml
```