// Copyright 2016 The etcd-operator Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package k8sutil

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"time"

	api "github.com/on2itsecurity/etcd-operator/pkg/apis/etcd/v1beta2"
	"github.com/on2itsecurity/etcd-operator/pkg/util/etcdutil"
	"github.com/on2itsecurity/etcd-operator/pkg/util/retryutil"
	"github.com/pborman/uuid"

	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	policyv1beta1 "k8s.io/api/policy/v1beta1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	utilrand "k8s.io/apimachinery/pkg/util/rand"
	"k8s.io/apimachinery/pkg/util/strategicpatch"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp" // for gcp auth
	"k8s.io/client-go/rest"
	"k8s.io/client-go/util/retry"
)

const (
	// EtcdClientPort is the client port on client service and etcd nodes.
	EtcdClientPort = 2379

	etcdVolumeMountDir       = "/var/etcd"
	dataDir                  = etcdVolumeMountDir + "/data"
	backupFile               = "/var/etcd/latest.backup"
	etcdVersionAnnotationKey = "etcd.version"
	peerTLSDir               = "/etc/etcdtls/member/peer-tls"
	peerTLSVolume            = "member-peer-tls"
	serverTLSDir             = "/etc/etcdtls/member/server-tls"
	serverTLSVolume          = "member-server-tls"
	operatorEtcdTLSDir       = "/etc/etcdtls/operator/etcd-tls"
	operatorEtcdTLSVolume    = "etcd-client-tls"

	randomSuffixLength = 10
	// k8s object name has a maximum length
	MaxNameLength = 63 - randomSuffixLength - 1

	defaultBusyboxImage = "busybox:1.28.0-glibc"

	// AnnotationScope annotation name for defining instance scope. Used for specifying cluster wide clusters.
	AnnotationScope = "etcd.database.coreos.com/scope"
	//AnnotationClusterWide annotation value for cluster wide clusters.
	AnnotationClusterWide = "clusterwide"

	// defaultDNSTimeout is the default maximum allowed time for the init container of the etcd pod
	// to reverse DNS lookup its IP. The default behavior is to wait forever and has a value of 0.
	defaultDNSTimeout = int64(0)
)

func GetEtcdVersion(pod *v1.Pod) string {
	return pod.Annotations[etcdVersionAnnotationKey]
}

func SetEtcdVersion(pod *v1.Pod, version string) {
	pod.Annotations[etcdVersionAnnotationKey] = version
}

func GetPodNames(pods []*v1.Pod) []string {
	if len(pods) == 0 {
		return nil
	}
	res := []string{}
	for _, p := range pods {
		res = append(res, p.Name)
	}
	return res
}

// PVCNameFromMember the way we get PVC name from the member name
func PVCNameFromMember(memberName string) string {
	return memberName
}

func makeRestoreInitContainers(backupURL *url.URL, token, repo, version string, m *etcdutil.Member) []v1.Container {
	return []v1.Container{
		{
			Name:  "fetch-backup",
			Image: "curlimages/curl",
			Command: []string{
				"/bin/ash", "-ec",
				fmt.Sprintf(`
httpcode=$(curl --write-out %%\{http_code\} --silent --output %[1]s %[2]s)
if [[ "$httpcode" != "200" ]]; then
	echo "http status code: ${httpcode}" >> /dev/termination-log
	cat %[1]s >> /dev/termination-log
	exit 1
fi
					`, backupFile, backupURL.String()),
			},
			VolumeMounts: etcdVolumeMounts(),
		},
		{
			Name:  "restore-datadir",
			Image: ImageName(repo, version),
			Command: []string{
				"/bin/sh", "-ec",
				fmt.Sprintf("ETCDCTL_API=3 etcdctl snapshot restore %[1]s"+
					" --name %[2]s"+
					" --initial-cluster %[2]s=%[3]s"+
					" --initial-cluster-token %[4]s"+
					" --initial-advertise-peer-urls %[3]s"+
					" --data-dir %[5]s 2>/dev/termination-log", backupFile, m.Name, m.PeerURL(), token, dataDir),
			},
			VolumeMounts: etcdVolumeMounts(),
		},
	}
}

