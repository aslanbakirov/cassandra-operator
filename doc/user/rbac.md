# Cassandra Operator RBAC Setup

If you kubernetes cluster is RBAC enabled, the users need to setup RBAC rules for cassandra operator. This doc serves a tutorial for it.

## Setup

First of all, we will have to create a cluster role. This role will define which actions can be done on which resources.

```bash
kubectl create -f example/cluster-role.yaml
```

After you created cluster role, we are going to bind this role to specific namespace for specific service-account.

```bash
kubectl create -f example/cluster-role-bindings.yaml
```

This will create a `clusterrolebinding` object which will bind the cluster role defined in previous step to specific service-account (In our case, we are using `default` service account just for testing, but probably it will be different in your case.) for specific namespace (In our example it is `test` namespace, you can make this binding cluster scoped, as well.)

After this RBAC setup, you can start using cassandra operator.