// backend/services/s3_service.go

package services

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"gorm.io/gorm"

	"github.com/shuttlersit/ads-player/backend/models"
)

var S3BaseURL string

// S3Service provides methods for interacting with AWS S3
type S3Service interface {
	FetchAdvertisementsFromS3() ([]models.Advertisement, error)
	GetVideoURL(videoKey string) (string, error)
	GetVideoContent(videoKey string) ([]byte, error)
	UploadFile(videokey string, file io.Reader) error
	DownloadFile(videoKey string) (bytes.Reader, error)
	GeneratePresignedURL(key string) (string, error)
}

type S3ServiceDBModel struct {
	DB *gorm.DB
}

// DefaultS3Service is the default implementation of S3Service
type DefaultS3Service struct {
	// Add any necessary configuration or dependencies
	BucketURL string // URL of the S3 bucket
	Region    string // Region of the S3 bucket
	AccessKey string
	SecretKey string
}

// NewS3Service creates a new instance of S3ServiceImpl
func NewDefaultS3Service(bucketURL, region, accessKey, secretKey string) *DefaultS3Service {
	// Perform any setup or initialization
	return &DefaultS3Service{
		BucketURL: bucketURL,
		Region:    region,
		AccessKey: accessKey,
		SecretKey: secretKey,
		// Initialize other fields as needed
	}
}

// GeneratePresignedURL generates a presigned URL for the specified S3 key
func (s3 *DefaultS3Service) GeneratePresignedURL(key string) (string, error) {
	// Implement logic to generate a presigned URL for the specified S3 key
	// Example: Use an S3 SDK to generate a presigned URL
	// Replace the following line with your actual implementation
	return "https://example.com/presigned-url", nil
}

// UploadFile uploads a file to the S3 bucket with the specified key
func (s3 *DefaultS3Service) UploadFile(key string, file io.Reader) error {
	// Placeholder for S3 file upload logic (replace with actual code)
	return nil
}

// DownloadFile downloads a file from S3 and returns a bytes.Reader for playback
func (s3 *DefaultS3Service) DownloadFile(videoKey string) (bytes.Reader, error) {
	// Replace with the actual implementation to download the file from S3
	videoURL := fmt.Sprintf("%s/%s", s3.BucketURL, videoKey)
	resp, err := http.Get(videoURL)
	if err != nil {
		log.Printf("Error downloading file from S3: %v\n", err)
		return bytes.Reader{}, err
	}
	defer resp.Body.Close()

	fileContent, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading file content: %v\n", err)
		return bytes.Reader{}, err
	}

	return *bytes.NewReader(fileContent), nil
}

// FetchAdvertisementsFromS3 fetches new advertisements from AWS S3
func (s3 *DefaultS3Service) FetchAdvertisementsFromS3() ([]models.Advertisement, error) {
	fmt.Println("Fetching advertisements from AWS S3...")

	// Replace this with your actual logic to fetch advertisements from AWS S3
	// For example, use an AWS SDK to list objects in a bucket and retrieve metadata

	// Mocking data for demonstration
	newAdvertisements := []models.Advertisement{
		{ID: 101, Name: "New Ad 1", ContentURL: "s3://bucket/ad1.mp4", Duration: 15},
		{ID: 102, Name: "New Ad 2", ContentURL: "s3://bucket/ad2.mp4", Duration: 20},
	}

	fmt.Println("Advertisements fetched from AWS S3.")
	return newAdvertisements, nil
}

// GetVideoContent fetches video content from S3 based on the video's key
func (s3 *DefaultS3Service) GetVideoContent(videoKey string) ([]byte, error) {
	// Replace with the actual implementation to download video content from S3
	videoURL := fmt.Sprintf("%s/%s", s3.BucketURL, videoKey)
	resp, err := http.Get(videoURL)
	if err != nil {
		log.Printf("Error fetching video content from S3: %v\n", err)
		return nil, err
	}
	defer resp.Body.Close()

	videoContent, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading video content: %v\n", err)
		return nil, err
	}

	return videoContent, nil
}

// GetVideoURL fetches the URL of the video content from S3 based on the video's key
func (s3 *DefaultS3Service) GetVideoURL(videoKey string) (string, error) {
	// Replace with the actual implementation to generate a signed URL or construct the URL
	videoURL := fmt.Sprintf("%s/%s", s3.BucketURL, videoKey)
	return videoURL, nil
}

