package run

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"products-frontend/lib/config"

	"golang.org/x/oauth2/google"
)

// Client contains an authenticated http client
// and a map-based cache for the results of the API
//
// If you reuse the client to make a second call,
// for the same url, it will not call the Cloud Run
// API anymore.
type Client struct {
	client *http.Client
	cache  map[key]string
}

// The key for the cache
type key struct {
	Region, Name string
}

// NewClient initializes the client
func NewClient() (*Client, error) {
	// Build the authenticated http.Client
	// On your local environment it will use the
	// gcloud application-default credentials,
	// on Cloud Run it will query the metadata service to
	// get credentials.
	client, err := google.DefaultClient(context.Background())
	if err != nil {
		return nil, fmt.Errorf("setting up http client: %v", err)
	}
	return &Client{
		client,
		make(map[key]string)}, nil
}

func (c *Client) getServiceUrlDirect(region, name string) (string, error) {
	// Build the URL to call. There is a separate
	// API endpoint for each Cloud Run region.
	apiUrl := fmt.Sprintf("https://%s-run.googleapis.com/"+
		"apis/serving.knative.dev/v1/"+
		"namespaces/%s/"+
		"services/%v",
		region,
		config.ProjectID(),
		name)

	// Call the API and handle errors
	resp, err := c.client.Get(apiUrl)
	if err != nil {
		return "", fmt.Errorf("calling Run API: %v", err)
	}
	// Close the response when the surrounding
	// function returns
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return "", fmt.Errorf("service %s not found "+
			"in region %s, "+
			"project %s",
			name,
			region,
			config.ProjectID)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("status %v from Run API", resp.StatusCode)
	}

	// Parse the JSON-formatted response
	service := serviceItem{}
	json.NewDecoder(resp.Body).Decode(&service)

	if err != nil {
		return "", fmt.Errorf("reading response: %v", err)
	}

	// Populate the cache with the result and return
	serviceUrl := service.Status.Url
	return serviceUrl, nil
}

func (c *Client) DeleteSelf(region, name string) error {
	// Build the URL to call. There is a separate
	// API endpoint for each Cloud Run region.
	apiUrl := fmt.Sprintf("https://%s-run.googleapis.com/"+
		"apis/serving.knative.dev/v1/"+
		"namespaces/%s/"+
		"services/%v",
		region,
		config.ProjectID(),
		name)

	requestURL, _ := url.Parse(apiUrl)
	req := &http.Request{
		Method: "DELETE",
		URL:    requestURL,
	}

	// Call the API and handle errors
	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("calling Run API: %v", err)
	}
	// Close the response when the surrounding
	// function returns
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("service %s not found "+
			"in region %s, "+
			"project %s",
			name,
			region,
			config.ProjectID)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("status %v from Run API", resp.StatusCode)
	}
	return nil

}

// GetServiceUrl calls the Cloud Run management API to
// find the URL of a Cloud Run function.
//
// This is an example URL:
// https://product-rwrmxiaqmq-ew.a.run.app
//
// The format is https://[NAME]-[HASH]-[REGION].a.run.app
// We know the name and region of a function,
// but the hash is determined when a function is first deployed.
func (c *Client) GetServiceUrl(region, name string) (string, error) {

	// If the cache contains the result, return that instead of
	// making an API call
	if value, ok := c.cache[key{region, name}]; ok {
		log.Printf("returning from cache %v, %v => %v",
			region, name, value)
		return value, nil
	}

	// Get the service url from the api
	serviceUrl, err := c.getServiceUrlDirect(region, name)
	if err != nil {
		return "", err
	}

	// Populate the cache with the result and return
	c.cache[key{region, name}] = serviceUrl
	return serviceUrl, nil

}

// serviceItem models the response from the Cloud Run API.
type serviceItem struct {
	Status *serviceStatus `json:status`
}

// serviceStatus contains the URL of the service
type serviceStatus struct {
	Url string `json:url`
}
