package repository

import (
	"io"

	"github.com/kumojin/go-blob-repository/models"
)

// BlobRepository is for interacting with blobs (files, images, etc.)
type BlobRepository interface {
	Upload(metadata models.BlobMetadata, in io.Reader) (*models.UploadURL, error)
}