func ImageName(repo, version string) string {
	return fmt.Sprintf("%s:%s", repo, version)
}

// imageNameBusybox returns the default image for busybox init container, or the image specified in the PodPolicy
func imageNameBusybox(policy *api.PodPolicy) string {
	if policy != nil && len(policy.BusyboxImage) > 0 {
		return policy.BusyboxImage
	}
	return defaultBusyboxImage
}

func PodWithNodeSelector(p *v1.Pod, ns map[string]string) *v1.Pod {
	p.Spec.NodeSelector = ns
	return p
}

func CreateClientService(ctx context.Context, kubecli kubernetes.Interface, clusterName, ns string, owner metav1.OwnerReference, tls bool, policy *api.ServicePolicy) error {

	var EtcdClientPortName string
	if tls {
		EtcdClientPortName = "https-client"
	} else {
		EtcdClientPortName = "http-client"
	}

	ports := []v1.ServicePort{{
		Name:       EtcdClientPortName,
		Port:       EtcdClientPort,
		TargetPort: intstr.FromInt(EtcdClientPort),
		Protocol:   v1.ProtocolTCP,
	}}
	return createService(ctx, kubecli, ClientServiceName(clusterName, policy), clusterName, ns, "", ports, owner, false, policy)
}

func ClientServiceName(clusterName string, policy *api.ServicePolicy) string {
	if policy != nil && len(policy.Name) > 0 {
		return policy.Name
	}
	return clusterName + "-client"
}

func CreatePeerService(ctx context.Context, kubecli kubernetes.Interface, clusterName, ns string, owner metav1.OwnerReference, tls bool, policy *api.ServicePolicy) error {

	var EtcdClientPortName string
	if tls {
		EtcdClientPortName = "https-client"
	} else {
		EtcdClientPortName = "http-client"
	}

	ports := []v1.ServicePort{{
		Name:       EtcdClientPortName,
		Port:       EtcdClientPort,
		TargetPort: intstr.FromInt(EtcdClientPort),
		Protocol:   v1.ProtocolTCP,
	}, {
		Name:       "peer",
		Port:       2380,
		TargetPort: intstr.FromInt(2380),
		Protocol:   v1.ProtocolTCP,
	}}

	return createService(ctx, kubecli, clusterName, clusterName, ns, v1.ClusterIPNone, ports, owner, true, nil)
}

func createService(ctx context.Context, kubecli kubernetes.Interface, svcName, clusterName, ns, clusterIP string, ports []v1.ServicePort, owner metav1.OwnerReference, publishNotReadyAddresses bool, policy *api.ServicePolicy) error {
	svc := newEtcdServiceManifest(svcName, clusterName, clusterIP, ports, publishNotReadyAddresses)

	applyServicePolicy(svc, policy)
	addOwnerRefToObject(svc.GetObjectMeta(), owner)
	_, err := kubecli.CoreV1().Services(ns).Create(ctx, svc, metav1.CreateOptions{})
	if err != nil && !apierrors.IsAlreadyExists(err) {
		return err
	}
	return nil
}

// CreateAndWaitPod creates a pod and waits until it is running
func CreateAndWaitPod(ctx context.Context, kubecli kubernetes.Interface, ns string, pod *v1.Pod, timeout time.Duration) (*v1.Pod, error) {
	_, err := kubecli.CoreV1().Pods(ns).Create(ctx, pod, metav1.CreateOptions{})
	if err != nil {
		return nil, err
	}

	interval := 5 * time.Second
	var retPod *v1.Pod
	err = retryutil.Retry(interval, int(timeout/(interval)), func() (bool, error) {
		retPod, err = kubecli.CoreV1().Pods(ns).Get(ctx, pod.Name, metav1.GetOptions{})
		if err != nil {
			return false, err
		}
		switch retPod.Status.Phase {
		case v1.PodRunning:
			return true, nil
		case v1.PodPending:
			return false, nil
		default:
			return false, fmt.Errorf("unexpected pod status.phase: %v", retPod.Status.Phase)
		}
	})

	if err != nil {
		if retryutil.IsRetryFailure(err) {
			return nil, fmt.Errorf("failed to wait pod running, it is still pending: %v", err)
		}
		return nil, fmt.Errorf("failed to wait pod running: %v", err)
	}

	return retPod, nil
}

