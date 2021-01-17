package service

import (
	"bytes"
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/revel/revel"
	"time"
)

type MinioClient struct {
	Endpoint string
	AccessKey string
	SecretKey string
	Bucket string
	minioClient *minio.Client
}

var client *MinioClient


func GetMinioClient() (mc *MinioClient, err error) {
	if client != nil {
		return client, nil
	}
	endpoint, _ := revel.Config.String("minio.Endpoint")
	accessKey, _ := revel.Config.String("minio.AccessKey")
	secretKey, _ := revel.Config.String("minio.SecretKey")
	bucket, _ := revel.Config.String("minio.Bucket")
	minioClient, err  := minio.New(endpoint, &minio.Options{
		Creds: credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: false,
	})
	if err != nil {
		fmt.Println("初始化 Minio Client 错误！", err)
		return nil, err
	}
	client = &MinioClient{
		Endpoint: endpoint,
		AccessKey: accessKey,
		SecretKey: secretKey,
		minioClient: minioClient,
	}
	found, err := minioClient.BucketExists(context.Background(), bucket)
	if err != nil {
		fmt.Println("检查 Bucket 出错！", err)
		return client, nil
	}
	// bucket 还没有存在，则创建
	if !found {
		err := minioClient.MakeBucket(
			context.Background(),
			bucket,
			minio.MakeBucketOptions{Region: "",ObjectLocking: true})
		if err != nil {
			fmt.Println("创建 Bucket 错误！", err)
		}
	}
	return client, nil
}

func (mc MinioClient)PutObject(data []byte, object string, contentType string)(path string, b bool)  {
	reader := bytes.NewReader(data)
	uploadInfo, err := mc.minioClient.PutObject(
		context.Background(),
		mc.Bucket,
		object,
		reader,
		reader.Size(),
		minio.PutObjectOptions{
			ContentType: contentType,
		})
	if err != nil {
		revel.AppLog.Error("上传数据到 Minio 错误！")
		return object, false
	}
	revel.AppLog.Info("上传数据到 Minio 成功！", uploadInfo)
	return object, true
}


func (mc MinioClient)StatObject(object string) *minio.ObjectInfo  {
	objInfo, err := mc.minioClient.StatObject(context.Background(), mc.Bucket, object, minio.StatObjectOptions{})
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return &objInfo
}

func (mc MinioClient)GetObject(object string) (bytes []byte, contentType string, b bool) {
	objectInfo := mc.StatObject(object)
	if objectInfo == nil {
		return nil, "", false
	}
	obj, err := mc.minioClient.GetObject(context.Background(), mc.Bucket, object, minio.GetObjectOptions{})
	if err != nil {
		return nil, "", false
	}
	_, er := obj.Read(bytes)
	if er != nil {
		return nil, "", false
	}
	return bytes, objectInfo.ContentType, true
}


func (mc MinioClient)CopyObject(dstObject string, srcObject string) (b bool, info *minio.UploadInfo) {
	// Source object
	srcOpts := minio.CopySrcOptions{
		Bucket: mc.Bucket,
		Object: srcObject,
	}
	// Destination object
	dstOpts := minio.CopyDestOptions{
		Bucket: mc.Bucket,
		Object: dstObject,
	}
	// Copy object call
	uploadInfo, err := mc.minioClient.CopyObject(context.Background(), dstOpts, srcOpts)
	if err != nil {
		fmt.Println(err)
		return false, nil
	}
	return true, &uploadInfo
}


func (mc MinioClient)DeleteObject(object string) (b bool) {
	err := mc.minioClient.RemoveObject(context.Background(), mc.Bucket, object, minio.RemoveObjectOptions{})
	if err != nil {
		return false
	}
	return true
}


func (mc MinioClient)PresignedGetObject(object string, expireAtAfterSeconds time.Duration) (url string) {
	presignedURL, err := mc.minioClient.PresignedGetObject(context.Background(), mc.Bucket, object, time.Second * expireAtAfterSeconds, nil)
	if err != nil {
		return mc.Endpoint + "/" + object
	}
	//fmt.Println(err)
	//fmt.Println(presignedURL.Scheme)
	//fmt.Println(presignedURL.String())
	return presignedURL.String()
}


func InitMinioClient() {
	GetMinioClient()
}