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

package absfactory

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/service"
	api "github.com/on2itsecurity/etcd-operator/pkg/apis/etcd/v1beta2"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type ABSClient struct {
	ServiceClient *service.Client
	BlobClient    *azblob.Client
}

func NewClientFromSecret(ctx context.Context, kubecli kubernetes.Interface, namespace, absSecret string) (*ABSClient, error) {
	se, err := kubecli.CoreV1().Secrets(namespace).Get(ctx, absSecret, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get k8s secret: %v", err)
	}

	accountName := string(se.Data[api.AzureSecretStorageAccount])
	accountKey := string(se.Data[api.AzureSecretStorageKey])
	cred, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create shared key credential: %v", err)
	}

	url := fmt.Sprintf("https://%s.blob.core.windows.net/", accountName)

	svcClient, err := service.NewClientWithSharedKeyCredential(url, cred, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create service client: %v", err)
	}

	blobClient, err := azblob.NewClientWithSharedKeyCredential(url, cred, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create blob client: %v", err)
	}

	return &ABSClient{
		ServiceClient: svcClient,
		BlobClient:    blobClient,
	}, nil

}

// ABSReader wraps an azblob.Client for reading operations.
type ABSReader struct {
	client *azblob.Client
}

func NewABSReader(client *azblob.Client) *ABSReader {
	return &ABSReader{client: client}
}
