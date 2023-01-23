# Changelog

## [1.2.2](https://github.com/on2itsecurity/etcd-operator/compare/v1.2.1...v1.2.2) (2023-01-23)

### Must Reads
* Upgrade to 1.2.2 before upgrading etcd to 3.5.7 ([14c27ce](https://github.com/on2itsecurity/etcd-operator/commit/14c27ceb08b2bd496546fbcddcd0e20ee6958e58))
* Dropped etcd 3.2 support ([14c27ce](https://github.com/on2itsecurity/etcd-operator/commit/14c27ceb08b2bd496546fbcddcd0e20ee6958e58))
* Next version will not create crds, use manifests from `example/crd` for new installs ([945971e](https://github.com/on2itsecurity/etcd-operator/commit/945971ed5a3ab5d659ab1db4375c7924db8c5f37))

### Bug Fixes
* openapi validation fixes ([b9f8821](https://github.com/on2itsecurity/etcd-operator/commit/b9f8821a7c23d9d7050eeab52ba197ea2c3fd7e9))
* do not continue making a pod when an error occurs. closes issue #80 ([2498fd5](https://github.com/on2itsecurity/etcd-operator/commit/2498fd5d8127d86167e4d573e684a1c550a7da61))

### Features
* crds are now generated as yaml with openapi support, applying it manually means etcd-operator can be started with `-create-crd=false` and https://github.com/on2itsecurity/etcd-operator/blob/f378e398dce7bc4546e8ade3c156c0d7d6e33850/example/rbac/cluster-role-template.yaml#L14-L19 can be removed  ([5a2ef43](https://github.com/on2itsecurity/etcd-operator/commit/5a2ef4f34594643370956570365b15c4a640ed73))
* pvc support in restore ([614ea03](https://github.com/on2itsecurity/etcd-operator/commit/614ea035d00f8d4788a2bde4fb7401ad0c1b34b9))
* build against Golang 1.19 ([6750543](https://github.com/on2itsecurity/etcd-operator/commit/6750543af27f095ef491087e20122398997de23f))
* Kubernetes 1.25 compatibility & test against ([c8d20c1](https://github.com/on2itsecurity/etcd-operator/commit/c8d20c16b9fe4b08b4c48ba7ed4f4361f001145c))
* Test on Kubernetes 1.26 ([4e88f0a](https://github.com/on2itsecurity/etcd-operator/commit/4e88f0a77e208253894136df616b78336d4152c8))
* Use Kubernetes 1.26 Libraries ([4532bff](https://github.com/on2itsecurity/etcd-operator/commit/4532bff0b574a879f3c558c0386e72411333023d))
* Default to etcd 3.5.6 + test against ([24cfdf0](https://github.com/on2itsecurity/etcd-operator/commit/24cfdf0b2d188201b6cda606e9d7ada763c4e941))
* Upgrade alpine to 3.17 ([36b086d](https://github.com/on2itsecurity/etcd-operator/commit/36b086d389fb3674dd49b4508073b44b541ec652))
* Build without linking to musl ([1042088](https://github.com/on2itsecurity/etcd-operator/commit/1042088a4f41e5828ef23579a7fa6e191cd1131f))
* Add notice for create-crd=true ([81032db](https://github.com/on2itsecurity/etcd-operator/commit/81032dbf8879e0ec8f0b2367b71f1f758c992126))

## [1.2.1](https://github.com/on2itsecurity/etcd-operator/compare/v1.2.0...v1.2.1) (2022-05-19)

### Bug Fixes
* try to update the podDisrusptionBudget after it was initially created ([8e66977](https://github.com/on2itsecurity/etcd-operator/commit/8e6697706e32255701c09c98678588af52bd740d))
* remove last etcd 2.3.8 dependency ([37891d0](https://github.com/on2itsecurity/etcd-operator/commit/37891d0c3abeefd5d84a72428a65163ec0799508))


### Features
* retry with client-go native routine (backoff) ([c2fdf81](https://github.com/on2itsecurity/etcd-operator/commit/c2fdf8180e99e671e680b792a7d576ac2ec73c9d))
* up Kubernetes api to 1.24.0 ([c259924](https://github.com/on2itsecurity/etcd-operator/commit/c25992439dc811a48934a5bd130767ab82261fc4))
* up default ETCD to 3.5.4 ([e914d26](https://github.com/on2itsecurity/etcd-operator/commit/e914d2654702b104694adf54aacb48ec9b0b7b15))
* test against kubernetes 1.24.0-rc1 ([2b9b4e7](https://github.com/on2itsecurity/etcd-operator/commit/2b9b4e7597971e7cceed3fa71de92ea675bd1b2e))

## [1.2.0](https://github.com/on2itsecurity/etcd-operator/compare/v1.1.1...v1.2.0) (2022-04-07)

### Bug Fixes
* ready check circular dependency with rolling deployment and leader election ([645c78d](https://github.com/on2itsecurity/etcd-operator/commit/645c78d26b0efbd3339a8a63b28ee0359ca52576))
* Make etcd-operator use the right amount of threads in containers ([01d0dd6](https://github.com/on2itsecurity/etcd-operator/commit/01d0dd63c62609e6f8413d1f58f7d693dcd7ea4a))

## Features
* Kubernetes 1.24 API / ETCD 3.5.2 API ([c1424cb](https://github.com/on2itsecurity/etcd-operator/commit/c1424cb90bd25ec822e239b730f2349bd751df8b))
* Using PublishNotReadyAddresses for services ([e73da47](https://github.com/on2itsecurity/etcd-operator/commit/e73da475e09d1c1546b9dd2e49a19b8dc7c6c9c0))
* Default to etcd 3.4.18 and official gcr.io/development/etcd-io ([c5807b3](https://github.com/on2itsecurity/etcd-operator/commit/c5807b3ef9b1820e087a8e3213e2ed8c75702936))
* Build docker images for arm64 ([69d10e1](https://github.com/on2itsecurity/etcd-operator/commit/69d10e16f867f181830ef5f88970a320d050553a))
* add priorityClassName support to etcd crd ([d4d528a](https://github.com/on2itsecurity/etcd-operator/commit/d4d528aa7707a851fa0960e555a6ac2753e6e859))
* podDisruptionBudget creates and maintains the policy to protect the etcd cluster from disruptive kubernetes actions ([f378e39](https://github.com/on2itsecurity/etcd-operator/commit/f378e398dce7bc4546e8ade3c156c0d7d6e33850))

## [1.1.1](https://gitlab.com/cbws/tools/etcd-operator/compare/v1.1.0...v1.1.1) (2021-02-18)


### Bug Fixes

* lock golang version to older version ([33def29](https://gitlab.com/cbws/tools/etcd-operator/commit/33def29c56e38eab69c817050d44be4324d792b8))

# [1.1.0](https://gitlab.com/cbws/tools/etcd-operator/compare/v1.0.0...v1.1.0) (2021-01-23)


### Features

* **backup-operator:** add prometheus-metrics. ([6dece02](https://gitlab.com/cbws/tools/etcd-operator/commit/6dece028c3595252bdc9f769cf9333b6b3a4847b))

# [1.0.0](https://gitlab.com/cbws/tools/etcd-operator/compare/v0.11.0...v1.0.0) (2021-01-22)


### Features

* **tag:** remove auto-adding v to the etcd image version ([125c141](https://gitlab.com/cbws/tools/etcd-operator/commit/125c141793b8a19d6192cd503bac5b96084fa420))


### BREAKING CHANGES

* **tag:** image names now need the full tag

# [0.11.0](https://gitlab.com/cbws/tools/etcd-operator/compare/v0.10.1...v0.11.0) (2021-01-21)


### Features

* upgrade k8s and etcd libraries ([d8f44fb](https://gitlab.com/cbws/tools/etcd-operator/commit/d8f44fbcaa79f695fd6d9372ae0c30c82e8bdc69))

## [0.10.1](https://gitlab.com/cbws/tools/etcd-operator/compare/v0.10.0...v0.10.1) (2021-01-21)


### Bug Fixes

* try and fix github.com/aliyun/aliyun-oss-go-sdk dependency update ([5a4025b](https://gitlab.com/cbws/tools/etcd-operator/commit/5a4025b2b11714a3d273135be50cef056df71622))

## [Unreleased]

### Added

- Added `spec.Pod.ClusterDomain` to explicitly set the cluster domain used for the etcd member URLs. [#2082](https://github.com/coreos/etcd-operator/pull/2082)

### Changed

### Removed

### Fixed

### Deprecated

### Security

## [Release 0.10.0]

### Added
* Alibaba Cloud Object Storage Service (OSS) backend - coreos#2065 (@stormgbs)
* Cluster wide restore - coreos#2074 (@kapouille)
* Cluster wide backup - coreos#2072 (@kapouille)
* Support specifying cluster domain for clientURLs - coreos#2028 (@mikkeloscar )
* Add tmpfs support - coreos#2103 (@chriscoester)

### Fixed
* Don't expose unready pods via client service - coreos#2063 (@ckoehn)

## [Release 0.9.4]

### Added

- Added `spec.BackupSource.S3.ForcePathStyle` to `EtcdBackup` to force path style s3 uploads. [#2036](https://github.com/coreos/etcd-operator/pull/2036)
- Added `spec.RestoreSource.S3.ForcePathStyle` to `EtcdRestore` to force path style s3 downloads. [#2036](https://github.com/coreos/etcd-operator/pull/2036)

### Changed

- Update Go version to 1.11.5
- Update k8s to 1.12.6
- EtcdBackup: Support periodic backups. This change added 3 new fields to EtcdBackup schema, 2 variables in spec, 1 variables in status.
  - in spec.backupPolicy
    - maxBackup: maximum number of backups to keep.
    - backupIntervalInSecond: how often to perform backup operation.
  - in status
    - LastSuccessDate: last time to succeed in taking backup

### Fixed

- Fixed a bug where `same CR names` in different namespaces with cluster-wide operator were not working as expected [#2026](https://github.com/coreos/etcd-operator/pull/2026)
- Fixed a bug where cluster names could exceed 63 octets the maximum defined by [RFC 1035 section 2.3.4](https://tools.ietf.org/html/rfc1035) resulting in hanging pods [2027](https://github.com/coreos/etcd-operator/pull/2027).

## [Release 0.9.3]

### Added

- Added `spec.pod.DNSTimeoutInSecond` to `EtcdCluster` that allows setting a maximum allowed time for the init container of the etcd pod to reverse DNS lookup its IP given the hostname.

### Changed

- Update Go version to 1.11.2
- Update k8s to 1.11.4
- k8s codegen updates are longer performed via container. Go dependencies are now vendored
  and updates are performed with shell script locally.

### Fixed

- Fixed leaking http connections while verifying backup snapshots. [#1976](https://github.com/coreos/etcd-operator/pull/1976)

## [Release 0.9.2]

### Added

- Added the field `spec.pod.securityContext` to `EtcdCluster` that allows setting a specific [PodSecurityContext](https://kubernetes.io/docs/tasks/configure-pod-container/security-context/#set-the-security-context-for-a-pod) for the etcd pods. [#1949](https://github.com/coreos/etcd-operator/pull/1949)

### Changed

- Update Go version to 1.10
- Build `gcr.io/coreos-k8s-scale-testing/etcd-operator-builder:0.4.1-2` container
  with Go 1.10 and dep 0.4.1
- etcd pod containers no longer run with a non-root security context by default. This setting can be configured per cluster via the PodPolicy.

## [Release 0.9.1]

### Added

- Added optional flag `--cluster-wide` to etcd-operator to allow it to manage etcd clusters across all namespaces. [#1777](https://github.com/coreos/etcd-operator/pull/1777)
- Added support for annotation `etcd.database.coreos.com/scope: clusterwide` in `EtcdCluster` to allow it to be managed by a cluster wide operator. [#1777](https://github.com/coreos/etcd-operator/pull/1777)
- Added the field `spec.pod.busyboxImage` to the `PodPolicy` of the `EtcdCluster` to allow overriding the default busybox image used for the etcd pod's init container. [#1928](https://github.com/coreos/etcd-operator/pull/1928)

### Fixed

- Fixed a bug where the informer watch stream would timeout after 30s of not receiving an event. [#1936](https://github.com/coreos/etcd-operator/pull/1936)

## [Release 0.9.0]

Same as v0.8.4. The version is bumping to 0.9.0 due to adding a new ABS backup API into etcd-backup-operator.

## [Release 0.8.4]

### Added

- Added ABS support for backup and restore
- Added tag to initContainer to enable offline deploys
- Enabled configurable backup timeout in backup operator

### Changed

- Set 30s default request timeout for kube client
- Change check-dns init container image to busybox:1.28.0-glibc to fix nslookup failure in some environment.

### Removed

- Removed self-hosted code

## [Release 0.8.3]

### Added

- Added the option to use PersistentVolume as non-stable storage for etcd pods. This feature is still alpha and subject to change in future releases [#1861](https://github.com/coreos/etcd-operator/pull/1861)

### Changed

- Changed etcd pod member names to be unique by having a random suffix instead of a sequence number. This change is backward compatible and should not affect operator upgrade.
    Previously the etcd pod names would look like:
    ```
    NAME                            READY     STATUS    RESTARTS   AGE
    example-etcd-cluster-0000       1/1       Running   0          1m
    example-etcd-cluster-0001       1/1       Running   0          1m
    example-etcd-cluster-0002       1/1       Running   0          1m
    ```
    After this change:
    ```
    NAME                                  READY     STATUS    RESTARTS   AGE
    example-etcd-cluster-2885zjw9he       1/1       Running   0          1m
    example-etcd-cluster-gghrmbeid4       1/1       Running   0          1m
    example-etcd-cluster-w5q9sn37fd       1/1       Running   0          1m
    ```

### Fixed

- Fixed a bug where the restore operator would fail to restore the seed member because recreating an etcd pod with the same name as a recently deleted one would conflict as the older pod and its resources, like the DNS name, might still not be deleted. [#1825](https://github.com/coreos/etcd-operator/issues/1825)


## [Release 0.8.2]

### Added

- Add support for backup and restore from custom S3 endpoint.

### Changed

- All etcd pod containers now run as non-root.


## [Release 0.8.1]

### Changed

- etcd-restore-operator will create a service for itself as the backup storage proxy. Delete the service in deployment yaml.

### Fixed

- Fix etcd-restore-operator wouldn't report error and keep looping if EtcdRestore name is different than EtcdCluster name.


## [Release 0.8.0]

**Important Changes**

Both etcd backup operator and etcd restore operator have changed their CR definition.
Please follow the latest backup/restore CR definition for future backup and restore.

### Added

- Add readiness probe to etcd pod. The readiness state will be reflected on `status.members.ready/unready`.
- TLS etcd cluster support in backup/restore-operator.
- Add spec validation in restore operator.
- Add BackupStorageType to EtcdRestore.RestoreSpec to indicate type of the backup storage which is used as RestoreSource and validation of BackupStorageType in restore operator.
- Add EtcdClusterRef to EtcdRestore.RestoreSpec to reference an EtcdCluster resource whose metadata and spec will be used to create the new restored EtcdCluster CR.
- Add create-crd flag to etcd backup operator allowing user to disable automatic backup CRD creation.
- Add create-crd flag to etcd restore operator allowing user to disable automatic restore CRD creation.
- Add EtcdVersion and EtcdRevision to EtcdBackup.BackupStatus.
- BackupStatus: Add detailed error when backup fails.

### Changed

- Rename BackupCRStatus to BackupStatus.
- EtcdBackup: BackupSpec passes in S3BackupSource.Path as the S3 path to save the backup.
- EtcdBackup: BackupSpec spec uses etcd endpoints to retrieve snapshot.
- Change default etcd version to `3.2.13`.

### Removed

- EtcdBackup: BackupSpec removed ClusterName field in favor of etcd endpoints.
- EtcdCluster: ClusterSpec removed deprecated BaseImage field.


## [Release 0.7.2]

> Note: This is a bug fix release.

When we bump the etcd version to 3.2, the images were only available on gcr.io . But now it is added on quay.io .
We'd better use quay.io and keep it compatible to work for 3.1 versions of etcd.

## [Release 0.7.1]

### Added

- TLS etcd cluster feature for EtcdBackup
- Log collector program for collecting logs in e2e test.
- ClusterSpec: In PodPolicy, add generic `Affinity` field to substitute bool `AntiAffinity` field.
- ClusterSpec: Add `Repository` field to substitute `BaseImage` field.

### Changed

- Default base image is changed to `gcr.io/etcd-development/etcd`, default etcd version is `3.2.11`.
- Migrate dependency management tooling from glide to dep.
- Containerize e2e test in a pod instead of running on raw jenkin slave.

### Removed

- ClusterSpec: Remove `PodPolicy.AutomountServiceAccountToken` field.
  No etcd pod will have service account token mounted.

### Fixed

- Ignore Terminating pods when polling etcd pods.

### Deprecated

- ClusterSpec: `BaseImage` is deprecated. It will be automatically converted to `Repository` in this release.
- ClusterSpec: In PodPolicy, `AntiAffinity` is deprecated. It will be automatically converted to `Affinity.PodAntiAffinity`
  terms with label selector on given cluster name and topology key on node in this release.

### Security

- All operator images by default uses user `etcd-operator` instead of root.


## [Release 0.7.0]

Existing backup and restore features in EtcdCluster API won’t be supported after 0.7.0 release.
See [Decoupling Backup and Restore Logic from Etcd Operator](https://github.com/coreos/etcd-operator/issues/1626) for more detail.

If applicable then see the [upgrade guide](./doc/user/upgrade/upgrade_guide.md) on how to upgrade from `v0.6.1` to `v0.7.0` .

### Added

- Add `ServiceName` and `ClientPort` into ClusterStatus.
- Add etcd backup operator for backing up an etcd cluster.
- Add etcd restore operator for restoring an etcd cluster.

### Removed

- Remove `pv-provisioner` flag from etcd operator.
- Remove etcd cluster Backup feature from etcd operator.
- Remove etcd cluster Restore from etcd operator.


## [Release 0.6.1]

The operator will no longer create a storage class specified by `--pv-provisioner` by default. If applicable then see the [upgrade guide](./doc/user/upgrade/upgrade_guide.md) on how to upgrade from `v0.6.0` to `v0.6.1` .

### Added

- backup binary supports serving backup defined by backupSpec. In addition, when backupSpec
is specified, backup binary changes to serve http backup requests only mode.
- Add operator flag `--create-crd`. By default it is `true` and operator will create EtcdCluster CRD.
  It can be set to `false` and operator won't create EtcdCluster CRD.
- Add operator flag `--create-storage-class`. By default it is `false` and operator won't create default storage class.
  It can be set to `true` and operator will create default storage class.

### Changed

- An EtcdCluster CR with an invalid spec will not be marked as failed. Any changes that result in an invalid spec will be ignored and logged by the operator.

### Fixed

- Fix the problem that operator might keep failing on version conflict updating CR status.

### Deprecated

- The operator flag `--pv-provisioner` is depercated. We recommend to use per cluster storageClass.


## [Release 0.6.0]

**BREAKING CHANGE**: operator level S3 backup is removed. See [upgrade](./doc/user/upgrade/upgrade_guide.md) on how to upgrade from 0.5.x to 0.6.0.

### Added

- Add cluster events into EtcdCluster custom resource. See `doc/user/conditions_and_events.md` .

### Changed

- Redefine status.conditions. See `doc/user/conditions_and_events.md` .

### Removed

- Remove operator level S3 flag.
- Remove analytics flag. Disable Google analytics.


## [Release 0.5.2]

### Added

- Expose `/metrics` endpoint at port 8080
- Add cluster S3 spec `prefix` feature. Let user choose a prefix under the bucket.
- Add `automountServiceAccountToken` to pod policy. Let users disable automounting of the Kubernetes access token into etcd-operator controlled pods.
- Cluster backups can now be saved using Azure Blob Storage (ABS).

### Deprecated

- Deprecate operator S3 flag. Add warning note for using it in this release. The flag will be removed in 0.6.0 release.


## [Release 0.5.1]

Upgrade notice for TLS cluster users:
If you are using TLS-enabled etcd cluster, the SAN domain has been changed. See [TLS docs](./doc/user/cluster_tls.md).
Before upgrading operator, you need to rotate certs on each secrets to allow both the old and new domains.
Then restart each etcd pod -- the simplest way is to "upgrade" cluster version.
Finally, it is safe to upgrade operator. It's highly recommended to save a backup before upgrade.

### Added

- A new `StorageClass` spec field, allowing more granular control over how etcd clusters are backed up to PVs.

### Changed

- Default timeout for snapshots done by backup sidecar increased from 5 seconds to 1 minute

### Fixed

- Fix periodFullGC only executed once problem.
- [GH-1021] Use the cluster domain provided by kubelet instead of hardcoded `.cluster.local` .


## [Release 0.5.0]

**BREAKING CHANGE**: The cluster object will now be defined via a [Custom Resource Definition(CRD)](https://kubernetes.io/docs/tasks/access-kubernetes-api/extend-api-custom-resource-definitions/) instead of a [Third Party Resource(TPR)](https://kubernetes.io/docs/tasks/access-kubernetes-api/extend-api-third-party-resource/). See the `Changed` section below for details.

### Changed

- With k8s 1.7 and onwards TPRs have been deprecated and are replaced with CRD. See the k8s 1.7 [blogpost](http://blog.kubernetes.io/2017/06/kubernetes-1.7-security-hardening-stateful-application-extensibility-updates.html) or [release notes](https://github.com/kubernetes/kubernetes/blob/master/CHANGELOG.md#major-themes) for more details. For this release a live migration of the cluster spec from TPR to CRD is not supported. To preserve the cluster state during the upgrade you will need to create a backup of the cluster and recreate the cluster from the backup after upgrading the operator. See the [upgrade guide](https://github.com/coreos/etcd-operator/blob/master/doc/user/upgrade/upgrade_guide.md) for more detailed steps on how to do that.

- Changes in the cluster object's type metadata:
  - The `apiVersion` field has been changed from `etcd.coreos.com/v1beta1` to `etcd.database.coreos.com/v1beta2`
  - The `kind` field has been changed from `Cluster` to `EtcdCluster`


## [Release 0.4.2]

### Added

- [GH-1232](https://github.com/coreos/etcd-operator/pull/1232) the operator can now log critical actions like pod creation/deletion to a user specified path via the optional flag `debug-logfile-path`. The logs will only be generated if the cluster is self hosted and the flag is set. This can be used in conjunction with a persistent volume to persist the critical actions to disk for later inspection.

### Changed

- enable alpha feature "tolerate unready endpoints" on etcd client and peer service

### Fixed

- Fix append-hosts init-container not run on some restart cases.


## [Release 0.4.1]

This is a bug-fix release. We have done a lot of testing against k8s 1.7
and making it stable on 1.7 .

### Added

- New self-hosted field `SkipBootMemberRemoval` allows users to skip the
  auto-deletion of the boot etcd member.

### Fixed

- Make sure etcd pod's FQDN is resolvable before running etcd commands .


## [Release 0.4.0]

**BREAKING CHANGE**: Re-naming of TLS spec and TLS secrets' fields.

TLS spec:
- member's `clientSecret` is changed to `serverSecret`

TLS secrets:
- member's `peerSecret`'s fields change:
  - peer-crt.pem -> peer.crt
  - peer-key.pem -> peer.key
  - peer-ca-crt.pem -> peer-ca.crt
- member's `clientSecret` is changed to `serverSecret`, its fields change:
  - client-crt.pem -> server.crt
  - client-key.pem -> server.key
  - client-ca-crt.pem -> server-ca.crt
- `operatorSecret`'s fields change:
  - etcd-crt.pem -> etcd-client.crt
  - etcd-key.pem -> etcd-client.key
  - etcd-ca-crt.pem -> etcd-client-ca.crt


**BREAKING CHANGE**: Backup spec: `CleanupBackupsOnClusterDelete` field is renamed to `AutoDelete`.

Previous spec like this one

```yaml
spec:
  backup:
    storageType: "PersistentVolume"
    ...
    cleanupBackupsOnClusterDelete: true
```

should be changed to

```yaml
spec:
  backup:
    storageType: "PersistentVolume"
    ...
    autoDelete: true
```


## [Release 0.3.3]

### Added

- Adds ability for users to specify base image for etcd pods in a cluster.
  Default base image is `quay.io/coreos/etcd-operator`.

### Fixed

- [GH-1138] Fixed operator stucks in managing selfhosted cluster when there are not enough nodes to start new etcd member.
- [GH-1196] Fixed etcd operator could not start S3 backup sidecar if given non-root user.


## [Release 0.3.2]

Bug fix release to fix self-hosted etcd issue [GH-1171] .

## [Release 0.3.1]


**Notes for self-hosted etcd**:
The newly introduced TLS feature for self hosted etcd is a breaking change.
Existing self hosted etcd cluster MUST be recreated for updating to this release.

### Added

- Self-hosted etcd: if `etcd-hosts.checkpoint` file exists under `${datadir}/`,
  etcd pod will restore the hosts mapping from it before etcd bootstraps.
- Add static TLS support for self-hosted etcd mode.
- The operator will now post Kubernetes [events](https://kubernetes.io/docs/resources-reference/v1.6/#event-v1-core). To allow this the necessary RBAC rule for the resource `events` must be added to the clusterrole. See the [rbac guide](https://github.com/coreos/etcd-operator/blob/master/doc/user/rbac.md#create-clusterrole) to see how to set up RBAC rules for the operator. If the rbac rule for 'events' is not present then the operator will continue to function normally but will also print out an error message on the failure to post an event.
- Add revision field in backup status.
- Support getting a specific backup with version and revision from the backup service.

### Changed

- Self-hosted etcd: use FQDN for client/peer URL.
- Updated RBAC rules for resources `thirdpartyresources` and `storageclasses` to all verbs `*` . We loose granularity early so that we have more flexibility to use other methods (e.g. Get) later.

### Removed

- Update default etcd version to 3.1.8

### Fixed

- [GH-1108] selfHosted: fix backup unable to talk to etcd pods

### Deprecated

### Security


## [Release 0.3.0]

### Upgrade Notice

Check https://github.com/coreos/etcd-operator/blob/master/doc/user/upgrade/upgrade_guide.md#v02x-to-v03x

### Added

- Added support for backup policy to be dynamically added, updated
- Added per cluster policy support for S3.

### Changed

- Backup sidecar deployment created with `Recreate` strategy.
- Spec.Backup.MaxBackups meaning change: 0 means unlimited backups; < 0 will be rejected.

### Removed

### Fixed

- [GH-1068] Backup sidecar deployment stuck at upgrading

### Deprecated

### Security


## [Release 0.2.6]

### Upgrade Notice

- Once operator is upgraded, all backup-enabled cluster will go through an upgrade process that
  deletes backup sidecar's ReplicaSet and creates new Deployment for sidecar.
  If upgrading failed for any reason, cluster TPR's `status.phase` will be FAILED.
  Recreate of the cluster TPR is needed on failure case.

### Added

- PodPolicy provides `EtcdEnv` option to add custom env to the etcd process.
- PodPolicy provides `Labels` option to add custom labels to the etcd pod.
- TLS feature: user can now create TLS-secured cluster via operator.
  See [TLS guide](https://github.com/coreos/etcd-operator/blob/master/doc/user/cluster_tls.md).

### Changed

- Self-hosted etcd pod's anti-affinity label selector is changed to select `{"app": "etcd"}`.
  That is, no two etcd pods should sit on the same node, even if they belongs to different clusters.
- Using Deployment to manage backup sidecar instead of ReplicaSet.
- S3 backup path is changed to `${BUCKET_NAME}/v1/${NAMESPACE}/${CLUSTER_NAME}/`.

### Removed

### Fixed

- Fixed an issue where liveness probes failed when authentication was enabled. [#1957](https://github.com/coreos/etcd-operator/issues/1957)

### Deprecated

### Security


## [Release 0.2.5]

### Added

- Add "none" PV provisioner option. If operator flag "pv-provisioner" is set to "none",
  operator won’t create any storage class and users couldn’t make use of operator’s PV backup feature.
- Add headless service `${clusterName}` which selects etcd pods of given cluster.
- Pod Tolerations.

### Changed

- TLSSpec json tag changed as `omitempty`
- Time related fields in spec, i.e. TransitionTime and CreationTime, is changed to type `string`.
  This should be backward compatible and no effect on operator upgrade.
- Update default etcd version to 3.1.4
- Self-hosted etcd pod is started with "--metrics extensive" flag.
  This is only available in etcd 3.1+.
- Change client LB service name to `${clusterName}-client`.
- Add hostname and subdomain to etcd pods, which makes them have A records formatted in `${memberName}.${clusterName}.${namespace}.svc.cluster.local` .
  For more info, see https://kubernetes.io/docs/concepts/services-networking/dns-pod-service/ .
  We also change PeerURL of etcd members to use such hostnames.

### Removed

- Individual etcd member's services were removed. Use hostname and subdomain of etcd pod instead.

### Fixed

- [GH-910] Operator keeps updating status even if there is no change.

### Deprecated

### Security


## [Release 0.2.4]

### Added

### Changed

### Removed

### Fixed

- [GH-900] Fix looping of reconcile skip due to unready members

### Deprecated

### Security


## [Release 0.2.3]

### Added

### Changed

### Removed

### Fixed

- [GH-890] Fix a race that when majority of members went down cluster couldn't recover
- Fix self-hosted cluster reboot case

### Deprecated

### Security


## [Release 0.2.2]

### Added

- Backup creation time is added in backup status.
- Total size of backups time is added in backup service status.
- Cluster members that are ready and unready to serve requests are tracked via the ClusterStatus fields `Members.Ready` and `Members.Unready`

### Changed

- PodPolicy `resourceRequirements` field is renamed to `resources`
- Default etcd version is changed to `3.1.2`
- Self-hosted etcd pod uses hostPath with path `/var/etcd/$ns-$member`

### Removed

### Fixed

- [GH-851] Fixed a race that caused nil pointer access panic
- [GH-823] Fixed backup service status not shown in TPR status

### Deprecated

### Security


## [Release 0.2.1]
### Added

- Experimental client for interacting with backup service
- The operator panics itself when it gets stuck unexpectedly. It relies on Kubernetes to
get restarted.
- Add resource requirements field in `Spec.Pod` . Users can specify resource requirements for the
etcd container with this new field.
- Add status endpoint to backup sidecar service.
- Service account of the etcd operator pod is passed to backup pod.
- Add backup service status into cluster status.

### Changed

- Example deployments pin to the released version of the operator image
- Downward API of pod's namespace and name is required to start etcd operator pod.
  See `example/deployment.yaml` .

### Removed

- Drop etcd operator command line flags: "masterHost", "cert-file", "key-file", "ca-file".

### Fixed

### Deprecated

### Security
