/*
Copyright 2021 The etcd-operator Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by client-gen. DO NOT EDIT.

package v1beta2

import (
	"context"
	"time"

	v1beta2 "github.com/on2itsecurity/etcd-operator/pkg/apis/etcd/v1beta2"
	scheme "github.com/on2itsecurity/etcd-operator/pkg/generated/clientset/versioned/scheme"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// EtcdClustersGetter has a method to return a EtcdClusterInterface.
// A group's client should implement this interface.
type EtcdClustersGetter interface {
	EtcdClusters(namespace string) EtcdClusterInterface
}

// EtcdClusterInterface has methods to work with EtcdCluster resources.
type EtcdClusterInterface interface {
	Create(ctx context.Context, etcdCluster *v1beta2.EtcdCluster, opts v1.CreateOptions) (*v1beta2.EtcdCluster, error)
	Update(ctx context.Context, etcdCluster *v1beta2.EtcdCluster, opts v1.UpdateOptions) (*v1beta2.EtcdCluster, error)
	UpdateStatus(ctx context.Context, etcdCluster *v1beta2.EtcdCluster, opts v1.UpdateOptions) (*v1beta2.EtcdCluster, error)
	Delete(ctx context.Context, name string, opts v1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error
	Get(ctx context.Context, name string, opts v1.GetOptions) (*v1beta2.EtcdCluster, error)
	List(ctx context.Context, opts v1.ListOptions) (*v1beta2.EtcdClusterList, error)
	Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1beta2.EtcdCluster, err error)
	EtcdClusterExpansion
}

// etcdClusters implements EtcdClusterInterface
type etcdClusters struct {
	client rest.Interface
	ns     string
}

// newEtcdClusters returns a EtcdClusters
func newEtcdClusters(c *EtcdV1beta2Client, namespace string) *etcdClusters {
	return &etcdClusters{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the etcdCluster, and returns the corresponding etcdCluster object, and an error if there is any.
func (c *etcdClusters) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1beta2.EtcdCluster, err error) {
	result = &v1beta2.EtcdCluster{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("etcdclusters").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of EtcdClusters that match those selectors.
func (c *etcdClusters) List(ctx context.Context, opts v1.ListOptions) (result *v1beta2.EtcdClusterList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1beta2.EtcdClusterList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("etcdclusters").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested etcdClusters.
func (c *etcdClusters) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("etcdclusters").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

// Create takes the representation of a etcdCluster and creates it.  Returns the server's representation of the etcdCluster, and an error, if there is any.
func (c *etcdClusters) Create(ctx context.Context, etcdCluster *v1beta2.EtcdCluster, opts v1.CreateOptions) (result *v1beta2.EtcdCluster, err error) {
	result = &v1beta2.EtcdCluster{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("etcdclusters").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(etcdCluster).
		Do(ctx).
		Into(result)
	return
}

// Update takes the representation of a etcdCluster and updates it. Returns the server's representation of the etcdCluster, and an error, if there is any.
func (c *etcdClusters) Update(ctx context.Context, etcdCluster *v1beta2.EtcdCluster, opts v1.UpdateOptions) (result *v1beta2.EtcdCluster, err error) {
	result = &v1beta2.EtcdCluster{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("etcdclusters").
		Name(etcdCluster.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(etcdCluster).
		Do(ctx).
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *etcdClusters) UpdateStatus(ctx context.Context, etcdCluster *v1beta2.EtcdCluster, opts v1.UpdateOptions) (result *v1beta2.EtcdCluster, err error) {
	result = &v1beta2.EtcdCluster{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("etcdclusters").
		Name(etcdCluster.Name).
		SubResource("status").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(etcdCluster).
		Do(ctx).
		Into(result)
	return
}

// Delete takes name of the etcdCluster and deletes it. Returns an error if one occurs.
func (c *etcdClusters) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("etcdclusters").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *etcdClusters) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	var timeout time.Duration
	if listOpts.TimeoutSeconds != nil {
		timeout = time.Duration(*listOpts.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Namespace(c.ns).
		Resource("etcdclusters").
		VersionedParams(&listOpts, scheme.ParameterCodec).
		Timeout(timeout).
		Body(&opts).
		Do(ctx).
		Error()
}

// Patch applies the patch and returns the patched etcdCluster.
func (c *etcdClusters) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1beta2.EtcdCluster, err error) {
	result = &v1beta2.EtcdCluster{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("etcdclusters").
		Name(name).
		SubResource(subresources...).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}