func newEtcdServiceManifest(svcName, clusterName, clusterIP string, ports []v1.ServicePort, publishNotReadyAddresses bool) *v1.Service {
	labels := LabelsForCluster(clusterName)
	svc := &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:   svcName,
			Labels: labels,
		},
		Spec: v1.ServiceSpec{
			Ports:                    ports,
			Selector:                 labels,
			ClusterIP:                clusterIP,
			PublishNotReadyAddresses: publishNotReadyAddresses,
		},
	}
	return svc
}

func newEtcdPodDisruptionBudgetManifest(pdbName, clusterName string, minAvailable *intstr.IntOrString) *policyv1beta1.PodDisruptionBudget {
	labels := LabelsForCluster(clusterName)

	pdb := &policyv1beta1.PodDisruptionBudget{
		ObjectMeta: metav1.ObjectMeta{
			Name:   pdbName,
			Labels: labels,
		},
		Spec: policyv1beta1.PodDisruptionBudgetSpec{
			MinAvailable: minAvailable,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
		},
	}

	return pdb
}

// AddEtcdVolumeToPod abstract the process of appending volume spec to pod spec
func AddEtcdVolumeToPod(pod *v1.Pod, pvc *v1.PersistentVolumeClaim, tmpfs bool) {
	vol := v1.Volume{Name: etcdVolumeName}
	if pvc != nil {
		vol.VolumeSource = v1.VolumeSource{
			PersistentVolumeClaim: &v1.PersistentVolumeClaimVolumeSource{ClaimName: pvc.Name},
		}
		// Addresses case C from https://github.com/on2itsecurity/etcd-operator/blob/master/doc/design/persistent_volumes_etcd_data.md
		// When PVC is used, make the pod auto recover in case of failure
		pod.Spec.RestartPolicy = v1.RestartPolicyAlways
	} else {
		vol.VolumeSource = v1.VolumeSource{EmptyDir: &v1.EmptyDirVolumeSource{}}
		if tmpfs {
			vol.VolumeSource.EmptyDir.Medium = v1.StorageMediumMemory
		}
	}
	pod.Spec.Volumes = append(pod.Spec.Volumes, vol)
}

func addRecoveryToPod(pod *v1.Pod, token string, m *etcdutil.Member, cs api.ClusterSpec, backupURL *url.URL) {
	pod.Spec.InitContainers = append(pod.Spec.InitContainers,
		makeRestoreInitContainers(backupURL, token, cs.Repository, cs.Version, m)...)
}

func addOwnerRefToObject(o metav1.Object, r metav1.OwnerReference) {
	o.SetOwnerReferences(append(o.GetOwnerReferences(), r))
}

// NewSeedMemberPod returns a Pod manifest for a seed member.
// It's special that it has new token, and might need recovery init containers
func NewSeedMemberPod(ctx context.Context, kubecli kubernetes.Interface, clusterName, clusterNamespace string, ms etcdutil.MemberSet, m *etcdutil.Member, cs api.ClusterSpec, owner metav1.OwnerReference, backupURL *url.URL) (*v1.Pod, error) {
	token := uuid.New()
	pod, err := newEtcdPod(ctx, kubecli, m, ms.PeerURLPairs(), clusterName, clusterNamespace, "new", token, cs)
	if err != nil {
		return pod, err
	}
	// TODO: PVC datadir support for restore process
	AddEtcdVolumeToPod(pod, nil, cs.Pod.Tmpfs)
	if backupURL != nil {
		addRecoveryToPod(pod, token, m, cs, backupURL)
	}
	applyPodPolicy(clusterName, pod, cs.Pod)
	addOwnerRefToObject(pod.GetObjectMeta(), owner)
	return pod, err
}

