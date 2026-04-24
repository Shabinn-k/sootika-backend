package uploads

import (
	"context"
	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"mime/multipart"
	"os"
	"strings"
	"time"
)

type CloudinaryResult struct {
	URL      string
	PublicID string
}

func UploadImageFile(file multipart.File, filename string) (*CloudinaryResult, error) {
	cld, err := cloudinary.NewFromParams(
		os.Getenv("CLOUD_NAME"),
		os.Getenv("API_KEY"),
		os.Getenv("API_SECRET"),
	)
	if err != nil {
		return nil, err
	}
	publicID := strings.TrimSuffix(filename, ".jpg")
	publicID = strings.TrimSuffix(publicID, ".png")
	publicID = strings.TrimSuffix(publicID, ".jpeg")
	publicID = strings.TrimSuffix(publicID, ".webp")
	publicID = publicID + "_" + time.Now().Format("20060102150405")

	resp, err := cld.Upload.Upload(context.Background(), file, uploader.UploadParams{
		PublicID: publicID,
		Folder:   os.Getenv("CLOUDINARY_UPLOAD_FOLDER"),
	})
	if err != nil {
		return nil, err
	}
	return &CloudinaryResult{
		URL:      resp.SecureURL,
		PublicID: resp.PublicID,
	}, nil
}

func UploadImage(filePath string) (string, error) {
	cld, err := cloudinary.NewFromParams(
		os.Getenv("CLOUD_NAME"),
		os.Getenv("API_KEY"),
		os.Getenv("API_SECRET"),
	)
	if err != nil {
		return "", err
	}
	resp, err := cld.Upload.Upload(context.Background(), filePath, uploader.UploadParams{
		Folder: os.Getenv("CLOUDINARY_UPLOAD_FOLDER"),
	})
	if err != nil {
		return "", err
	}
	return resp.SecureURL, nil
}

func DeleteImage(publicID string) error {
	cld, err := cloudinary.NewFromParams(
		os.Getenv("CLOUD_NAME"),
		os.Getenv("API_KEY"),
		os.Getenv("API_SECRET"),
	)
	if err != nil {
		return err
	}
	_, err = cld.Upload.Destroy(context.Background(), uploader.DestroyParams{
		PublicID: publicID,
	})
	return err
}
