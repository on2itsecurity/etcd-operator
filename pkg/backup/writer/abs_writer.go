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

package writer

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/service"
	"github.com/on2itsecurity/etcd-operator/pkg/backup/util"
	"github.com/pborman/uuid"
)

type absWriter struct {
	client *service.Client
}

var _ Writer = &absWriter{}

type ABSClient struct {
	Client *service.Client
}

func NewABSWriter(client *service.Client) Writer {
	return &absWriter{client}
}

const AzureBlobBlockChunkLimitInBytes = 100 * 1024 * 1024

func (absw *absWriter) Write(ctx context.Context, path string, r io.Reader) (int64, error) {
	container, key, err := util.ParseBucketAndKey(path)
	if err != nil {
		return 0, err
	}

	containerClient := absw.client.NewContainerClient(container)
	if _, err = containerClient.GetProperties(ctx, nil); err != nil {
		return 0, fmt.Errorf("container %v does not exist or is inaccessible: %v", container, err)
	}

	blobClient := containerClient.NewBlockBlobClient(key)
	blockIDs := []string{}
	total := int64(0)
	buf := make([]byte, AzureBlobBlockChunkLimitInBytes)

	for {
		n, readErr := io.ReadFull(r, buf)
		if readErr != nil && readErr != io.ErrUnexpectedEOF && readErr != io.EOF {
			return 0, readErr
		}
		if n == 0 {
			break
		}

		blockID := base64.StdEncoding.EncodeToString([]byte(uuid.New()))
		_, err := blobClient.StageBlock(ctx, blockID, streamingBytes(buf[:n]), nil)
		if err != nil {
			return 0, err
		}
		blockIDs = append(blockIDs, blockID)
		total += int64(n)

		if readErr == io.ErrUnexpectedEOF || readErr == io.EOF {
			break
		}
	}

	_, err = blobClient.CommitBlockList(ctx, blockIDs, nil)
	if err != nil {
		return 0, err
	}

	props, err := blobClient.GetProperties(ctx, nil)
	if err != nil {
		return 0, err
	}
	return *props.ContentLength, nil
}

func (absw *absWriter) List(ctx context.Context, basePath string) ([]string, error) {
	container, prefix, err := util.ParseBucketAndKey(basePath)
	if err != nil {
		return nil, err
	}

	containerClient := absw.client.NewContainerClient(container)
	pager := containerClient.NewListBlobsFlatPager(&azblob.ListBlobsFlatOptions{
		Prefix: &prefix,
	})

	var blobKeys []string
	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		for _, blob := range page.Segment.BlobItems {
			blobKeys = append(blobKeys, container+"/"+*blob.Name)
		}
	}
	return blobKeys, nil
}

func (absw *absWriter) Delete(ctx context.Context, path string) error {
	container, key, err := util.ParseBucketAndKey(path)
	if err != nil {
		return err
	}

	containerClient := absw.client.NewContainerClient(container)
	blobClient := containerClient.NewBlockBlobClient(key)

	_, err = blobClient.Delete(ctx, nil)
	return err
}

func streamingBytes(b []byte) io.ReadSeekCloser {
	return &readSeekNopCloser{data: b}
}

type readSeekNopCloser struct {
	data []byte
	pos  int64
}

func (r *readSeekNopCloser) Read(p []byte) (int, error) {
	if r.pos >= int64(len(r.data)) {
		return 0, io.EOF
	}
	n := copy(p, r.data[r.pos:])
	r.pos += int64(n)
	return n, nil
}

func (r *readSeekNopCloser) Seek(offset int64, whence int) (int64, error) {
	var abs int64
	switch whence {
	case io.SeekStart:
		abs = offset
	case io.SeekCurrent:
		abs = r.pos + offset
	case io.SeekEnd:
		abs = int64(len(r.data)) + offset
	default:
		return 0, fmt.Errorf("invalid seek whence: %d", whence)
	}
	if abs < 0 {
		return 0, fmt.Errorf("negative position")
	}
	r.pos = abs
	return abs, nil
}

func (r *readSeekNopCloser) Close() error { return nil }
