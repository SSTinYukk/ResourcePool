package utils

import (
	"context"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"

	"g/front/backend/config"
)

// MinioUtils Minio工具结构体
type MinioUtils struct {
	Client     *minio.Client
	BucketName string
}

// NewMinioUtils 创建Minio工具实例
func NewMinioUtils(client *minio.Client) *MinioUtils {
	return &MinioUtils{
		Client:     client,
		BucketName: config.GetEnv("MINIO_BUCKET", "pool"),
	}
}

// UploadFile 上传文件到Minio
func (m *MinioUtils) UploadFile(file *multipart.FileHeader, directory string) (string, int64, error) {
	// 打开文件
	src, err := file.Open()
	if err != nil {
		return "", 0, err
	}
	defer src.Close()

	// 生成唯一文件名
	ext := filepath.Ext(file.Filename)
	fileName := fmt.Sprintf("%s/%s%s", directory, uuid.New().String(), ext)

	// 上传文件到Minio
	info, err := m.Client.PutObject(
		context.Background(),
		m.BucketName,
		fileName,
		src,
		file.Size,
		minio.PutObjectOptions{ContentType: file.Header.Get("Content-Type")},
	)

	if err != nil {
		return "", 0, err
	}

	log.Printf("文件上传成功: %s, 大小: %d", fileName, info.Size)
	return fileName, info.Size, nil
}

// GetFileURL 获取文件的临时URL
func (m *MinioUtils) GetFileURL(fileName string, expires time.Duration) (string, error) {
	// 检查文件是否存在
	_, err := m.Client.StatObject(context.Background(), m.BucketName, fileName, minio.StatObjectOptions{})
	if err != nil {
		return "", fmt.Errorf("文件不存在: %v", err)
	}

	// 如果是avatars存储桶，直接返回公开URL
	if m.BucketName == "avatars" {
		return fmt.Sprintf("%s/%s/%s", config.GetEnv("MINIO_ENDPOINT", "47.121.210.209:9000"), m.BucketName, fileName), nil
	}

	// 其他存储桶生成临时URL
	presignedURL, err := m.Client.PresignedGetObject(
		context.Background(),
		m.BucketName,
		fileName,
		expires,
		nil,
	)

	if err != nil {
		return "", err
	}

	return presignedURL.String(), nil
}

// DownloadFile 下载文件
func (m *MinioUtils) DownloadFile(fileName string) (io.ReadCloser, error) {
	// 获取文件对象
	object, err := m.Client.GetObject(
		context.Background(),
		m.BucketName,
		fileName,
		minio.GetObjectOptions{},
	)

	if err != nil {
		return nil, err
	}

	return object, nil
}

// DeleteFile 删除文件
func (m *MinioUtils) DeleteFile(fileName string) error {
	// 删除文件
	err := m.Client.RemoveObject(
		context.Background(),
		m.BucketName,
		fileName,
		minio.RemoveObjectOptions{},
	)

	if err != nil {
		return err
	}

	log.Printf("文件删除成功: %s", fileName)
	return nil
}
