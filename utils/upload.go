package utils

import (
	"context"
	"log"
	"os"
	"strings"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

var cld *cloudinary.Cloudinary

func ConfigureCloudinary() {
	cloudName := os.Getenv("CLOUDINARY_CLOUD_NAME")
	apiKey := os.Getenv("CLOUDINARY_API_KEY")
	apiSecret := os.Getenv("CLOUDINARY_API_SECRET")

	cloudinaryInstance, err := cloudinary.NewFromParams(cloudName, apiKey, apiSecret)
	if err != nil {
		log.Fatal("Can't connect to cloudinary", err)
	}

	cld = cloudinaryInstance
}

func UploadFile(file string) (string, error) {
	var ctx = context.Background()
	resp, err := cld.Upload.Upload(ctx, file, uploader.UploadParams{})

	return resp.SecureURL, err
}

func DeleteFileFromUrl(url string) error {
	urlParts := strings.Split(url, "/")
	fileNameParts := strings.Split(urlParts[len(urlParts)-1], ".")
	publicId := fileNameParts[0]
	var ctx = context.Background()
	_, err := cld.Upload.Destroy(ctx, uploader.DestroyParams{
		PublicID: publicId,
	})

	return err
}