// NewEtcdPodPVC create PVC object from etcd pod's PVC spec
func NewEtcdPodPVC(m *etcdutil.Member, pvcSpec v1.PersistentVolumeClaimSpec, clusterName, namespace string, owner metav1.OwnerReference) *v1.PersistentVolumeClaim {
	pvc := &v1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      PVCNameFromMember(m.Name),
			Namespace: namespace,
			Labels:    LabelsForCluster(clusterName),
		},
		Spec: pvcSpec,
	}
	addOwnerRefToObject(pvc.GetObjectMeta(), owner)
	return pvc
}

func newEtcdPod(ctx context.Context, kubecli kubernetes.Interface, m *etcdutil.Member, initialCluster []string, clusterName, clusterNamespace, state, token string, cs api.ClusterSpec) (*v1.Pod, error) {
	commands := fmt.Sprintf("/usr/local/bin/etcd --data-dir=%s --name=%s --initial-advertise-peer-urls=%s "+
		"--listen-peer-urls=%s --listen-client-urls=%s --advertise-client-urls=%s "+
		"--initial-cluster=%s --initial-cluster-state=%s",
		dataDir, m.Name, m.PeerURL(), m.ListenPeerURL(), m.ListenClientURL(), m.ClientURL(), strings.Join(initialCluster, ","), state)
	if m.SecurePeer {
		secret, err := kubecli.CoreV1().Secrets(clusterNamespace).Get(ctx, cs.TLS.Static.Member.PeerSecret, metav1.GetOptions{})
		if err != nil {
			return nil, err
		}
		if secret.Type == v1.SecretTypeTLS {
			commands += fmt.Sprintf(" --peer-client-cert-auth=true --peer-trusted-ca-file=%[1]s/ca.crt --peer-cert-file=%[1]s/tls.crt --peer-key-file=%[1]s/tls.key", peerTLSDir)
		} else {
			commands += fmt.Sprintf(" --peer-client-cert-auth=true --peer-trusted-ca-file=%[1]s/peer-ca.crt --peer-cert-file=%[1]s/peer.crt --peer-key-file=%[1]s/peer.key", peerTLSDir)
		}
	}
	if m.SecureClient {
		secret, err := kubecli.CoreV1().Secrets(clusterNamespace).Get(ctx, cs.TLS.Static.Member.ServerSecret, metav1.GetOptions{})
		if err != nil {
			return nil, err
		}
		if secret.Type == v1.SecretTypeTLS {
			commands += fmt.Sprintf(" --client-cert-auth=true --trusted-ca-file=%[1]s/ca.crt --cert-file=%[1]s/tls.crt --key-file=%[1]s/tls.key", serverTLSDir)
		} else {
			commands += fmt.Sprintf(" --client-cert-auth=true --trusted-ca-file=%[1]s/server-ca.crt --cert-file=%[1]s/server.crt --key-file=%[1]s/server.key", serverTLSDir)
		}
	}
	if state == "new" {
		commands = fmt.Sprintf("%s --initial-cluster-token=%s", commands, token)
	}

	labels := map[string]string{
		"app":          "etcd",
		"etcd_node":    m.Name,
		"etcd_cluster": clusterName,
	}

	isTLSSecret := false
	if cs.TLS.IsSecureClient() {
		secret, err := kubecli.CoreV1().Secrets(clusterNamespace).Get(ctx, cs.TLS.Static.OperatorSecret, metav1.GetOptions{})
		if err != nil {
			return nil, err
		}
		isTLSSecret = secret.Type == v1.SecretTypeTLS
	}
	livenessProbe := newEtcdProbe(cs.TLS.IsSecureClient(), isTLSSecret)
	readinessProbe := newEtcdProbe(cs.TLS.IsSecureClient(), isTLSSecret)
	readinessProbe.InitialDelaySeconds = 1
	readinessProbe.InitialDelaySeconds = 1
	readinessProbe.TimeoutSeconds = 5
	readinessProbe.PeriodSeconds = 5
	readinessProbe.FailureThreshold = 3

	container := containerWithProbes(
		etcdContainer(strings.Split(commands, " "), cs.Repository, cs.Version),
		livenessProbe,
		readinessProbe)

	volumes := []v1.Volume{}

	if m.SecurePeer {
		container.VolumeMounts = append(container.VolumeMounts, v1.VolumeMount{
			MountPath: peerTLSDir,
			Name:      peerTLSVolume,
		})
		volumes = append(volumes, v1.Volume{Name: peerTLSVolume, VolumeSource: v1.VolumeSource{
			Secret: &v1.SecretVolumeSource{SecretName: cs.TLS.Static.Member.PeerSecret},
		}})
	}
	if m.SecureClient {
		container.VolumeMounts = append(container.VolumeMounts, v1.VolumeMount{
			MountPath: serverTLSDir,
			Name:      serverTLSVolume,
		}, v1.VolumeMount{
			MountPath: operatorEtcdTLSDir,
			Name:      operatorEtcdTLSVolume,
		})
		volumes = append(volumes, v1.Volume{Name: serverTLSVolume, VolumeSource: v1.VolumeSource{
			Secret: &v1.SecretVolumeSource{SecretName: cs.TLS.Static.Member.ServerSecret},
		}}, v1.Volume{Name: operatorEtcdTLSVolume, VolumeSource: v1.VolumeSource{
			Secret: &v1.SecretVolumeSource{SecretName: cs.TLS.Static.OperatorSecret},
		}})
	}

	DNSTimeout := defaultDNSTimeout
	if cs.Pod != nil {
		DNSTimeout = cs.Pod.DNSTimeoutInSecond
	}
	pod := &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:        m.Name,
			Labels:      labels,
			Annotations: map[string]string{},
		},
		Spec: v1.PodSpec{
			InitContainers: []v1.Container{{
				// busybox:latest uses uclibc which contains a bug that sometimes prevents name resolution
				// More info: https://github.com/docker-library/busybox/issues/27
				//Image default: "busybox:1.28.0-glibc",
				Image: imageNameBusybox(cs.Pod),
				Name:  "check-dns",
				// In etcd 3.2, TLS listener will do a reverse-DNS lookup for pod IP -> hostname.
				// If DNS entry is not warmed up, it will return empty result and peer connection will be rejected.
				// In some cases the DNS is not created correctly so we need to time out after a given period.
				Command: []string{"/bin/sh", "-c", fmt.Sprintf(`
					TIMEOUT_READY=%d
					while ( ! nslookup %s )
					do
						# If TIMEOUT_READY is 0 we should never time out and exit
						TIMEOUT_READY=$(( TIMEOUT_READY-1 ))
                        if [ $TIMEOUT_READY -eq 0 ];
				        then
				            echo "Timed out waiting for DNS entry"
				            exit 1
				        fi
						sleep 1
					done`, DNSTimeout, m.Addr())},
			}},
			Containers:    []v1.Container{container},
			RestartPolicy: v1.RestartPolicyNever,
			Volumes:       volumes,
			// DNS A record: `[m.Name].[clusterName].Namespace.svc`
			// For example, etcd-795649v9kq in default namespace will have DNS name
			// `etcd-795649v9kq.etcd.default.svc`.
			Hostname:                     m.Name,
			Subdomain:                    clusterName,
			AutomountServiceAccountToken: func(b bool) *bool { return &b }(false),
			SecurityContext:              podSecurityContext(cs.Pod),
		},
	}
	SetEtcdVersion(pod, cs.Version)
	return pod, nil
}

