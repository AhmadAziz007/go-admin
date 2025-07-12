package service

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinioService struct {
	client     *minio.Client
	bucketName string
	endpoint   string
	usessl     bool
}

func NewMinioService(endpoint, accessKey, secretKey, bucketName string, usessl bool) (*MinioService, error) {
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: usessl,
	})
	if err != nil {
		return nil, err
	}

	return &MinioService{
		client:     client,
		bucketName: bucketName,
		endpoint:   endpoint,
		usessl: usessl,
	}, nil
}

func (m *MinioService) UploadFile(file io.Reader, fileSize int64, contentType string) (string, error) {
	ctx := context.Background()

	// Generate unique filename
	objectName := fmt.Sprintf("%d", time.Now().UnixNano())
	if strings.Contains(contentType, "image/jpeg") {
		objectName += ".jpg"
	} else if strings.Contains(contentType, "image/png") {
		objectName += ".png"
	} else {
		objectName += ".bin"
	}

	// Create bucket if not exists
	exists, err := m.client.BucketExists(ctx, m.bucketName)
	if err == nil && !exists {
		err = m.client.MakeBucket(ctx, m.bucketName, minio.MakeBucketOptions{})
	}
	if err != nil {
		return "", err
	}

	// Upload file
	_, err = m.client.PutObject(
		ctx,
		m.bucketName,
		objectName,
		file,
		fileSize,
		minio.PutObjectOptions{ContentType: contentType},
	)
	if err != nil {
		return "", err
	}

	protocol := "http"
	if m.usessl {
		protocol = "https"
	}
	return fmt.Sprintf("%s://%s/%s/%s", protocol, m.endpoint, m.bucketName, objectName), nil
}

func (m *MinioService) GetFile(imgUrl string) ([]byte, error) {
	ctx := context.Background()

	// Parse URL to get object name
	u, err := url.Parse(imgUrl)
	if err != nil {
		return nil, err
	}
	objectName := path.Base(u.Path)

	// Get object
	obj, err := m.client.GetObject(
		ctx,
		m.bucketName,
		objectName,
		minio.GetObjectOptions{},
	)
	if err != nil {
		return nil, err
	}
	defer obj.Close()

	return io.ReadAll(obj)
}

func (m *MinioService) DeleteFile(imgUrl string) error {
	ctx := context.Background()

	// Parse URL to get object name
	u, err := url.Parse(imgUrl)
	if err != nil {
		return err
	}
	objectName := path.Base(u.Path)

	return m.client.RemoveObject(ctx, m.bucketName, objectName, minio.RemoveObjectOptions{})
}