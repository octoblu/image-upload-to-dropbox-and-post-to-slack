package uploader_test

import "io"
import "github.com/dropbox/dropbox-sdk-go-unofficial/files"
import "github.com/dropbox/dropbox-sdk-go-unofficial/sharing"

type FakeClient struct {
	CreateSharedLinkWithSettingsSpy struct {
		CallCount                 int
		LastCalledWith            *sharing.CreateSharedLinkWithSettingsArg
		ReturnsError              error
		ReturnsSharedLinkMetadata *sharing.SharedLinkMetadata
	}

	ListSharedLinksSpy struct {
		CallCount                    int
		LastCalledWith               *sharing.ListSharedLinksArg
		ReturnsListSharedLinksResult *sharing.ListSharedLinksResult
	}

	UploadSpy struct {
		CallCount                int
		LastCalledWithCommitInfo *files.CommitInfo
		LastCalledWithContent    io.Reader
		ReturnsError             error
		ReturnsFileMetadata      *files.FileMetadata
	}
}

func NewFakeClient() *FakeClient {
	return &FakeClient{}
}

func (client *FakeClient) CreateSharedLinkWithSettings(arg *sharing.CreateSharedLinkWithSettingsArg) (res *sharing.SharedLinkMetadata, err error) {
	spy := &client.CreateSharedLinkWithSettingsSpy

	spy.CallCount++
	spy.LastCalledWith = arg
	return spy.ReturnsSharedLinkMetadata, spy.ReturnsError
}

func (client *FakeClient) ListSharedLinks(arg *sharing.ListSharedLinksArg) (res *sharing.ListSharedLinksResult, err error) {
	spy := &client.ListSharedLinksSpy

	spy.CallCount++
	spy.LastCalledWith = arg
	return spy.ReturnsListSharedLinksResult, nil
}

func (client *FakeClient) Upload(commitInfo *files.CommitInfo, content io.Reader) (res *files.FileMetadata, err error) {
	spy := &client.UploadSpy

	spy.CallCount++
	spy.LastCalledWithCommitInfo = commitInfo
	spy.LastCalledWithContent = content
	return spy.ReturnsFileMetadata, spy.ReturnsError
}