func podSecurityContext(podPolicy *api.PodPolicy) *v1.PodSecurityContext {
	if podPolicy == nil {
		return nil
	}
	return podPolicy.SecurityContext
}

func NewEtcdPod(ctx context.Context, kubecli kubernetes.Interface, m *etcdutil.Member, initialCluster []string, clusterName, clusterNamespace, state, token string, cs api.ClusterSpec, owner metav1.OwnerReference) (*v1.Pod, error) {
	pod, err := newEtcdPod(ctx, kubecli, m, initialCluster, clusterName, clusterNamespace, state, token, cs)
	applyPodPolicy(clusterName, pod, cs.Pod)
	addOwnerRefToObject(pod.GetObjectMeta(), owner)
	return pod, err
}

func MustNewKubeClient() kubernetes.Interface {
	cfg, err := InClusterConfig()
	if err != nil {
		panic(err)
	}
	return kubernetes.NewForConfigOrDie(cfg)
}

func InClusterConfig() (*rest.Config, error) {
	cfg, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

func IsKubernetesResourceAlreadyExistError(err error) bool {
	return apierrors.IsAlreadyExists(err)
}

func IsKubernetesResourceNotFoundError(err error) bool {
	return apierrors.IsNotFound(err)
}

// We are using internal api types for cluster related.
func ClusterListOpt(clusterName string) metav1.ListOptions {
	return metav1.ListOptions{
		LabelSelector: labels.SelectorFromSet(LabelsForCluster(clusterName)).String(),
	}
}

func LabelsForCluster(clusterName string) map[string]string {
	return map[string]string{
		"etcd_cluster": clusterName,
		"app":          "etcd",
	}
}

func CreatePatch(o, n, datastruct interface{}) ([]byte, error) {
	oldData, err := json.Marshal(o)
	if err != nil {
		return nil, err
	}
	newData, err := json.Marshal(n)
	if err != nil {
		return nil, err
	}
	return strategicpatch.CreateTwoWayMergePatch(oldData, newData, datastruct)
}

func PatchDeployment(ctx context.Context, kubecli kubernetes.Interface, namespace, name string, updateFunc func(*appsv1.Deployment)) error {
	od, err := kubecli.AppsV1().Deployments(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return err
	}
	nd := od.DeepCopy()
	updateFunc(nd)
	patchData, err := CreatePatch(od, nd, appsv1.Deployment{})
	if err != nil {
		return err
	}
	_, err = kubecli.AppsV1().Deployments(namespace).Patch(ctx, name, types.StrategicMergePatchType, patchData, metav1.PatchOptions{})
	return err
}

func UpdateOrCreatePodDisruptionBudget(ctx context.Context, kubecli kubernetes.Interface, namespace, name string, minAvailable int, owner metav1.OwnerReference) error {
	pdb := newEtcdPodDisruptionBudgetManifest(name, name, &intstr.IntOrString{IntVal: int32(minAvailable)})
	addOwnerRefToObject(pdb.GetObjectMeta(), owner)

	return retry.RetryOnConflict(retry.DefaultRetry, func() error {
		od, err := kubecli.PolicyV1beta1().PodDisruptionBudgets(namespace).Get(ctx, name, metav1.GetOptions{})
		if IsKubernetesResourceNotFoundError(err) {
			_, err := kubecli.PolicyV1beta1().PodDisruptionBudgets(namespace).Create(ctx, pdb, metav1.CreateOptions{})
			return err
		}
		if err != nil {
			return err
		}

		pdb.Status = od.Status
		pdb.ObjectMeta.ResourceVersion = od.ObjectMeta.ResourceVersion // Optimistic locking
		_, err = kubecli.PolicyV1beta1().PodDisruptionBudgets(namespace).Update(ctx, pdb, metav1.UpdateOptions{})

		return err
	})
}

func CascadeDeleteOptions(gracePeriodSeconds int64) *metav1.DeleteOptions {
	return &metav1.DeleteOptions{
		GracePeriodSeconds: func(t int64) *int64 { return &t }(gracePeriodSeconds),
		PropagationPolicy: func() *metav1.DeletionPropagation {
			foreground := metav1.DeletePropagationForeground
			return &foreground
		}(),
	}
}

// mergeLabels merges l2 into l1. Conflicting label will be skipped.
func mergeLabels(l1, l2 map[string]string) {
	for k, v := range l2 {
		if _, ok := l1[k]; ok {
			continue
		}
		l1[k] = v
	}
}

func UniqueMemberName(clusterName string) string {
	suffix := utilrand.String(randomSuffixLength)
	if len(clusterName) > MaxNameLength {
		clusterName = clusterName[:MaxNameLength]
	}
	return clusterName + "-" + suffix
}
