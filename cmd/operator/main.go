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

package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/on2itsecurity/etcd-operator/pkg/chaos"
	"github.com/on2itsecurity/etcd-operator/pkg/client"
	"github.com/on2itsecurity/etcd-operator/pkg/controller"
	"github.com/on2itsecurity/etcd-operator/pkg/util/constants"
	"github.com/on2itsecurity/etcd-operator/pkg/util/k8sutil"
	"github.com/on2itsecurity/etcd-operator/pkg/util/probe"
	"github.com/on2itsecurity/etcd-operator/pkg/util/retryutil"
	"github.com/on2itsecurity/etcd-operator/version"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	_ "github.com/KimMachineGun/automemlimit"
	"github.com/sirupsen/logrus"
	_ "go.uber.org/automaxprocs"

	"golang.org/x/time/rate"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/leaderelection"
	"k8s.io/client-go/tools/leaderelection/resourcelock"
)

var (
	namespace  string
	name       string
	listenAddr string
	gcInterval time.Duration

	chaosLevel int

	printVersion bool

	createCRD bool

	clusterWide bool
)

func init() {
	flag.StringVar(&listenAddr, "listen-addr", "0.0.0.0:8080", "The address on which the HTTP server will listen to")
	// chaos level will be removed once we have a formal tool to inject failures.
	flag.IntVar(&chaosLevel, "chaos-level", -1, "DO NOT USE IN PRODUCTION - level of chaos injected into the etcd clusters created by the operator.")
	flag.BoolVar(&printVersion, "version", false, "Show version and quit")
	flag.BoolVar(&createCRD, "create-crd", true, "The operator will not create the EtcdCluster CRD when this flag is set to false.")
	flag.DurationVar(&gcInterval, "gc-interval", 10*time.Minute, "GC interval")
	flag.BoolVar(&clusterWide, "cluster-wide", false, "Enable operator to watch clusters in all namespaces")
	flag.Parse()
}

func main() {
	namespace = os.Getenv(constants.EnvOperatorPodNamespace)
	if len(namespace) == 0 {
		logrus.Fatalf("must set env (%s)", constants.EnvOperatorPodNamespace)
	}
	name = os.Getenv(constants.EnvOperatorPodName)
	if len(name) == 0 {
		logrus.Fatalf("must set env (%s)", constants.EnvOperatorPodName)
	}

	if printVersion {
		fmt.Println("etcd-operator Version:", version.Version)
		fmt.Println("Git SHA:", version.GitSHA)
		fmt.Println("Go Version:", runtime.Version())
		fmt.Printf("Go OS/Arch: %s/%s\n", runtime.GOOS, runtime.GOARCH)
		os.Exit(0)
	}

	logrus.Infof("etcd-operator Version: %v", version.Version)
	logrus.Infof("Git SHA: %s", version.GitSHA)
	logrus.Infof("Go Version: %s", runtime.Version())
	logrus.Infof("Go OS/Arch: %s/%s", runtime.GOOS, runtime.GOARCH)

	id, err := os.Hostname()
	if err != nil {
		logrus.Fatalf("failed to get hostname: %v", err)
	}
	if createCRD {
		logrus.Infof("NOTICE: create-crd is set to true by default, this is going to change in a future versions and create-crd will be deprecated. Please use the crds in example/crd.")
	}
	kubecli := k8sutil.MustNewKubeClient()

	http.HandleFunc(probe.HTTPReadyzEndpoint, probe.ReadyzHandler)
	http.Handle("/metrics", promhttp.Handler())
	go http.ListenAndServe(listenAddr, nil)

	rl := &resourcelock.LeaseLock{
		LeaseMeta: metav1.ObjectMeta{
			Name:      "etcd-operator",
			Namespace: namespace,
		},
		Client: kubecli.CoordinationV1(),
		LockConfig: resourcelock.ResourceLockConfig{
			Identity: id,
		},
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	probe.SetReady()
	defer probe.UnsetReady()

	leaderelection.RunOrDie(ctx, leaderelection.LeaderElectionConfig{
		Lock:          rl,
		LeaseDuration: 15 * time.Second,
		RenewDeadline: 10 * time.Second,
		RetryPeriod:   2 * time.Second,
		Callbacks: leaderelection.LeaderCallbacks{
			OnStartedLeading: run,
			OnStoppedLeading: func() {
				logrus.Fatalf("leader election lost")
			},
		},
	})

	panic("unreachable")
}

func run(ctx context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	cfg := newControllerConfig(ctx)

	startChaos(context.Background(), cfg.KubeCli, cfg.Namespace, chaosLevel)

	c := controller.New(cfg)
	err := c.Start()
	logrus.Fatalf("controller Start() failed: %v", err)
}

func newControllerConfig(ctx context.Context) controller.Config {
	kubecli := k8sutil.MustNewKubeClient()

	serviceAccount, err := getMyPodServiceAccount(ctx, kubecli)
	if err != nil {
		logrus.Fatalf("fail to get my pod's service account: %v", err)
	}

	cfg := controller.Config{
		Namespace:      namespace,
		ClusterWide:    clusterWide,
		ServiceAccount: serviceAccount,
		KubeCli:        kubecli,
		KubeExtCli:     k8sutil.MustNewKubeExtClient(),
		EtcdCRCli:      client.MustNewInCluster(),
		CreateCRD:      createCRD,
	}

	return cfg
}

func getMyPodServiceAccount(ctx context.Context, kubecli kubernetes.Interface) (string, error) {
	var sa string
	err := retryutil.Retry(5*time.Second, 100, func() (bool, error) {
		pod, err := kubecli.CoreV1().Pods(namespace).Get(ctx, name, metav1.GetOptions{})
		if err != nil {
			logrus.Errorf("fail to get operator pod (%s): %v", name, err)
			return false, nil
		}
		sa = pod.Spec.ServiceAccountName
		return true, nil
	})
	return sa, err
}

func startChaos(ctx context.Context, kubecli kubernetes.Interface, ns string, chaosLevel int) {
	m := chaos.NewMonkeys(kubecli)
	ls := labels.SelectorFromSet(map[string]string{"app": "etcd"})

	switch chaosLevel {
	case 1:
		logrus.Info("chaos level = 1: randomly kill one etcd pod every 30 seconds at 50%")
		c := &chaos.CrashConfig{
			Namespace: ns,
			Selector:  ls,

			KillRate:        rate.Every(30 * time.Second),
			KillProbability: 0.5,
			KillMax:         1,
		}
		go func() {
			time.Sleep(60 * time.Second) // don't start until quorum up
			m.CrushPods(ctx, c)
		}()

	case 2:
		logrus.Info("chaos level = 2: randomly kill at most five etcd pods every 30 seconds at 50%")
		c := &chaos.CrashConfig{
			Namespace: ns,
			Selector:  ls,

			KillRate:        rate.Every(30 * time.Second),
			KillProbability: 0.5,
			KillMax:         5,
		}

		go m.CrushPods(ctx, c)

	default:
	}
}
