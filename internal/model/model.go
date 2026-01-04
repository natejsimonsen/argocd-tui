package model

import (
	"example.com/main/services/argocd"
	"github.com/sirupsen/logrus"
)

type AppModel struct {
	ArgoCDService        *argocd.Service
	Logger               *logrus.Logger
	Applications         []argocd.ApplicationItem
	SelectedAppName      string
	MainFilter           string
	AppFilter            string
	SelectedAppResources []argocd.ApplicationNode
	ScrollOffset         int
	PrevIndex            int
	PrevText             string
}

func NewAppModel(logger *logrus.Logger, svc *argocd.Service) *AppModel {
	return &AppModel{
		ArgoCDService: svc,
		Logger:        logger,
		PrevIndex:     0,
	}
}

func (m *AppModel) LoadApplications() {
	result := m.ArgoCDService.ListApplications()
	m.Applications = result.Items
	m.PrevText = result.Items[0].Metadata.Name
}

func (m *AppModel) LoadResources(appName string) {
	m.SelectedAppName = appName
	m.SelectedAppResources = m.ArgoCDService.GetResourceTree(m.SelectedAppName)
}
