package argocd

type LoginToken struct {
	Token string `json:"token"`
}

type ListApplicationsResponse struct {
	Items    []ApplicationItems     `json:"items"`
	Metadata map[string]interface{} `json:"metadata"`
}

type ApplicationItems struct {
	Metadata  map[string]interface{} `json:"metadata"`
	Operation map[string]interface{} `json:"operation"`
	Spec      map[string]interface{} `json:"spec"`
	Status    map[string]interface{} `json:"status"`
}
