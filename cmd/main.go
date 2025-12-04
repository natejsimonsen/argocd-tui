package main

import (
	"example.com/main/internal/controller"
	"example.com/main/internal/model"
	"example.com/main/internal/view"
	"example.com/main/services/argocd"
	"example.com/main/services/config"
	"example.com/main/services/logger"
	"github.com/rivo/tview"
)

func main() {
	app := tview.NewApplication()
	l := logger.SetupLogger()
	argocdSvc := argocd.NewService(l)
	appModel := model.NewAppModel(l, argocdSvc)
	config := config.NewConfig()
	appView := view.NewAppView(app, config)
	appController := controller.NewAppController(
		appModel,
		appView,
	)

	err := appController.Start()
	if err != nil {
		l.Fatalf("Could not start controller: %v", err)
	}
}
