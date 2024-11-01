package azure

import (
	"context"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blob"
	"io"
	"mime"
	"net/url"

	"github.com/kumojin/go-blob-repository/models"
	"github.com/kumojin/go-blob-repository/repository"
)

type defaultBlobRepository struct {
	config models.BlobRepositoryConfiguration
}

// NewBlobRepository returns a new instance of the blob repository
func NewBlobRepository(config models.BlobRepositoryConfiguration) repository.BlobRepository {
	return defaultBlobRepository{config: config}
}

func (r defaultBlobRepository) Upload(metadata models.BlobMetadata, in io.Reader) (*models.UploadURL, error) {
	fileName := BuildFileNameFromMetadata(metadata)

	cdnURL, err := r.GetCDNURL(fileName)
	if err != nil {
		return nil, err
	}

	credentials, err := azblob.NewSharedKeyCredential(r.config.AccountName, r.config.APIKey)
	if err != nil {
		return nil, err
	}

	client, err := azblob.NewClientWithSharedKeyCredential(r.config.Endpoint(), credentials, nil)
	if err != nil {
		return nil, err
	}

	noCache := "no-cache"
	_, err = client.UploadStream(context.Background(), r.config.Bucket, fileName, in, &azblob.UploadStreamOptions{
		BlockSize:   2 * 1024 * 1024,
		Concurrency: 3,
		HTTPHeaders: &blob.HTTPHeaders{
			BlobContentType:  &metadata.ContentType,
			BlobCacheControl: &noCache,
		},
	})
	if err != nil {
		return nil, err
	}

	blobUrl, _ := r.GetBlobURL(fileName)

	urls := models.UploadURL{
		BlobURL: *blobUrl,
		CDNUrl:  *cdnURL,
	}
	return &urls, err
}

func (r defaultBlobRepository) GetBlobURL(blobName string) (*url.URL, error) {
	return url.Parse(fmt.Sprint(r.config.Endpoint(), r.config.Bucket, "/", blobName))
}

func (r defaultBlobRepository) GetCDNURL(blobName string) (*url.URL, error) {
	return url.Parse(fmt.Sprint(r.config.CDNEndpoint(), r.config.Bucket, "/", blobName))
}

func BuildFileNameFromMetadata(metadata models.BlobMetadata) string {
	return fmt.Sprintf("%s%s", metadata.Name, getFirstKnownExtensionTypeOrFallbackToNothing(metadata.ContentType))
}

func getFirstKnownExtensionTypeOrFallbackToNothing(contentType string) string {
	types, err := mime.ExtensionsByType(contentType)
	if err != nil || len(types) <= 0 {
		return ""
	}

	return types[0]
}
