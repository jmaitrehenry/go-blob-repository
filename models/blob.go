package models

import (
	"fmt"
	"net/url"
)

// BlobMetadata contains blob metadata
type BlobMetadata struct {
	ContentType string
	Name        string
}

// BlobAccessConfiguration contains blob access configuration
type BlobAccessConfiguration struct {
	APIKey      string
	AccountName string
}

// BlobRepositoryConfiguration contains blob repository configuration
type BlobRepositoryConfiguration struct {
	BlobAccessConfiguration
	Bucket string
}

// UploadURL contains both Blob and CDN URLs
// for a stored file.
type UploadURL struct {
	BlobURL url.URL
	CDNUrl  url.URL
}

// Endpoint returns the URL of the blob repository
func (brc BlobRepositoryConfiguration) Endpoint() string {
	return fmt.Sprintf("https://%s.blob.core.windows.net/", brc.AccountName)
}

// CDNEndpoint returns the CDN URL of the blob repository
func (brc BlobRepositoryConfiguration) CDNEndpoint() string {
	return fmt.Sprintf("https://%s.azureedge.net/", brc.AccountName)
}