/*

// Implement the GetVideoURL method
func (s3 *DefaultS3Service) GetVideoURL(videoKey string) (string, error) {
	// Replace with the actual implementation to generate a signed URL or construct the URL
	// Use the AWS SDK or any other library you prefer for S3 operations
	// Example using AWS SDK:
	// sess := session.Must(session.NewSession(&aws.Config{
	// 	Region: aws.String(s3.Region),
	// }))
	// svc := s3manager.NewDownloader(sess)
	// req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
	// 	Bucket: aws.String(s3.BucketURL),
	// 	Key:    aws.String(videoKey),
	// })
	// urlStr, err := req.Presign(time.Minute * 15) // Presign URL for 15 minutes
	// if err != nil {
	// 	return "", err
	// }
	// return urlStr, nil

	// Placeholder return (replace with actual logic)
	return fmt.Sprintf("%s/%s", s3.BucketURL, videoKey), nil
}

// Implement the GetVideoContent method
func (s3 *DefaultS3Service) GetVideoContent(videoKey string) ([]byte, error) {
	// Replace with the actual implementation to download video content from S3
	// Use the AWS SDK or any other library you prefer for S3 operations
	// Example using AWS SDK:
	// sess := session.Must(session.NewSession(&aws.Config{
	// 	Region: aws.String(s3.Region),
	// }))
	// svc := s3manager.NewDownloader(sess)
	// buf := aws.NewWriteAtBuffer([]byte{})
	// _, err := svc.Download(buf,
	// 	&s3.GetObjectInput{
	// 		Bucket: aws.String(s3.BucketURL),
	// 		Key:    aws.String(videoKey),
	// 	})
	// if err != nil {
	// 	return nil, err
	// }
	// return buf.Bytes(), nil

	// Placeholder return (replace with actual logic)
	return []byte(fmt.Sprintf("Video content for key %s", videoKey)), nil
}

// Implement the DownloadFile method
func (s3 *DefaultS3Service) DownloadFile(videoKey string) (io.Reader, error) {
	// Replace with the actual implementation to download the file from S3
	// Use the AWS SDK or any other library you prefer for S3 operations
	// Example using AWS SDK:
	// sess := session.Must(session.NewSession(&aws.Config{
	// 	Region: aws.String(s3.Region),
	// }))
	// svc := s3manager.NewDownloader(sess)
	// buf := aws.NewWriteAtBuffer([]byte{})
	// _, err := svc.Download(buf,
	// 	&s3.GetObjectInput{
	// 		Bucket: aws.String(s3.BucketURL),
	// 		Key:    aws.String(videoKey),
	// 	})
	// if err != nil {
	// 	return nil, err
	// }
	// return bytes.NewReader(buf.Bytes()), nil

	// Placeholder return (replace with actual logic)
	return strings.NewReader(fmt.Sprintf("File content for key %s", videoKey)), nil
}

// Implement the UploadFile method
func (s3 *S3ServiceImpl) UploadFile(key string, content io.Reader) error {
	// Replace with the actual implementation to upload a file to S3
	// Use the AWS SDK or any other library you prefer for S3 operations
	// Example using AWS SDK:
	// sess := session.Must(session.NewSession(&aws.Config{
	// 	Region: aws.String(s3.Region),
	// }))
	// uploader := s3manager.NewUploader(sess)
	// _, err := uploader.Upload(&s3manager.UploadInput{
	// 	Bucket: aws.String(s3.BucketURL),
	// 	Key:    aws.String(key),
	// 	Body:   content,
	// })
	// return err

	// Placeholder return (replace with actual logic)
	return nil
}

// Implement the DeleteFile method
func (s3 *S3ServiceImpl) DeleteFile(key string) error {
	// Replace with the actual implementation to delete a file from S3
	// Use the AWS SDK or any other library you prefer for S3 operations
	// Example using AWS SDK:
	// sess := session.Must(session.NewSession(&aws.Config{
	// 	Region: aws.String(s3.Region),
	// }))
	// svc := s3.New(sess)
	// _, err := svc.DeleteObject(&s3.DeleteObjectInput{
	// 	Bucket: aws.String(s3.BucketURL),
	// 	Key:    aws.String(key),
	// })
	// return err

	// Placeholder return (replace with actual logic)
	return nil
}

*/
