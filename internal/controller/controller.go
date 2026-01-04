package controller

import (
	"strings"

	"example.com/main/internal/model"
	"example.com/main/internal/view"
	"example.com/main/services/argocd"
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
	// AppTable Commands
	c.CommandModel.Add(
		model.KeyStroke{Rune: 'g'},
		model.Global,
		"Goes to the top of the list",
		func(ctx model.Context) {
			c.View.ScrollTo(0)
		},
	)

	c.CommandModel.Add(
		model.KeyStroke{Rune: 'G'},
		model.Global,
		"Goes to the bottom of the list",
		func(ctx model.Context) {
			c.View.ScrollTo(-1)
		},
	)

	c.CommandModel.Add(
		model.KeyStroke{Rune: 'j'},
		model.Global,
		"Scrolls down one row",
		func(ctx model.Context) {
			c.View.Scroll(-1)
		},
	)

	c.CommandModel.Add(
		model.KeyStroke{Rune: 'k'},
		model.Global,
		"Scrolls up one row",
		func(ctx model.Context) {
			c.View.Scroll(1)
		},
	)

	// App Commands
	c.CommandModel.Add(
		model.KeyStroke{Rune: 'q'},
		model.Global,
		"Quits the application",
		func(ctx model.Context) {
			c.View.App.Stop()
		},
	)

	c.CommandModel.Add(
		model.KeyStroke{Rune: '?'},
		model.Global,
		"Toggles the help page",
		func(ctx model.Context) {
			c.View.ToggleHelp(c.CommandModel.Commands)
		},
	)

	c.CommandModel.Add(
		model.KeyStroke{Rune: '/'},
		model.Global,
		"Toggle Search Bar",
		func(ctx model.Context) {
			c.View.ToggleCommandBar()
		},
	)

	// Help Page Commands
	c.CommandModel.Add(
		model.KeyStroke{Key: tcell.KeyEsc},
		model.Help,
		"Exit Search Page",
		func(ctx model.Context) {
			c.View.ToggleHelp(c.CommandModel.Commands)
		},
	)

	// Command Bar Commands
	c.CommandModel.Add(
		model.KeyStroke{Rune: '/'},
		model.CommandBar,
		"Toggle Search Bar",
		func(ctx model.Context) {
			c.View.ToggleCommandBar()
		},
	)

	c.CommandModel.Add(
		model.KeyStroke{Key: tcell.KeyEsc},
		model.CommandBar,
		"Toggle Search Bar",
		func(ctx model.Context) {
			c.View.ToggleCommandBar()
		},
	)

	c.CommandModel.Add(
		model.KeyStroke{Key: tcell.KeyEnter},
		model.CommandBar,
		"Search for substrings in the currently focused pane",
		func(ctx model.Context) {
			c.Model.MainFilter = c.View.SearchInput.GetText()
			c.View.ToggleCommandBar()
			c.View.SetSearchTitle(c.Model.MainFilter)
			c.Model.SelectedAppResources = c.FilterContent(c.Model.MainFilter)
			c.View.UpdateMainContent(c.Model.SelectedAppResources)
		},
	)
}

func (c *AppController) SetupEventHandlers() {
	c.AddCommands()

	c.View.AppTable.SetSelectionChangedFunc(func(row int, col int) {
		name := c.View.AppTable.GetCell(row, col)
		c.Model.LoadResources(name.Text)
		c.View.UpdateMainContent(c.Model.SelectedAppResources)
	})

	// apptable commands
	c.View.AppTable.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyRune {
			if cmd, ok := c.CommandModel.Commands[model.AppTable][model.KeyStroke{Rune: event.Rune()}]; ok {
				cmd.Handler()
				return nil
			}
		}

		return event
	})

	// help page cmds
	c.View.HelpPage.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyRune {
			// runes
			if cmd, ok := c.CommandModel.Commands[model.Help][model.KeyStroke{Rune: event.Rune()}]; ok {
				cmd.Handler()
				return nil
			}
		}

		// keys
		if cmd, ok := c.CommandModel.Commands[model.Help][model.KeyStroke{Key: event.Key()}]; ok {
			cmd.Handler()
			return nil
		}

		return event
	})

	c.View.HelpPage.SetBlurFunc(func() {
		c.View.RemoveHelp()
	})

	// main page cmds
	c.View.MainPage.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyRune {
			if cmd, ok := c.CommandModel.Commands[model.MainPage][model.KeyStroke{Rune: event.Rune()}]; ok {
				cmd.Handler()
				return nil
			}
		}

		return event
	})

	// global cmds
	c.View.App.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if c.View.CommandBar.HasFocus() {
			return event
		}

		if event.Key() == tcell.KeyRune {
			// global commands
			if cmd, ok := c.CommandModel.Commands[model.Global][model.KeyStroke{Rune: event.Rune()}]; ok {
				cmd.Handler()
				return nil
			}
		}

		if event.Key() == tcell.KeyTab {
			if c.View.App.GetFocus() == c.View.MainTable {
				c.View.App.SetFocus(c.View.AppTable)
				return nil
			}
			c.View.App.SetFocus(c.View.MainTable)
			return nil
		}

		if event.Key() == tcell.KeyEsc {
			if c.Model.MainFilter != "" {
				c.Model.MainFilter = ""
				c.View.ClearSearch()
				return nil
			}

			return event
		}

		return event
	})

	// command bar cmds
	c.View.CommandBar.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		// handle runes
		if event.Key() == tcell.KeyRune {
			if cmd, ok := c.CommandModel.Commands[model.CommandBar][model.KeyStroke{Rune: event.Rune()}]; ok {
				cmd.Handler()
				return nil
			}
		}

		// handle tcell.Key
		if cmd, ok := c.CommandModel.Commands[model.CommandBar][model.KeyStroke{Key: event.Key()}]; ok {
			cmd.Handler()
			return event
		}

		// if no event found, bubble event
		return event
	})
}

func (c *AppController) FilterContent(search string) []argocd.ApplicationNode {
	var filteredResources []argocd.ApplicationNode

	for _, app := range c.Model.SelectedAppResources {
		if strings.Contains(app.Name, c.Model.MainFilter) {
			filteredResources = append(filteredResources, app)
		}
	}

	return filteredResources
}

func (c *AppController) Start() error {
	c.SetupEventHandlers()
	c.Model.LoadApplications()
	c.View.UpdateAppTable(c.Model.Applications)
	c.View.App.SetRoot(c.View.Pages, true)
	return c.View.App.Run()
}
