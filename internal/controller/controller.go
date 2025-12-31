package controller

import (
	"strings"

	"example.com/main/internal/model"
	"example.com/main/internal/view"
	"example.com/main/services/argocd"
	"example.com/main/services/utils"
	"github.com/gdamore/tcell/v2"
)

type AppController struct {
	Model        *model.AppModel
	CommandModel *model.CommandModel
	View         *view.AppView
}

func NewAppController(m *model.AppModel, cm *model.CommandModel, v *view.AppView) *AppController {
	return &AppController{
		Model:        m,
		CommandModel: cm,
		View:         v,
	}
}

func (c *AppController) AddCommands() {
	// AppList Commands
	c.CommandModel.Add(
		'g',
		model.Global,
		"Goes to the top of the list",
		func(ctx model.Context) {
			c.View.ScrollTo(0)
		},
	)

	c.CommandModel.Add(
		'G',
		model.Global,
		"Goes to the bottom of the list",
		func(ctx model.Context) {
			c.View.ScrollTo(-1)
		},
	)

	c.CommandModel.Add(
		'j',
		model.Global,
		"Scrolls down one row",
		func(ctx model.Context) {
			c.View.Scroll(-1)
		},
	)

	c.CommandModel.Add(
		'k',
		model.Global,
		"Scrolls up one row",
		func(ctx model.Context) {
			c.View.Scroll(1)
		},
	)

	// App Commands
	c.CommandModel.Add(
		'q',
		model.App,
		"Quits the application",
		func(ctx model.Context) {
			c.View.App.Stop()
		},
	)

	c.CommandModel.Add(
		'?',
		model.App,
		"Toggles the help page",
		func(ctx model.Context) {
			c.View.ToggleHelp(c.CommandModel.Commands)
		},
	)

	// Help Commands

	// Main Page Commands

	c.CommandModel.Add(
		'l',
		model.MainPage,
		"WIP for horizontal scrolling",
		func(ctx model.Context) {
			c.View.HorizontallyScrollMainTable(1)
		},
	)

	c.CommandModel.Add(
		'h',
		model.MainPage,
		"WIP for horizontal scrolling",
		func(ctx model.Context) {
			c.View.HorizontallyScrollMainTable(-1)
		},
	)

	c.CommandModel.Add(
		'D',
		model.MainPage,
		"Page Down",
		func(ctx model.Context) {
			c.View.PageMainContent(1)
		},
	)

	c.CommandModel.Add(
		'U',
		model.MainPage,
		"Page Up",
		func(ctx model.Context) {
			c.View.PageMainContent(-1)
		},
	)

	c.CommandModel.Add(
		'/',
		model.MainPage,
		"Toggle Search Bar",
		func(ctx model.Context) {
			c.View.ToggleCommandBar()
		},
	)
}

func (c *AppController) SetupEventHandlers() {
	c.AddCommands()

	// application on change
	c.View.AppList.SetChangedFunc(func(index int, mainText, secondaryText string, shortcut rune) {
		appName := utils.StripTags(mainText)

		c.Model.LoadResources(appName)

		c.View.UpdateMainContent(c.Model.SelectedAppResources)
		c.View.UpdateTitles(index, c.Model.PrevIndex, appName, c.Model.PrevText)

		c.Model.PrevText = mainText
		c.Model.PrevIndex = index
	})

	// applist commands
	c.View.AppList.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyRune {
			if cmd, ok := c.CommandModel.Commands[model.AppList][event.Rune()]; ok {
				cmd.Handler()
				return nil
			}
		}

		return event
	})

	c.View.HelpPage.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyRune {
			if cmd, ok := c.CommandModel.Commands[model.Help][event.Rune()]; ok {
				cmd.Handler()
				return nil
			}
		}

		return event
	})

	c.View.SideBar.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyRune {
			if cmd, ok := c.CommandModel.Commands[model.MainPage][event.Rune()]; ok {
				cmd.Handler()
				return nil
			}
		}

		if event.Key() == tcell.KeyPgDn {
			c.View.PageMainContent(1)
			return nil
		}

		if event.Key() == tcell.KeyPgUp {
			c.View.PageMainContent(-1)
			return nil
		}

		if event.Key() == tcell.KeyTab {
			if c.View.App.GetFocus() == c.View.MainTable {
				c.View.App.SetFocus(c.View.AppList)
				return nil
			}
			c.View.App.SetFocus(c.View.MainTable)
			return nil
		}

		if event.Key() == tcell.KeyBacktab {
			c.View.App.SetFocus(c.View.AppList)
			return nil
		}

		if event.Key() == tcell.KeyEsc {
			if c.Model.SearchString != "" {
				c.Model.LoadResources(c.Model.SelectedAppName)
				c.View.SetSearchTitle("")
				c.View.UpdateMainContent(c.Model.SelectedAppResources)
			}

			c.Model.SearchString = ""
			c.View.ClearSearch()
			return nil
		}

		return event
	})

	c.View.App.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyRune {
			// app commands
			if cmd, ok := c.CommandModel.Commands[model.App][event.Rune()]; ok {
				cmd.Handler()
				return nil
			}

			// global commands
			if cmd, ok := c.CommandModel.Commands[model.Global][event.Rune()]; ok {
				cmd.Handler()
				return nil
			}
		}
		return event
	})

	c.View.CommandBar.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyRune {
			if event.Rune() == '/' {
				c.View.ToggleCommandBar()
				return nil
			}
		}

		if event.Key() == tcell.KeyEsc {
			c.View.ToggleCommandBar()
			return nil
		}

		if event.Key() == tcell.KeyEnter {
			c.Model.SearchString = c.View.SearchInput.GetText()
			c.View.ToggleCommandBar()
			c.View.SetSearchTitle(c.Model.SearchString)
			c.Model.SelectedAppResources = c.FilterContent(c.Model.SearchString)
			c.View.UpdateMainContent(c.Model.SelectedAppResources)
			return nil
		}

		return event
	})
}

func (c *AppController) FilterContent(search string) []argocd.ApplicationNode {
	var filteredResources []argocd.ApplicationNode

	for _, app := range c.Model.SelectedAppResources {
		if strings.Contains(app.Name, c.Model.SearchString) {
			filteredResources = append(filteredResources, app)
		}
	}

	return filteredResources
}

func (c *AppController) Start() error {
	c.SetupEventHandlers()
	c.Model.LoadApplications()
	c.View.UpdateAppList(c.Model.Applications)
	c.View.App.SetRoot(c.View.Pages, true)
	return c.View.App.Run()
}
