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
	SelectedAppResources []argocd.ApplicationNode
	SelectedIndex        int
}

func NewAppModel(logger *logrus.Logger, svc *argocd.Service) *AppModel {
	return &AppModel{
		ArgoCDService: svc,
		Logger:        logger,
		SelectedIndex: -1,
	}
}

func (m *AppModel) LoadApplications() {
	result := m.ArgoCDService.ListApplications()
	m.Applications = result.Items
}

func (m *AppModel) LoadResources(appName string) {
	m.SelectedAppName = appName
	m.SelectedAppResources = m.ArgoCDService.GetResourceTree(m.SelectedAppName)
}
