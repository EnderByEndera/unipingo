package utils

import (
	"context"
	"fmt"
	"log"
	"melodie-site/server/config"
	"mime/multipart"
	"strings"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type BucketNames struct {
	Files string
}

type OSSHandler struct {
	Buckets BucketNames
	Client  *minio.Client
}

var ossHandler *OSSHandler

func GetOSSHandler() *OSSHandler {
	if ossHandler == nil {
		ossHandler = Connect()
		ossHandler.Init()
	}
	return ossHandler
}

func (handler *OSSHandler) Init() {
	handler.Buckets = BucketNames{Files: "files"}
	handler.EnsureBucket(handler.Buckets.Files)
}

// 确保bucket存在
func (handler *OSSHandler) EnsureBucket(bucketName string) {
	ctx := context.Background()
	err := handler.Client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
	if err != nil {
		// Check to see if we already own this bucket (which happens if you run this twice)
		exists, errBucketExists := handler.Client.BucketExists(ctx, bucketName)
		if errBucketExists == nil && exists {
			log.Printf("We already own %s\n", bucketName)
		} else {
			log.Fatalln(err)
		}
	} else {
		log.Printf("Successfully created %s\n", bucketName)
	}
}

func Connect() *OSSHandler {
	cfg := config.GetConfig()
	endpoint := fmt.Sprintf("%s:%s", cfg.ADDRESSES.OSS_ADDR, cfg.ADDRESSES.OSS_PORT)
	accessKeyID := cfg.INFRASTRUCTURE_USER.NAME
	secretAccessKey := cfg.INFRASTRUCTURE_USER.PASSWORD
	useSSL := false

	// Initialize minio client object.
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatalln(err)
	}
	return &OSSHandler{Client: minioClient}
}

func PutFile(bucketName, objectName, filePath, contentType string) error {
	client := GetOSSHandler().Client
	ctx := context.Background()
	// Upload the zip file with FPutObject
	info, err := client.FPutObject(ctx, bucketName, objectName, filePath, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		return err
	}

	log.Printf("Successfully uploaded %s of size %d\n", objectName, info.Size)
	return nil
}

func PutObject(bucketName, objectName, contentType string, file multipart.File, size int64) (err error) {
	client := GetOSSHandler().Client
	t0 := time.Now()
	uploadInfo, err := client.PutObject(context.Background(),
		bucketName, objectName, file, size,
		minio.PutObjectOptions{ContentType: contentType})
	fmt.Println("since", time.Since(t0))
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Successfully uploaded bytes: ", uploadInfo)
	return
}

func (handler *OSSHandler) GetFileAndSave(objectName, saveAsFile string) (err error) {
	err = handler.Client.FGetObject(context.Background(), handler.Buckets.Files, objectName,
		saveAsFile, minio.GetObjectOptions{})
	if err != nil {
		return
	}
	return
}

func SplitExt(fileName string) (ext string) {
	splitted := strings.Split(fileName, ".")
	if len(splitted) > 1 {
		ext = splitted[len(splitted)-1]
	}
	return
}
