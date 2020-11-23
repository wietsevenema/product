package product

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"cloud.google.com/go/compute/metadata"
	"google.golang.org/api/idtoken"
)

type Client struct {
	serviceURL string
	http       *http.Client
}

func NewClient(baseURL string) (*Client, error) {

	result := &Client{
		serviceURL: baseURL,
		http:       nil,
	}

	if metadata.OnGCE() {
		client, err := idtoken.NewClient(context.Background(), baseURL)
		if err != nil {
			return nil, err
		}
		result.http = client
		return result, nil
	}
	result.http = &http.Client{}
	return result, nil

}

func (c *Client) GetProducts(ctx context.Context) (*[]Product, error) {
	url := fmt.Sprintf("%s/random/", c.serviceURL)
	req, err := http.NewRequest("GET", url, nil)
	req = req.WithContext(ctx)
	if err != nil {
		return nil, err
	}
	resp, err := c.http.Do(req) // Adds ID Token to request
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)

	products := &[]Product{}
	err = decoder.Decode(products)
	if err != nil {
		return nil, err
	}
	return products, nil

}
