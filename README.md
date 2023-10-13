# go-blob-repository

This library provides an implementation of a Blob repository that uses Azure
Storage underneath. The repository component is built in a way that allows
easily configuring it with the right Azure credentials. It is particularly
useful if you are integrating into an app/library using the clean/hexagonal
architecture as it abstracts Azure's code behing a repository.

To instantiate the blob repository, use the NewBlobRepository function. The
account name, API key, and bucket are configurable.

```go
    repo := NewBlobRepository(models.BlobRepositoryConfiguration{
        BlobAccessConfiguration: models.BlobAccessConfiguration{
            AccountName: "<test-account>",
            APIKey: "<api-key>",
        },
        Bucket: "<container-name>",
    }).(defaultBlobRepository)
```

Once the repository is instantiated, you can call the `Upload` function to
upload a document into the storage:

```go
    var content io.Reader

    meta := models.BlobMetadata{
        ContentType: "image/png",
        Name:        "some.picture",
    }

    uploadedURL, err := u.blobRepo.Upload(meta, content)
```
