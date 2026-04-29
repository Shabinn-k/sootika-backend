package uploads

import (
	"context"
	"fmt"
	"mime/multipart"
	"os"
	"strings"
	"time"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

type CloudinaryResult struct {
	URL      string
	PublicID string
}

func UploadImageFile(file multipart.File, filename string) (*CloudinaryResult, error) {
	// Get credentials from environment
	cloudName := os.Getenv("CLOUD_NAME")
	apiKey := os.Getenv("API_KEY")
	apiSecret := os.Getenv("API_SECRET")
	
	// Debug output
	fmt.Println("=== Cloudinary Debug ===")
	fmt.Println("Cloud Name:", cloudName)
	fmt.Println("API Key:", apiKey)
	fmt.Println("API Secret length:", len(apiSecret))
	
	if cloudName == "" {
		return nil, fmt.Errorf("CLOUD_NAME environment variable is not set")
	}
	if apiKey == "" {
		return nil, fmt.Errorf("API_KEY environment variable is not set")
	}
	if apiSecret == "" {
		return nil, fmt.Errorf("API_SECRET environment variable is not set")
	}
	
	// Create Cloudinary instance
	cld, err := cloudinary.NewFromParams(cloudName, apiKey, apiSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to create Cloudinary client: %v", err)
	}
	
	// Generate public ID
	publicID := strings.TrimSuffix(filename, ".jpg")
	publicID = strings.TrimSuffix(publicID, ".png")
	publicID = strings.TrimSuffix(publicID, ".jpeg")
	publicID = strings.TrimSuffix(publicID, ".webp")
	publicID = strings.TrimSuffix(publicID, ".gif")
	publicID = strings.TrimSuffix(publicID, ".bmp")
	publicID = publicID + "_" + time.Now().Format("20060102150405")
	
	folder := os.Getenv("CLOUDINARY_UPLOAD_FOLDER")
	if folder == "" {
		folder = "sootika_products"
	}
	
	fmt.Println("Uploading to folder:", folder)
	fmt.Println("Public ID:", publicID)
	
	// Upload to Cloudinary
	ctx := context.Background()
	resp, err := cld.Upload.Upload(ctx, file, uploader.UploadParams{
		PublicID: publicID,
		Folder:   folder,
	})
	if err != nil {
		return nil, fmt.Errorf("Cloudinary upload failed: %v", err)
	}
	
	fmt.Println("Upload successful! URL:", resp.SecureURL)
	
	return &CloudinaryResult{
		URL:      resp.SecureURL,
		PublicID: resp.PublicID,
	}, nil
}

func DeleteImage(publicID string) error {
	cloudName := os.Getenv("CLOUD_NAME")
	apiKey := os.Getenv("API_KEY")
	apiSecret := os.Getenv("API_SECRET")
	
	if cloudName == "" || apiKey == "" || apiSecret == "" {
		return fmt.Errorf("Cloudinary credentials missing")
	}
	
	cld, err := cloudinary.NewFromParams(cloudName, apiKey, apiSecret)
	if err != nil {
		return err
	}
	
	ctx := context.Background()
	_, err = cld.Upload.Destroy(ctx, uploader.DestroyParams{
		PublicID: publicID,
	})
	return err
}