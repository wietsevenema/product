package config

import (
	"cloud.google.com/go/compute/metadata"
	"fmt"
	"log"
	"os/exec"
	"strings"
)

var projectID string
var region string
var functionServiceAccount string

func ProjectID() string {
	if projectID == "" {
		loadProjectID()
	}
	return projectID
}

func loadProjectID() {
	var err error
	if metadata.OnGCE() {
		projectID, err = metadata.ProjectID()
		if err != nil {
			log.Fatal(err)
		}
		return
	}
	projectID, err = gCloud("config", "get-value", "project")
	if err != nil {
		panic(err)
	}
}

func Region() string {
	if region == "" {
		loadRegion()
	}
	return region
}

func loadRegion() {
	var err error
	if metadata.OnGCE() {

		// Cloud Run is a regional resource,
		// the zone is reported with the suffix
		// '-1' instead of the usual '-a', '-b',
		// or '-c'.
		// Example: europe-west1-1
		zone, err := metadata.Zone()
		if err != nil {
			log.Fatal(err)
		}
		region = zone[:len(zone)-2]
		return
	}
	region, err = gCloud("config", "get-value", "run/region")
	if err != nil {
		panic(err)
	}
}

func ServiceAccount() string {
	if functionServiceAccount == "" {
		loadServiceAccount()
	}
	return functionServiceAccount
}

func loadServiceAccount() {
	var err error
	if metadata.OnGCE(){
		functionServiceAccount, err = metadata.Email("default")
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	projectNumber, err := gCloud("projects",
		"describe", ProjectID(), "--format", "value(projectNumber)")
	if err != nil {
		panic(err)
	}

	// Assumes Compute Engine default service account.
	functionServiceAccount = fmt.Sprintf(
		"%s-compute@developer.gserviceaccount.com", projectNumber)

}

func gCloud(args ...string) (string, error) {
	cmd := exec.Command("gcloud", args...)

	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("command failed with %s, %s\n", out, err)
	}
	result := strings.TrimSpace(string(out))
	log.Print("gcloud ", strings.Join(args, " "), " => ", result)
	return result, nil
}