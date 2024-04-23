package resolvers

import (
	"graphql-api/internal/image"
	"graphql-api/pkg/data/models"
	"github.com/graphql-go/graphql"
)

func GetImageResolve(params graphql.ResolveParams) (interface{}, error) {

	// Fetch Image from rest api
	image_url, err := image.RandomImageUrl(150,150)
	image := models.ImageModel{ImageUrl: image_url}
	if err != nil {
		return nil, err
	}
	return image, nil
}