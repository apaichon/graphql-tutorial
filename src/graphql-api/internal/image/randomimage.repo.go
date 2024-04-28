package image

import (
	"errors"
	"net/http"
	"fmt"
)

func RandomImageUrl(h, w int) (string, error) {

	// The URL of the REST API endpoint
	url := fmt.Sprintf("https://picsum.photos/%v/%v", w, h)

	// Create an HTTP client
	client := &http.Client{}

	// Create an HTTP GET request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	// Send the request and get the response
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Check if the response is successful
	if resp.StatusCode != http.StatusOK {
		return "", errors.New("Request failed with status:" + resp.Status)
	}

	image_url := fmt.Sprintf("%s",resp.Request.URL)
	return image_url, err
}