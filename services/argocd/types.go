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

type ApplicationItems struct {
	Metadata  map[string]any `json:"metadata"`
	Operation map[string]any `json:"operation"`
	Spec      map[string]any `json:"spec"`
	Status    map[string]any `json:"status"`
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
