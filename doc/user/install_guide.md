# Installation guide

## Set up RBAC

Set up basic [RBAC rules][rbac-rules] for etcd operator:

```bash
$ example/rbac/create_role.sh
```

## Set up service account

Set up basic service account for etcd operator:

```bash
$ kubectl create -f example/serviceaccount.yaml
```

## Set up CRDs

Set up custom resource definitions for the etcd operator

```bash
$ kubectl create -f example/crd/etcd.database.coreos.com_etcdclusters.yaml
```

## Install etcd operator

Create a deployment for etcd operator:

```bash
$ kubectl create -f example/deployment.yaml
```

## Uninstall etcd operator

Note that the etcd clusters managed by etcd operator will **NOT** be deleted even if the operator is uninstalled.

This is an intentional design to prevent accidental operator failure from killing all the etcd clusters.

To delete all clusters, delete all cluster CR objects before uninstalling the operator.

Clean up etcd operator:

```bash
kubectl delete -f example/deployment.yaml
kubectl delete -f example/serviceaccount.yaml
kubectl delete -f example/crd/etcd.database.coreos.com_etcdclusters.yaml
kubectl delete clusterrole etcd-operator
kubectl delete clusterrolebinding etcd-operator
```

## Installation using Helm

etcd operator is available as a [Helm chart][etcd-helm] for this fork. Follow the instructions on the chart to install etcd operator on clusters.
[Phil Porada][pgporada] is the active maintainer.


[rbac-rules]: rbac.md
[etcd-helm]: https://github.com/pgporada/etcd-operator-helm-chart
[pgporada]: https://github.com/pgporada