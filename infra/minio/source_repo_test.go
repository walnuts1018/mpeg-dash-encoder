package minio

import (
	"context"
	"strings"
	"testing"

	"github.com/minio/minio-go/v7"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("SourceClient", Ordered, func() {
	client := NewSourceClient(sourceClientBucketName, minioClient)

	ctx := context.Background()

	sourceFiles := map[string]string{
		"test1": "thisismusicsourcefile1",
		"test2": "thisismusicsourcefile2",
	}

	for k, v := range sourceFiles {
		_, err := minioClient.PutObject(ctx, sourceClientBucketName, k, strings.NewReader(v), int64(len(v)), minio.PutObjectOptions{})
		Expect(err).NotTo(HaveOccurred())
	}

	It("Normal", func() {
		uploadedFiles := client.ListUploadedFiles(ctx)
		for file, err := range uploadedFiles {
			Expect(err).NotTo(HaveOccurred())
			Expect(file.ID).To(BeKeyOf(sourceFiles))
			Expect(file.Tags).To(BeEmpty())
		}

		err := client.SetObjectTags(ctx, "test1", map[string]string{"tag1": "value1"})
		Expect(err).NotTo(HaveOccurred())

		uploadedFiles = client.ListUploadedFiles(ctx)
		for file, err := range uploadedFiles {
			Expect(err).NotTo(HaveOccurred())
			if file.ID == "test1" {
				Expect(file.Tags).To(Equal(map[string]string{"tag1": "value1"}))
			} else {
				Expect(file.Tags).To(BeEmpty())
			}
		}

		err = client.SetObjectTags(ctx, "test1", map[string]string{"tag1": "updatedvalue1", "tag2": "value2"})
		Expect(err).NotTo(HaveOccurred())

		uploadedFiles = client.ListUploadedFiles(ctx)
		for file, err := range uploadedFiles {
			Expect(err).NotTo(HaveOccurred())
			if file.ID == "test1" {
				Expect(file.Tags).To(Equal(map[string]string{"tag1": "updatedvalue1", "tag2": "value2"}))
			} else {
				Expect(file.Tags).To(BeEmpty())
			}
		}

		err = client.RemoveObjectTags(ctx, "test1")
		Expect(err).NotTo(HaveOccurred())

		uploadedFiles = client.ListUploadedFiles(ctx)
		for file, err := range uploadedFiles {
			Expect(err).NotTo(HaveOccurred())
			Expect(file.Tags).To(BeEmpty())
		}

		err = client.RemoveObjectTags(ctx, "test1")
		Expect(err).NotTo(HaveOccurred())

		err = client.DeleteSourceContent(ctx, "test1")
		Expect(err).NotTo(HaveOccurred())

		uploadedFiles = client.ListUploadedFiles(ctx)
		for file, err := range uploadedFiles {
			Expect(err).NotTo(HaveOccurred())
			Expect(file.ID).To(Equal("test2"))
		}

		err = client.DeleteSourceContent(ctx, "test2")
		Expect(err).NotTo(HaveOccurred())
	})
})

func TestSourceRepo(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "SourceRepo Suite")
}
