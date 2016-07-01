package uploader

import (
	"io"

	"github.com/dropbox/dropbox-sdk-go-unofficial"
	"github.com/dropbox/dropbox-sdk-go-unofficial/files"
	"github.com/dropbox/dropbox-sdk-go-unofficial/sharing"
)

// Uploader defines the interface for uploading a file
type Uploader interface {
	// Upload a file to dropbox at the given remote filepath
	Upload(filepath string, content io.Reader) (string, error)
}

// Client defines the interface of the client the Uploader will use
type Client interface {
	CreateSharedLinkWithSettings(arg *sharing.CreateSharedLinkWithSettingsArg) (res *sharing.SharedLinkMetadata, err error)
	Upload(arg *files.CommitInfo, content io.Reader) (res *files.FileMetadata, err error)
}

type dropBoxUploader struct {
	client Client
}

// New constructs a new Uploader instance using the dropbox client
func New(accessToken string) Uploader {
	client := dropbox.Client(accessToken, dropbox.Options{})
	return &dropBoxUploader{client}
}

// NewWithClient constructs a new Uploader instance using the given client
func NewWithClient(client Client) Uploader {
	return &dropBoxUploader{client}
}

func (uploader *dropBoxUploader) Upload(filepath string, content io.Reader) (string, error) {
	commitInfo := files.NewCommitInfo(filepath)
	commitInfo.Mode = &files.WriteMode{Update: "overwrite"}
	fileMetadata, err := uploader.client.Upload(commitInfo, content)
	if err != nil {
		return "", err
	}

	settings := sharing.NewCreateSharedLinkWithSettingsArg(fileMetadata.PathLower)
	sharedLinkMetadata, err := uploader.client.CreateSharedLinkWithSettings(settings)
	if err != nil {
		return "", err
	}

	return sharedLinkMetadata.File.Url, nil
}
