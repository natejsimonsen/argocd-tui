package argocd

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

const ARGOCD_SERVER = "http://localhost:8080"

var token string

type Service struct {
	Logger *logrus.Logger
}

func getClient() *http.Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{
		Transport: tr,
		Timeout:   10 * time.Second,
	}

	return client
}

func login() {
	client := getClient()

	loginBody := map[string]string{
		"password": os.Getenv("ARGOCD_PASSWORD"),
		"username": os.Getenv("ARGOCD_USERNAME"),
	}

	jsonLoginBody, err := json.Marshal(loginBody)
	if err != nil {
		log.Fatalf("Error marshaling JSON: %v", err)
	}

	loginBodyReader := bytes.NewBuffer(jsonLoginBody)

	response, err := client.Post(fmt.Sprintf("%s/api/v1/session", ARGOCD_SERVER), "application/json", loginBodyReader)
	if err != nil {
		log.Fatalf("Error performing POST request for exchange token: %v", err)
		return
	}

	defer response.Body.Close()

	loginResBytes, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %v", err)
		return
	}

	var loginToken LoginToken

	err = json.Unmarshal(loginResBytes, &loginToken)
	if err != nil {
		log.Fatalf("Error unmarshaling json: %v", err)
		return
	}

	token = loginToken.Token
}

func (s *Service) ListApplications() ListApplicationsResponse {
	if len(token) == 0 {
		login()
	}

	client := getClient()

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/applications", ARGOCD_SERVER), nil)
	if err != nil {
		log.Fatalf("Error creating request: %var", err)
		return ListApplicationsResponse{}
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error executing request: %v", err)
		return ListApplicationsResponse{}
	}

	defer resp.Body.Close()

	var result ListApplicationsResponse

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		log.Fatalf("Error decoding json: %v", err)
	}

	return result
}

func (s *Service) GetResourceTree(application string) []ApplicationNode {
	if len(token) == 0 {
		login()
	}

	client := getClient()

	req, err := http.NewRequest(
		"GET",
		fmt.Sprintf("%s/api/v1/applications/%s/resource-tree", ARGOCD_SERVER, application),
		nil,
	)
	if err != nil {
		s.Logger.Fatalf("Error making request : %v", err)
		return nil
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Add("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		s.Logger.Fatalf("Error executing request: %v", err)
		return nil
	}

	defer resp.Body.Close()

	var result ResourceTreeResponse

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		s.Logger.Fatalf("Error with json decoder: %v", err)
	}

	return result.Nodes
}

func (s *Service) GetApplicationManifests(application string) []ApplicationManifest {
	if len(token) == 0 {
		login()
	}

	client := getClient()

	req, err := http.NewRequest(
		"GET",
		fmt.Sprintf("%s/api/v1/applications/%s/manifests", ARGOCD_SERVER, application),
		nil,
	)
	if err != nil {
		s.Logger.Fatalf("Error creating request: %v", err)
		return nil
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error executing request: %v", err)
		return nil
	}

	defer resp.Body.Close()

	var manifests ApplicationManifestsResponse

	err = json.NewDecoder(resp.Body).Decode(&manifests)
	if err != nil {
		log.Fatalf("Error decoding json: %v", err)
		return nil
	}

	appManifests := []ApplicationManifest{}

	for _, manifest := range manifests.Manifests {
		s.Logger.Debug(manifest)
		var appManifest ApplicationManifest
		err := json.Unmarshal([]byte(manifest), &appManifest)
		if err != nil {
			s.Logger.Fatalf("Error unmarshaling json: %v", err)
			return nil
		}
		appManifests = append(appManifests, appManifest)
	}

	return appManifests
}
