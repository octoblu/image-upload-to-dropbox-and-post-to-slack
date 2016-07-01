package uploader_test

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"

	"github.com/dropbox/dropbox-sdk-go-unofficial/files"
	"github.com/dropbox/dropbox-sdk-go-unofficial/sharing"
	"github.com/octoblu/image-upload-to-dropbox-and-post-to-slack/uploader"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func sampleImage() io.Reader {
	data, err := base64.StdEncoding.DecodeString("iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAAAAAA6fptVAAAACklEQVR4nGP6DwABBQECz6AuzQAAAABJRU5ErkJggg==")
	Expect(err).To(BeNil())
	return bytes.NewReader(data)
}

var _ = Describe("Uploader", func() {
	var sut uploader.Uploader
	var fakeClient *FakeClient
	var err error

	Describe("Constructing a new Uploader using an access token", func() {
		BeforeEach(func() {
			sut = uploader.New("access-token")
		})

		It("Should exist", func() {
			Expect(sut).NotTo(BeNil())
		})
	})

	Describe("sut.Upload(filepath, content)", func() {
		Describe("when fakeClient.Upload and fakeClient.CreateSharedLinkWithSettings both go well", func() {
			var url string

			BeforeEach(func() {
				fileMetadata := &files.FileMetadata{PathLower: "/failures/example-2016-01-02.png"}

				fileLinkMetadata := &sharing.FileLinkMetadata{Url: "https://dropbox.biz/failures/example-2016-01-02.png"}
				sharedLinkMetadata := &sharing.SharedLinkMetadata{File: fileLinkMetadata}

				fakeClient = NewFakeClient()
				fakeClient.UploadSpy.ReturnsFileMetadata = fileMetadata
				fakeClient.CreateSharedLinkWithSettingsSpy.ReturnsSharedLinkMetadata = sharedLinkMetadata

				sut = uploader.NewWithClient(fakeClient)
				url, err = sut.Upload("failures/example-2016-01-02.png", sampleImage())
			})

			It("Should have called fakeClient.Upload", func() {
				Expect(fakeClient.UploadSpy.CallCount).To(Equal(1))
			})

			It("Should have called fakeClient.Upload with CommitInfo that included the filepath", func() {
				commitInfo := fakeClient.UploadSpy.LastCalledWithCommitInfo
				Expect(commitInfo.Path).To(Equal("failures/example-2016-01-02.png"))
				Expect(commitInfo.Mode.Update).To(Equal("overwrite"))
			})

			It("Should have called fakeClient.CreateSharedLinkWithSettings", func() {
				spy := fakeClient.CreateSharedLinkWithSettingsSpy
				Expect(spy.CallCount).To(Equal(1))
			})

			It("Should have called fakeClient.CreateSharedLinkWithSettings with the remote filepath", func() {
				spy := fakeClient.CreateSharedLinkWithSettingsSpy
				Expect(spy.LastCalledWith.Path).To(Equal("/failures/example-2016-01-02.png"))
			})

			It("should return the url", func() {
				Expect(url).To(Equal("https://dropbox.biz/failures/example-2016-01-02.png"))
			})

			It("Should have a nil error", func() {
				Expect(err).To(BeNil())
			})
		})

		Describe("when fakeClient.Upload returns an error", func() {
			var url string

			BeforeEach(func() {
				fakeClient = NewFakeClient()
				fakeClient.UploadSpy.ReturnsError = fmt.Errorf("Error uploading.")
				sut = uploader.NewWithClient(fakeClient)
				url, err = sut.Upload("failures/example-2016-01-02.png", sampleImage())
			})

			It("Should have an empty url", func() {
				Expect(url).To(Equal(""))
			})

			It("Should have an error", func() {
				Expect(err).NotTo(BeNil())
				Expect(err.Error()).To(Equal("Error uploading."))
			})

			It("Should have called fakeClient.Upload", func() {
				Expect(fakeClient.UploadSpy.CallCount).To(Equal(1))
			})

			It("Should have called fakeClient.Upload with CommitInfo that included the filepath", func() {
				Expect(fakeClient.UploadSpy.LastCalledWithCommitInfo.Path).To(Equal("failures/example-2016-01-02.png"))
				Expect(fakeClient.UploadSpy.LastCalledWithCommitInfo.Mode.Update).To(Equal("overwrite"))
			})

			It("Should have called fakeClient.Upload with the content passed in", func() {
				Expect(fakeClient.UploadSpy.LastCalledWithContent).To(Equal(sampleImage()))
			})
		})
	})
})
