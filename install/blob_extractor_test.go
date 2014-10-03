package install_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	boshlog "github.com/cloudfoundry/bosh-agent/logger"

	fakeblobstore "github.com/cloudfoundry/bosh-agent/blobstore/fakes"
	fakesys "github.com/cloudfoundry/bosh-agent/system/fakes"
	testfakes "github.com/cloudfoundry/bosh-micro-cli/testutils/fakes"

	. "github.com/cloudfoundry/bosh-micro-cli/install"
)

var _ = Describe("BlobExtractor", func() {
	var (
		blobExtractor BlobExtractor
		blobstore     *fakeblobstore.FakeBlobstore
		targetDir     string
		fakeExtractor *testfakes.FakeMultiResponseExtractor
		logger        boshlog.Logger
		fs            *fakesys.FakeFileSystem

		blobID   string
		blobSHA1 string
	)
	BeforeEach(func() {
		blobstore = fakeblobstore.NewFakeBlobstore()
		targetDir = "fake-target-dir"
		fakeExtractor = testfakes.NewFakeMultiResponseExtractor()
		logger = boshlog.NewLogger(boshlog.LevelNone)
		fs = fakesys.NewFakeFileSystem()
		blobID = "fake-blob-id"
		blobSHA1 = "fake-sha1"

		blobExtractor = NewBlobExtractor(fs, fakeExtractor, blobstore, logger)
	})

	Context("when the specified blobID exists in the blobstore", func() {
		BeforeEach(func() {
			blobstore.GetFileName = "fake-blob-file"
		})

		It("creates the installed package dir if it does not exist", func() {
			Expect(fs.FileExists(targetDir)).To(BeFalse())
			err := blobExtractor.Extract(blobID, blobSHA1, targetDir)
			Expect(err).ToNot(HaveOccurred())
			Expect(fs.FileExists(targetDir)).To(BeTrue())
		})

		It("extracts the blob into the target dir", func() {
			err := blobExtractor.Extract(blobID, blobSHA1, targetDir)
			Expect(err).ToNot(HaveOccurred())
			Expect(fakeExtractor.DecompressedFiles()).To(ContainElement("fake-target-dir/fake-blob-file"))
		})

		It("cleans up the blob file", func() {
			err := blobExtractor.Extract(blobID, blobSHA1, targetDir)
			Expect(err).ToNot(HaveOccurred())
			Expect(blobstore.CleanUpFileName).To(Equal("fake-blob-file"))
		})

		Context("when getting the blob from the blobstore errors", func() {
			BeforeEach(func() {
				blobstore.GetError = errors.New("fake-error")
			})

			It("returns an error", func() {
				err := blobExtractor.Extract(blobID, blobSHA1, targetDir)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("fake-error"))
			})
		})

		Context("when creating the target dir fails", func() {
			It("return an error", func() {
				fs.MkdirAllError = errors.New("fake-error")
				err := blobExtractor.Extract(blobID, blobSHA1, targetDir)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Creating target dir"))
				Expect(err.Error()).To(ContainSubstring("fake-error"))
			})
		})

		Context("when extracting the blob fails", func() {
			BeforeEach(func() {
				fakeExtractor.SetDecompressBehavior(
					"fake-blob-file",
					"fake-target-dir",
					errors.New("fake-error"))
			})

			It("returns an error", func() {
				err := blobExtractor.Extract(blobID, blobSHA1, targetDir)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Extracting compiled package"))
				Expect(err.Error()).To(ContainSubstring("fake-error"))
			})

			It("cleans up the target dir if it was created by this blobExtractor", func() {
				Expect(fs.FileExists(targetDir)).To(BeFalse())
				err := blobExtractor.Extract(blobID, blobSHA1, targetDir)
				Expect(err).To(HaveOccurred())
				Expect(fs.FileExists(targetDir)).To(BeFalse())
			})
		})

		Context("when cleaning up the downloaded blob errors", func() {
			BeforeEach(func() {
				blobstore.CleanUpErr = errors.New("fake-error")
			})

			It("does not return the error", func() {
				err := blobExtractor.Extract(blobID, blobSHA1, targetDir)
				Expect(err).ToNot(HaveOccurred())
			})
		})
	})
})
