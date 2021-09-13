package azure

import (
	"context"
	"fmt"
	"io"
	"mime"
	"net/url"

	"github.com/kumojin/go-blob-repository/models"
	"github.com/kumojin/go-blob-repository/repository"

	"github.com/Azure/azure-storage-blob-go/azblob"
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
	blobURL, err := r.GetBlobURL(fileName)
	if err != nil {
		return nil, err
	}

	cdnURL, err := r.GetCDNURL(fileName)
	if err != nil {
		return nil, err
	}

	credentials, err := azblob.NewSharedKeyCredential(r.config.AccountName, r.config.APIKey)
	if err != nil {
		return nil, err
	}

	// Prepare objects: block URL (for uploading pieces one by one), context and upload options
	blockBlobURL := azblob.NewBlockBlobURL(*blobURL, azblob.NewPipeline(credentials, azblob.PipelineOptions{}))
	o := azblob.UploadStreamToBlockBlobOptions{
		BlobHTTPHeaders: azblob.BlobHTTPHeaders{
			ContentType:  metadata.ContentType,
			CacheControl: "no-cache",
		},
		BufferSize: 2 * 1024 * 1024, // Size of the rotating buffers that are used when uploading
		MaxBuffers: 3,               // Number of rotating buffers that are used when uploading
	}

	ctx := context.Background()
	_, err = azblob.UploadStreamToBlockBlob(ctx, in, blockBlobURL, o)
	urls := models.UploadURL{
		BlobURL: *blobURL,
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
