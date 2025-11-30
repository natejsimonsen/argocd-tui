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

type Service struct {
	Logger *logrus.Logger
	Client *http.Client
	Token  string
}

func NewService(logger *logrus.Logger) *Service {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	client := &http.Client{
		Transport: tr,
		Timeout:   10 * time.Second,
	}

	token, err := Login(*client)
	if err != nil {
		logger.Fatalf("Error logging in: %v", err)
	}

	svc := Service{
		Logger: logger,
		Client: client,
		Token:  token,
	}

	return &svc
}

func (s *Service) Get(path string) (*http.Response, error) {
	req, err := http.NewRequest(
		"GET",
		fmt.Sprintf("%s/api/v1/%s", ARGOCD_SERVER, path),
		nil,
	)
	if err != nil {
		s.Logger.Fatalf("Error creating request: %v", err)
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.Token))
	req.Header.Set("Accept", "application/json")

	resp, err := s.Client.Do(req)
	if err != nil {
		log.Fatalf("Error executing request: %v", err)
		return nil, err
	}

	return resp, nil
}

func Login(client http.Client) (string, error) {
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
		return "", err
	}

	defer response.Body.Close()

	loginResBytes, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %v", err)
		return "", err
	}

	var loginToken LoginToken

	err = json.Unmarshal(loginResBytes, &loginToken)
	if err != nil {
		log.Fatalf("Error unmarshaling json: %v", err)
		return "", err
	}

	return loginToken.Token, nil
}

func (s *Service) ListApplications() ListApplicationsResponse {
	resp, err := s.Get("applications")
	if err != nil {
		log.Fatalf("Could not get applications: %v", err)
	}

	defer resp.Body.Close()

	var result ListApplicationsResponse

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		log.Fatalf("Error decoding json: %v", err)
	}

	debugResult, err := json.Marshal(result)
	s.Logger.Debug(string(debugResult))

	return result
}

func (s *Service) GetResourceTree(application string) []ApplicationNode {
	resp, err := s.Get(fmt.Sprintf("applications/%s/resource-tree", application))
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
