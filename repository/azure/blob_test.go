package azure

import (
	"testing"

	"github.com/kumojin/go-blob-repository/models"

	"github.com/stretchr/testify/assert"
)

func Test_defaultBlobRepository_getBlobURL(t *testing.T) {
	br := getBlobRepository()

	url, err := br.GetBlobURL("test-blob")
	assert.NoError(t, err)
	assert.Equal(t, "https://test-account.blob.core.windows.net/test-bucket/test-blob", url.String())
}

func Test_defaultBlobRepository_GetCDNURL(t *testing.T) {
	br := getBlobRepository()

	url, err := br.GetCDNURL("test-blob")
	assert.NoError(t, err)
	assert.Equal(t, "https://test-account.azureedge.net/test-bucket/test-blob", url.String())
}

func TestBuildFileNameFromMetadata(t *testing.T) {
	file := BuildFileNameFromMetadata(models.BlobMetadata{
		ContentType: "application/pdf",
		Name:        "file",
	})

	assert.Equal(t, "file.pdf", file)

	file = BuildFileNameFromMetadata(models.BlobMetadata{
		ContentType: "application/pdf",
		Name:        "file.pdf",
	})

	assert.Equal(t, "file.pdf.pdf", file)

	file = BuildFileNameFromMetadata(models.BlobMetadata{
		ContentType: "image/png",
		Name:        "file.pdf",
	})

	assert.Equal(t, "file.pdf.png", file)

	file = BuildFileNameFromMetadata(models.BlobMetadata{
		ContentType: "type/subtype",
		Name:        "file.unknown",
	})
	assert.Equal(t, "file.unknown", file)
}

func getBlobRepository() defaultBlobRepository {
	return NewBlobRepository(models.BlobRepositoryConfiguration{
		BlobAccessConfiguration: models.BlobAccessConfiguration{
			AccountName: "test-account",
		},
		Bucket: "test-bucket",
	}).(defaultBlobRepository)
}
