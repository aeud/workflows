package flows

import (
	"context"
	"io"
	"os"
	"strings"

	"cloud.google.com/go/storage"
)

func gcsPathToBucketObject(path string) (string, string) {
	elements := strings.Split(strings.Replace(path, "gs://", "", 1), "/")
	return elements[0], strings.Join(elements[1:], "/")
}

func OpenFromPath(path string) (io.ReadCloser, error) {
	if strings.Contains(path, "gs://") {
		bucketRef, objectRef := gcsPathToBucketObject(path)
		ctx := context.Background()
		client, err := storage.NewClient(ctx)
		if err != nil {
			return nil, err
		}
		bucket := client.Bucket(bucketRef)
		object := bucket.Object(objectRef)
		return object.NewReader(ctx)
	}
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	return file, nil
}
