package main

import (
	"fmt"
	"example.com/main/services/argocd"
)

func main() {
	result := argocd.ListApplications()
	for _, item := range result.Items {
		fmt.Println(item.Metadata["name"])
	}
}
