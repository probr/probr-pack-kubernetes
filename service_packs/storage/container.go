package storage

import (
	"context"
	"fmt"
	"net/url"

	"github.com/Azure/azure-storage-blob-go/azblob"
)

// CreateContainer creates a new container with the specified name in the specified account
func CreateContainer(ctx context.Context, name, rgName, containerName string) (azblob.ContainerURL, error) {
	u := getContainerURL(ctx, name, rgName, containerName)
	_, err := u.Create(
		ctx,
		azblob.Metadata{},
		azblob.PublicAccessContainer)
	return u, err
}

// GetContainer gets info about an existing container.
func GetContainer(ctx context.Context, accountName, accountGroupName, containerName string) (azblob.ContainerURL, error) {
	u := getContainerURL(ctx, accountName, accountGroupName, containerName)

	_, err := u.GetProperties(ctx, azblob.LeaseAccessConditions{})
	//TODO do we really want to return u, or the properties?
	return u, err
}

// DeleteContainer deletes the named container.
func DeleteContainer(ctx context.Context, accountName, accountGroupName, containerName string) error {
	u := getContainerURL(ctx, accountName, accountGroupName, containerName)

	_, err := u.Delete(ctx, azblob.ContainerAccessConditions{})
	return err
}

// ListBlobs lists blobs on the specified container
func ListBlobs(ctx context.Context, accountName, accountGroupName, containerName string) (*azblob.ListBlobsFlatSegmentResponse, error) {
	u := getContainerURL(ctx, accountName, accountGroupName, containerName)
	return u.ListBlobsFlatSegment(
		ctx,
		azblob.Marker{},
		azblob.ListBlobsSegmentOptions{
			Details: azblob.BlobListingDetails{
				Snapshots: true,
			},
		})
}

func getContainerURL(ctx context.Context, accountName, accountGroupName, containerName string) azblob.ContainerURL {
	key := AccountPrimaryKey(ctx, accountName, accountGroupName)
	creds, _ := azblob.NewSharedKeyCredential(accountName, key)
	p := azblob.NewPipeline(creds, azblob.PipelineOptions{})
	u, _ := url.Parse(fmt.Sprintf(`https://%s.blob.core.windows.net`, accountName))
	service := azblob.NewServiceURL(*u, p)
	return service.NewContainerURL(containerName)
}
