package argocd

import (
	"time"
)

type LoginToken struct {
	Token string `json:"token"`
}

type ListApplicationsResponse struct {
	Items    []ApplicationItems `json:"items"`
	Metadata map[string]any     `json:"metadata"`
}

type ApplicationMetadata struct {
	Name   string `json:"name"`
	Status string `json:"status"`
}

type ApplicationHealthStatus string

const (
	StatusHealthy     ApplicationHealthStatus = "Healthy"
	StatusMissing     ApplicationHealthStatus = "Missing"
	StatusProgressing ApplicationHealthStatus = "Progressing"
	StatusUnknown     ApplicationHealthStatus = "Unknown"
	StatusDegraded    ApplicationHealthStatus = "Degraded"
)

type ApplicationStatus struct {
	Health struct {
		Status ApplicationHealthStatus `json:"status"`
	} `json:"health"`
}

type ApplicationItems struct {
	Metadata  ApplicationMetadata `json:"metadata"`
	Operation map[string]any      `json:"operation"`
	Spec      map[string]any      `json:"spec"`
	Status    ApplicationStatus   `json:"status"`
}

type ParentRef struct {
	Group     string `json:"group"`
	Version   string `json:"version"`
	Kind      string `json:"kind"`
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
	UID       string `json:"uid"`
}

type InfoItem struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type NetworkingInfo struct {
	Labels map[string]string `json:"labels"`
}

type Health struct {
	Status string `json:"status"`
}

type ResourceTreeResponse struct {
	Nodes []ApplicationNode `json:"nodes"`
}

type ApplicationNode struct {
	Version         string         `json:"version"`
	Kind            string         `json:"kind"`
	Namespace       string         `json:"namespace"`
	Name            string         `json:"name"`
	UID             string         `json:"uid"`
	ParentRefs      []ParentRef    `json:"parentRefs"`
	Info            []InfoItem     `json:"info"`
	NetworkingInfo  NetworkingInfo `json:"networkingInfo"`
	ResourceVersion string         `json:"resourceVersion"`
	Images          []string       `json:"images"`
	Health          Health         `json:"health"`
	CreatedAt       time.Time      `json:"createdAt"`
}
