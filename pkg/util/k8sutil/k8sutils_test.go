// Copyright 2017 The etcd-operator Authors
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
	"strconv"
	"testing"

	api "github.com/on2itsecurity/etcd-operator/pkg/apis/etcd/v1beta2"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

func TestDefaultBusyboxImageName(t *testing.T) {
	policy := &api.PodPolicy{}
	image := imageNameBusybox(policy)
	expected := defaultBusyboxImage
	if image != expected {
		t.Errorf("expect image=%s, get=%s", expected, image)
	}
}

func TestDefaultNilBusyboxImageName(t *testing.T) {
	image := imageNameBusybox(nil)
	expected := defaultBusyboxImage
	if image != expected {
		t.Errorf("expect image=%s, get=%s", expected, image)
	}
}

func TestSetBusyboxImageName(t *testing.T) {
	policy := &api.PodPolicy{
		BusyboxImage: "myRepo/busybox:1.3.2",
	}
	image := imageNameBusybox(policy)
	expected := "myRepo/busybox:1.3.2"
	if image != expected {
		t.Errorf("expect image=%s, get=%s", expected, image)
	}
}

func TestClientServiceNameNilPolicy(t *testing.T) {
	clusterName := "clusterName"
	var policy *api.ServicePolicy = nil

	svcName := ClientServiceName(clusterName, policy)
	expected := "clusterName-client"
	if svcName != expected {
		t.Errorf("expect svcName=%s, got=%s", expected, svcName)
	}
}

func TestClientServiceNameWithEmptyNamePolicy(t *testing.T) {
	clusterName := "clusterName"
	policy := &api.ServicePolicy{
		Name: "",
	}
	svcName := ClientServiceName(clusterName, policy)
	expected := "clusterName-client"
	if svcName != expected {
		t.Errorf("expect svcName=%s, got=%s", expected, svcName)
	}
}

func TestClientServiceNameWithNamePolicy(t *testing.T) {
	clusterName := "clusterName"
	policy := &api.ServicePolicy{
		Name: "clusterNameCustom",
	}

	svcName := ClientServiceName(clusterName, policy)
	expected := policy.Name
	if svcName != expected {
		t.Errorf("expect svcName=%s, got=%s", expected, svcName)
	}
}

func TestGoGCEnvFromResources(t *testing.T) {
	resources := v1.ResourceRequirements{
		Limits: v1.ResourceList{
			v1.ResourceMemory: resource.MustParse("1Gi"),
		},
	}

	envVar := goGCEnvFromResources(resources)
	expected := v1.EnvVar{
		Name:  "GOMEMLIMIT",
		Value: strconv.FormatInt(1024*1024*1024*goGCMemLimitPercentage/100, 10),
	}
	if envVar != expected {
		t.Errorf("expect envVar=%v, got=%v", expected, envVar)
	}
}
func TestGoGCEnvFromResourcesNoLimit(t *testing.T) {
	resources := v1.ResourceRequirements{
		Requests: v1.ResourceList{
			v1.ResourceMemory: resource.MustParse("1Gi"),
		},
	}

	envVar := goGCEnvFromResources(resources)
	expected := v1.EnvVar{}
	if envVar != expected {
		t.Errorf("expect envVar=%v, got=%v", expected, envVar)
	}
}
