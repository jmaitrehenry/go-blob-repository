package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBlobRepositoryConfiguration_Endpoint(t *testing.T) {
	blob := BlobRepositoryConfiguration{
		BlobAccessConfiguration: BlobAccessConfiguration{AccountName: "account"},
	}
	endpoint := blob.Endpoint()
	assert.Equal(t, "https://account.blob.core.windows.net/", endpoint)
}

func TestBlobRepositoryConfiguration_CDNEndpoint(t *testing.T) {
	blob := BlobRepositoryConfiguration{
		BlobAccessConfiguration: BlobAccessConfiguration{AccountName: "account"},
	}
	endpoint := blob.CDNEndpoint()
	assert.Equal(t, "https://account.azureedge.net/", endpoint)
}
