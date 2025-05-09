package config

import (
	"context"
	"log"
	"os"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// GetEnv 获取环境变量，如果不存在则返回默认值
func GetEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// InitMinioClient 初始化Minio客户端
func InitMinioClient() (*minio.Client, error) {
	endpoint := GetEnv("MINIO_ENDPOINT", "47.121.210.209:9000")
	accessKey := GetEnv("MINIO_ACCESS_KEY", "minioadmin")
	secretKey := GetEnv("MINIO_SECRET_KEY", "minioadmin")
	useSSL := false

	// 初始化Minio客户端
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, err
	}

	return client, nil
}

// EnsureBucketExists 确保bucket存在
func EnsureBucketExists(client *minio.Client, bucketName string) {
	ctx := context.Background()
	exists, err := client.BucketExists(ctx, bucketName)
	if err != nil {
		log.Fatalf("检查bucket失败: %v", err)
	}

	if !exists {
		err = client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
		if err != nil {
			log.Fatalf("创建bucket失败: %v", err)
		}
		log.Printf("成功创建bucket: %s", bucketName)
	} else {
		log.Printf("bucket已存在: %s", bucketName)
	}
}
